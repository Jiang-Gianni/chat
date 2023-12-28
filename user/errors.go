package user

type userError string

const (
	UsernameTakenError      = userError("username already taken")
	InvalidCredentialsError = userError("invalid credentials")
)

var _ error = (*userError)(nil)

func (e userError) Error() string {
	return string(e)
}
