package main

import "net/http"

type Endpoint struct {
	Pattern string
	Handler http.HandlerFunc
}

func (fh *functionHelper) Endpoints() []Endpoint {
	return []Endpoint{
		{
			Pattern: "api/v1/upload",
			Handler: fh.FileProcessor,
		},
		{
			Pattern: "api/v1/update-status",
			Handler: fh.UpdateStatus,
		},
		{
			Pattern: "api/v1/get-images",
			Handler: fh.GetImages,
		},
		{
			Pattern: "api/v1/download-image",
			Handler: fh.DownloadImage,
		},
		{
			Pattern: "api/v1/ws",
			Handler: fh.HandleWebSocket,
		},
	}
}
