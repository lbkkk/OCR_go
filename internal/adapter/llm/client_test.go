package llm

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestComplete(t *testing.T) {
	var captured chatRequest
	var gotAuth string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &captured); err != nil {
			t.Errorf("server: unmarshal request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"choices":[{"message":{"content":"# Title"}}]}`)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "secret-key", "qwen-test")
	got, err := client.Complete(context.Background(), []Message{
		System("you are helpful"),
		User("convert this"),
	})
	if err != nil {
		t.Fatalf("Complete returned error: %v", err)
	}

	if got != "# Title" {
		t.Errorf("content = %q, want %q", got, "# Title")
	}
	if gotAuth != "Bearer secret-key" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer secret-key")
	}
	if captured.Model != "qwen-test" {
		t.Errorf("model = %q, want %q", captured.Model, "qwen-test")
	}
	if len(captured.Messages) != 2 {
		t.Fatalf("messages len = %d, want 2", len(captured.Messages))
	}
}

func TestToWireMessageVision(t *testing.T) {
	msg := Message{
		Role:   RoleUser,
		Text:   "describe",
		Images: []Image{{Data: []byte("PNG"), Format: "png"}},
	}

	wire := toWireMessage(msg)

	parts, ok := wire.Content.([]contentPart)
	if !ok {
		t.Fatalf("Content type = %T, want []contentPart", wire.Content)
	}
	if len(parts) != 2 {
		t.Fatalf("parts len = %d, want 2 (text + image)", len(parts))
	}
	if parts[0].Type != "text" || parts[0].Text != "describe" {
		t.Errorf("part[0] = %+v, want text part", parts[0])
	}
	if parts[1].Type != "image_url" || parts[1].ImageURL == nil {
		t.Fatalf("part[1] = %+v, want image_url part", parts[1])
	}
	if !strings.HasPrefix(parts[1].ImageURL.URL, "data:image/png;base64,") {
		t.Errorf("image url = %q, want data URL prefix", parts[1].ImageURL.URL)
	}
}

func TestToWireMessageTextOnly(t *testing.T) {
	wire := toWireMessage(User("hello"))

	content, ok := wire.Content.(string)
	if !ok {
		t.Fatalf("Content type = %T, want string", wire.Content)
	}
	if content != "hello" {
		t.Errorf("content = %q, want %q", content, "hello")
	}
}
