package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/lanzafame/bobblehat/sense/screen"
	"github.com/lanzafame/bobblehat/sense/screen/color"
	"github.com/lanzafame/bobblehat/sense/stick"
)

func main() {
	fmt.Println("Hello, world")

	//red
	r := color.New(255, 0, 0)
	//maroon
	m := color.New(100, 0, 0)
	//pink
	p := color.New(255, 192, 203)
	//gray
	g := color.New(80, 80, 80)
	//off
	o := color.Black

	fb := screen.NewFrameBuffer()

	// the grid indices run Top-Right to Bottom-Left
	grid := [8][8]color.Pixel565{
		{o, o, o, o, o, o, o, o},
		{o, g, g, o, o, g, g, o},
		{g, m, m, g, g, p, p, g},
		{g, m, m, r, r, r, p, g},
		{g, m, m, m, r, r, r, g},
		{o, g, m, m, m, r, g, o},
		{o, o, g, m, m, g, o, o},
		{o, o, o, g, g, o, o, o},
	}

	setGrid(grid, fb)
	screen.Draw(fb)

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
				// reset the grid
				fmt.Println("⏎")

				grid = [8][8]color.Pixel565{
					{o, o, o, o, o, o, o, o},
					{o, g, g, o, o, g, g, o},
					{g, m, m, g, g, p, p, g},
					{g, m, m, r, r, r, p, g},
					{g, m, m, m, r, r, r, g},
					{o, g, m, m, m, r, g, o},
					{o, o, g, m, m, g, o, o},
					{o, o, o, g, g, o, o, o},
				}

				setGrid(grid, fb)
				screen.Draw(fb)
			case stick.Up:
				fmt.Println("↑")

				newGrid := shiftGrid(grid, 0, 1)
				setGrid(newGrid, fb)
				screen.Draw(fb)
				grid = newGrid
			case stick.Down:
				fmt.Println("↓")

				newGrid := shiftGrid(grid, 0, -1)
				setGrid(newGrid, fb)
				screen.Draw(fb)
				grid = newGrid
			case stick.Left:
				fmt.Println("←")

				newGrid := shiftGrid(grid, 1, 0)
				setGrid(newGrid, fb)
				screen.Draw(fb)
				grid = newGrid
			case stick.Right:
				fmt.Println("→")
				newGrid := shiftGrid(grid, -1, 0)
				setGrid(newGrid, fb)
				screen.Draw(fb)
				grid = newGrid
			}
		}
	}

}

func setGrid(grid [8][8]color.Pixel565, fb *screen.FrameBuffer) {
	for i, row := range grid {
		for j, color := range row {
			fb.SetPixel(i, j, color)
		}
	}
}

func shiftGrid(grid [8][8]color.Pixel565, x int, y int) [8][8]color.Pixel565 {
	newGrid := grid
	for i, row := range grid {
		for j, _ := range row {
			newGrid[i][j] = grid[(i+x+8)%8][(j+y+8)%8]
		}
	}
	return newGrid
}
