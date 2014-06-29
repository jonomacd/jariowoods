package logic

import (
	"fmt"
	"math/rand"
	"time"
)

type Cell struct {
	Actors []Actor
	X      int
	Y      int
}

type Actor interface {
	State() string
	SetState(string)

	Type() string
	SetType(string)
}

type DefaultActor struct {
	CurrentState string
	CurrentType  string
}

func (a *DefaultActor) State() string {
	return a.CurrentState
}

func (a *DefaultActor) SetState(state string) {
	a.CurrentState = state

}

func (a *DefaultActor) Type() string {
	return a.CurrentType
}

func (a *DefaultActor) SetType(t string) {
	a.CurrentType = t
}

const (
	Width      = 6
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
			board[ii][jj].Actors = append(board[ii][jj].Actors, &DefaultActor{
				CurrentState: "rest",
				CurrentType:  "monster",
			})
		}
	}

	return board
}

func DropNew(board Board, newType string) Board {
	drop := rand.Intn(Width)

	if len(board[0][drop].Actors) == 0 {
		board[0][drop].Actors = append(board[0][drop].Actors, &DefaultActor{
			CurrentState: "falling",
			CurrentType:  newType,
		})
	}
	// TODO set state to dropping
	// TODO decide what it is I actually want to drop

	return board
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func PrintBoard(board Board) {
	for ii, _ := range board {
		for jj, _ := range board[ii] {
			for _, a := range board[ii][jj].Actors {
				fmt.Printf("%+v   ", a.Type())
			}
			if len(board[ii][jj].Actors) == 0 {
				fmt.Printf("[     ]   ")
			}

		}
		fmt.Println("")
	}

	fmt.Printf("\n____\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
}

func RunBoard(b Board, control chan string) {
	for {
		time.Sleep(300 * time.Millisecond)
		change := false
		for ii := len(b) - 1; ii >= 0; ii-- {
			for jj := len(b[ii]) - 1; jj >= 0; jj-- {
				for kk, _ := range b[ii][jj].Actors {
					if b[ii][jj].Actors[kk].State() == "falling" {
						if ii < len(b)-1 {
							if len(b[ii+1][jj].Actors) <= 0 {
								fmt.Printf("%v, %v, %v\n", ii, jj, kk)
								tmp := b[ii][jj].Actors[kk]
								b[ii][jj].Actors = append(b[ii][jj].Actors[:kk], b[ii][jj].Actors[kk+1:]...)
								b[ii+1][jj].Actors = append(b[ii+1][jj].Actors, tmp)
								change = true
							} else {
								b[ii][jj].Actors[kk].SetState("rest")
							}
						}

					} else {
						b[ii][jj].Actors[kk].SetState("rest")
					}
				}
			}
		}

		select {
		case c := <-control:
			switch c {
			case "right":
				for ii := len(b) - 1; ii >= 0; ii-- {
					for jj := len(b[ii]) - 1; jj >= 0; jj-- {
						for kk, _ := range b[ii][jj].Actors {
							if b[ii][jj].Actors[kk].Type() == "player1" && b[ii][jj].Actors[kk].State() != "falling" {
								if jj < len(b[ii]) {
									if len(b[ii][jj+1].Actors) == 0 {
										tmp := b[ii][jj].Actors[kk]
										b[ii][jj].Actors = append(b[ii][jj].Actors[:kk], b[ii][jj].Actors[kk+1:]...)
										b[ii][jj+1].Actors = append(b[ii][jj+1].Actors, tmp)
									}
								}
							}
						}
					}
				}
			case "left":
			}
		default:
		}

		PrintBoard(b)
		if !change {
			DropNew(b, "monster")
			PrintBoard(b)
		}
	}
}
