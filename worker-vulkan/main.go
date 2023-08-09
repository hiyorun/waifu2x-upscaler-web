package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	beanstalk "github.com/beanstalkd/go-beanstalk"
	"github.com/google/uuid"
)

type workerMessage struct {
	UUID     uuid.UUID
	Scale    int
	Noise    int
	FileName string
	Status   string
}

type workerHelper struct {
	beanstalkAddr string
	backendAddr   string
	sharedFolder  string
}

func main() {
	beanstalkAddr := flag.String("beanstalk", "127.0.0.1:11300", "The beanstalk server address")
	backendAddr := flag.String("backend", "http://127.0.0.1:8080", "The backend server address")
	sharedFolder := flag.String("sharedFolder", "./", "Shared folder location")

	flag.Parse()

	conn, err := beanstalk.Dial("tcp", *beanstalkAddr)
	if err != nil {
		fmt.Println("Error connecting to Beanstalkd:", err)
		return
	}
	defer conn.Close()

	wh := &workerHelper{
		beanstalkAddr: *beanstalkAddr,
		backendAddr:   *backendAddr,
		sharedFolder:  *sharedFolder,
	}

	log.Println("Ready to take jobs")

	go wh.worker(conn)

	// Wait indefinitely
	select {}
}

func (wh *workerHelper) worker(conn *beanstalk.Conn) {
	for {
		id, body, err := conn.Reserve(5 * time.Second)
		if err != nil {
			continue
		}

		var task workerMessage
		err = json.Unmarshal(body, &task)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			continue
		}

		fmt.Printf("Processing %s scale %dx, and noise %d for session %s\n", task.FileName, task.Scale, task.Noise, task.UUID)
		wh.sendStatus(task.UUID, task.FileName, "processing")
		err = wh.startUpscale(task.FileName, task.Noise, task.Scale)
		if err != nil {
			wh.sendStatus(task.UUID, task.FileName, "failed")
			continue
		}
		log.Println("Done")
		wh.sendStatus(task.UUID, task.FileName, "done")

		conn.Delete(id)
	}
}

func (wh *workerHelper) sendStatus(uuid uuid.UUID, filename string, status string) {
	data := struct {
		UUID     string `json:"uuid"`
		FileName string `json:"filename"`
		Status   string `json:"status"`
	}{
		UUID:     uuid.String(),
		FileName: filename,
		Status:   status,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	response, err := http.Post(fmt.Sprint(wh.backendAddr, "/api/v1/update-status"), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
}

func (wh *workerHelper) startUpscale(filename string, noise, scale int) error {
	cmd := exec.Command("waifu2x-ncnn-vulkan", "-i", fmt.Sprint(wh.sharedFolder, "/temp-images/", filename), "-n", fmt.Sprint(noise), "-s", fmt.Sprint(scale), "-o", fmt.Sprint(wh.sharedFolder, "/upscaled-images/", filename))
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	if err := cmd.Run(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
