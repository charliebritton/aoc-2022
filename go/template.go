package main

import (
	"bufio"
	"os"
)

func part1(fileName string) {

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {

		line := fileScanner.Text()

		// Todo implement me

	}

}

func part2(fileName string) {

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {

		// line := fileScanner.Text()

		// Todo implement me

	}

}

func main() {
	part1("dayNt")
	// part1("dayNa")
	// part2("dayNt")
	// part2("dayNa")
}
