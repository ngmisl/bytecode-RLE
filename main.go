package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// ReadFile reads the content of the file at the given path
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes the content to the file at the given path
func WriteFile(path string, data string) error {
	return os.WriteFile(path, []byte(data), 0644)
}

// RLEncode compresses the input string using Run-Length Encoding
func RLEncode(input string) string {
	var encoded strings.Builder
	count := 1
	for i := 1; i < len(input); i++ {
		if input[i] == input[i-1] {
			count++
		} else {
			encoded.WriteString(string(input[i-1]) + "-" + strconv.Itoa(count) + "-")
			count = 1
		}
	}
	encoded.WriteString(string(input[len(input)-1]) + "-" + strconv.Itoa(count))
	return encoded.String()
}

// RLDecode decompresses the RLE encoded string
func RLDecode(input string) string {
	var decoded strings.Builder
	parts := strings.Split(input, "-")
	for i := 0; i < len(parts); i += 2 {
		if i+1 < len(parts) {
			char := parts[i]
			count, err := strconv.Atoi(parts[i+1])
			if err != nil {
				log.Fatalf("Invalid RLE string: %v", err)
			}
			decoded.WriteString(strings.Repeat(char, count))
		}
	}
	return decoded.String()
}

func main() {
	// Read data from the file
	filePath := "data.txt"
	data, err := ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Trim the data to remove any surrounding white space or new lines
	data = strings.TrimSpace(data)

	// Print original data
	fmt.Println("Original Data:")
	fmt.Println(data)

	// Compress data using RLE
	encodedData := RLEncode(data)
	fmt.Println("\nCompressed Data:")
	fmt.Println(encodedData)

	// Calculate sizes
	originalSize := len(data)
	compressedSize := len(encodedData)
	compressionRatio := float64(originalSize) / float64(compressedSize)

	// Print size statistics
	fmt.Printf("\nOriginal Size: %d bytes\n", originalSize)
	fmt.Printf("Compressed Size: %d bytes\n", compressedSize)
	fmt.Printf("Compression Ratio: %.2f\n", compressionRatio)

	// Save compressed data to file if ratio is better than 1
	if compressionRatio > 1 {
		compressedFilePath := "compressed.txt"
		err = WriteFile(compressedFilePath, encodedData)
		if err != nil {
			log.Fatalf("Failed to write compressed file: %v", err)
		}
		fmt.Printf("\nCompressed data saved to %s\n", compressedFilePath)
	} else {
		fmt.Println("\nCompression ratio is not better than 1. Compressed data not saved.")
	}

	// Decompress the data to verify
	decodedData := RLDecode(encodedData)
	fmt.Println("\nDecompressed Data:")
	fmt.Println(decodedData)

	// Check if decompressed data matches original data
	if data == decodedData {
		fmt.Println("\nDecompression successful, original data matches decompressed data.")
	} else {
		fmt.Println("\nDecompression failed, original data does not match decompressed data.")
	}
}
