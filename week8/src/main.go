package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func readFile(file string) <-chan string {
	content := make(chan string)
	go func() {
		readFile, err := os.Open(file)

		if err != nil {
			fmt.Println(err)
		}
		fileScanner := bufio.NewScanner(readFile)

		fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			//fmt.Println(fileScanner.Text())
			content <- fileScanner.Text()
		}

		readFile.Close()
		close(content)
	}()
	return content
}

func unique(f func() <-chan string) <-chan string {
	content := f()
	uniqueChan := make(chan string)
	go func() {
		defer close(uniqueChan)
		prev := ""
		for c := range content {
			if prev == "" {
				uniqueChan <- c
				prev = c
			} else if prev != c {
				uniqueChan <- c
				prev = c
			}
		}
	}()
	return uniqueChan
}

func uniqueWithCount(f func() <-chan string) <-chan struct {
	int
	string
} {
	content := f()
	uniqueChan := make(chan struct {
		int
		string
	})
	go func() {
		defer close(uniqueChan)
		prev := ""
		count := 0
		for c := range content {
			if prev == "" {
				prev = c
				count = 1
			} else if prev == c {
				count += 1
			} else if prev != c {
				uniqueChan <- struct {
					int
					string
				}{count, prev}
				prev = c
				count = 1
			}
		}
		uniqueChan <- struct {
			int
			string
		}{count, prev}
	}()
	return uniqueChan
}

type options []string

func (o options) Len() int {
	return len(o)
}

func (o options) Less(i, j int) bool {
	return o[i] > o[j]
}

func (o options) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func parseCommmandLine(optionList options) struct {
	options
	string
} {
	commands := make([]string, 0)
	var filename string
	for _, str := range optionList {
		if str[0] == '-' {
			commands = append(commands, str)
		} else {
			filename = str
		}
	}
	return struct {
		options
		string
	}{options: commands, string: filename}
}

func main() {
	options := parseCommmandLine(os.Args[1:])
	//option := "./resource/countries.txt"
	//option := "./resource/test.txt"
	if options.Len() > 0 {
		for _, option := range options.options {
			if option == "-" {
				uniqueChan := unique(func() <-chan string {
					return readContentFromPipe()
				})
				if options.string != "" {
					writeToFile(uniqueChan, options.string)
				}
			} else if option == "-c" {
				println(options.string)
				if options.string != "" {
					responseChannel := uniqueWithCount(func() <-chan string {
						return readFile(options.string)
					})
					printFromResponseChannel(responseChannel)
				}
			} else if option == "-d" {
				if options.string != "" {
					responseChannel := uniqueWithDuplicate(func() <-chan string {
						return readFile(options.string)
					})
					printFromResponseChannelString(responseChannel)
				}
			} else if option == "-u" {
				if options.string != "" {
					responseChannel := onlyUnique(func() <-chan string {
						return readFile(options.string)
					})
					printFromResponseChannelString(responseChannel)
				}
			}
		}
	} else {
		channel := unique(func() <-chan string {
			return readFile(options.string)
		})
		printFromResponseChannelString(channel)
	}
	//fmt.Println("Hello world")
	//fmt.Println("I am here")
}

func onlyUnique(f func() <-chan string) <-chan string {
	uniqueString := make(chan string)
	content := f()
	go func() {
		defer close(uniqueString)
		prev := ""
		count := 0
		for c := range content {
			if prev == "" {
				prev = c
				count = 1
			} else if prev == c {
				count += 1
			} else if prev != c {
				if count == 1 {
					uniqueString <- prev
				}
				prev = c
				count = 1
			}
		}
		if count == 1 {
			uniqueString <- prev
		}
	}()
	return uniqueString
}

func uniqueWithDuplicate(f func() <-chan string) <-chan string {
	duplicateString := make(chan string)
	content := f()
	go func() {
		defer close(duplicateString)
		prev := ""
		count := 0
		for c := range content {
			if prev == "" {
				prev = c
				count = 1
			} else if prev == c {
				count += 1
			} else if prev != c {
				if count > 1 {
					duplicateString <- prev
				}
				prev = c
				count = 1
			}
		}
		if count > 1 {
			duplicateString <- prev
		}
	}()
	return duplicateString
}

func printFromResponseChannel(channel <-chan struct {
	int
	string
}) {
	for val := range channel {
		fmt.Printf("%d %s\n", val.int, val.string)
	}
}

func printFromResponseChannelString(channel <-chan string) {
	for val := range channel {
		fmt.Printf("%s\n", val)
	}
}

func writeToFile(uniqueChan <-chan string, file string) {
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		println("Unable to open file")
		return
	}
	for str := range uniqueChan {
		f.WriteString(str)
		f.WriteString("\n")
	}
}

func readContentFromPipe() <-chan string {
	//stat, _ := os.Stdin.Stat()
	//if (stat.Mode() & os.ModeCharDevice) == 0 {

	chann := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			//stdin = append(stdin, scanner.Bytes()...)
			chann <- string(scanner.Bytes())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		close(chann)
	}()
	return chann

	//}
}

func readFromPipe() {

}
