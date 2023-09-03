package main

import (
	"bytes"
	"fmt"
	"strings"

	svg "github.com/ajstarks/svgo"
)

func MergeSvgFiles(selectedImages *[][]byte) *bytes.Buffer {
	buffer := bytes.Buffer{}
	canvas := svg.New(&buffer)
	canvas.Start(270.0, 100.0)
	xCoord := 0
	for index, image := range *selectedImages {
		imageWithoutTags := RemoveTagsFromImage(string(image))
		canvas.Writer.Write([]byte(fmt.Sprintf("<g xmlns=\"http://www.w3.org/2000/svg\" id=\"id%d:id%d\" transform=\"translate(%d.0,0.0)\">", index, index, xCoord)))
		canvas.Writer.Write([]byte(imageWithoutTags))
		canvas.Writer.Write([]byte("</g>\r\n"))
		xCoord += 45
	}
	canvas.End()
	return &buffer
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

func generateSVG(number int, images *[][]byte) *bytes.Buffer {
	selectedImages, err := SelectImagesForRepr(number, images)
	if err != nil {
		panic(err)
	}
	buffer := MergeSvgFiles(selectedImages)
	return buffer
}
