package citizen

type Statement struct {
	Command string
	State   bool
}

// check if command belongs to user
func Ordeal(command string, user Citizen) bool {
	_, ok := user.CommandsSet[command]
	return ok
}
