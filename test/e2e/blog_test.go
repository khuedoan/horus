package e2e

import (
	"net/http"
	"testing"
)

func TestBlog(t *testing.T) {
	resp, err := http.Get("https://khuedoan.com") // TODO get domain name automatically
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code to be 200, but got %d", resp.StatusCode)
	}
}
