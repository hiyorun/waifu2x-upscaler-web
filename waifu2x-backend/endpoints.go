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
	}
}
