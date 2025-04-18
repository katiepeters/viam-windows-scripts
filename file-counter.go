package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	// Define command-line flags
	dirPath := flag.String("dir", ".", "Directory to monitor for file count")
	outputCSV := flag.String("output", "file_count.csv", "Path to output CSV file")
	flag.Parse()

	// Ensure the directory exists
	if _, err := os.Stat(*dirPath); os.IsNotExist(err) {
		fmt.Printf("Error: Directory '%s' does not exist\n", *dirPath)
		os.Exit(1)
	}

	// Create or open the CSV file
	file, err := os.OpenFile(*outputCSV, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Check if the file is new and write header if needed
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		os.Exit(1)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if fileInfo.Size() == 0 {
		// Write header to CSV
		if err := writer.Write([]string{"Timestamp", "FileCount"}); err != nil {
			fmt.Printf("Error writing CSV header: %v\n", err)
			os.Exit(1)
		}
		writer.Flush()
	}

	fmt.Printf("Monitoring directory: %s\n", *dirPath)
	fmt.Printf("Recording to: %s\n", *outputCSV)
	fmt.Println("Press Ctrl+C to stop...")

	// Use a ticker to execute every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			// Get the current number of files
			entries, err := os.ReadDir(*dirPath)
			if err != nil {
				fmt.Printf("Error reading directory: %v\n", err)
				continue
			}

			// Count only files, not directories
			fileCount := 0
			for _, entry := range entries {
				if !entry.IsDir() {
					fileCount++
				}
			}

			// Format timestamp (YYYY-MM-DD HH:MM:SS)
			timestamp := t.Format("2006-01-02 15:04:05")

			// Write to CSV
			if err := writer.Write([]string{timestamp, fmt.Sprintf("%d", fileCount)}); err != nil {
				fmt.Printf("Error writing to CSV: %v\n", err)
				continue
			}
			writer.Flush()

			// Display current count (optional)
			fmt.Printf("[%s] Found %d files in %s\n", timestamp, fileCount, *dirPath)
		}
	}
}
