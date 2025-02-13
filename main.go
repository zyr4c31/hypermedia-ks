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

	"github.com/skip2/go-qrcode"
)

type SignalRequest struct {
	Input    bool   `json:"input"`
	Message  string `json:"message"`
	Messages string `json:"messages"`
}

type SignalMessages struct {
	Messages string `json:"messages"`
}

type SignalResponses struct {
	Responses string `json:"responses"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func main() {
	sm := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))
	sp := http.StripPrefix("/assets/", fs)

	sm.Handle("GET /assets/", sp)

	sm.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		templ := home()
		templ.Render(context.Background(), w)
	})

	var bulkMessage string
	aiResponses := make(chan string)

	sm.HandleFunc("GET /messages", func(w http.ResponseWriter, r *http.Request) {
		flusher := w.(http.Flusher)

		addSSEHeaders(w)

		clientGone := r.Context().Done()

		for {
			select {
			case <-clientGone:
				return
			case <-aiResponses:
				for response := range aiResponses {

					responses := SignalResponses{
						Responses: response,
					}

					responseBytes, _ := json.Marshal(responses)

					w.Write([]byte("event: datastar-merge-signals\n"))
					w.Write([]byte("retry: 1000\n"))
					w.Write([]byte(fmt.Sprintf(`data: signals %v`, string(responseBytes))))
					w.Write([]byte("\n\n\n"))
					flusher.Flush()
					break
				}
			}
		}
	})

	sm.HandleFunc("POST /chat", func(w http.ResponseWriter, r *http.Request) {

		reqBody, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		var postBody SignalRequest

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

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		resBytes := make([]byte, 512)

		for {
			_, err := res.Body.Read(resBytes)
			if err == io.EOF {
				break
			}

			resBytesReader := bytes.NewReader(resBytes)
			decoder := json.NewDecoder(resBytesReader)
			for {
				var message string
				tok, err := decoder.Token()
				if err == io.EOF {
					return
				}
				if tok == "" {
					return
				}
				if tok == "content" {
					decoder.Decode(&message)
					fmt.Printf("%v", message)

					bulkMessage += message
					aiResponses <- bulkMessage
					break
				}
			}

		}
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

	checkbox := make(chan bool)

	sm.HandleFunc("GET /checkbox", func(w http.ResponseWriter, r *http.Request) {
		rc := http.NewResponseController(w)

		addSSEHeaders(w)

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

	sm.HandleFunc("POST /checkbox", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		var postBody SignalRequest
		json.Unmarshal(body, &postBody)

		checkbox <- postBody.Input
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: sm,
	}

	fmt.Printf("server.Addr: http://192.168.3.112%v\n", server.Addr)
	// fmt.Printf("server.Addr: http://192.168.0.245%v\n", server.Addr)

	generateQR()

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func generateQR() {
	qrCode, _ := qrcode.New("http://192.168.3.112:8080", qrcode.Medium)
	fileName := fmt.Sprintf("assets/%v.png", "qrcode")
	qrCode.WriteFile(256, fileName)
	qrcode := qrCode.ToString(false)
	fmt.Printf("%v", qrcode)
}

func addSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
}

func mergeSignals(w http.ResponseWriter) {
}
