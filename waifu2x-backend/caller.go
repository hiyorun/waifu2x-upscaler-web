package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func FileProcessor(w http.ResponseWriter, r *http.Request) {
	ScaleFactor, _ := strconv.Atoi(r.FormValue("scale"))
	NoiseFactor, _ := strconv.Atoi(r.FormValue("noise"))
	_, handler, err := r.FormFile("imageFile")
	if err != nil {
		log.Println(err)
		return
	}
	fileId, err := uploadFile(r)
	if err != nil {
		log.Println(err)
		return
	}
	filename, err := startUpscale(fileId, handler.Filename, NoiseFactor, ScaleFactor)
	if err != nil {
		log.Println(err)
		return
	}

	returnFile, err := os.ReadFile("./upscaled-images/" + filename)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Connection", "keep-alive")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Disposition", "attachment; filename="+filename)
	w.Write(returnFile)
}
