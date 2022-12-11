package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Need to include the equal to show that it is the same
type WorryLevel = big.Int

type Monkey struct {
	Items           []WorryLevel
	Operation       func(int) int     // From example, test always of the form int <operation> int, so function takes old regardless
	BigOp           func(*WorryLevel) // Type alias doesn't seem to work for pointers
	DivisibleTest   int
	DivisibleTrue   int // TODO might make sense to refactor to a pointer to another monkey but unlikely
	DivisibleFalse  int
	ItemInspections int // Number of items the monkey has inspected
}

var monkeys []Monkey

// func part1() {

// 	f, err := os.OpenFile("../data/day11a", os.O_RDONLY, 0755)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fileScanner := bufio.NewScanner(f)
// 	fileScanner.Split(bufio.ScanLines)

// 	parseAllMonkeys(fileScanner)

// 	// Says to do 20 rounds
// 	for roundNum := 1; roundNum <= 20; roundNum++ {

// 		for j := range monkeys {

// 			for i := range monkeys[j].Items {

// 				monkeys[j].Items[i] = WorryLevel(monkeys[j].Operation(int(monkeys[j].Items[i])))
// 				monkeys[j].Items[i] = monkeys[j].Items[i] / 3

// 				if monkeys[j].Items[i]%WorryLevel(monkeys[j].DivisibleTest) == 0 {
// 					monkeys[monkeys[j].DivisibleTrue].Items = append(monkeys[monkeys[j].DivisibleTrue].Items, monkeys[j].Items[i])
// 				} else {
// 					monkeys[monkeys[j].DivisibleFalse].Items = append(monkeys[monkeys[j].DivisibleFalse].Items, monkeys[j].Items[i])
// 				}

// 				monkeys[j].ItemInspections++

// 			}

// 			// Clear all items from this monkey
// 			monkeys[j].Items = monkeys[j].Items[:0]

// 		}

// 		// fmt.Printf("\nAfter round %d, the monkeys are holding items with these worry levels:\n", roundNum)

// 		// for num, m := range monkeys {

// 		// 	// fmt.Printf("Monkey %d: %v\n", num, m.Items)

// 		// }

// 	}

// 	fmt.Println("TOTAL INSPECTIONS:")

// 	var insps []int

// 	for _, m := range monkeys {

// 		// fmt.Printf("Monkey %d inspected items %d times.\n", num, m.ItemInspections)
// 		insps = append(insps, m.ItemInspections)

// 	}

// 	sort.IntSlice(insps).Sort()

// 	insps = insps[len(insps)-2:]

// 	total := 1
// 	for _, item := range insps {
// 		total *= item
// 	}

// 	fmt.Println(total)

// }

