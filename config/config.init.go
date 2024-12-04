package config

import (
	"html/template"
	"net/http"
	"fmt"
	"os"
	"io"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Host string
	Port int
}

type Database struct {
	Username, Password, Host string
	Port                     int
	Dbname                   string
}

func DbConfig() Database {
	database := Database{
		Username: "root",
		Password: "kama1983",
		Port:     3306,
		Host:     "localhost",
		Dbname:   "lawoffice",
	}
	return database
}

func ServerConfig() Server {
	server := Server{
		Host: "localhost",
		Port: 8000,
	}
	return server
}

func RenderTemplate(c echo.Context, templateName string, data interface{}) error {

	t, err := template.ParseFiles("views/dashboard.html" ,"views/" + templateName)

	if err != nil {

		return c.String(http.StatusNotFound, err.Error())

	}

	return t.Execute(c.Response().Writer, data)

}

// Function to check if the file is an image
func isImageFile(file *os.File) bool {

	// Read the first 512 bytes to determine the content type

	buffer := make([]byte, 512)

	if _, err := file.Read(buffer); err != nil {

		return false

	}

	contentType := http.DetectContentType(buffer)

	// Check if the content type starts with "image/"

	return contentType[:6] == "image/"

}

// Handler for the file upload

func uploadFile(c echo.Context) error {

	// Get the file from the form input

	file, err := c.FormFile("file")

	if err != nil {
		return c.String(http.StatusBadRequest, "Unable to get file from form data")
	}
	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Unable to open file")
	}

	defer src.Close()

	// Create a temporary file to check its content type

	tempFile, err := os.CreateTemp("", "upload-")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Unable to create temporary file")
	}
	defer os.Remove(tempFile.Name()) // Clean up the temporary file
	// Copy the uploaded file to the temporary file
	if _, err := io.Copy(tempFile, src); err != nil {
		return c.String(http.StatusInternalServerError, "Unable to save file")
	}

	// Reset the file pointer to the beginning of the file
	if _, err := tempFile.Seek(0, 0); err != nil {
		return c.String(http.StatusInternalServerError, "Unable to seek in temporary file")
	}

	// Check if the uploaded file is an image
	if !isImageFile(tempFile) {
		return c.String(http.StatusBadRequest, "Uploaded file is not a valid image")
	}

	// Create a new file in the current directory with the original name

	dst, err := os.Create(file.Filename)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Unable to create file")
	}
	defer dst.Close()
	// Copy the uploaded file to the new file
	if _, err := tempFile.Seek(0, 0); err != nil {
		return c.String(http.StatusInternalServerError, "Unable to seek in temporary file")
	}

	if _, err := io.Copy(dst, tempFile); err != nil {

		return c.String(http.StatusInternalServerError, "Unable to save file")

	}

	return c.String(http.StatusOK, fmt.Sprintf("File uploaded successfully: %s", file.Filename))

}