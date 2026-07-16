package protocol

// Shared types for protocol handlers.

// TankIDName packet fields
const (
	ActionTankIDName     = "tankIDName"
	ActionEnterGame      = "enter_game"
	ActionJoinRequest    = "join_request"
	ActionInput          = "input"
	ActionDrop           = "drop"
	ActionTrash          = "trash"
	ActionWrench         = "wrench"
	ActionRefreshItems   = "refresh_item_data"
	ActionItemFavorite   = "itemfavourite"
	ActionInventoryFav   = "inventoryfav"
	ActionStore          = "store"
	ActionBuy            = "buy"
	ActionSetSkin        = "setskin"
	ActionRespawn        = "respawn"
	ActionQuit           = "quit"
	ActionInfo           = "info"
	ActionFriends        = "friends"
)

// OnSuperMain packets we send to client
var (
	OnSuperMainStartAcceptLogon = "OnSuperMainStartAcceptLogonHrdxs47254722215a"
	OnRequestWorldSelectMenu    = "OnRequestWorldSelectMenu"
	OnSendToServer              = "OnSendToServer"
	OnConsoleMessage            = "onConsoleMessage"
	OnDialogRequest             = "onDialogRequest"
	OnSetBux                    = "OnSetBux"
	OnSetClothing               = "OnSetClothing"
	OnSetInventorySize          = "OnSetInventorySize"
	OnSetGems                   = "OnSetGems"
	OnSetLevel                  = "OnSetLevel"
	OnSetHome                   = "OnSetHome"
)

// World constants
const (
	DefaultWorld = "START"
	HomeWorld    = "HOME"
	WorldWidth   = 100
	WorldHeight  = 60
)
