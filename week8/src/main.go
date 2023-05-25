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

func main() {
	filename := os.Args[1]
	//filename := "./resource/countries.txt"
	//filename := "./resource/test.txt"
	if filename == "-" {
		uniqueChan := unique(func() <-chan string {
			return readContentFromPipe()
		})
		if len(os.Args) > 2 {
			writeToFile(uniqueChan, os.Args[2])
		}
	} else {
		unique(func() <-chan string {
			return readFile(filename)
		})
	}
	//fmt.Println("Hello world")
	//fmt.Println("I am here")
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
