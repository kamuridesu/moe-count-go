package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func LoadAllImages(basePath string) (*[][]byte, error) {
	files, e := os.ReadDir(basePath)
	if e != nil {
		log.Printf("Error! Could not read files from dir %s!", basePath)
		return nil, e
	}
	var contents [][]byte
	for _, file := range files {
		content, err := os.ReadFile(filepath.Join(basePath, file.Name()))
		if err != nil {
			log.Printf("Error! Could not read some file from dir, %e", err)
		}
		contents = append(contents, content)
	}
	return &contents, nil
}

func SelectImagesForRepr(number int, imagesList *[][]byte) (*[][]byte, error) {
	numberString := fmt.Sprintf("%06s", strconv.Itoa(number))
	var selectedImages [][]byte
	for _, _rune := range numberString {
		intRepr, err := strconv.Atoi(string(_rune))
		if err != nil {
			log.Print("Error! Cannot convert string to integer representation!")
			return nil, err
		}
		selectedImages = append(selectedImages, (*imagesList)[intRepr])
	}
	return &selectedImages, nil
}
