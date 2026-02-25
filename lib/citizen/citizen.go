package citizen

type Role string

const (
	Admin  Role = "admin"
	Peon   Role = "peon"
	Baddy  Role = "baddy"
	Hokage Role = "hokage"
)

type Citizen struct {
	UserID      int64
	Username    string
	Firstname   string
	Lastname    string
	Role        Role
	ChatRole    map[int64]Role
	CommandsSet map[string]bool
	GrokToks    int64
}

func NewCitizen(UID int64, username, firstname, lastname string) Citizen {
	return Citizen{
		UserID:      UID,
		Username:    username,
		Firstname:   firstname,
		Lastname:    lastname,
		GrokToks:    10,
		Role:        Peon,
		ChatRole:    make(map[int64]Role),
		CommandsSet: make(map[string]bool),
	}
}

func (c Citizen) IsAdmin() bool {
	return c.Role == Admin || c.Role == Hokage
}

func (c Citizen) IsBanned() bool {
	return c.Role == Baddy
}
