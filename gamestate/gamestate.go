package gamestate

import (
	//"fmt"
	character "github.com/jonomacd/jariowoods/characters"
	"github.com/jonomacd/playjunk/image"
	"github.com/jonomacd/playjunk/object"
	"github.com/skelterjohn/geom"
	"strconv"
)

const (
	Width      = 6
	Height     = 12
	SquareSize = 56
	FillHeight = 4
)

type Board [][]*Square

type Square struct {
	Contains []object.Object
}

func NewBoard() Board {
	var squares Board
	squares = Board(make([][]*Square, Width))

	for ii, _ := range squares {
		squares[ii] = make([]*Square, Height)
	}
	return squares
}

func (b Board) InitialFill() {
	for ii, _ := range b {
		for jj := 0; jj < FillHeight; jj++ {
			b[ii][jj] = &Square{
				Contains: make([]object.Object, 1),
			}
			mc := &character.MainCharacter{}
			mc.Character = &character.Character{}
			mc.IdC = "monster" + strconv.Itoa(ii) + strconv.Itoa(jj)
			mc.CoordC = &geom.Coord{X: float64(ii * SquareSize), Y: float64(jj * SquareSize)}
			mc.ImageC = image.Images["resources/fakewario.gif"]
			mc.SizeC = &mc.ImageC.Size
			mc.ZC = 1
			mc.PreviousLoc = mc.SizeC
			mc.DirtyC = true
			b[ii][jj].Contains[0] = mc
		}
	}
}

func (b Board) GetObjects() []object.Object {

	objects := make([]object.Object, 0)

	for ii, _ := range b {
		for jj, _ := range b[ii] {
			if b[ii][jj] != nil {
				objects = append(objects, b[ii][jj].Contains...)
			}
		}
	}

	return objects
}
