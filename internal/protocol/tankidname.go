package protocol

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"vibetopia/internal/enet"
)

// TankIDName handles the initial login handshake packet from the client.
// This is the first packet sent after ENet connects, before enter_game.
func TankIDName(ctx *Context, db *sql.DB) {
	tankName := ctx.Packet.Get(ActionTankIDName)
	tankPass := ctx.Packet.Get("tankIDPass")
	requestedName := ctx.Packet.Get("requestedName")

	if tankName == "" {
		tankName = requestedName
	}
	if tankName == "" {
		tankName = "guest"
	}

	log.Printf("[tankIDName] peer=%d growID=%s", ctx.PeerID, tankName)
	ctx.Player.GrowID = tankName

	// Look up player in DB
	var playerID int64
	var level, gems int
	err := db.QueryRow(
		"SELECT id, level, gems FROM players WHERE LOWER(growid) = LOWER($1)",
		tankName,
	).Scan(&playerID, &level, &gems)
	if err != nil {
		log.Printf("[tankIDName] player not found: %s, using defaults", tankName)
		level = 1
		gems = 0
	}

	// Get inventory count
	var invCount int
	db.QueryRow("SELECT COUNT(*) FROM player_inventory WHERE player_id = $1", playerID).Scan(&invCount)

	// Send OnSuperMainStartAcceptLogon — the initial setup packet
	startLogon := enet.Packet{
		"action":  OnSuperMainStartAcceptLogon,
		"tankIDName": tankName,
		"tankIDPass": tankPass,
		"country":    "id",
		"user":       tankName,
		"_token":     ctx.Packet.Get("_token"),
		"protocol":   "225",
		"server":     "VIBETOPIA",
		"version":    "5.50",
		"meta":       "vibetopia",
	}
	ctx.Host.Send(ctx.PeerID, startLogon.Serialize())

	// Auto-trigger enter_game (same logic as Gurotopia auto-enter patch)
	// This bypasses the need for client to send action|enter_game from WebView
	EnterGame(ctx, db)
}

// EnterGame handles the post-login setup — sends world select menu and player state.
func EnterGame(ctx *Context, db *sql.DB) {
	growID := ctx.Player.GrowID
	log.Printf("[enter_game] peer=%d growID=%s", ctx.PeerID, growID)

	var playerID int64
	var level, gems int
	var displayName string
	err := db.QueryRow(
		"SELECT id, level, gems, display FROM players WHERE LOWER(growid) = LOWER($1)",
		growID,
	).Scan(&playerID, &level, &gems, &displayName)
	if err != nil {
		log.Printf("[enter_game] player not found: %s", growID)
		level = 1
		gems = 100 // Starter gems
		displayName = growID
	}

	// Update last login
	db.Exec("UPDATE players SET last_login = NOW() WHERE id = $1", playerID)

	// 1. Console message
	ctx.Host.LogOn(ctx.PeerID, "Welcome to VIBETOPIA, `"+displayName+"`!")

	// 2. Send inventory
	sendInventoryState(ctx, db, playerID)

	// 3. Currency (gems)
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":  OnSetBux,
		"OnSetBux": formatInt(gems),
	}.Serialize())

	// 4. SetHasGrowID
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":       "SetHasGrowID",
		"SetHasGrowID": growID,
	}.Serialize())

	// 5. Today's date
	now := time.Now()
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":       "OnTodaysDate",
		"OnTodaysDate": now.Format("02/01/2006"),
	}.Serialize())

	// 6. WORLD SELECT MENU — THE critical packet
	// Without this, client stays at "Connecting..."
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":                   OnRequestWorldSelectMenu,
		OnRequestWorldSelectMenu:   "0",
	}.Serialize())

	// 7. Gazette (news)
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":          "OnRequestGazette",
		"OnRequestGazette": "",
	}.Serialize())

	// 8. Feature flags
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":                  "OnSetFeatureEnableFlags",
		"OnSetFeatureEnableFlags": "0",
	}.Serialize())

	// 9. Ping
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action": "PACKET_PING_REQUEST",
		"hash":   "1",
		"meta":   "vibetopia",
	}.Serialize())
}

// JoinRequest handles the world join action.
func JoinRequest(ctx *Context, db *sql.DB) {
	worldName := ctx.Packet.Get("worldName")
	if worldName == "" {
		worldName = ctx.Packet.Get("name")
	}
	if worldName == "" {
		worldName = DefaultWorld
	}

	log.Printf("[join_request] peer=%d world=%s", ctx.PeerID, worldName)
	ctx.Player.World = worldName
	ctx.Player.InWorld = true

	// Send OnSendToServer — game join success
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":         OnSendToServer,
		"port":           "17091",
		"server":         "VIBETOPIA",
		"OnSendToServer": "1",
	}.Serialize())
}

// sendInventoryState sends the player's inventory to the client.
func sendInventoryState(ctx *Context, db *sql.DB, playerID int64) {
	rows, err := db.Query(
		"SELECT item_id, quantity FROM player_inventory WHERE player_id = $1 ORDER BY item_id",
		playerID,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	var items []string
	for rows.Next() {
		var itemID, qty int
		rows.Scan(&itemID, &qty)
		items = append(items, formatInt(itemID)+"|"+formatInt(qty))
	}

	invStr := ""
	if len(items) > 0 {
		invStr = strings.Join(items, "|")
	}

	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":               "send_inventory_state",
		"send_inventory_state": invStr,
	}.Serialize())

	// Backpack size
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":               OnSetInventorySize,
		OnSetInventorySize:     formatInt(len(items) + 32), // 32 extra slots
	}.Serialize())
}

// Input handles chat messages and commands.
func Input(ctx *Context) {
	text := ctx.Packet.Get("text")
	if text == "" {
		return
	}
	log.Printf("[input] peer=%d text=%s", ctx.PeerID, text)

	if len(text) > 0 && text[0] == '/' {
		handleCommand(ctx, text[1:])
	} else {
		// Broadcast chat to all players in world
		ctx.Host.SendAll(enet.Packet{
			"action":     "OnTalkBubble",
			"text":       text,
			"name":       ctx.Player.GrowID,
			"country":    "id",
			"append":     "0",
			"set_state":  "0",
			"color":      "`o",
		}.Serialize())
	}
}

func handleCommand(ctx *Context, cmd string) {
	switch {
	case cmd == "help":
		ctx.Host.LogOn(ctx.PeerID, "Commands: /help, /pos, /skin, /id")
	case cmd == "pos":
		ctx.Host.LogOn(ctx.PeerID, "You're in world `"+ctx.Player.World+"`")
	default:
		ctx.Host.LogOn(ctx.PeerID, "Unknown command: `"+cmd+"`. Try /help")
	}
}

func formatInt(n int) string {
	return fmt.Sprintf("%d", n)
}
