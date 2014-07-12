package logic

type PlayerLocation struct {
	Id string
	X  int
	Y  int
}

// This is a shortcut to find where the player is, otherwise we
// would need to loop through EVERYTHING!
var PlayerLocations map[string]*PlayerLocation

func init() {
	PlayerLocations = make(map[string]*PlayerLocation)
}

type Player struct {
	Id          string
	Facing      int
	Carrying    []Actor
	PlayerState string
}

func (p *Player) State() string {
	return p.PlayerState
}

func (p *Player) SetState(s string) {
	p.PlayerState = s
}
