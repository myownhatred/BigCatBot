package citizen

type Role string

const (
	Admin Role = "admin"
	Peon  Role = "peon"
	Baddy Role = "baddy"
)

type Citizen struct {
	UserID      int64
	Username    string
	Firstname   string
	Lastname    string
	Role        Role
	ChatRole    map[int64]Role
	CommandsSet map[string]bool
}
