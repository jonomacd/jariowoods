package characters

import (
	"github.com/jonomacd/playjunk/image"
	"github.com/jonomacd/playjunk/object"
	"github.com/skelterjohn/geom"
)

type Character struct {
	CoordC      *geom.Coord
	SizeC       *geom.Rect
	AnimateC    bool
	PanelC      *object.Panel
	ZC          int
	ImageC      *image.Image
	AlphaC      int
	DirtyC      bool
	VisibleC    bool
	IdC         string
	PreviousLoc *geom.Rect
}

func (self *Character) Id() string {
	return self.IdC
}

func (self *Character) Coord() *geom.Coord {
	return self.CoordC
}

func (self *Character) SetCoord(coord *geom.Coord) {

	self.DirtyC = true
	self.PreviousLoc = &geom.Rect{}
	self.PreviousLoc.Min = geom.Coord{X: self.CoordC.X, Y: self.CoordC.Y}
	self.PreviousLoc.Max = geom.Coord{X: self.CoordC.X + self.Size().Width(), Y: self.CoordC.Y + self.Size().Height()}

	self.CoordC = coord
}

func (self *Character) Size() *geom.Rect {
	return self.SizeC
}

func (self *Character) Panel() *object.Panel {
	return self.PanelC
}

func (self *Character) Animate() bool {
	return self.AnimateC
}

func (self *Character) Z() int {
	return self.ZC
}

func (self *Character) Image() *image.Image {
	return self.ImageC
}

func (self *Character) Alpha() int {
	return self.AlphaC
}

func (self *Character) Dirty() bool {

	return self.DirtyC
}

func (self *Character) Equals(o object.Object) bool {
	return o.Coord().Equals(self.Coord()) &&
		o.Size().Equals(self.Size()) &&
		o.Panel().Equals(o.Panel())
}

func (self *Character) Visible() bool {
	return self.VisibleC
}

func (self *Character) Previous() *geom.Rect {
	return self.PreviousLoc
}

func (self *Character) ClearDirty() {
	self.DirtyC = false
}

func (self *Character) AddToPanel(p *object.Panel) {
	self.PanelC = p
}
