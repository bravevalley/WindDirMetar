package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	dataFile  = regexp.MustCompile(`sample`)
	validLine = regexp.MustCompile(`^JFK*`)
	windMetar = regexp.MustCompile(`\d*KT`)
	windDir = regexp.MustCompile(`\d{3}`)
)

// Get the text

// 2. Seperate the text and pass the stream lined text

func textProcessing(textProcessSender <-chan string, sendProcessedText chan<- string) {
	for rawStringData := range textProcessSender {

		lines := strings.Split(rawStringData, "\n")
		for _, singleLine := range lines {
	
			// JFK,2012-04-24 07:51,42.98,33.80,69.80,230.00,15.00,0.00,29.45,997.
			// 20,10.00,M,BKN,M,M,M,2700.00,M,M,M,M,M,M,M,M,M,M,35.01,KJFK 240751Z
			// 23015KT 10SM BKN027 06/01 A2945 RMK SLP972 T0061 0011,M
	
			// Ensure it is not a comment or heading
			if !validLine.MatchString(singleLine) {
				continue
			}
	
			// Send the data a Line at a time
			sendProcessedText <- singleLine
		}
	}

	// Close the channel feeding the processed string
	close(sendProcessedText)
	
}

// 3. Extract the information from stream line

func extractWindInfo(processTextReceiver <-chan string, extractedWindChan chan<- string) {
	// Get incoming data from the previous process

	// 70,10.00,M,FEW,OVC,M,M,4000.00,9500.00,M,M,-PL,M,M,M,M,M,M,31.71,KJFK 282251Z 07013KT

	for text := range processTextReceiver {

		// Match and extract the Wind Data needed for our calculation
		// [07013KT]
		windMetaData := windMetar.FindAllString(text, 1)

		// Extract the string from the wind metar []string
		for _, windData := range windMetaData {

			// Send the extracted wind information to the next thread - 07013KT
			extractedWindChan <- windData
		}
	}

	// Close the channel feeding the extracted wind data
	close(extractedWindChan)
}

func extractWindDir(extractedWindReceiver <-chan string, windDirChan chan<- int) {
	for windData := range extractedWindReceiver {
		winDirExtract := windDir.FindAllString(windData, 1)
		for _, windDirData := range winDirExtract {
			windDirDataInt, err := strconv.Atoi(windDirData)
			if err != nil {
				fmt.Println("Error parsing String to Interger")
			}
			fmt.Printf("%v\n\n", windDirDataInt)
		
				
			windDirChan <- windDirDataInt
		}
			
		
	}
}

// Collate information from the extracted information

func main() {
	textProcessSender := make(chan string)
	textProcessedSender := make(chan string)
	extractedWindSender := make(chan string)
	extractedWindDirSender := make(chan int)

	start := time.Now()

	go textProcessing(textProcessSender, textProcessedSender)

	go extractWindInfo(textProcessedSender, extractedWindSender)

	go extractWindDir(extractedWindSender, extractedWindDirSender)
	// 1. Get the text

	// Define the absolute path to the Directory where the data
	// is located
	absPath, err := filepath.Abs("./")

	// Don't just check errors, handle them gracefully. - Rob Pike
	if err != nil {
		fmt.Println("Unable to define path.")
	}

	// Read the Directory based on the file path generated above
	workingDir, err := os.ReadDir(absPath)
	if err != nil {
		fmt.Printf("Unable to read, check if Dir exists")
	}

	// Loop over the files gotten from reading the filepath
	for _, files := range workingDir {

		// Check which of the file matches the name of our data file
		if dataFile.MatchString(files.Name()) {

			// Define the path to the data file
			// NOTE: I could have defined the path and read the file
			// without reading the Dir first, I chose this route because
			// I want to learn more.

			pathfile, err := filepath.Abs(files.Name())
			if err != nil {
				fmt.Println("Unable to define path.")
			}

			// Read the data in the data file
			file, err := os.ReadFile(pathfile)
			if err != nil {
				fmt.Println("Error reading file")
			}

			text := string(file)

			textProcessSender <- text
		}

	}

	close(textProcessSender)

	for v := range extractedWindDirSender {
		fmt.Printf("%v", v)
		
	} 
	end := time.Since(start)

	fmt.Println(end)

}
