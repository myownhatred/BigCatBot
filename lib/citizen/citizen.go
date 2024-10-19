package citizen

type Role string

const (
	Admin Role = "admin"
	Peon  Role = "peon"
	Baddy Role = "baddy"
)

type Citizen struct {
	UserID    int64
	Username  string
	Firstname string
	Lastname  string
	Role      Role
}
