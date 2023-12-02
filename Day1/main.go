package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// 54239

// 55343

var maxChannels int = 4

var digitWords []string

var minWordLen = 3
var maxWordLen = 5

func main() {
	digitWords = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	puzzleLines := fileToLines("input.txt")

	inboundVals := lineSender(puzzleLines)

	resultChannelSet := make([]chan int, maxChannels)

	for i, _ := range resultChannelSet {
		// fan out
		resultChannelSet[i] = lineProcessor(inboundVals)
	}

	var total int

	//fan in
	for num := range merge(resultChannelSet) {
		total += num
	}

	fmt.Println("result: ", total)
}

func fileToLines(path string) []string {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("cannot read file: ", err)
	}

	output := strings.Split(string(fileBytes), "\n")

	return output
}

// feeds our input channel with values
func lineSender(lines []string) chan string {
	out := make(chan string)
	go func() {
		for _, line := range lines {
			out <- line
		}
		close(out)
	}()
	return out
}

// take in an input string channel
// on each string emission
// find the first and then last numerical digit
// combine and form a 2 digit integer
func lineProcessor(line chan string) chan int {
	out := make(chan int)
	go func() {
		for n := range line {
			out <- combineDigits(n)
		}
		close(out)
	}()
	return out
}

func combineDigits(line string) int {
	var firstDigit int
	var lastDigit int

	foundFirstDigit := false
	foundLastDigit := false

	firstDigitIndex := -1
	lastDigitIndex := -1

	for i := 0; i < len(line); i++ {
		val, err := strconv.Atoi(string(line[i]))
		if err == nil {
			firstDigit = val * 10 // this is in the tens place after all
			if i < 3 {
				foundFirstDigit = true
			} else {
				firstDigitIndex = i
			}
			break
		}
	}

	for i := len(line) - 1; i > -1; i-- {
		val, err := strconv.Atoi(string(line[i]))
		if err == nil {
			lastDigit = val // ones place
			if i > len(line)-4 {
				foundLastDigit = true
			} else {
				lastDigitIndex = i
			}
			break
		}
	}

	if !foundFirstDigit {
		val := slideForwardOverString(line, firstDigitIndex)
		if val != 0 {
			firstDigit = val * 10
		}
	}

	if !foundLastDigit {
		val := slideBackwardOverString(line, lastDigitIndex)
		if val != 0 {
			lastDigit = val
		}
	}

	combinedDigits := firstDigit + lastDigit

	return combinedDigits
}

func slideForwardOverString(data string, potentialDigitIndex int) int {
	for i := 0; i < len(data); i++ {
		for j := minWordLen; j < maxWordLen+1; j++ {
			if j+i <= len(data) {
				evalString := data[i : i+j]
				result := checkForWrittenDigit(evalString)
				if len(result) > 0 {
					if potentialDigitIndex == -1 || potentialDigitIndex >= i+j {
						return literalToInteger(result)
					}
				}
			} else {
				break // our string is to small now
			}
		}
	}

	return 0
}

func slideBackwardOverString(data string, potentialDigitIndex int) int {
	for i := len(data) - 1; i > minWordLen-1; i-- {
		for wordSize := minWordLen; wordSize < maxWordLen+1; wordSize++ {
			if i-wordSize > -1 {
				evalString := data[i-wordSize : i+1]
				result := checkForWrittenDigit(evalString)
				if len(result) > 0 {
					if potentialDigitIndex == -1 || i-wordSize >= potentialDigitIndex {
						return literalToInteger(result)
					}
				}
			} else {
				break // our string is to small now
			}
		}
	}
	return 0
}

func checkForWrittenDigit(s string) string {
	for i, _ := range digitWords {
		if strings.Contains(s, digitWords[i]) {
			return digitWords[i]
		}
	}
	return ""
}

func literalToInteger(s string) int {
	switch s {
	case "one":
		return 1
	case "two":
		return 2
	case "three":
		return 3
	case "four":
		return 4
	case "five":
		return 5
	case "six":
		return 6
	case "seven":
		return 7
	case "eight":
		return 8
	case "nine":
		return 9
	}
	return 0
}

// combine our processing channels output back into a single input source
func merge(channels []chan int) chan int {
	out := make(chan int)

	var fanInWaitGroup sync.WaitGroup
	fanInWaitGroup.Add(len(channels))

	for _, c := range channels {
		go func(ch chan int) {
			for n := range ch {
				out <- n
			}
			fanInWaitGroup.Done()
		}(c)
	}

	go func() {
		fanInWaitGroup.Wait()
		close(out)
	}()

	return out
}
