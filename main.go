package main

import (
	"bufio"
	uuid2 "github.com/google/uuid"
	"log"
	"os"
	"strings"
	"time"
)

const icsDateLayout = "20060102T150405Z"
const triggerDuration = "-P0DT12H0M0S" // 12 hours before

func main() {
	filename := getFilenameFromArgs()

	processFile(filename)
}

func processFile(filename string) {
	log.Println("processing ICS file")

	readFile, err := os.Open(filename)

	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	var fileTextLines []string
	var summary string
	var lastEventDate time.Time

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		if strings.HasPrefix(line, "SUMMARY") {
			summary = parseSummary(line)
		}

		if strings.HasPrefix(line, "DTSTART") {
			lastEventDate = parseEventDate(line)
		}

		if strings.EqualFold(line, "END:VEVENT") {
			fileTextLines = appendTrigger(fileTextLines, summary)
		}

		if strings.EqualFold(line, "END:VCALENDAR") {
			fileTextLines = appendDownloadNewCalendarDatesEvent(fileTextLines, lastEventDate)
		}

		fileTextLines = append(fileTextLines, line)
	}

	_ = readFile.Close()

	for _, eachline := range fileTextLines {
		log.Println(eachline)
	}

	log.Println("finished")
}

func appendTrigger(fileTextLines []string, triggerText string) []string {
	fileTextLines = append(fileTextLines, "BEGIN:VALARM")
	fileTextLines = append(fileTextLines, "ACTION:DISPLAY")
	fileTextLines = append(fileTextLines, "DESCRIPTION:"+triggerText)
	fileTextLines = append(fileTextLines, "TRIGGER:"+triggerDuration)
	fileTextLines = append(fileTextLines, "END:VALARM")
	return fileTextLines
}

func appendDownloadNewCalendarDatesEvent(fileTextLines []string, eventDate time.Time) []string {
	log.Println("adding reminder for downloading new calendar file")

	uuid := uuid2.New().String()

	startDate := eventDate.Add(time.Hour * -(12 + 24))
	endDate := startDate.Add(time.Hour)

	startDateString := startDate.Format(icsDateLayout)
	endDateString := endDate.Format(icsDateLayout)

	timestamp := time.Now().Format(icsDateLayout)

	fileTextLines = append(fileTextLines, "BEGIN:VEVENT")
	fileTextLines = append(fileTextLines, "UID:"+uuid)
	fileTextLines = append(fileTextLines, "DTSTART:"+startDateString)
	fileTextLines = append(fileTextLines, "SEQUENCE:0")
	fileTextLines = append(fileTextLines, "TRANSP:OPAQUE")
	fileTextLines = append(fileTextLines, "DTEND:"+endDateString)
	fileTextLines = append(fileTextLines, "SUMMARY:Neuen Abfallkalender herunterladen")
	fileTextLines = append(fileTextLines, "CLASS:PUBLIC")
	fileTextLines = append(fileTextLines, "DTSTAMP:"+timestamp)
	fileTextLines = append(fileTextLines, "BEGIN:VALARM")
	fileTextLines = append(fileTextLines, "ACTION:DISPLAY")
	fileTextLines = append(fileTextLines, "DESCRIPTION:Neuen Abfallkalender herunterladen")
	fileTextLines = append(fileTextLines, "TRIGGER:-P0DT0H0M1S")
	fileTextLines = append(fileTextLines, "END:VALARM")
	fileTextLines = append(fileTextLines, "END:VEVENT")
	return fileTextLines
}

func parseEventDate(line string) time.Time {
	lineTokens := strings.Split(line, ":")
	dateString := lineTokens[1]
	return parseDate(dateString)
}

func parseSummary(line string) string {
	lineTokens := strings.Split(line, ":")
	return lineTokens[1]
}

func parseDate(dateString string) time.Time {
	date, err := time.Parse(icsDateLayout, dateString)

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
