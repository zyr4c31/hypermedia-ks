package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Signal struct {
	Input   bool   `json:"input"`
	Message string `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Message   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	DoneReason         string `json:"done_reason"`
	Done               bool   `json:"done"`
	TotalDuration      int    `json:"total_duration"`
	LoadDuration       int    `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int    `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int    `json:"eval_duration"`
}

func main() {
	sm := http.NewServeMux()

	sm.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})
	sm.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		templ := home()
		templ.Render(context.Background(), w)
	})

	var subscribers []int

	sm.HandleFunc("GET /sub", func(w http.ResponseWriter, r *http.Request) {
		for _, sub := range subscribers {
			fmt.Printf("sub: %v\n", sub)
		}
	})

	checkbox := make(chan bool)

	sm.HandleFunc("GET /checkbox", func(w http.ResponseWriter, r *http.Request) {
		rc := http.NewResponseController(w)

		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")

		clientGone := r.Context().Done()

		for {
			select {
			case <-clientGone:
				return
			case <-checkbox:
				for box := range checkbox {
					w.Write([]byte("event: datastar-merge-signals\n"))
					w.Write([]byte("retry: 1000\n"))
					w.Write([]byte(fmt.Sprintf(`data: signals {input: %v}`, box)))
					w.Write([]byte("\n\n\n"))
					rc.Flush()
				}
			}
		}
	})

	aiResponse := make(chan string, 1)
	var aiMessages []string

	sm.HandleFunc("GET /messages", func(w http.ResponseWriter, r *http.Request) {
		flusher := w.(http.Flusher)

		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")

		clientGone := r.Context().Done()

		for {
			select {
			case <-clientGone:
				return
			case <-aiResponse:
				for response := range aiResponse {
					fmt.Printf("response: %v\n", response)
					aiMessages = append(aiMessages, response)

					w.Write([]byte("event: datastar-merge-signals\n"))
					w.Write([]byte("retry: 1000\n"))
					w.Write([]byte(fmt.Sprintf(`data: signals {messages: %v}`, response)))
					w.Write([]byte("\n\n\n"))
					flusher.Flush()
				}
			}
		}
	})

	sm.HandleFunc("POST /chat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")

		flusher := w.(http.Flusher)

		reqBody, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		var postBody Signal

		json.Unmarshal(reqBody, &postBody)

		endpoint := "http://localhost:11434/api/chat"

		message := Message{
			Role:    "user",
			Content: postBody.Message,
		}

		request := Request{
			Model:    "deepseek-r1:1.5b",
			Messages: []Message{message},
		}

		jsonRequestBody, _ := json.Marshal(request)

		req, _ := http.NewRequest("POST", endpoint, bytes.NewReader(jsonRequestBody))

		res, _ := http.DefaultClient.Do(req)

		resBytes := make([]byte, 512)

		var response Response

		for {
			_, err := res.Body.Read(resBytes)
			if err == io.EOF {
				break
			}

			fmt.Printf("resBytes: %v\n", string(resBytes))

			isValid := json.Valid(resBytes)
			fmt.Printf("isValid: %v\n", isValid)
			err = json.Unmarshal(resBytes, &response)
			fmt.Printf("err: %v\n", err)

			w.Write([]byte("event: datastar-merge-fragments\n"))
			w.Write([]byte("retry: 1000\n"))
			w.Write([]byte(fmt.Sprintf(`data: fragments <p id="messages" >%v</p>`, string(resBytes))))
			w.Write([]byte("\n\n\n"))
			flusher.Flush()
		}

	})

	sm.HandleFunc("POST /checkbox", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		var postBody Signal
		json.Unmarshal(body, &postBody)

		checkbox <- postBody.Input
	})

	sm.HandleFunc("POST /user", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		var signal Signal

		json.Unmarshal(body, &signal)

		flusher := w.(http.Flusher)

		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")

		w.Write([]byte("event: datastar-merge-fragments\n"))
		w.Write([]byte("retry: 1000\n"))
		w.Write([]byte("\n\n\n"))
		flusher.Flush()
		time.Sleep(time.Second)
	})

	clients := 0

	sm.HandleFunc("GET /clients", func(w http.ResponseWriter, r *http.Request) {
		clients++
		flusher := w.(http.Flusher)

		clientDisconnect := r.Context().Done()

		for {
			select {
			case <-clientDisconnect:
				clients--
				return
			default:
				var builder strings.Builder
				builder.WriteString("event: datastar-merge-fragments\n")
				builder.WriteString("retry: 1000\n")
				word := "viewer"
				if clients > 1 {
					word = "viewers"
				}
				builder.WriteString(fmt.Sprintf(`data: fragments <span id="clients">%v %v</span>`, clients, word))
				builder.WriteString("\n\n\n")
				w.Write([]byte(builder.String()))
				flusher.Flush()
				time.Sleep(time.Second)
			}
		}
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: sm,
	}

	fmt.Printf("server.Addr: http://192.168.3.112%v\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
