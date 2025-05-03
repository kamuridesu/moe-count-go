package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestMergeSvgFiles_EmptyImagesList(t *testing.T) {
	selectedImages := [][]byte{}

	buffer := MergeSvgFiles(&selectedImages)

	if len(buffer) != 114 {
		t.Errorf("Expected default 114 len str, got %d", len(buffer))
	}
}

func TestMergeSvgFiles_OneImage(t *testing.T) {
	image := []byte("image")
	selectedImages := [][]byte{image}

	buffer := MergeSvgFiles(&selectedImages)

	if len(buffer) != 205 {
		t.Errorf("Expected default 205 len str, got %d", len(buffer))
	}
}

func TestMergeSvgFiles_MultipleImages(t *testing.T) {
	images := imagesListBuilder(3)
	buffer := MergeSvgFiles(&images)
	for index := range images {
		expected := fmt.Sprintf(`<g xmlns="http://www.w3.org/2000/svg" id="id%d:id%d" transform="translate(%d.0,0.0)">image%d.svg`, index, index, index*45, index)
		if !strings.Contains(buffer, expected) {
			t.Errorf("Expected %s got %s", expected, buffer)
		}
	}
}

func TestGenerateSVG(t *testing.T) {
	counter := 3
	images := imagesListBuilder(3)
	buffer := generateSVG(counter, &images)
	if len(buffer) != 698 {
		t.Errorf("Expected default 698 len str, got %d", len(buffer))
	}
}
