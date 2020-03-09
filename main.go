package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const icsDateLayout = "20060102T150405Z"
const triggerDuration = "-P0DT12H0M0S" // 12 hours before
const outputFilename = "enhanced.ics"

func main() {
	filename := getFilenameFromArgs()

	log.Println("processing ICS file")

	icsFileLines, eventCount := parseEvents(filename)

	log.Printf("processed %v events", eventCount)

	writeOutputFile(icsFileLines)

	log.Printf("finished, written %v events to enhanced ICS file", eventCount+1)
}

func parseEvents(filename string) ([]string, int) {
	readFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	defer readFile.Close()

	var icsFileLines []string
	var summary string
	var lastEventDate time.Time
	var eventCount int

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		if isSummaryLine(line) {
			summary = parseSummary(line)
		} else if isEventStartLine(line) {
			eventCount += 1
		} else if isStartDateLine(line) {
			lastEventDate = parseEventDate(line)
		} else if isEventEndLine(line) {
			icsFileLines = appendTrigger(icsFileLines, summary)
		} else if isCalendarEndLine(line) {
			icsFileLines = appendDownloadNewCalendarFileEvent(icsFileLines, lastEventDate)
		}

		icsFileLines = append(icsFileLines, line)
	}
	return icsFileLines, eventCount
}

func writeOutputFile(icsFileLines []string) {
	file, err := os.Create(outputFilename)
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}

	defer file.Close()

	joinedString := strings.Join(icsFileLines, "\n")

	_, err = io.WriteString(file, joinedString)
	if err != nil {
		log.Fatalf("failed to write file content: %s", err)
	}

	err = file.Sync()
	if err != nil {
		log.Fatalf("failed to write enhanced file: %s", err)
	}
}

func appendTrigger(icsFileLines []string, triggerText string) []string {
	icsFileLines = append(icsFileLines, "BEGIN:VALARM")
	icsFileLines = append(icsFileLines, "ACTION:DISPLAY")
	icsFileLines = append(icsFileLines, "DESCRIPTION:"+triggerText)
	icsFileLines = append(icsFileLines, "TRIGGER:"+triggerDuration)
	icsFileLines = append(icsFileLines, "END:VALARM")

	return icsFileLines
}

func appendDownloadNewCalendarFileEvent(icsFileLines []string, eventDate time.Time) []string {
	log.Println("adding reminder for downloading new calendar file")

	startDate := eventDate.Add(time.Hour * -(12 + 24))
	endDate := startDate.Add(time.Hour)

	startDateString := startDate.Format(icsDateLayout)
	endDateString := endDate.Format(icsDateLayout)

	timestamp := time.Now().Format(icsDateLayout)

	uid := strconv.FormatInt(startDate.Unix(), 10)

	icsFileLines = append(icsFileLines, "BEGIN:VEVENT")
	icsFileLines = append(icsFileLines, "UID:"+uid)
	icsFileLines = append(icsFileLines, "DTSTART:"+startDateString)
	icsFileLines = append(icsFileLines, "SEQUENCE:0")
	icsFileLines = append(icsFileLines, "TRANSP:OPAQUE")
	icsFileLines = append(icsFileLines, "DTEND:"+endDateString)
	icsFileLines = append(icsFileLines, "SUMMARY:Neuen Abfallkalender herunterladen")
	icsFileLines = append(icsFileLines, "CLASS:PUBLIC")
	icsFileLines = append(icsFileLines, "DTSTAMP:"+timestamp)
	icsFileLines = append(icsFileLines, "BEGIN:VALARM")
	icsFileLines = append(icsFileLines, "ACTION:DISPLAY")
	icsFileLines = append(icsFileLines, "DESCRIPTION:Neuen Abfallkalender herunterladen")
	icsFileLines = append(icsFileLines, "TRIGGER:-P0DT0H0M1S")
	icsFileLines = append(icsFileLines, "END:VALARM")
	icsFileLines = append(icsFileLines, "END:VEVENT")

	return icsFileLines
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
		log.Fatal("filename is missing")
	}

	filename := argsWithoutProg[0]

	return filename
}

func isCalendarEndLine(line string) bool {
	return strings.EqualFold(line, "END:VCALENDAR")
}

func isEventEndLine(line string) bool {
	return strings.EqualFold(line, "END:VEVENT")
}

func isStartDateLine(line string) bool {
	return strings.HasPrefix(line, "DTSTART")
}

func isEventStartLine(line string) bool {
	return strings.HasPrefix(line, "BEGIN:VEVENT")
}

func isSummaryLine(line string) bool {
	return strings.HasPrefix(line, "SUMMARY")
}
