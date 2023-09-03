package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestMergeSvgFiles_EmptyImagesList(t *testing.T) {
	selectedImages := [][]byte{}

	buffer := MergeSvgFiles(&selectedImages)

	if len(buffer.Bytes()) != 175 {
		t.Errorf("Expected default 175 bytes buffer, got %d bytes", len(buffer.Bytes()))
	}
}

func TestMergeSvgFiles_OneImage(t *testing.T) {
	image := []byte("image")
	selectedImages := [][]byte{image}

	buffer := MergeSvgFiles(&selectedImages)

	if len(buffer.Bytes()) != 270 {
		t.Errorf("Expected default 270 bytes buffer, got %d bytes", len(buffer.Bytes()))
	}
}

func TestMergeSvgFiles_MultipleImages(t *testing.T) {
	// Arrange
	images := imagesListBuilder(3)

	// Act
	buffer := MergeSvgFiles(&images)
	x := 0
	for index := range images {
		expected := fmt.Sprintf(`<g xmlns="http://www.w3.org/2000/svg" id="id%d:id%d" transform="translate(%d.0,0.0)">image%d.svg`, index, index, x, index)
		if !strings.Contains(buffer.String(), expected) {
			t.Errorf("Expected %s got %s", expected, buffer.String())
		}
		x += 45
	}
}
