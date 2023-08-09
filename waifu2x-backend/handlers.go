package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type workerMessage struct {
	UUID     uuid.UUID
	Scale    int
	Noise    int
	FileName string
	Status   string
	Name     string
}

type processedImage struct {
	UUID   uuid.UUID `json:"uuid"`
	Status string    `json:"status"`
	Name   string    `json:"name"`
}

type Data struct {
	UUID     uuid.UUID `json:"uuid"`
	FileName string    `json:"filename"`
	Status   string    `json:"status"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
		Name:     handler.Filename,
	}

	// Insert data into the table
	fmt.Println("Uploading", task.UUID, task.FileName, task.Scale, task.Noise)
	insertProcess, err := fh.db.Prepare("INSERT INTO file_processes (uuid, filename, scale, noise,status,name) VALUES (?,?,?,?,?,?)")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer insertProcess.Close()

	_, err = insertProcess.Exec(task.UUID, task.FileName, task.Scale, task.Noise, task.Status, task.Name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fh.triggerUpdate(task.UUID)

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

func (fh *functionHelper) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	var data Data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Update from worker: Filename=%s, Status=%s\n", data.FileName, data.Status)
	fileStatus, err := fh.db.Prepare("UPDATE file_processes	SET status=? WHERE filename=?")
	if err != nil {
		log.Println("Update file status", err)
		return
	}
	defer fileStatus.Close()

	_, err = fileStatus.Exec(data.Status, data.FileName)
	if err != nil {
		log.Println(err)
		return
	}
	fh.triggerUpdate(data.UUID)
}

func (fh *functionHelper) GetImages(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")

	// Query all projects from the database
	rows, err := fh.db.Query("SELECT uuid,name,status FROM file_processes WHERE uuid=?", uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
	}

	defer rows.Close()

	var projects []processedImage
	for rows.Next() {
		var project processedImage
		err := rows.Scan(&project.UUID, &project.Name, &project.Status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		projects = append(projects, project)
	}

	// Convert projects slice to JSON
	jsonData, err := json.Marshal(projects)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	// Set response headers and write JSON data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (fh *functionHelper) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	fh.webSocket = conn
	defer fh.webSocket.Close()

	for {
		messageType, message, err := fh.webSocket.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received message: %s\n", message)

		err = fh.webSocket.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}

func (fh *functionHelper) triggerUpdate(uuid uuid.UUID) {
	// Send a message to the connected WebSocket client
	err := fh.webSocket.WriteMessage(websocket.TextMessage, []byte(uuid.String()))
	if err != nil {
		log.Println(err)
	}
}