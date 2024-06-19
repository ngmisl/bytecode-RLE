package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// ReadFile reads the content of the file at the given path
func ReadFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes the content to the file at the given path
func WriteFile(path string, data string) error {
	return ioutil.WriteFile(path, []byte(data), 0644)
}

// RLEncode compresses the input string using Run-Length Encoding, excluding blocks of digits
func RLEncode(input string) string {
	var encoded strings.Builder
	count := 1
	for i := 1; i < len(input); i++ {
		if isDigit(input[i]) && isDigit(input[i-1]) {
			continue
		} else if input[i] == input[i-1] && !isDigit(input[i]) {
			count++
		} else {
			if isDigit(input[i-1]) {
				start := i - 1
				for start > 0 && isDigit(input[start-1]) {
					start--
				}
				encoded.WriteString("|" + input[start:i] + "|")
			} else {
				encoded.WriteString(string(input[i-1]) + strconv.Itoa(count))
			}
			count = 1
		}
	}
	if isDigit(input[len(input)-1]) {
		start := len(input) - 1
		for start > 0 && isDigit(input[start-1]) {
			start--
		}
		encoded.WriteString("|" + input[start:] + "|")
	} else {
		encoded.WriteString(string(input[len(input)-1]) + strconv.Itoa(count))
	}
	return encoded.String()
}

// RLDecode decompresses the RLE encoded string
func RLDecode(input string) string {
	var decoded strings.Builder
	i := 0
	for i < len(input) {
		if input[i] == '|' {
			end := strings.IndexByte(input[i+1:], '|') + i + 1
			decoded.WriteString(input[i+1 : end])
			i = end + 1
		} else {
			char := input[i]
			countStart := i + 1
			for countStart < len(input) && isDigit(input[countStart]) {
				countStart++
			}
			count, err := strconv.Atoi(input[i+1 : countStart])
			if err != nil {
				log.Fatalf("Invalid RLE string: %v", err)
			}
			decoded.WriteString(strings.Repeat(string(char), count))
			i = countStart
		}
	}
	return decoded.String()
}

// isDigit checks if a byte is a digit
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
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
