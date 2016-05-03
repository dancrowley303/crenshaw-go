package util

import "github.com/nsf/termbox-go"

var width int
var height int
var xPos int
var yPos int

var screenMap = make(map[int][]rune)

func incrementLine() {
	xPos = 0
	if yPos == height-1 {
		for i := 0; i < height-1; i++ {
			tmp := make([]rune, width)
			copy(tmp, screenMap[i+1])
			screenMap[i] = tmp
		}
		screenMap[height-1] = make([]rune, width)
	} else {
		yPos++
	}
}

func init() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	width, height = termbox.Size()

	for i := 0; i < height; i++ {
		screenMap[i] = make([]rune, width)
	}

}

// Read reads a single character from stdin into a rune
func Read() (out rune) {
	for {
		if ev := termbox.PollEvent(); ev.Type == termbox.EventKey {
			switch ev.Key {
			case termbox.KeyCtrlZ:
				out = 0x1A
			case termbox.KeySpace:
				out = 0x20
			case termbox.KeyTab:
				out = 0x09
			case termbox.KeyEnter:
				out = 0x0D
			default:
				out = ev.Ch
			}
			break
		}
	}
	return
}

// WriteBlankLine Writes a blank line to stdout
func WriteBlankLine() {
	Write(string(0x0D))
}

// WriteLine Writes a line of content to stdout
func WriteLine(output string) {
	Write(output + string(0x0D))
}

// Write writes a string to stdout
func Write(output string) {

	for _, r := range output {

		if xPos >= width || r == 0x0D {
			incrementLine()
			if r == 0x0D {
				drawScreen()
				continue
			}
		}

		screenMap[yPos][xPos] = r
		xPos++
	}
	drawScreen()
}

func drawScreen() {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			termbox.SetCell(x, y, screenMap[y][x], termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}
