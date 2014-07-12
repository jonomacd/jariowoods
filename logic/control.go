package logic

func RightAction(b Board, playerId string) Board {
	return move(b, playerId, 1)
}

func LeftAction(b Board, playerId string) Board {
	return move(b, playerId, -1)
}

func move(b Board, playerId string, direction int) Board {
	x := PlayerLocations[playerId].X
	y := PlayerLocations[playerId].Y

	// Find the player in the cell
	for kk, _ := range b[x][y].Actors {

		player, ok := b[x][y].Actors[kk].(*Player)
		if !ok {
			continue
		}

		//Ensure this is the correct player and they are not falling
		if player.Id != playerId {
			continue
		}

		//check which way we are facing
		if player.Facing != direction {
			b[x][y].Actors[kk].(*Player).Facing *= -1
			return b
		}

		if direction == 1 {
			// Moving right, so check wall on right
			if y >= len(b[x])-1 {
				return moveUp(b, player, x, y, kk, direction, playerId)
			}
		} else if y <= 0 {
			return moveUp(b, player, x, y, kk, direction, playerId)
		}

		// Is the next cell occupied?
		if len(b[x][y+direction].Actors) != 0 {
			return moveUp(b, player, x, y, kk, direction, playerId)
		}

		if player.State() == "falling" || player.State() == "fallingfast" {
			continue
		}

		b[x][y+direction].Actors = append(b[x][y+direction].Actors, b[x][y].Actors[kk])
		b[x][y].Actors = make([]Actor, 0)
		PlayerLocations[playerId].Y += direction
		b[x][y+direction].Actors[kk].(*Player).PlayerState = "rest"
	}

	return b
}

func moveUp(b Board, player *Player, x, y, kk, direction int, playerId string) Board {
	if x == 0 {
		return b
	}
	moveCell := x - 1

	if len(b[moveCell][y].Actors) != 0 {
		return b
	}

	b[moveCell][y].Actors = append(b[moveCell][y].Actors, b[x][y].Actors[kk])
	b[x][y].Actors = make([]Actor, 0)
	PlayerLocations[playerId].X -= 1

	b[moveCell][y].Actors[kk].(*Player).PlayerState = "climbing"

	return b

}

func UpAction(b Board, playerId string) Board {
	x := PlayerLocations[playerId].X
	y := PlayerLocations[playerId].Y

	// Find the player in the cell
	for kk, _ := range b[x][y].Actors {

		player, ok := b[x][y].Actors[kk].(*Player)
		if !ok {
			continue
		}

		//Ensure this is the correct player and they are not falling
		if player.Id != playerId {
			continue
		}

		if x == 0 {
			return b
		}

		check := x - 1
		for {
			if len(b[check][y].Actors) != 0 {
				check--
				if check != 0 {
					continue
				}
			}
			if check == 0 && len(b[check][y].Actors) != 0 {
				b = shiftVertical(b, x, y, kk, check)
				PlayerLocations[playerId].X = check
			} else {
				b = shiftVertical(b, x, y, kk, check+1)
				PlayerLocations[playerId].X = check + 1
			}
			break
		}

	}

	return b

}

func shiftVertical(b Board, x, y, kk, check int) Board {
	if x == 0 {
		return b
	}
	if x == check {
		return b
	}

	tmp := b[x][y].Actors[kk]

	for ii := x; ii > check; ii-- {
		b[ii][y].Actors = b[ii-1][y].Actors
	}

	b[check][y].Actors = make([]Actor, 1)
	b[check][y].Actors[0] = tmp

	return b
}

func DownAction() {

}

func AAction() {

}

func BAction() {

}
