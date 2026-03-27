package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/genai"
)

const defaultModel = "gemini-2.5-flash"

type chatRequest struct {
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func main() {
	addr := getenv("LISTEN_ADDR", ":9090")
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is required")
	}
	model := getenv("GEMINI_MODEL", defaultModel)

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("genai client: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handleChat(w, r, client, model)
	})

	log.Printf("backend listening on %s (model=%s)", addr, model)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func handleChat(w http.ResponseWriter, r *http.Request, client *genai.Client, model string) {
	var req chatRequest
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if len(req.Messages) == 0 {
		http.Error(w, "messages required", http.StatusBadRequest)
		return
	}

	contents := make([]*genai.Content, 0, len(req.Messages))
	for _, m := range req.Messages {
		var role genai.Role
		switch m.Role {
		case "user":
			role = genai.RoleUser
		case "model":
			role = genai.RoleModel
		default:
			http.Error(w, "role must be user or model", http.StatusBadRequest)
			return
		}
		contents = append(contents, genai.NewContentFromText(m.Content, role))
	}

	ctx := r.Context()
	resp, err := client.Models.GenerateContent(ctx, model, contents, nil)
	if err != nil {
		log.Printf("gemini GenerateContent: %v", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	text := resp.Text()
	if text == "" {
		http.Error(w, "empty model response", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"reply": text})
}
