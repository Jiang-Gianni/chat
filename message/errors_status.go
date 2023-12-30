package message

type messageError string

const (
	RoomIDError = messageError("room_id is missing from metadata")
)

var _ (error) = (*messageError)(nil)

func (e messageError) Error() string {
	return string(e)
}
