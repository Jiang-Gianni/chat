package config

import "strconv"

const (
	ChatEndpoint              = "/chat"
	ChatParamEndpoint         = "/chat/{roomID}"
	ChatRedirectParamEndpoint = "/chat/redirect/{roomID}"
	ChatWsEndParampoint       = "/chat/{roomID}/ws"
	LoginEndpoint             = "/login"
	RegisterEndpoint          = "/register"
	RoomEndpoint              = "/room"
	DiscardEndpoint           = "/discard"
	DeniedEndpoint            = "/denied"
	LogoutEndpoint            = "/logout"
	IndexEndpoint             = "/"
)

var ChatRoomIDEndpoint = func(roomID int) string {
	if roomID > 0 {
		return "/chat/" + strconv.Itoa(roomID)
	}
	return "/chat"
}

var ChatRedirectRoomIDEndpoint = func(roomID int) string {
	if roomID > 0 {
		return "/chat/redirect/" + strconv.Itoa(roomID)
	}
	return "/chat"
}

var ChatWsEndpoint = func(roomID int) string {
	if roomID > 0 {
		return "/chat/" + strconv.Itoa(roomID) + "/ws"
	}
	return "/chat"
}
