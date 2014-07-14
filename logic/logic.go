package logic

import (
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

var FrameRate time.Duration = time.Millisecond * 16

type Cell struct {
	Actors []Actor
	X      int
	Y      int
}

type Actor interface {
	State() string
	SetState(string)
}

const (
	Width      = 16
	Height     = 12
	SquareSize = 56
)

var FillHeight int = 4

type Board [][]*Cell

func NewBoard() Board {
	board := make([][]*Cell, Height)
	for ii, _ := range board {
		board[ii] = make([]*Cell, Width)
		for jj, _ := range board[ii] {
			board[ii][jj] = &Cell{
				X:      jj,
				Y:      ii,
				Actors: make([]Actor, 0),
			}
		}
	}

	return board
}

func FillBoard(board Board, fillHeight int) Board {

	if fillHeight <= 0 {
		fillHeight = Height - FillHeight
	} else {
		fillHeight = Height - fillHeight
	}

	if fillHeight > Height {
		fillHeight = 0
	}

	for ii := Height - 1; ii >= fillHeight; ii-- {
		for jj, _ := range board[ii] {
			color := "red"
			switch rand.Intn(4) {
			case 0:
				color = "green"
			case 1:
				color = "blue"
			case 2:
				color = "yellow"
			}
			board[ii][jj].Actors = append(board[ii][jj].Actors, &Monster{
				MonsterState: "rest",
				Color:        color,
			})
		}
	}

	return board
}

func DropNew(board Board, newType Actor) Board {
	drop := rand.Intn(Width)

	if len(board[0][drop].Actors) == 0 {
		switch v := newType.(type) {
		case *Monster:
			color := "red"
			switch rand.Intn(4) {
			case 0:
				color = "green"
			case 1:
				color = "blue"
			case 2:
				color = "yellow"
			}
			newType.(*Monster).Color = color
		case *Bomb:
			color := "red"
			switch rand.Intn(4) {
			case 0:
				color = "green"
			case 1:
				color = "blue"
			case 2:
				color = "yellow"
			}
			newType.(*Bomb).Color = color
		case *Player:
			PlayerLocations[newType.(*Player).Id] = &PlayerLocation{
				Id: v.Id,
				X:  0,
				Y:  drop,
			}
			newType.(*Player).Facing = 1
		}
		newType.SetState("falling")

		board[0][drop].Actors = append(board[0][drop].Actors, newType)
	}

	return board
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func PrintBoard(board Board) {

	for ii, _ := range board {
		for jj, _ := range board[ii] {
			for _, a := range board[ii][jj].Actors {
				rn := '*'
				color := termbox.ColorDefault
				if p, ok := a.(*Player); ok {
					if p.Facing == 1 {
						rn = 'p'
						if p.State() == "climbing" {
							rn = '౿'
						}
					} else if p.Facing == -1 {
						rn = '9'
						if p.State() == "climbing" {
							rn = 'ڡ'
						}
					}
					if p.State() == "fallingfast" {
						rn = 'ٽ'
					}
				}
				if m, ok := a.(*Monster); ok {
					switch m.Color {
					case "red":
						color = termbox.ColorRed
					case "blue":
						color = termbox.ColorBlue
					case "yellow":
						color = termbox.ColorCyan
					case "green":
						color = termbox.ColorGreen
					}
				}

				if m, ok := a.(*Bomb); ok {
					switch m.Color {
					case "red":
						color = termbox.ColorRed
					case "blue":
						color = termbox.ColorBlue
					case "yellow":
						color = termbox.ColorCyan
					case "green":
						color = termbox.ColorGreen
					}
					rn = 'o'
				}

				termbox.SetCell(jj, ii, rn, color, termbox.ColorDefault)
			}
			if len(board[ii][jj].Actors) == 0 {
				termbox.SetCell(jj, ii, '-', termbox.ColorDefault, termbox.ColorDefault)

			}

		}
		termbox.Flush()
	}

}

func RunBoard(b Board, control chan string) {
	monsterFallTick := time.NewTicker(500 * time.Millisecond)
	fastFallTick := time.NewTicker(100 * time.Millisecond)
	dropped := true
	sinceLastFrame := time.Now()
	slow := true
	for {

		// Pre Input actions
		// check for blowups (TODO this does not work!)
		checkAllForExplosions(b)

		select {
		case <-monsterFallTick.C:
			drop := false
			for ii := len(b) - 1; ii >= 0; ii-- {
				for jj := len(b[ii]) - 1; jj >= 0; jj-- {
					for kk, _ := range b[ii][jj].Actors {
						if b[ii][jj].Actors[kk].State() == "falling" {
							if ii < len(b)-1 {
								if len(b[ii+1][jj].Actors) <= 0 {
									if _, ok := b[ii][jj].Actors[kk].(*Player); ok {
										PlayerLocations[b[ii][jj].Actors[kk].(*Player).Id].X = ii + 1
									}
									tmp := b[ii][jj].Actors[kk]
									b[ii][jj].Actors = append(b[ii][jj].Actors[:kk], b[ii][jj].Actors[kk+1:]...)
									b[ii+1][jj].Actors = append(b[ii+1][jj].Actors, tmp)
									drop = true
								} else {
									b[ii][jj].Actors[kk].SetState("rest")
								}
							} else {
								b[ii][jj].Actors[kk].SetState("rest")
							}
						}
					}
				}
			}
			if !drop {
				dropped = false
			}
		case <-fastFallTick.C:
			// TODO CLEAN THIS UP
			for ii := len(b) - 1; ii >= 0; ii-- {
				for jj := len(b[ii]) - 1; jj >= 0; jj-- {
					for kk, _ := range b[ii][jj].Actors {
						if b[ii][jj].Actors[kk].State() == "rest" {
							if ii < len(b)-1 {
								if len(b[ii+1][jj].Actors) == 0 {
									b[ii][jj].Actors[kk].SetState("fallingfast")
								}
							}
						} else if b[ii][jj].Actors[kk].State() == "climbing" {
							b[ii][jj].Actors[kk].SetState("limbo")
						} else if b[ii][jj].Actors[kk].State() == "limbo" {
							b[ii][jj].Actors[kk].SetState("fallingfast")
						}

						if b[ii][jj].Actors[kk].State() == "fallingfast" {
							if ii < len(b)-1 {
								if len(b[ii+1][jj].Actors) == 0 {
									if _, ok := b[ii][jj].Actors[kk].(*Player); ok {
										PlayerLocations[b[ii][jj].Actors[kk].(*Player).Id].X = ii + 1
									}
									tmp := b[ii][jj].Actors[kk]
									b[ii][jj].Actors = append(b[ii][jj].Actors[:kk], b[ii][jj].Actors[kk+1:]...)
									b[ii+1][jj].Actors = append(b[ii+1][jj].Actors, tmp)
								} else {
									b[ii][jj].Actors[kk].SetState("rest")
								}
							} else {
								b[ii][jj].Actors[kk].SetState("rest")
							}
						}
					}
				}
			}
		default:
		}

		// Input Actions
		select {
		case c := <-control:
			switch c {
			case "right":
				RightAction(b, "player1")
			case "left":
				LeftAction(b, "player1")
			case "up":
				UpAction(b, "player1")
			case "down":
				if slow {
					monsterFallTick.Stop()
					monsterFallTick = time.NewTicker(100 * time.Millisecond)
					slow = false
				} else {
					monsterFallTick.Stop()
					monsterFallTick = time.NewTicker(500 * time.Millisecond)
					slow = true
				}
			case "a":
				pickUpOne(b, "player1")
			case "b":
				pickUpAll(b, "player1")
			}

		default:
		}

		// Post Input Actions
		if !dropped {
			n := rand.Intn(10)
			if n >= 4 {
				DropNew(b, &Monster{})
			} else {
				DropNew(b, &Bomb{})
			}
			dropped = true
		}

		// DRAW
		DelayDraw(sinceLastFrame)
		PrintBoard(b)
		now := time.Now()

		sinceLastFrame = now
	}
}

func DelayDraw(sinceLastFrame time.Time) {
	now := time.Now()
	nextFrame := sinceLastFrame.Add(FrameRate)
	if !nextFrame.Before(now) {
		// Wait until time to draw
		time.Sleep(nextFrame.Sub(now))
	}
}

func checkAllForExplosions(b Board) {
	for xx := 0; xx < len(b); xx++ {
		for yy := 0; yy < len(b[xx]); yy++ {
			for kk, _ := range b[xx][yy].Actors {
				if bomb, ok := b[xx][yy].Actors[kk].(*Bomb); ok {
					c := bomb.Color

					// set up to explode structure
					toExplode := make([][]Cell, 4)
					for ll, _ := range toExplode {
						toExplode[ll] = make([]Cell, 1)
						toExplode[ll][0] = Cell{
							X: xx,
							Y: yy,
						}
					}
					for direction := 0; direction < 8; direction++ {
						toExplode = checkLineForExplosions(b, xx, yy, direction, toExplode, c)
					}
					for direction, _ := range toExplode {
						if len(toExplode[direction]) >= 3 {
							for _, cell := range toExplode[direction] {
								b[cell.X][cell.Y].Actors = make([]Actor, 0)
							}
						}
					}

				}
			}
		}
	}
}

func checkLineForExplosions(b Board, xx, yy, direction int, toExplode [][]Cell, c string) [][]Cell {

	ym := yy
	ymStopFunc := func(checky int) bool {
		return true
	}
	ymUpdateFunc := func(updatey int) int {
		return updatey
	}

	xm := xx
	xmStopFunc := func(checkx int) bool {
		return true
	}
	xmUpdateFunc := func(updatex int) int {
		return updatex
	}

	if direction == 0 {
		// Left
		ym = yy - 1
		ymStopFunc = func(checky int) bool {
			return checky >= 0
		}
		ymUpdateFunc = func(updatey int) int {
			return updatey - 1
		}

	} else if direction == 1 {
		// Right
		ym = yy + 1
		ymStopFunc = func(checky int) bool {
			return checky < len(b[0])
		}
		ymUpdateFunc = func(updatey int) int {
			return updatey + 1
		}
	} else if direction == 2 {
		// Up
		xm = xx - 1
		xmStopFunc = func(checkx int) bool {
			return checkx >= 0
		}
		xmUpdateFunc = func(updatex int) int {
			return updatex - 1
		}
	} else if direction == 3 {
		// Down
		xm = xm + 1
		xmStopFunc = func(checkx int) bool {
			return checkx < len(b)
		}
		xmUpdateFunc = func(updatex int) int {
			return updatex + 1
		}
	} else if direction == 4 {
		// Left up
		ym = yy - 1
		xm = xx - 1
		ymStopFunc = func(checky int) bool {
			return checky >= 0
		}
		ymUpdateFunc = func(updatey int) int {
			return updatey - 1
		}
		xmStopFunc = func(checky int) bool {
			return checky >= 0
		}
		xmUpdateFunc = func(updatey int) int {
			return updatey - 1
		}

	} else if direction == 5 {
		// Right Up
		ym = yy + 1
		xm = xx + 1
		ymStopFunc = func(checky int) bool {
			return checky < len(b[0])
		}
		ymUpdateFunc = func(updatey int) int {
			return updatey + 1
		}
		xmStopFunc = func(checky int) bool {
			return checky < len(b)
		}
		xmUpdateFunc = func(updatey int) int {
			return updatey + 1
		}
	} else if direction == 6 {
		// Up
		xm = xx - 1
		ym = yy + 1
		xmStopFunc = func(checkx int) bool {
			return checkx >= 0
		}
		xmUpdateFunc = func(updatex int) int {
			return updatex - 1
		}

		ymStopFunc = func(checky int) bool {
			return checky < len(b[0])
		}
		ymUpdateFunc = func(updatey int) int {
			return updatey + 1
		}
	} else if direction == 7 {
		// Down
		xm = xm + 1
		ym = yy - 1
		xmStopFunc = func(checkx int) bool {
			return checkx < len(b)
		}
		xmUpdateFunc = func(updatex int) int {
			return updatex + 1
		}
		ymStopFunc = func(checky int) bool {
			return checky >= 0
		}
		ymUpdateFunc = func(updatey int) int {
			return updatey - 1
		}
	}

	direction = (direction / 2)

loop:
	for {

		if !xmStopFunc(xm) || !ymStopFunc(ym) {
			break loop
		}

		if len(b[xm][ym].Actors) > 0 {
			for nn, _ := range b[xm][ym].Actors {
				switch b[xm][ym].Actors[nn].(type) {
				case *Bomb:
					if b[xm][ym].Actors[nn].(*Bomb).Color == c {
						toExplode[direction] = append(toExplode[direction], Cell{
							X: xm,
							Y: ym,
						})
					} else {
						break loop
					}
				case *Monster:
					if b[xm][ym].Actors[nn].(*Monster).Color == c {
						toExplode[direction] = append(toExplode[direction], Cell{
							X: xm,
							Y: ym,
						})
					} else {
						break loop
					}
				default:
					break loop
				}
			}
		} else {
			break loop
		}

		xm = xmUpdateFunc(xm)
		ym = ymUpdateFunc(ym)

	}

	return toExplode
}
