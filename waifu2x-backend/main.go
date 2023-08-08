// package main

// import (
// 	"fmt"
// 	"net/http"
// )

// func main() {
// 	fmt.Println("Upscale server activated")
// 	setupRoutes()
// }

// func setupRoutes() {
// 	http.HandleFunc("/upload", uploadFile)
// 	http.ListenAndServe(":5173", nil)
// }

package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/upscale", FileProcessor)

	// Start the HTTP server!
	fmt.Println("HTTP server listening on 42070")
	if err := http.ListenAndServe("0.0.0.0:42070", nil); err != nil {
		fmt.Println(err.Error())
	}
}
