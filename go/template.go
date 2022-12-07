package main

import (
	"bufio"
	"os"
)

func main() {

	f, err := os.OpenFile("../data/dayNt", os.O_RDONLY, 0755)
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
