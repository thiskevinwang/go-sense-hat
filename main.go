package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/lanzafame/bobblehat/sense/screen"
	"github.com/lanzafame/bobblehat/sense/screen/color"
	"github.com/lanzafame/bobblehat/sense/screen/texture"
	"github.com/lanzafame/bobblehat/sense/stick"
)

func main() {
	fmt.Println("Hello, world")

	fb := screen.NewFrameBuffer()
	fb.SetPixel(0, 5, color.Red)
	fb.SetPixel(1, 2, color.Green)
	fb.SetPixel(3, 4, color.White)
	fb.SetPixel(5, 5, color.Blue)
	fb.SetPixel(7, 7, color.Magenta)
	screen.Draw(fb)

	tx := texture.New(16, 16)
	tx.SetPixel(8, 8, color.White)

	// Capture Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			fmt.Println(sig)
			screen.Clear()
			os.Exit(1)
		}
	}()

	// Process joystick interactions
	input, err := stick.Open("/dev/input/event0")
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case e := <-input.Events:
			switch e.Code {
			case stick.Enter:
				fmt.Println("⏎")
			case stick.Up:
				fmt.Println("↑")
			case stick.Down:
				fmt.Println("↓")
			case stick.Left:
				fmt.Println("←")
			case stick.Right:
				fmt.Println("→")
			}
		}
	}

}
