package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	// Set up the HTTP server
	http.HandleFunc("/", handler)
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Define the target URL
	targetURL := "https://discord.com/api/webhooks/1235232270944702576/gmrZHO2D1Pz6aRx2MaGmgp3-qX00b2dM-Gb_5fBNmYWWYPwHrEpH3ZbnuLntupWWjZ1S" + r.URL.Path

	// Create a new HTTP request that replicates the original request
	outReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Error creating request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy all headers from the original request
	copyHeaders(outReq.Header, r.Header)

	// Make the request to the target URL
	client := &http.Client{}
	resp, err := client.Do(outReq)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy all headers from the response
	copyHeaders(w.Header(), resp.Header)

	// Write the status code from the response
	w.WriteHeader(resp.StatusCode)

	// Copy the response body to the response writer
	io.Copy(w, resp.Body)
}

// copyHeaders copies headers from source to destination, needed to preserve header values
func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
