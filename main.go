package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	dataFile = regexp.MustCompile(`sample`)
)

// Get the text

// Seperate the text and pass the stream lined text

// Extract the information from stream line

// Collate information from the extracted information

func main() {

	// Get the text

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

			// Informal unit test
			fmt.Println("We found him")

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

			fmt.Println(string(file))
		}
	}

}
