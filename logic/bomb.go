package logic

type Bomb struct {
	BombState string
	Color     string
}

func (a *Bomb) State() string {
	return a.BombState
}

func (a *Bomb) SetState(state string) {
	a.BombState = state

}