func part2() {

	f, err := os.OpenFile("../data/day11t", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	parseAllMonkeys(fileScanner)

	// Says to do 20 rounds
	for roundNum := 1; roundNum <= 10000; roundNum++ {

		for j := range monkeys {

			for i := range monkeys[j].Items {
				monkeys[j].ItemInspections++

				monkeys[j].BigOp(&monkeys[j].Items[i])
				// monkeys[j].Items[i] = monkeys[j].Items[i] / 3

				modRes := big.Int{}
				mod1 := big.Int(monkeys[j].Items[i])
				mod2 := big.NewInt(int64(monkeys[j].DivisibleTest))

				modRes.Mod(&mod1, mod2)
				if modRes.Int64() == 0 {
					monkeys[monkeys[j].DivisibleTrue].Items = append(monkeys[monkeys[j].DivisibleTrue].Items, monkeys[j].Items[i])
				} else {
					monkeys[monkeys[j].DivisibleFalse].Items = append(monkeys[monkeys[j].DivisibleFalse].Items, monkeys[j].Items[i])
				}

			}

			// Clear all items from this monkey
			monkeys[j].Items = monkeys[j].Items[:0]

		}

		if roundNum == 1 || roundNum == 20 || roundNum%1000 == 0 {
			fmt.Printf("\n\n== After round %d ==\n", roundNum)

			for num, m := range monkeys {

				fmt.Printf("Monkey %d currently has %v\n", num, m.Items)
				fmt.Printf("Monkey %d inspected items %d times.\n", num, m.ItemInspections)

			}
		} else {
			fmt.Printf("\rRound %d", roundNum)
		}

	}

	fmt.Println("TOTAL INSPECTIONS:")

	var insps []int

	for num, m := range monkeys {

		fmt.Printf("Monkey %d inspected items %d times.\n", num, m.ItemInspections)
		insps = append(insps, m.ItemInspections)

	}

	sort.IntSlice(insps).Sort()

	insps = insps[len(insps)-2:]

	total := 1
	for _, item := range insps {
		total *= item
	}

	fmt.Println(total)

}

func parseAllMonkeys(fileScanner *bufio.Scanner) {

	var monkeyBuf []string

	// fmt.Println("Starting to parse monkeys...")
	for fileScanner.Scan() {

		line := fileScanner.Text()

		// End of monkey, clear the buffer, parseMonkey and save to list of monkeys
		if line == "" {

			monkeys = append(monkeys, *parseMonkey(monkeyBuf))
			monkeyBuf = []string{}
			continue

		}

		monkeyBuf = append(monkeyBuf, line)

	}

	// End of monkey, clear the buffer, parseMonkey and save to list of monkeys
	monkeys = append(monkeys, *parseMonkey(monkeyBuf))

	// fmt.Printf("Finished parsing monkeys. There are %d in total\n", len(monkeys))
	// fmt.Printf("Monkeys: %v\n", monkeys)

}

// Parse an input of a Monkey
func parseMonkey(info []string) *Monkey {

	divNum, trueMonkey, falseMonkey := parseDivisible(info)

	return &Monkey{
		Items:           parseItems(info),
		Operation:       parseOperation(info),
		BigOp:           parseOperationBig(info),
		DivisibleTest:   divNum,
		DivisibleTrue:   trueMonkey,
		DivisibleFalse:  falseMonkey,
		ItemInspections: 0,
	}

}

func main() {
	// part1()
	monkeys = monkeys[:0]
	part2()
}

func parseItems(info []string) []WorryLevel {

	itemsString := strings.Split(info[1], ":")[1][1:]
	itemsStringSplit := strings.Split(itemsString, ", ")

	var items []WorryLevel

	for _, i := range itemsStringSplit {

		itemInt, _ := strconv.Atoi(i)
		items = append(items, *big.NewInt(int64(itemInt)))

	}

	return items

}

func parseOperation(info []string) func(int) int {

	opString := strings.Split(info[2], "= ")[1]

	tokens := strings.Split(opString, " ")

	// Both tokens do something on an old
	if tokens[0] == "old" && tokens[2] == "old" {

		switch tokens[1] {
		case "+":
			return func(old int) int { return old + old }
		case "-":
			return func(old int) int { return old - old }
		case "*":
			return func(old int) int { return old * old }
		case "/":
			return func(old int) int { return old / old }
		}

	} else if tokens[0] == "old" {

		secondInt, _ := strconv.Atoi(tokens[2])

		switch tokens[1] {
		case "+":
			return func(old int) int { return old + secondInt }
		case "-":
			return func(old int) int { return old - secondInt }
		case "*":
			return func(old int) int { return old * secondInt }
		case "/":
			return func(old int) int { return old / secondInt }
		}

	} else if tokens[2] == "old" {

		firstInt, _ := strconv.Atoi(tokens[0])

		switch tokens[1] {
		case "+":
			return func(old int) int { return firstInt + old }
		case "-":
			return func(old int) int { return firstInt - old }
		case "*":
			return func(old int) int { return firstInt * old }
		case "/":
			return func(old int) int { return firstInt / old }
		}

	} else {

		firstInt, _ := strconv.Atoi(tokens[0])
		secondInt, _ := strconv.Atoi(tokens[2])

		switch tokens[1] {
		case "+":
			return func(old int) int { return firstInt + secondInt }
		case "-":
			return func(old int) int { return firstInt - secondInt }
		case "*":
			return func(old int) int { return firstInt * secondInt }
		case "/":
			return func(old int) int { return firstInt / secondInt }
		}

	}

	return nil

}

func parseDivisible(info []string) (int, int, int) {

	divNum, _ := strconv.Atoi(strings.Split(info[3], "by ")[1])
	trueMonkey, _ := strconv.Atoi(strings.Split(info[4], "monkey ")[1])
	falseMonkey, _ := strconv.Atoi(strings.Split(info[5], "monkey ")[1])
	return divNum, trueMonkey, falseMonkey

}

// Takes the integer and updates without writing to new location
func parseOperationBig(info []string) func(*WorryLevel) {

	opString := strings.Split(info[2], "= ")[1]

	tokens := strings.Split(opString, " ")

	// Both tokens do something on an old
	if tokens[0] == "old" && tokens[2] == "old" {

		switch tokens[1] {
		case "+":
			return func(old *big.Int) { old.Add(old, old) }
		case "-":
			return func(old *big.Int) { old.Sub(old, old) }
		case "*":
			return func(old *big.Int) { old.Mul(old, old) }
		case "/":
			return func(old *big.Int) { old.Div(old, old) }
		}

	} else if tokens[0] == "old" {

		secondInt, _ := strconv.Atoi(tokens[2])
		secondBigInt := big.NewInt(int64(secondInt))

		switch tokens[1] {
		case "+":
			return func(old *big.Int) { old.Add(old, secondBigInt) }
		case "-":
			return func(old *big.Int) { old.Sub(old, secondBigInt) }
		case "*":
			return func(old *big.Int) { old.Mul(old, secondBigInt) }
		case "/":
			return func(old *big.Int) { old.Div(old, secondBigInt) }
		}

	} else if tokens[2] == "old" {

		firstInt, _ := strconv.Atoi(tokens[0])
		firstBigInt := big.NewInt(int64(firstInt))

		switch tokens[1] {
		case "+":
			return func(old *big.Int) { old.Add(firstBigInt, old) }
		case "-":
			return func(old *big.Int) { old.Sub(firstBigInt, old) }
		case "*":
			return func(old *big.Int) { old.Mul(firstBigInt, old) }
		case "/":
			return func(old *big.Int) { old.Div(firstBigInt, old) }
		}

	} else {

		firstInt, _ := strconv.Atoi(tokens[0])
		firstBigInt := big.NewInt(int64(firstInt))

		secondInt, _ := strconv.Atoi(tokens[2])
		secondBigInt := big.NewInt(int64(secondInt))

		switch tokens[1] {
		case "+":
			return func(old *big.Int) { old.Add(firstBigInt, secondBigInt) }
		case "-":
			return func(old *big.Int) { old.Sub(firstBigInt, secondBigInt) }
		case "*":
			return func(old *big.Int) { old.Mul(firstBigInt, secondBigInt) }
		case "/":
			return func(old *big.Int) { old.Div(firstBigInt, secondBigInt) }
		}

	}

	return nil

}
