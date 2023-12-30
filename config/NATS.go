package config

import "github.com/nats-io/nats.go"

const NATS_URL = nats.DefaultURL
const NATSUserLogin = "userService.login"
const NATSUserRegister = "userService.register"
const NATSRoomCreate = "roomService.create"

var NATSMessageStreamRoom = func(roomID string) string {
	return "messageService.stream." + roomID
}
