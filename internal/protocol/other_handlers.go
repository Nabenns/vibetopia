package protocol

import "vibetopia/internal/enet"

// Stub handlers for the remaining Growtopia actions.
// These will be expanded in later phases (Items, World, Trade, etc.)

func Drop(ctx *Context) {
	ctx.Host.LogOn(ctx.PeerID, "`5Drop: not yet implemented. VIBETOPIA is under construction.")
}

func Trash(ctx *Context) {
	ctx.Host.LogOn(ctx.PeerID, "`5Trash: not yet implemented.")
}

func Wrench(ctx *Context) {
	ctx.Host.LogOn(ctx.PeerID, "`5World editing coming soon.")
}

func RefreshItemData(ctx *Context) {
	// Return empty — prevents client from freezing
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":             "refresh_item_data",
		"refresh_item_data":  "",
	}.Serialize())
}

func ItemFavorite(ctx *Context) {
	ctx.Host.LogOn(ctx.PeerID, "`5Favorites coming soon.")
}

func InventoryFav(ctx *Context) {
	// Return empty favorites list so client doesn't freeze
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":        "inventoryfav",
		"inventoryfav":  "0|0",
	}.Serialize())
}

func Store(ctx *Context) {
	// Send empty store — client expects this
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action": "store",
		"store":  "",
	}.Serialize())
}

func Buy(ctx *Context) {
	ctx.Host.LogOn(ctx.PeerID, "`5Store coming soon.")
}

func SetSkin(ctx *Context) {
	ctx.Host.LogOn(ctx.PeerID, "Skin applied.")
}

func Respawn(ctx *Context) {
	// Send respawn to self
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":  "respawn",
		"respawn": "1",
	}.Serialize())
}

func Quit(ctx *Context) {
	ctx.Host.Disconnect(ctx.PeerID)
}

func Info(ctx *Context) {
	ctx.Host.LogOn(ctx.PeerID, "`5VIBETOPIA v0.1.0 — vibe coded GTPS from scratch")
}

func Friends(ctx *Context) {
	ctx.Host.Send(ctx.PeerID, enet.Packet{
		"action":  "friends",
		"friends": "",
	}.Serialize())
}
