package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const icsDateFormat = "20060102T150405Z"
const triggerDuration = "-P0DT12H0M0S" // 12 hours before

func main() {
	filename := getFilenameFromArgs()

	processFile(filename)
}

func processFile(filename string) {
	readFile, err := os.Open(filename)

	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	var fileTextLines []string
	var summary string

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		if strings.HasPrefix(line, "SUMMARY") {
			summary = parseSummary(line)
		}

		if strings.EqualFold(line, "END:VEVENT") {
			// lineTokens := strings.Split(line, ":")
			// startDateString := lineTokens[1]
			// startDate := parseDate(startDateString)
			// reminderDate := startDate.Add(time.Hour * -11)

			fileTextLines = appendTrigger(fileTextLines, summary)
		}

		fileTextLines = append(fileTextLines, line)
	}

	_ = readFile.Close()

	for _, eachline := range fileTextLines {
		fmt.Println(eachline)
	}
}

func appendTrigger(fileTextLines []string, triggerText string) []string {
	fileTextLines = append(fileTextLines, "BEGIN:VALARM")
	fileTextLines = append(fileTextLines, "ACTION:DISPLAY")
	fileTextLines = append(fileTextLines, "DESCRIPTION:"+triggerText)
	fileTextLines = append(fileTextLines, "TRIGGER:"+triggerDuration)
	fileTextLines = append(fileTextLines, "END:VALARM")
	return fileTextLines
}

func parseSummary(line string) string {
	lineTokens := strings.Split(line, ":")
	return lineTokens[1]
}

func parseDate(dateString string) time.Time {
	date, err := time.Parse(icsDateFormat, dateString)

	if err != nil {
		log.Fatalf("failed to parse date: %s", err)
	}

	return date
}

func getFilenameFromArgs() string {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) < 1 {
		log.Fatalf("filename is missing")
	}

	filename := argsWithoutProg[0]
	return filename
}
