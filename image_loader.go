package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func RemoveTagsFromImage(svg string) string {
	__svg := ""
	for _, line := range strings.Split(strings.ReplaceAll(svg, "\r\n", "\n"), "\n") {
		if (!strings.HasPrefix(line, "<?xml")) &&
			(!strings.HasPrefix(line, "<!DOCTYPE")) &&
			(!strings.HasSuffix(line, ".dtd\">")) &&
			(!strings.HasPrefix(line, "<svg")) &&
			(!strings.HasSuffix(line, "\"1.1\">")) &&
			(!strings.HasPrefix(line, "</svg")) {
			__svg += line + "\r\n"
		}
	}
	return __svg
}

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
		contents = append(contents, []byte(RemoveTagsFromImage(string(content))))
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
