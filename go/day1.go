package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {

	f, err := os.OpenFile("../data/day1a", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	type elfCalories struct {
		cals  int
		elfId int
	}
	var totalArr = []elfCalories{}
	var counter = 0
	var elf = 1

	for fileScanner.Scan() {

		line := fileScanner.Text()

		if line == "" {
			fmt.Printf("Finished elf %03d who ate %05d Calories\n", len(totalArr), counter)
			totalArr = append(totalArr, elfCalories{counter, elf})
			counter = 0
			elf++
		} else {
			val, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}

			counter += val

		}
	}

	sort.Slice(totalArr, func(i, j int) bool {
		return totalArr[i].cals < totalArr[j].cals
	})

	fmt.Println(totalArr)

	fmt.Printf("Most to eat was elf %d with %d calories", totalArr[len(totalArr)-1].elfId, totalArr[len(totalArr)-1].cals)

	topThree := totalArr[len(totalArr)-3:]
	totalCals := 0
	for _, el := range topThree {
		totalCals += el.cals
	}

	fmt.Printf("Most to eat was elf %d with %d calories", totalArr[len(totalArr)-1].elfId, totalArr[len(totalArr)-1].cals)

	fmt.Printf("total top 3: %d\n", totalCals)

}
