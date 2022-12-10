package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var registers = make(map[rune]int, 0)
var signalStrength = make(map[rune]int, 0)
var clk = 1

func main() {

	f, err := os.OpenFile("../data/day10a", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	// x is initially 1
	registers['x'] = 1

	printPositionOfSprite()

	for fileScanner.Scan() {

		line := fileScanner.Text()

		tokens := strings.Split(line, " ")

		if tokens[0] == "noop" {

			fmt.Printf("Start cycle  % 2d: begin executing noop\n", clk)

			// Draw during cycle
			drawCRT()

			fmt.Printf("End of cycle % 2d: finish executing %s\n", clk, line)
			clk++

		} else {

			// We read the instruction between the start of the cycle
			i, _ := strconv.Atoi(tokens[1])
			fmt.Printf("Start cycle  % 2d: begin executing %s\n", clk, line)

			drawCRT()

			clk++
			fmt.Println()

			drawCRT()
			registers['x'] += i

			fmt.Printf("End of cycle % 2d: finish executing %s (Register X is now %d)\n", clk, line, registers['x'])

			clk++
		}

		printPositionOfSprite()

	}

	renderImage()

	// registers['x'] += opInQueue

	// fmt.Printf("\nRegisters at end: %v, signal strength sum: %v", registers, signalStrength)

	// for i := 0; i <= 41; i++ {
	// 	registers['x'] = i
	// 	printPositionOfSprite()
	// }

}

var crtString = ""

func drawCRT() {

	pos := (clk - 1) % 40

	startRowIndex := clk - pos

	if drawBright(pos) {
		crtString += "#"
	} else {
		crtString += "."
	}

	currentRow := crtString[startRowIndex-1 : startRowIndex+pos]

	fmt.Printf("During cycle % 2d: CRT draws pixel in position %d\nCurrent CRT row: %s (len: %d)\n", clk, pos, currentRow, len(currentRow))

}

func modCycle() {

	printPositionOfSprite()
	if (clk-20)%40 == 0 {

		for k, v := range registers {

			signalStrength[k] += v * clk
			// fmt.Printf("==> %03d-20 %% 40 == 0, calc signalStrength to be %d * %d = %d\n", clk, clk, v, clk*v)
		}

	} else {
		// fmt.Printf("==> Condition not met, just increasing the clock\n")
	}

	clk++
}

const WIDTH = 40

func printPositionOfSprite() {

	pos := registers['x']

	numHash := 3

	// x = 1  ###
	// x = 0  ##
	// x = -1 #
	// x = -2 or less nothing

	// x = 37 #x#.|
	// x = 38  #x#|
	// x = 39   #x|#
	// x = 40    #|x#

	// None at all
	if pos < -2 || pos >= WIDTH+1 {
		numHash = 0
	} else if pos == -1 || pos == WIDTH {
		numHash = 1
	} else if pos == 0 || pos == WIDTH-1 {
		numHash = 2
	}

	stringBuf := "Sprite position: "

	toLeft := pos - 1
	if toLeft < 0 {
		toLeft = 0
	}
	stringBuf += strings.Repeat(".", toLeft)

	stringBuf += strings.Repeat("#", numHash)

	toRight := WIDTH - numHash - toLeft
	if toRight < 0 {
		toRight = 0
	}

	stringBuf += strings.Repeat(".", toRight)

	stringBuf += fmt.Sprintf(" (x = %d)\n", registers['x'])

	fmt.Println(stringBuf)

}

func drawBright(pos int) bool {

	currentX := registers['x']

	difference := int(math.Abs(float64(currentX) - float64(pos)))

	if difference < 2 {
		return true
	} else {
		return false
	}

}

func renderImage() {

	for i, c := range crtString {

		if i%40 == 0 {
			fmt.Println()
		}
		fmt.Printf("%c", c)

	}

}
