package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

func startUpscale(id uuid.UUID, file string, noise, scale int) (string, error) {
	fileName := fmt.Sprintf("./temp-images/%s-%s", id.String(), file)
	output := fmt.Sprintf("%s-[%d%%][%dx]%s", id.String(), scale*100, noise, file)
	log.Println(output)
	cmd := exec.Command("waifu2x-converter-cpp", "-i", fmt.Sprint(fileName), "--noise-level", fmt.Sprint(noise), "--scale-ratio", fmt.Sprint(scale), "-o", "./upscaled-images/"+output)
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
