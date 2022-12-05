package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Play int

const (
	Rock = iota + 1
	Paper
	Scissors
)

type Strat int

const (
	Win = iota + 1
	Draw
	Lose
)

func calcScore(opp Play, us Play) int {

	if opp == us {
		return 3 + int(us)
	}

	// Rock
	if opp == Rock {
		if us == Paper {
			return 6 + int(us)
		} else {
			return 0 + int(us)
		}
	}

	// Paper
	if opp == Paper {
		if us == Scissors {
			return 6 + int(us)
		} else {
			return 0 + int(us)
		}
	}

	if us == Rock {
		return 6 + int(us)
	} else {
		return 0 + int(us)
	}
}

func calcStrat(opp Play, strat Strat) int {

	if strat == Win {
		switch opp {
		case Rock:
			return calcScore(opp, Paper)
		case Paper:
			return calcScore(opp, Scissors)
		case Scissors:
			return calcScore(opp, Rock)
		}
	} else if strat == Draw {
		return calcScore(opp, opp)
	}

	switch opp {
	case Rock:
		return calcScore(opp, Scissors)
	case Paper:
		return calcScore(opp, Rock)
	case Scissors:
		return calcScore(opp, Paper)
	}

	return -1

}

func main() {

	f, err := os.OpenFile("../data/day2a", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	var totalScore = 0
	var stratScore = 0

	for fileScanner.Scan() {

		line := fileScanner.Text()

		strat := strings.Split(line, " ")

		them := strat[0]
		us := strat[1]

		var themEnc Play
		var usEnc Play
		var usStrat Strat

		switch them {
		case "A":
			themEnc = Rock
		case "B":
			themEnc = Paper
		case "C":
			themEnc = Scissors
		}

		switch us {
		case "X":
			usEnc = Rock
			usStrat = Lose
		case "Y":
			usEnc = Paper
			usStrat = Draw
		case "Z":
			usEnc = Scissors
			usStrat = Win
		}

		totalScore += calcScore(themEnc, usEnc)
		stratScore += calcStrat(themEnc, usStrat)

	}

	fmt.Printf("Intial: %d, strat: %d\n", totalScore, stratScore)

}
