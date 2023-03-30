package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	optionToHandlerFunc := make(map[string]func(string2 string) (int, error))
	optionToHandlerFunc["-c"] = getFileSize
	optionToHandlerFunc["-l"] = getLineCount
	optionToHandlerFunc["-w"] = getWordCount
	optionToHandlerFunc["-n"] = getCharacterCount
	results := make([]interface{}, 0)
	if len(os.Args) == 3 {
		option := os.Args[1]
		file := os.Args[2]
		i, err := optionToHandlerFunc[option](file)
		if err != nil {
			return
		}
		results = append(results, int(i))
		results = append(results, file)
	} else if len(os.Args) == 2 {
		file := os.Args[1]
		for _, function := range optionToHandlerFunc {
			i, err := function(file)
			if err != nil {
				return
			}
			results = append(results, i)
		}
		results = append(results, file)
	}
	for _, val := range results {
		if _, ok := val.(string); ok {
			valStr := fmt.Sprint(val)
			fmt.Print(valStr + " ")
		} else if _, ok := val.(int); ok {
			valStr := fmt.Sprint(val)
			fmt.Print(valStr + " ")
		}
	}
}

func getFileSize(filename string) (int, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	// get the size
	size := fi.Size()
	return int(size), nil
}

func getLineCount(filename string) (int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return -1, err
	}
	//fmt.Printf(string(data))
	temp := strings.Split(string(data), "\n")
	return len(temp), nil
}

func getWordCount(filename string) (int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return -1, err
	}
	temp := strings.Fields(string(data))
	return len(temp), nil
}

func getCharacterCount(filename string) (int, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		return -1, err
	}
	//fmt.Printf(string(data))
	temp := utf8.RuneCountInString(string(data))
	return temp, nil
}
