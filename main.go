package main

import (
	"bufio"
	uuid2 "github.com/google/uuid"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const icsDateLayout = "20060102T150405Z"
const triggerDuration = "-P0DT12H0M0S" // 12 hours before
const outputFilename = "enhanced.ics"  // 12 hours before

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

	defer readFile.Close()

	var icsTextLines []string
	var summary string
	var lastEventDate time.Time
	var eventCount int

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		if strings.HasPrefix(line, "SUMMARY") {
			summary = parseSummary(line)
		}

		if strings.HasPrefix(line, "BEGIN:VEVENT") {
			eventCount += 1
		}

		if strings.HasPrefix(line, "DTSTART") {
			lastEventDate = parseEventDate(line)
		}

		if strings.EqualFold(line, "END:VEVENT") {
			icsTextLines = appendTrigger(icsTextLines, summary)
		}

		if strings.EqualFold(line, "END:VCALENDAR") {
			icsTextLines = appendDownloadNewCalendarDatesEvent(icsTextLines, lastEventDate)
		}

		icsTextLines = append(icsTextLines, line)
	}

	log.Printf("processed %v events", eventCount)

	err = writeOutputFile(icsTextLines)
	if err != nil {
		log.Fatalf("failed to create enhanced file: %s", err)
	}

	log.Printf("written %v events to enhanced file", eventCount+1)

	log.Println("finished")
}

func writeOutputFile(icsTextLines []string) error {
	file, err := os.Create(outputFilename)
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}

	defer file.Close()

	joinedString := strings.Join(icsTextLines, "\n")

	_, err = io.WriteString(file, joinedString)
	if err != nil {
		log.Fatalf("failed to write file content: %s", err)
	}

	return file.Sync()
}

func appendTrigger(icsTextLines []string, triggerText string) []string {
	icsTextLines = append(icsTextLines, "BEGIN:VALARM")
	icsTextLines = append(icsTextLines, "ACTION:DISPLAY")
	icsTextLines = append(icsTextLines, "DESCRIPTION:"+triggerText)
	icsTextLines = append(icsTextLines, "TRIGGER:"+triggerDuration)
	icsTextLines = append(icsTextLines, "END:VALARM")

	return icsTextLines
}

func appendDownloadNewCalendarDatesEvent(icsTextLines []string, eventDate time.Time) []string {
	log.Println("adding reminder for downloading new calendar file")

	uuid := uuid2.New().String()

	startDate := eventDate.Add(time.Hour * -(12 + 24))
	endDate := startDate.Add(time.Hour)

	startDateString := startDate.Format(icsDateLayout)
	endDateString := endDate.Format(icsDateLayout)

	timestamp := time.Now().Format(icsDateLayout)

	icsTextLines = append(icsTextLines, "BEGIN:VEVENT")
	icsTextLines = append(icsTextLines, "UID:"+uuid)
	icsTextLines = append(icsTextLines, "DTSTART:"+startDateString)
	icsTextLines = append(icsTextLines, "SEQUENCE:0")
	icsTextLines = append(icsTextLines, "TRANSP:OPAQUE")
	icsTextLines = append(icsTextLines, "DTEND:"+endDateString)
	icsTextLines = append(icsTextLines, "SUMMARY:Neuen Abfallkalender herunterladen")
	icsTextLines = append(icsTextLines, "CLASS:PUBLIC")
	icsTextLines = append(icsTextLines, "DTSTAMP:"+timestamp)
	icsTextLines = append(icsTextLines, "BEGIN:VALARM")
	icsTextLines = append(icsTextLines, "ACTION:DISPLAY")
	icsTextLines = append(icsTextLines, "DESCRIPTION:Neuen Abfallkalender herunterladen")
	icsTextLines = append(icsTextLines, "TRIGGER:-P0DT0H0M1S")
	icsTextLines = append(icsTextLines, "END:VALARM")
	icsTextLines = append(icsTextLines, "END:VEVENT")

	return icsTextLines
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
