package logic

type PlayerState struct {
	Id string
	X  int
	Y  int
}

var PlayerStates map[string]*PlayerState

func init() {
	PlayerStates = make(map[string]*PlayerState)
}
