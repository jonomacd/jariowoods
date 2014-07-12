package main

import (
	"github.com/jonomacd/jariowoods/logic"
	termbox "github.com/nsf/termbox-go"
	"os"
)

func main() {
	b := logic.NewBoard()
	logic.PrintBoard(b)
	logic.FillBoard(b, 0)
	logic.PrintBoard(b)
	logic.DropNew(b, &logic.Player{
		Id: "player1",
	})
	logic.PrintBoard(b)
	cntr := make(chan string)
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	go func() {

		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:

				if ev.Key == termbox.KeyCtrlC {
					termbox.Close()
					os.Exit(0)
				} else if ev.Key == termbox.KeyArrowLeft {

					cntr <- "left"
				} else if ev.Key == termbox.KeyArrowRight {

					cntr <- "right"
				} else if ev.Key == termbox.KeyArrowUp {

					cntr <- "up"
				} else if ev.Key == termbox.KeyArrowDown {

					cntr <- "down"
				} else if ev.Ch == 'a' {

					cntr <- "a"
				} else if ev.Ch == 's' {

					cntr <- "b"
				}
			case termbox.EventError:
				panic(ev.Err)
			}
		}

	}()
	logic.RunBoard(b, cntr)
}
