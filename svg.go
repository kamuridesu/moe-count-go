package main

import (
	"bytes"
	"fmt"
)

const (
	SVG_HEADER_TEMPLATE = `<svg width="270" height="100" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">%s</svg>`
	G_HEADER_TEMPLATE   = `<g xmlns="http://www.w3.org/2000/svg" id="id%d:id%d" transform="translate(%d.0,0.0)">%s</g>`
)

func MergeSvgFiles(selectedImages *[][]byte) string {
	buffer := bytes.Buffer{}
	for index, image := range *selectedImages {
		buffer.WriteString(fmt.Sprintf(G_HEADER_TEMPLATE, index, index, index*45, string(image)))
	}
	newFile := fmt.Sprintf(SVG_HEADER_TEMPLATE, buffer.String())
	return newFile
}

func generateSVG(number int, images *[][]byte) string {
	selectedImages, err := SelectImagesForRepr(number, images)
	if err != nil {
		panic(err)
	}
	buffer := MergeSvgFiles(selectedImages)
	return buffer
}
