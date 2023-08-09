package main

import "net/http"

type Endpoint struct {
	Pattern string
	Handler http.HandlerFunc
}

func (fh *functionHelper) Endpoints() []Endpoint {
	return []Endpoint{
		{
			Pattern: "/upload",
			Handler: fh.FileProcessor,
		},
		{
			Pattern: "/update-status",
			Handler: fh.UpdateStatus,
		},
		{
			Pattern: "/get-images",
			Handler: fh.GetImages,
		},
		{
			Pattern: "/download-image",
			Handler: fh.DownloadImage,
		},
		{
			Pattern: "/ws",
			Handler: fh.HandleWebSocket,
		},
	}
}
