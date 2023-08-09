package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type workerMessage struct {
	UUID     uuid.UUID
	Scale    int
	Noise    int
	FileName string
	Status   string
}

func (fh *functionHelper) checkUUID(uuid uuid.UUID) {
	searchUUID, err := fh.db.Query("SELECT uuid FROM users WHERE uuid = ? LIMIT 1", uuid)
	if err != nil {
		log.Println("UUID search", err)
		return
	}
	defer searchUUID.Close()

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
	fh.checkUUID(userID)

	scaleFactor, _ := strconv.Atoi(r.FormValue("scale"))
	noiseFactor, _ := strconv.Atoi(r.FormValue("noise"))
	_, handler, err := r.FormFile("imageFile")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileId, err := fh.uploadFile(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task := workerMessage{
		UUID:     userID,
		Scale:    scaleFactor,
		Noise:    noiseFactor,
		FileName: fmt.Sprint(fileId, "-", handler.Filename),
		Status:   "pending",
	}

	// Insert data into the table
	fmt.Println("Uploading", task.UUID, task.FileName, task.Scale, task.Noise)
	insertProcess, err := fh.db.Prepare("INSERT INTO file_processes (uuid, filename, scale, noise,status) VALUES (?,?,?,?,?)")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer insertProcess.Close()

	_, err = insertProcess.Exec(task.UUID, task.FileName, task.Scale, task.Noise, task.Status)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	taskJson, err := json.Marshal(task)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := fh.beanstalk.Put((taskJson), 1, time.Second*3, time.Second*0)
	if err != nil {
		fmt.Println("Error enqueueing task:", err)
	} else {
		fmt.Printf("Enqueued task with ID %d: %s\n", id, fmt.Sprint(fileId, handler.Filename))
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
