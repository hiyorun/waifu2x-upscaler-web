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
	"sync"
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

	wh.worker(conn)

	// Wait indefinitely
	select {}
}

func (wh *workerHelper) touchBeanstalkJob(stopCh <-chan struct{}, wg *sync.WaitGroup, conn *beanstalk.Conn, jobID uint64) {
	defer wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Use the Touch function to reset the job's TTR (time-to-run)
			if err := conn.Touch(jobID); err != nil {
				fmt.Println("Error touching job:", err)
			}
			fmt.Println("Job is taking longer than usual, letting Beanstalk know")
		case <-stopCh:
			fmt.Println("touchBeanstalkJob() is stopping...")
			return
		}
	}
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
		// Create a channel to stop the touchBeanstalkJob goroutine when needed
		stopLetThemKnow := make(chan struct{})
		var wg sync.WaitGroup

		// Start the touchBeanstalkJob goroutine
		wg.Add(1)
		go wh.touchBeanstalkJob(stopLetThemKnow, &wg, conn, id)

		// Perform the upscale operation
		wh.sendStatus(task.UUID, task.FileName, "processing")
		err = wh.startUpscale(task.FileName, task.Noise, task.Scale)

		// Stop the touchBeanstalkJob goroutine
		close(stopLetThemKnow)
		wg.Wait() // Wait for the goroutine to finish

		// Handle the result of the upscale operation
		if err != nil {
			wh.sendStatus(task.UUID, task.FileName, "failed, retrying")
			conn.Release(id, 1, 3*time.Second)
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
	cmd := exec.Command("waifu2x-converter-cpp", "-i", fmt.Sprint(wh.sharedFolder, "/temp-images/", filename), "--noise-level", fmt.Sprint(noise), "--scale-ratio", fmt.Sprint(scale), "-o", fmt.Sprint(wh.sharedFolder, "/upscaled-images/", filename))
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	if err := cmd.Run(); err != nil {
		log.Println(err)
		return err
	}

	if err := fileExists(fmt.Sprint(wh.sharedFolder, "/upscaled-images/", filename)); err != nil {
		return err
	}

	return nil
}

func fileExists(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}
