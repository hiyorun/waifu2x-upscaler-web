package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
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
	sharedFolder  string
}

func main() {
	beanstalkAddr := flag.String("beanstalk", "127.0.0.1:11300", "The beanstalk server address")
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
		sharedFolder:  *sharedFolder,
	}

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
		wh.startUpscale(task.FileName, task.Noise, task.Scale)
		log.Println("Done")
		conn.Delete(id)
	}
}

func (wh *workerHelper) startUpscale(filename string, noise, scale int) (string, error) {
	output := fmt.Sprintf("[%d%%][%dx]%s", scale*100, noise, filename)
	log.Println(output)
	cmd := exec.Command("waifu2x-ncnn-vulkan", "-i", fmt.Sprint(wh.sharedFolder, "/temp-images/", filename), "-n", fmt.Sprint(noise), "-s", fmt.Sprint(scale), "-o", fmt.Sprint(wh.sharedFolder, "/upscaled-images/", output))
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	if err := cmd.Run(); err != nil {
		log.Println(err)
		return "", err
	}
	return output, nil
}
