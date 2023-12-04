package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var numerical = regexp.MustCompile("[0-9]+")

func main() {
	var sum int
	data := fileToLines("input.txt")

	// pad data so I don't have to if check on perimeters
	var pad string
	for i := 0; i < len(data[0]); i++ {
		pad += "."
	}

	data = append(data, pad)
	data = append([]string{pad}, data...)

	for i := range data {
		data[i] = "." + data[i] + "."
	}

	// start late and exit early to avoid perimeter pad we added
	for i := 1; i < len(data)-1; i++ {
		sum += SumOfValidLineNums(append([]string{}, data[i-1:i+2]...)) // make a slice copy to avoid reference passing
	}

	fmt.Println(sum)
}

func fileToLines(path string) []string {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("cannot read file: ", err)
	}

	output := strings.Split(string(fileBytes), "\n")

	return output
}

// check the middle line in a 3 line string slice
func SumOfValidLineNums(block []string) int {
	searchingForNums := true
	var sum int

	for searchingForNums {
		result := numerical.FindStringIndex(block[1])
		if result == nil {
			break
		}
		// for ease of reading
		startIndex := result[0]
		endIndex := result[1]
		if perimeterHasSymbol((block[0][startIndex-1:endIndex+1] +
			block[1][startIndex-1:startIndex] +
			block[1][endIndex:endIndex+1] +
			block[2][startIndex-1:endIndex+1])) {
			num, err := strconv.Atoi(block[1][startIndex:endIndex])
			if err == nil {
				sum += num
			}
		}
		for i := range block {
			// shrink the block
			block[i] = block[i][endIndex:]
		}
		if len(block) == 0 {
			searchingForNums = false
		}
	}

	return sum
}

func perimeterHasSymbol(s string) bool {
	for _, v := range s {
		if string(v) != "." && string(v) != "\r" {
			return true
		}
	}
	return false
}
