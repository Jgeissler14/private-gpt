package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`    // Add model field to the request
	Messages []ChatMessage `json:"messages"` // Fix the struct tag for Messages field
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest

	// Decode the form data into the ChatRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Model = "gpt-3.5-turbo"

	// Convert the payload to JSON
	reqJSON, err := json.Marshal(req) // Marshal the entire ChatRequest struct
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiKey := os.Getenv("open_ai_token")
	if apiKey == "" {
		http.Error(w, "OpenAI API key not found", http.StatusInternalServerError)
		return
	}

	apiURL := "https://api.openai.com/v1/chat/completions"

	// Send a POST request to the OpenAI API
	reqOpenAI, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqOpenAI.Header.Set("Content-Type", "application/json")
	reqOpenAI.Header.Set("Authorization", "Bearer "+apiKey)
	reqOpenAI.Header.Set("Content-Length", strconv.Itoa(len(reqJSON))) // Set the Content-Length header

	client := http.Client{}
	resp, err := client.Do(reqOpenAI)
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

	// Parse the response body
	var response struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the OpenAI response body to the ResponseWriter
	w.Write([]byte(response.Choices[0].Message.Content))
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
