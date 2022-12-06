package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Returns false if any matches, true if all unique
func subCheck(sub string) bool {

	// fmt.Printf("Checking if \"%c\" is contained in substring %s\n", sub[0], sub[1:])

	if len(sub) == 1 {
		// fmt.Printf("Unique")
		return true
	}

	if strings.Contains(sub[1:], string(sub[0])) {
		// fmt.Printf("Repeated Character\n\n")
		return false
	} else {
		return subCheck(sub[1:])
	}

}

func pick(line string, numUnique int) {
	for i := range line {

		// Wait until we have numUnique characters
		if i < numUnique-1 {
			continue
		}

		if subCheck(line[i-numUnique+1 : i+1]) {
			fmt.Printf("Token for %d uniques at %d\n", numUnique, i+1)
			break
		}

	}
}

var numUnique = 14

func main() {

	f, err := os.OpenFile("../data/day6a", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {

		line := fileScanner.Text()

		pick(line, 3)
		pick(line, 14)

	}

}
