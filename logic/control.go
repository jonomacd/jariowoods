package logic

func RightAction(b Board, playerId string) Board {

	x := PlayerStates[playerId].X
	y := PlayerStates[playerId].Y

	for kk, _ := range b[x][y].Actors {
		if b[x][y].Actors[kk].Type() == "player1" && b[x][y].Actors[kk].State() != "falling" {
			if y < len(b[x])-1 {
				if len(b[x][y+1].Actors) == 0 {
					tmp := b[x][y].Actors[kk]
					b[x][y].Actors = append(b[x][y].Actors[:kk], b[x][y].Actors[kk+1:]...)
					b[x][y+1].Actors = append(b[x][y+1].Actors, tmp)
					PlayerStates[playerId].Y++
				}
			}
		}

	}

	return b
}

func LeftAction(b Board, playerId string) Board {

	x := PlayerStates[playerId].X
	y := PlayerStates[playerId].Y

	for kk, _ := range b[x][y].Actors {
		if b[x][y].Actors[kk].Type() == "player1" && b[x][y].Actors[kk].State() != "falling" {
			if y > 0 {
				if len(b[x][y-1].Actors) == 0 {

					b[x][y-1].Actors = append(b[x][y-1].Actors, b[x][y].Actors[kk])
					b[x][y].Actors = make([]Actor, 0)
					PlayerStates[playerId].Y--
				}
			}
		}
	}

	return b
}

func UpAction() {

}

func DownAction() {

}

func AAction() {

}

func BAction() {

}
