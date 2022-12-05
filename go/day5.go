package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Take the strings
// Remove anything at index 0, 2-4, 7-9 etc

// Last is top
// E.g.
// [N]
// [Z]
//
//	1
//
// Would parse to 1 = {'Z', 'N'}
type Stack []string

func parseStacks(stacks []string) []Stack {

	// Modify the stacks so we only have the data and none of the surrounding output
	for i, stack := range stacks {
		stack = stack[1:]

		var tmp string
		for i, c := range stack {

			if i%4 == 0 {
				tmp += string(c)
			}

		}

		stacks[i] = tmp
		fmt.Printf("Stack %d: %s\n", i, stacks[i])

	}

	// Count total stacks
	numStacks := len(stacks[len(stacks)-1])
	fmt.Printf("There are %d stacks\n", numStacks)

	// Need a variable to store what we return
	var outStacks = make([]Stack, numStacks)

	// Number of stacks
	for i := 0; i < numStacks; i++ {

		// As we work like a stack, the first element will be the bottom of the stack
		stackIndex := len(stacks) - 2

		// Max number of items in a stack
		for j := 0; j < len(stacks)-1; j++ {

			if stacks[stackIndex][i] == ' ' {
				continue
			} else {
				outStacks[i] = append(outStacks[i], string(stacks[stackIndex][i]))
			}

			stackIndex--

		}

	}

	return outStacks

}

func parseMoveCommand(line string) (int, int, int) {

	data := strings.Split(line, " ")

	num, _ := strconv.Atoi(data[1])
	orig, _ := strconv.Atoi(data[3])
	dest, _ := strconv.Atoi(data[5])

	return num, orig - 1, dest - 1

}

var parsedStacks []Stack
var parsedStacks9001 []Stack

func move(orig int, dest int) {

	fmt.Printf("Before Move: %v\n", parsedStacks)
	tmp := parsedStacks[orig][len(parsedStacks[orig])-1]
	parsedStacks[orig] = parsedStacks[orig][:len(parsedStacks[orig])-1]
	parsedStacks[dest] = append(parsedStacks[dest], tmp)
	fmt.Printf("After Move: %v\n\n", parsedStacks)
}

func moveMult(orig int, dest int, num int) {

	fmt.Printf("MULT Before Move: %v\n", parsedStacks9001)

	// Slice of containers to move
	tmp := parsedStacks9001[orig][len(parsedStacks9001[orig])-num : len(parsedStacks9001[orig])]

	fmt.Printf("Temporary store: %v\n", tmp)

	// Clear off end of initial
	parsedStacks9001[orig] = parsedStacks9001[orig][:len(parsedStacks9001[orig])-num]

	// Add to end of new
	for _, val := range tmp {
		parsedStacks9001[dest] = append(parsedStacks9001[dest], val)
	}
	fmt.Printf("MULT After Move: %v\n\n", parsedStacks9001)

}

func main() {

	f, err := os.OpenFile("../data/day5a", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	var stacks []string
	parseStacksFlag := false

	for fileScanner.Scan() {

		line := fileScanner.Text()

		// Send current data to parseStacks only once
		if line == "" && !parseStacksFlag {
			parseStacksFlag = true

			parsedStacks = parseStacks(stacks)
			parsedStacks9001 = make([]Stack, len(parsedStacks))
			copy(parsedStacks9001, parsedStacks)

			fmt.Printf("Length of parsed stacks: %d\n", len(parsedStacks))
			fmt.Printf("Stacks are: %v\n", parsedStacks)

		}

		// Add to the list until we hit that newline
		if !parseStacksFlag {
			stacks = append(stacks, line)
			continue
		}

		if line == "" {
			continue
		}

		// Start to parse moves for the rest of the time

		numToMove, origin, dest := parseMoveCommand(line)

		fmt.Printf("Moving %d elements from %d to %d\n", numToMove, origin+1, dest+1)
		for i := 0; i < numToMove; i++ {
			move(origin, dest)
		}
		moveMult(origin, dest, numToMove)

	}

	fmt.Printf("Stack rearrangement complete.\n")

	var encMsg string
	for _, stack := range parsedStacks {

		encMsg += stack[len(stack)-1]

	}

	var encMsg2 string
	for _, stack := range parsedStacks9001 {

		encMsg2 += stack[len(stack)-1]

	}

	fmt.Printf("Encoded message is: %s, %s\n", encMsg, encMsg2)

}
