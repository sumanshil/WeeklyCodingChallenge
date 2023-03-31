package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {

	optionToHandlerFunc := make(map[string]func(data string) (int, error))
	optionToHandlerFunc["-l"] = getLineCount
	optionToHandlerFunc["-w"] = getWordCount
	optionToHandlerFunc["-m"] = getCharacterCount
	stat, _ := os.Stdin.Stat()
	results := make([]interface{}, 0)
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		stdin, _ := io.ReadAll(os.Stdin)
		lines := string(stdin)
		if len(os.Args) == 2 {
			option := os.Args[1]
			result, _ := optionToHandlerFunc[option](lines)
			results = append(results, result)
		} else {
			for _, function := range optionToHandlerFunc {
				i, err := function(lines)
				if err != nil {
					return
				}
				results = append(results, i)
			}
		}
	} else {
		if len(os.Args) == 3 {
			option := os.Args[1]
			file := os.Args[2]
			if option == "-c" {
				i, _ := getFileSize(file)
				results = append(results, i)
			} else {
				data, err := os.ReadFile(file)
				if err != nil {
					return
				}

				i, err := optionToHandlerFunc[option](string(data))
				if err != nil {
					return
				}
				results = append(results, int(i))
			}
			results = append(results, file)
		} else if len(os.Args) == 2 {
			file := os.Args[1]
			data, err := os.ReadFile(file)
			if err != nil {
				return
			}
			//for _, function := range optionToHandlerFunc {
			//	i, err := function(string(data))
			//	if err != nil {
			//		return
			//	}
			//	results = append(results, i)
			//}
			r, _ := getLineCount(string(data))
			results = append(results, r)
			r, _ = getWordCount(string(data))
			results = append(results, r)
			r, _ = getFileSize(file)
			results = append(results, r)
			results = append(results, file)
		}
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

func getLineCount(data string) (int, error) {
	temp := strings.Split(data, "\n")
	return len(temp), nil
}

func getWordCount(data string) (int, error) {
	temp := strings.Fields(data)
	return len(temp), nil
}

func getCharacterCount(data string) (int, error) {
	temp := utf8.RuneCountInString(data)
	return temp, nil
}
