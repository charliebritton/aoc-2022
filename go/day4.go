package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	f, err := os.OpenFile("../data/day4a", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	var total = 0
	var over = 0

	for fileScanner.Scan() {

		line := fileScanner.Text()

		secs := strings.Split(line, ",")
		e0s := strings.Split(secs[0], "-")
		e1s := strings.Split(secs[1], "-")

		var e0 []int
		var e1 []int

		for _, e := range e0s {
			i, _ := strconv.Atoi(e)
			e0 = append(e0, i)
		}

		for _, e := range e1s {
			i, _ := strconv.Atoi(e)
			e1 = append(e1, i)
		}

		fmt.Printf("e0[0]: %d\ne0[1]: %d\ne1[0]: %d\ne1[1]: %d\n", e0[0], e0[1], e1[0], e1[1])

		// ...456...  4-6
		// .....6...  6-6
		if e0[0] < e1[0] && e0[1] >= e1[1] {
			fmt.Printf("e0 wholly contains e1 near top end\n\n")
			total += 1
			over += 1
			// e1 wholly contained at bot end
		} else if e0[0] <= e1[0] && e0[1] > e1[1] {
			fmt.Printf("e0 wholly contains e1 near bottom end\n\n")
			total += 1
			over += 1
			// e1 wholly contained at top end
		} else if e1[0] < e0[0] && e1[1] >= e0[1] {
			fmt.Printf("e1 wholle contains e0 near top end\n\n")
			total += 1
			over += 1
			// e1 wholly contained at bot end
		} else if e1[0] <= e0[0] && e1[1] > e0[1] {
			fmt.Printf("e1 wholly containes e0 near bottom end\n\n")
			total += 1
			over += 1
		} else if e0[0] == e1[0] && e0[1] == e1[1] {
			fmt.Printf("e1 and e0 are equal\n\n")
			total += 1
			over += 1
			// e0 starts before e1 overlaps e1 or e1 starts before e0 overlaps e1
		} else if (e0[1] >= e1[0] && e0[0] <= e1[0]) || (e1[1] >= e0[0] && e1[0] <= e0[0]) {
			fmt.Printf("Partial overlap\n\n")
			over += 1
		} else {
			fmt.Printf("No overlap at all\n\n")
		}
	}
	fmt.Printf("TOTAL is: %d, %d", total, over)

}
