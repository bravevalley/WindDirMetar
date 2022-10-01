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
	WindDirections = map[string]int{
		"NE" : 45, 
		"E"  : 90,
		"SE" : 135,
		"S"  : 180,
		"SW" : 225,
		"W"  : 270,
		"NW" : 315,
		"N" : 360,
		

	}
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

	// Receive the processed windata from previous process
	for windData := range extractedWindReceiver {

		// Extract the wind data angles from the imported data
		winDirExtract := windDir.FindAllString(windData, 1)

		// Loop over each data 
		for _, windDirData := range winDirExtract {

			// Convert each wind data to Interger
			windDirDataInt, err := strconv.Atoi(windDirData)

			// To err is man
			if err != nil {
				fmt.Println("Error parsing String to Interger")
			}
		
			// Pass the data to the next thread
			windDirChan <- windDirDataInt
		}	
	}
	// Close the thread
	close(windDirChan)
}

func computeWindDir(windDirReceiver <-chan int, results chan<- map[string]int)  {
	// Retrieve the data via a channel from the previous process
	for windData := range windDirReceiver {

		// Increment the cardial map with each value within the scope
		switch {
		case windData < 45 : WindDirections["NE"]++
		case windData < 90 : WindDirections["E"]++
		case windData < 135 : WindDirections["SE"]++
		case windData < 180 : WindDirections["S"]++
		case windData < 225 : WindDirections["SW"]++
		case windData < 270: WindDirections["W"]++
		case windData < 315 : WindDirections["NW"]++
		case windData < 360 : WindDirections["N"]++
		}
	}

	// Send the result
	results <- WindDirections

	close(results)
}

// Collate information from the extracted information

func main() {
	textProcessSender := make(chan string)
	textProcessedSender := make(chan string)
	extractedWindSender := make(chan string)
	extractedWindDirSender := make(chan int)
	resultChannel := make(chan map[string]int)

	start := time.Now()

	go textProcessing(textProcessSender, textProcessedSender)

	go extractWindInfo(textProcessedSender, extractedWindSender)

	go extractWindDir(extractedWindSender, extractedWindDirSender)

	go computeWindDir(extractedWindDirSender, resultChannel)
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
	// Close the channel feeding the data from the file
	close(textProcessSender)


	// Retrieve result from the thread processing the data
	result := <- resultChannel

	// Print result to the console
	fmt.Println(result)
	end := time.Since(start)

	// Print the time elaspsed
	fmt.Println(end)

}
