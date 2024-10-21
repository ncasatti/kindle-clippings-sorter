package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const inputFolder = "./input"
const inputFile = "clippings.txt"
const outputFolder = "./output"

func main() {

	fmt.Println("Starting sorting clippings...")

	// Open the file
	inputF := inputFolder + "/" + inputFile
	fmt.Println("Opening file:", inputF)
	file, err := os.Open(inputF)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a map to hold quotes by book title
	books := make(map[string][]string)

	// Read line by line
	scanner := bufio.NewScanner(file)
	var currentBook string
	var currentQuote strings.Builder

	// Initialize the book title to true
	var title = true

	fmt.Println("Reading file...")
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		// Remove BOM if present
		line = strings.TrimPrefix(line, "\ufeff")

		// When separator is found, add the current quote to the book's list
		if line == "==========" {
			if currentBook != "" && currentQuote.Len() > 0 {
				books[currentBook] = append(books[currentBook], currentQuote.String())
				currentQuote.Reset()
			}
			title = true
			continue
		}

		// If we find a book title (lines not starting with "-" or empty), treat as new book
		if !strings.HasPrefix(line, "-") && line != "" && title {
			// Prints only one time the current book
			if line != currentBook {
				fmt.Println("New book:", line)
			}
			currentBook = line
			title = false
		} else if strings.HasPrefix(line, "-") || line == "" {
			// Skip metadata lines
			continue
		} else {
			// Collect the quote
			currentQuote.WriteString(line + "\n")
		}
	}

	// Write each book's quotes to a new file
	for book, quotes := range books {
		bookName := strings.ReplaceAll(book, " ", "_")
		fileName := outputFolder + "/" + bookName + ".txt"
		fmt.Println("Writing file:", fileName)

		// Ensure the output directory exists
		if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
			err = os.MkdirAll(outputFolder, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating output directory:", err)
				return
			}
		}
		// Create the file
		outFile, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating file:", err)
			continue
		}
		defer outFile.Close()

		// Write the quotes to the file
		for _, quote := range quotes {
			outFile.WriteString(quote + "\n")
		}

		fmt.Printf("File created for book: %s\n", book)
	}

}
