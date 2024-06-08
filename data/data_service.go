package data

type DataService interface {
	GetAllUsers() ([]User, error)
	CreatePendingUserSession(userId, id, cookieTokenValue, ipAddress, userAgent string) error
	VerifyPendingUserSession(id, cookieTokenValue string) (PendingUserSession, error)
	CreateUserSession(userId, id, cookieTokenValue string) error
	GetUserByEmail(email string) (User, error)
	DeletePendingUserSession(id string) error
	VerifyUserSession(id, token string) (User, error)
	GetUserById(id string) (User, error)
}
