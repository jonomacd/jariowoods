package main

import (
	//"bufio"
	"fmt"
	"github.com/jonomacd/jariowoods/logic"
	"os"
	"time"
)

func main() {
	b := logic.NewBoard()
	logic.PrintBoard(b)
	logic.FillBoard(b, 0)
	logic.PrintBoard(b)
	logic.DropNew(b, "player1")
	logic.PrintBoard(b)
	cntr := make(chan string)
	go func() {
		in := make([]byte, 40)
		for {
			in = make([]byte, 40)

			os.Stdin.Read(in)

			fmt.Println("cruft", string(in))
			if string(in) == "d" {
				cntr <- "right"
			}
			if string(in) == "a" {
				cntr <- "left"
			}

			time.Sleep(100 * time.Millisecond)
		}

	}()
	logic.RunBoard(b, cntr)
}
