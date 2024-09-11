package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	JPEG_SOI   = 0xFFD8
	JPEG_APP1  = 0xFFE1
	EXIF_MAGIC = "Exif\x00\x00"
)

func main() {
	// Prompt the user for the image path
	fmt.Print("Enter the path to the image: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	imagePath := scanner.Text()

	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("No jpg file at", (imagePath))
		return
	}
	defer file.Close()

	fileStat, _ := file.Stat()
	data := make([]byte, fileStat.Size())
	_, err = file.Read(data)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if binary.BigEndian.Uint16(data[0:2]) != JPEG_SOI {
		fmt.Println("Not a valid JPEG file")
		return
	}

	offset := 2
	for offset < len(data) {
		marker := binary.BigEndian.Uint16(data[offset : offset+2])
		size := binary.BigEndian.Uint16(data[offset+2 : offset+4])

		if marker == JPEG_APP1 {
			if string(data[offset+4:offset+10]) == EXIF_MAGIC {
				fmt.Println("EXIF data found")
				return
			}
		}
		offset += int(size) + 2
	}

	fmt.Println("No EXIF data found")
}
