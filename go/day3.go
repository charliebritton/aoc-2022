package main

import (
	"bufio"
	"fmt"
	"os"
)

func getPriority(r rune) int {

	// Other way round to the ASCII chart, lowercase comes first
	if int(r) > int('Z') {
		return int(r) - 96
	} else {
		return int(r) - 64 + 26
	}

}

func main() {

	f, err := os.OpenFile("../data/day3a", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	var total int
	var group3 []string
	var counter int
	var groupName rune
	var totalGrp int

	for fileScanner.Scan() {

		line := fileScanner.Text()

		lineFirst := line[:len(line)/2]
		lineSecond := line[len(line)/2:]

		group3 = append(group3, line)
		fmt.Printf("Counter at %d\n", counter)
		counter++
		if counter%3 == 0 && counter != 0 {

			fmt.Printf("Firing w/ %v\n", group3)

			for _, r := range group3[0] {

				for _, s := range group3[1] {

					for _, t := range group3[2] {

						if r == s && r == t {
							groupName = r
						}
					}

				}

			}

			group3 = []string{}
			fmt.Printf("==> Group Name is %c\n\n", groupName)
			totalGrp += getPriority(groupName)

		}

		var common rune

		for _, r := range lineFirst {

			for _, s := range lineSecond {

				if r == s {
					common = r
				}

			}

		}

		total += getPriority(common)
		fmt.Printf("First Half:  %s\nSecond Half: %s\nCommon: %c (ascii %d)\n\n", lineFirst, lineSecond, common, getPriority(common))

	}
	fmt.Printf("TOTAL IS: %d, %d", total, totalGrp)
}
