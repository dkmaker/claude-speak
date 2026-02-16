package elevenlabs_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dkmaker/claude-speak/internal/elevenlabs"
)

func TestSynthesize(t *testing.T) {
	fakeMP3 := []byte("fake-mp3-data")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/v1/text-to-speech/") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("xi-api-key") != "test-key" {
			t.Errorf("missing or wrong API key")
		}

		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), "hello world") {
			t.Errorf("body missing text: %s", body)
		}

		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write(fakeMP3)
	}))
	defer server.Close()

	client := elevenlabs.NewClient("test-key",
		elevenlabs.WithBaseURL(server.URL),
		elevenlabs.WithVoiceID("test-voice"),
		elevenlabs.WithModel("test-model"),
	)

	data, err := client.Synthesize("hello world")
	if err != nil {
		t.Fatalf("Synthesize: %v", err)
	}
	if string(data) != string(fakeMP3) {
		t.Errorf("got %q, want %q", data, fakeMP3)
	}
}

func TestSynthesizeNoAPIKey(t *testing.T) {
	client := elevenlabs.NewClient("")
	_, err := client.Synthesize("test")
	if err == nil {
		t.Error("expected error for empty API key")
	}
}
