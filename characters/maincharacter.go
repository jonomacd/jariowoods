package characters

type MainCharacter struct {
	*Character
	Holding []*Character
}

func (self *MainCharacter) Carry(c []*Character) error {

	self.Holding = append(self.Holding, c...)
	return nil

}

func (self *MainCharacter) DropOne() *Character {
	if len(self.Holding) == 0 {
		return nil
	}
	if len(self.Holding) == 1 {
		drop := self.Holding[1]
		self.Holding = make([]*Character, 0)
		return drop
	}

	drop := self.Holding[0]
	self.Holding = self.Holding[1:]
	return drop
}

func (self *MainCharacter) DropAll() []*Character {
	var drop []*Character
	copy(self.Holding, drop)
	self.Holding = make([]*Character, 0)
	return drop
}
