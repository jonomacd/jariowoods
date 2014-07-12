package logic

var blowUp1 rune = '҉'
var blowUp2 rune = '҈'

type Monster struct {
	MonsterState string
	Color        string
}

func (a *Monster) State() string {
	return a.MonsterState
}

func (a *Monster) SetState(state string) {
	a.MonsterState = state

}
