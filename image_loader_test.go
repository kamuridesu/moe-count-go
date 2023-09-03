package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAllImages_EmptyDir(t *testing.T) {
	basePath := "/tmp/empty"

	contents, err := LoadAllImages(basePath)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if contents != nil {
		t.Errorf("Expected 0 images, got %v", *contents)
	}
}

func TestLoadAllImages_OneImage(t *testing.T) {
	basePath := "one_image/"
	os.Mkdir(basePath, 0755)
	os.WriteFile(filepath.Join(basePath, "image.svg"), []byte{}, 0644)

	defer os.RemoveAll(basePath)

	contents, err := LoadAllImages(basePath)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(*contents) != 1 {
		t.Errorf("Expected 1 image, got %d", len(*contents))
	}
	if (*contents)[0] == nil {
		t.Errorf("Expected non-empty image content")
	}

}

func TestSelectImagesForRepr_OneNumber(t *testing.T) {
	imagesList := imagesListBuilder(0)
	number := 2

	selectedImages, err := SelectImagesForRepr(number, &imagesList)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(*selectedImages) != 6 {
		t.Errorf("Expected 6 images, got %d", len(*selectedImages))
	}
	image := ((*selectedImages)[5])
	if string(image) != "image2.svg" {
		t.Errorf("Expected last image to be image2.svg, got %v", string(image))
	}
}

func TestSelectImagesForRepr_MultipleNumbers(t *testing.T) {
	imagesList := imagesListBuilder(0)
	number := 123456

	selectedImages, err := SelectImagesForRepr(number, &imagesList)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(*selectedImages) != 6 {
		t.Errorf("Expected 6 images, got %d", len(*selectedImages))
	}
	for index, image := range *selectedImages {
		expectedImageName := fmt.Sprintf("image%d.svg", index+1)
		if string(image) != expectedImageName {
			t.Errorf("Expected image to be %s, got %v", expectedImageName, string(image))
		}
	}
}

func imagesListBuilder(size int) [][]byte {
	imagesList := [][]byte{}

	if size == 0 {
		size = 9
	}

	for i := 0; i < size+1; i++ {
		imagesList = append(imagesList, []byte(fmt.Sprintf("image%d.svg", i)))
	}
	return imagesList
}
