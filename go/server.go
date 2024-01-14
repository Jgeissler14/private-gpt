package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong!")
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the ChatRequest to JSON
	reqJson, err := json.Marshal(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a POST request to the GPT4All model
	resp, err := http.Post("http://localhost:4891/v1/chat/completions", "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response body to the ResponseWriter
	w.Write(body)
}

func main() {
	fs := http.FileServer(http.Dir("../htmx")) // serve files from the current directory
	http.Handle("/", fs)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong!"))
	})

	http.HandleFunc("/chat", chatHandler)

	http.ListenAndServe(":8080", nil)
}
