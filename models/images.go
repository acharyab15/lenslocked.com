package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ImageService interface {
	Create(galleryId uint, r io.Reader, filename string) error
	ByGalleryID(galleryID uint) ([]string, error)
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (is *imageService) Create(galleryID uint,
	r io.Reader, filename string) error {
	path, err := is.mkImageDir(galleryID)
	if err != nil {
		return err
	}
	// Create a destination file
	dst, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy reader data to the destination file
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) ByGalleryID(galleryID uint) ([]string, error) {
	path := is.imageDir(galleryID)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}
	// Add a leading "/" to all image file paths
	for i := range strings {
		strings[i] = "/" + strings[i]
	}
	return strings, nil
}

// Going to need this when we know it is already made
func (is *imageService) imageDir(galleryID uint) string {
	// filepath.Join will return a path like:
	// images/galleries/123
	// We use filepath.Join instead of building the path
	// manually because the slashes and other characters
	// could vary between operating systems.
	return filepath.Join("images", "galleries",
		fmt.Sprintf("%v", galleryID))
}

func (is *imageService) mkImageDir(galleryID uint) (string, error) {
	// Create the directory to contain our images

	galleryPath := is.imageDir(galleryID)
	// Create our directory (and any necessary parent dirs)
	// using 0755 permissions.
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}
