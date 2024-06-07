package data

type User struct {
	Id         string
	Email      string
	IsAdmin    bool
	InsertedAt string
}

type PendingUserSession struct {
	Id                string
	UserId            string
	HashedCookieToken string
	Salt              string
	IpAddress         string
	UserAgent         string
	InsertedAt        string
}

type UserSession struct {
	Id                string
	UserId            string
	HashedCookieToken string
	Salt              string
	InsertedAt        string
}

type Org struct {
	Id         string
	Name       string
	InsertedAt string
}

type OrgUser struct {
	OrgId  string
	UserId string
	Role   string
}

type OrgGroup struct {
	Id         string
	OrgId      string
	Name       string
	InsertedAt string
}

type OrgGroupUser struct {
	OrgGroupId string
	UserId     string
	Role       string
}

type OrgGroupMessage struct {
	OrgGroupId string
	UserId     string
	Content    string
	InsertedAt string
}
