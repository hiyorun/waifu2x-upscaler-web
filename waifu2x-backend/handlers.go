package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (fh *functionHelper) checkUUID(uuid uuid.UUID) {
	searchUUID, err := fh.db.Query("SELECT uuid FROM users WHERE uuid = ? LIMIT 1", uuid)
	defer searchUUID.Close()
	if err != nil {
		log.Println("UUID search", err)
		return
	}

	var id string
	for searchUUID.Next() {
		err = searchUUID.Scan(&id)
		log.Println(id)
		if err != nil {
			log.Println("Scan", err)
			return
		}
		if id == uuid.String() {
			log.Println("UUID match", id, uuid)
			return
		}
	}

	insertUUID, err := fh.db.Prepare("INSERT INTO users (uuid) VALUES (?)")
	if err != nil {
		log.Println("Prepare insert UUID", err)
		return
	}
	defer insertUUID.Close()

	_, err = insertUUID.Exec(uuid)
	if err != nil {
		log.Println(err)
		return
	}
}

func (fh *functionHelper) FileProcessor(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.FormValue("uuid"))
	if err != nil {
		log.Println(err)
		return
	}
	scaleFactor, _ := strconv.Atoi(r.FormValue("scale"))
	noiseFactor, _ := strconv.Atoi(r.FormValue("noise"))
	_, handler, err := r.FormFile("imageFile")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fh.checkUUID(userID)

	fileId, err := uploadFile(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Insert data into the table
	fmt.Println("Uploading", userID, fmt.Sprint(fileId, handler.Filename), scaleFactor, noiseFactor)
	insertProcess, err := fh.db.Prepare("INSERT INTO file_processes (uuid, filename, scale, noise) VALUES (?,?,?,?)")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer insertProcess.Close()

	_, err = insertProcess.Exec(userID, fmt.Sprint(fileId, handler.Filename), scaleFactor, noiseFactor)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	// filename, err := startUpscale(fileId, handler.Filename, NoiseFactor, ScaleFactor)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// returnFile, err := os.ReadFile("./upscaled-images/" + filename)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// w.Header().Set("Connection", "keep-alive")
	// w.Header().Add("Access-Control-Allow-Origin", "*")
	// w.Header().Add("Content-Disposition", "attachment; filename="+filename)
	// w.Write(returnFile)
}
