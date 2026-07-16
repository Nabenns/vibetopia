package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"vibetopia/internal/config"
	"vibetopia/internal/enet"
	"vibetopia/internal/protocol"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[game-server] db connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("[game-server] db ping: %v", err)
	}
	log.Println("[game-server] PostgreSQL connected")

	// Create ENet host
	host, err := enet.NewHost("0.0.0.0", 17091, 128)
	if err != nil {
		log.Fatalf("[game-server] enet host: %v", err)
	}
	defer host.Close()

	// Build protocol registry
	reg := protocol.NewRegistry()

	// Core handlers
	reg[protocol.ActionTankIDName] = func(ctx *protocol.Context) {
		protocol.TankIDName(ctx, db)
	}
	reg[protocol.ActionEnterGame] = func(ctx *protocol.Context) {
		protocol.EnterGame(ctx, db)
	}
	reg[protocol.ActionJoinRequest] = func(ctx *protocol.Context) {
		protocol.JoinRequest(ctx, db)
	}
	reg[protocol.ActionInput] = func(ctx *protocol.Context) {
		protocol.Input(ctx)
	}

	// Stub handlers (will be expanded in later phases)
	reg[protocol.ActionDrop] = func(ctx *protocol.Context) { protocol.Drop(ctx) }
	reg[protocol.ActionTrash] = func(ctx *protocol.Context) { protocol.Trash(ctx) }
	reg[protocol.ActionWrench] = func(ctx *protocol.Context) { protocol.Wrench(ctx) }
	reg[protocol.ActionRefreshItems] = func(ctx *protocol.Context) { protocol.RefreshItemData(ctx) }
	reg[protocol.ActionItemFavorite] = func(ctx *protocol.Context) { protocol.ItemFavorite(ctx) }
	reg[protocol.ActionInventoryFav] = func(ctx *protocol.Context) { protocol.InventoryFav(ctx) }
	reg[protocol.ActionStore] = func(ctx *protocol.Context) { protocol.Store(ctx) }
	reg[protocol.ActionBuy] = func(ctx *protocol.Context) { protocol.Buy(ctx) }
	reg[protocol.ActionSetSkin] = func(ctx *protocol.Context) { protocol.SetSkin(ctx) }
	reg[protocol.ActionRespawn] = func(ctx *protocol.Context) { protocol.Respawn(ctx) }
	reg[protocol.ActionQuit] = func(ctx *protocol.Context) { protocol.Quit(ctx) }
	reg[protocol.ActionInfo] = func(ctx *protocol.Context) { protocol.Info(ctx) }
	reg[protocol.ActionFriends] = func(ctx *protocol.Context) { protocol.Friends(ctx) }

	// Wire receive callback
	host.SetOnReceive(func(peerID uint32, data []byte) {
		reg.Handle(peerID, host, data)
	})

	host.SetOnConnect(func(peerID uint32) {
		log.Printf("[game-server] peer connected: %d", peerID)
	})

	host.SetOnDisconnect(func(peerID uint32) {
		log.Printf("[game-server] peer disconnected: %d", peerID)
	})

	fmt.Printf("VIBETOPIA game-server listening on UDP :%s\n", cfg.GamePort)
	log.Printf("[game-server] %d peers max, protocol v225 ready", 128)

	// Main loop
	for {
		host.Service(0)
	}
}
