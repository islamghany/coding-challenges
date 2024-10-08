package commands


type Commander struct {
	currentPath string
}


func NewCommander() Commander {
	return Commander{}
}