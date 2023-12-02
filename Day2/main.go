package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type GameData struct {
	Id     int
	Rounds []GameRound
}

type GameRound struct {
	GrabResult []ColorAmount
}

type ColorAmount struct {
	Amount int
	Color  ColorField
}

type ColorField int

const (
	red ColorField = iota
	green
	blue
	noColor
)

type BagMaximums struct {
	red, green, blue int
}

var inBagValues = BagMaximums{
	12, 13, 14,
}

func main() {
	var sumOfValidIds int
	var sumOfPowers int

	gameData := readGameData(fileToLines("input.txt"))

	for i, _ := range gameData {
		if isGameValid((gameData[i])) {
			sumOfValidIds += gameData[i].Id
		}
	}

	for i, _ := range gameData {
		sumOfPowers += powerOfMinimumValues(gameData[i])
	}

	fmt.Println("sum of ids: ", sumOfValidIds)
	fmt.Println("power: ", sumOfPowers)
}

func fileToLines(path string) []string {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("cannot read file: ", err)
	}

	output := strings.Split(string(fileBytes), "\n")

	return output
}

func readGameData(data []string) []GameData {
	gameData := make([]GameData, len(data))
	var err error
	for i, _ := range data {
		partialLine := strings.Split(strings.Split(data[i], "Game ")[1], ":")
		gameID := partialLine[0]
		gameData[i].Id, err = strconv.Atoi(gameID)
		if err != nil {
			fmt.Println("unexpected line format: ", data[i])
			return gameData
		}
		rounds := strings.Split(partialLine[1], ";")
		gameData[i].Rounds = make([]GameRound, len(rounds))
		for j, _ := range rounds {
			gameData[i].Rounds[j].GrabResult = getColorAmounts(rounds[j])
		}
	}
	return gameData
}

func getColorAmounts(s string) []ColorAmount {
	colorSets := strings.Split(s, ",")
	colors := make([]ColorAmount, len(colorSets))
	for i, _ := range colorSets {
		fields := strings.Fields(colorSets[i])
		colors[i].Amount, _ = strconv.Atoi(fields[0])
		colors[i].Color = getColorField(fields[1])
	}
	return colors
}

func getColorField(col string) ColorField {
	switch col {
	case "red":
		return red
	case "blue":
		return blue
	case "green":
		return green
	default:
		fmt.Println("couldn't find color! ", col)
		return noColor
	}
}

func powerOfMinimumValues(data GameData) int {
	maxValue := BagMaximums{}
	for i, _ := range data.Rounds {
		for j, _ := range data.Rounds[i].GrabResult {
			switch data.Rounds[i].GrabResult[j].Color {
			case red:
				if maxValue.red < data.Rounds[i].GrabResult[j].Amount {
					maxValue.red = data.Rounds[i].GrabResult[j].Amount
				}
			case green:
				if maxValue.green < data.Rounds[i].GrabResult[j].Amount {
					maxValue.green = data.Rounds[i].GrabResult[j].Amount
				}
			case blue:
				if maxValue.blue < data.Rounds[i].GrabResult[j].Amount {
					maxValue.blue = data.Rounds[i].GrabResult[j].Amount
				}
			}
		}
	}
	return (maxValue.red * maxValue.green * maxValue.blue)
}

func isGameValid(data GameData) bool {
	for i, _ := range data.Rounds {
		for j, _ := range data.Rounds[i].GrabResult {
			switch data.Rounds[i].GrabResult[j].Color {
			case red:
				if data.Rounds[i].GrabResult[j].Amount > inBagValues.red {
					return false
				}
			case green:
				if data.Rounds[i].GrabResult[j].Amount > inBagValues.green {
					return false
				}
			case blue:
				if data.Rounds[i].GrabResult[j].Amount > inBagValues.blue {
					return false
				}
			}
		}
	}
	return true
}
