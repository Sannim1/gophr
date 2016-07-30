package main

import (
	"image"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
)

// Image represents an image
type Image struct {
	ID          string
	UserID      string
	Name        string
	Location    string
	Size        int64
	CreatedAt   time.Time
	Description string
}

const (
	imageIDLength = 10

	thumbnailWidth  = 400
	thumbnailHeight = thumbnailWidth

	previewImageWidth = 800
)

// A map of accepted mime types and their corresponding file extensions
var mimeExtensions = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
}

func isValidExtension(extension string) bool {
	for _, validExtension := range mimeExtensions {
		if validExtension == extension {
			return true
		}
	}

	return false
}

// NewImage creates a new image for a specified user
func NewImage(user *User) *Image {
	return &Image{
		ID:     GenerateID("img", imageIDLength),
		UserID: user.ID,
	}
}

// StaticRoute generates a path from which the image can be retrieved
func (image *Image) StaticRoute() string {
	return "/im/" + image.Location
}

// StaticThumbnailRoute generates a URL from which a thumbnail of the image can be retrieved
func (image *Image) StaticThumbnailRoute() string {
	return "/im/thumbnail/" + image.Location
}

// StaticPreviewRoute generates a URL from a preview of the image can be retrieved
func (image *Image) StaticPreviewRoute() string {
	return "/im/preview/" + image.Location
}

// ShowRoute generates a path to the image's display page
func (image *Image) ShowRoute() string {
	return "/image/" + image.ID
}

// CreateFromURL creates and persists an image from specified URL
func (image *Image) CreateFromURL(imageURL string) error {
	// Get the response from the URL
	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}

	// Make sure there's a 200 response
	if response.StatusCode != http.StatusOK {
		return errImageURLInvalid
	}

	defer response.Body.Close()

	// Ascertain the type of file downloaded
	mimeType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return errInvalidImageType
	}

	// Get an extension for the file
	extension, extensionFound := mimeExtensions[mimeType]
	if !extensionFound {
		return errInvalidImageType
	}

	// Get a name from the URL
	image.Name = filepath.Base(imageURL)
	image.Location = image.ID + extension

	// Open a file at target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}

	defer savedFile.Close()

	// copy the entire response body to the output file
	imageSize, err := io.Copy(savedFile, response.Body)
	if err != nil {
		return err
	}

	image.Size = imageSize

	// Create the various resized versions of the image
	err = image.CreateResizedImages()
	if err != nil {
		return err
	}

	// save image object to the database
	return globalImageStore.Save(image)
}

// CreateFromFile creates an image from an uploaded file
func (image *Image) CreateFromFile(file multipart.File, headers *multipart.FileHeader) error {

	// Move the file to an appropriate location, with an appropriate name
	image.Name = headers.Filename
	image.Location = image.ID + filepath.Ext(image.Name)

	// check that the file has a valid extension
	fileExtension := filepath.Ext(image.Name)
	if !isValidExtension(fileExtension) {
		return errInvalidImageType
	}

	// Open a file at target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}

	defer savedFile.Close()

	// Copy uploaded file to the target location
	imageSize, err := io.Copy(savedFile, file)
	if err != nil {
		return err
	}

	image.Size = imageSize

	// Create the various resized versions of the image
	err = image.CreateResizedImages()
	if err != nil {
		return err
	}

	// save image object to the database
	return globalImageStore.Save(image)
}

// CreateResizedImages generates predefined thumbnail images for the specified image
func (image *Image) CreateResizedImages() error {
	// Genereate an image from a file
	srcImage, err := imaging.Open("./data/images/" + image.Location)
	if err != nil {
		return err
	}

	// Create a channel to receive errors on
	errorChannel := make(chan error)

	// Process each resize
	go image.resizePreview(errorChannel, srcImage)
	go image.resizeThumbnail(errorChannel, srcImage)

	// Wait for images to finish resizing
	for i := 0; i < 2; i++ {
		err = <-errorChannel
	}

	return err
}

func (image *Image) resizeThumbnail(errorChannel chan error, srcImage image.Image) {

	destinationImage := imaging.Thumbnail(srcImage, thumbnailWidth, thumbnailHeight, imaging.Lanczos)

	destinationImagePath := "./data/images/thumbnail/" + image.Location

	errorChannel <- imaging.Save(destinationImage, destinationImagePath)
}

func (image *Image) resizePreview(errorChannel chan error, srcImage image.Image) {
	originalImageSize := srcImage.Bounds().Size()
	aspectRatio := float64(originalImageSize.Y) / float64(originalImageSize.X)

	previewImageHeight := int(float64(previewImageWidth) * aspectRatio)

	destinationImage := imaging.Resize(srcImage, previewImageWidth, previewImageHeight, imaging.Lanczos)

	destinationImagePath := "./data/images/preview/" + image.Location
	errorChannel <- imaging.Save(destinationImage, destinationImagePath)
}

