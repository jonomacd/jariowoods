package logic

import (
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

var prune rune = 'p'
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
			board[ii][jj].Actors = append(board[ii][jj].Actors, &Monster{
				MonsterState: "rest",
				Color:        "blue",
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
			newType.(*Monster).Color = "red"
		case *Bomb:
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
				termbox.SetCell(jj, ii, rn, termbox.ColorDefault, termbox.ColorDefault)
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
							}
						}
					}
				}
			}
		default:
		}

		// Input Actions
		prune = 'p'
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
				prune = 'a'
			case "b":
				prune = 'b'
			}

		default:
		}

		// Post Input Actions
		if !dropped {
			DropNew(b, &Monster{})

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
