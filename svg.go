package main

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	SVG_HEADER_TEMPLATE = `<svg width="270" height="100" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">%s</svg>`
	G_HEADER_TEMPLATE   = `<g xmlns="http://www.w3.org/2000/svg" id="id%d:id%d" transform="translate(%d.0,0.0)">%s</g>`
)

func MergeSvgFiles(selectedImages *[][]byte) string {
	buffer := bytes.Buffer{}
	for index, image := range *selectedImages {
		buffer.WriteString(fmt.Sprintf(G_HEADER_TEMPLATE, index, index, index*45, RemoveTagsFromImage(string(image))))
	}
	newFile := fmt.Sprintf(SVG_HEADER_TEMPLATE, buffer.String())
	return newFile
}

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

func generateSVG(number int, images *[][]byte) string {
	selectedImages, err := SelectImagesForRepr(number, images)
	if err != nil {
		panic(err)
	}
	buffer := MergeSvgFiles(selectedImages)
	return buffer
}
