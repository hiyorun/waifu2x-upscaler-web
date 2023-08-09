package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func (fh *functionHelper) uploadFile(r *http.Request) (uuid.UUID, error) {
	fmt.Println("File Upload Endpoint Hit")

	// Size Limit
	// // Parse our multipart form, 10 << 20 specifies a maximum
	// // upload of 10 MB files.
	// r.ParseMultipartForm(10 << 20)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println(err)
		return uuid.Nil, err
	}

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("imageFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return uuid.Nil, err
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fileId := uuid.New()
	filePath, err := filepath.Abs(fmt.Sprintf(fh.sharedFolder, "./temp-images/%s-%s", fileId.String(), handler.Filename))
	if err != nil {
		fmt.Println("filepath", err)
		return uuid.Nil, err
	}
	tempFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return uuid.Nil, err
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return uuid.Nil, err
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	fmt.Println("Uploaded", fileId)
	return fileId, nil
}
