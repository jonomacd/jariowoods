package characters

const (
	EnemyZ = 3
)

type EnemyCharacter struct {
	*Character
	Type      string // Colour and bomb direction
	State     string // FallingInitial, FallingLater, Set, Carry,
	LiveState bool   // Live or Dying

}
