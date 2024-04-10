// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /prompt4Type.go
// Original timestamp: 2024/04/10 15:23

package helperFunctions

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Getting typed values from prompt
func GetStringValFromPrompt(prompt string) string {
	inputScanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s", prompt)
	inputScanner.Scan()
	nval := inputScanner.Text()
	value := ""

	if nval != "" {
		value = nval
	}
	return value
}

func GetIntValFromPrompt(prompt string) int {
	var err error
	value := 0
	inputScanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s", prompt)
	inputScanner.Scan()
	nval := inputScanner.Text()

	if nval != "" {
		value, err = strconv.Atoi(nval)
		if err != nil {
			value = 1
		}
	}
	return value
}

func GetBoolValFromPrompt(prompt string) bool {
	fmt.Printf("%s", prompt)
	bval := ""
	var value = false

	fmt.Scanln(&bval)
	if strings.HasPrefix(strings.ToLower(bval), "t") || bval == "1" {
		value = true
	}
	return value
}

func GetStringSliceFromPrompt(prompt string) []string {
	slice := []string{}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("%s\n", prompt)
	for {
		fmt.Println("Just press enter to end the loop")
		scanner.Scan()
		input := scanner.Text()
		if input == "" {
			break
		} else {
			slice = append(slice, input)
		}
	}
	return slice
}
