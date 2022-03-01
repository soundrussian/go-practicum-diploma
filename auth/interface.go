package auth

type Auth interface {
	Register(login string, password string) (*User, error)
	Authenticate(login string, password string) (*User, error)
	AuthToken(user *User) string
}
