package clients

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
)

var (
	conf          *oauth2.Config
	authCodeCh    = make(chan *string)
	serverCloseCh = make(chan struct{})
)

func StartAuthFlow(ctx context.Context, url string, state string) (string, error) {
	openBrowser(url)

	// Start a local web server to listen for the OAuth callback
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		callbackHandler(w, r, state)
	})

	server := &http.Server{Addr: ":7999"}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Wait for auth code
	authCode := <-authCodeCh

	if authCode == nil || len(*authCode) == 0 {
		return "", fmt.Errorf("Failed to Authenticate")
	}

	// Shutdown the server
	server.Shutdown(ctx)
	close(serverCloseCh)

	return *authCode, nil
}

func callbackHandler(w http.ResponseWriter, r *http.Request, state string) {
	if r.URL.Query().Get("state") != state {
		fmt.Fprintf(w, "Authentication failed. Please try again")
		authCodeCh <- nil
		return
	}
	code := r.URL.Query().Get("code")

	if code == "" {
		fmt.Fprintf(w, "Authentication was cancelled. You can close this window")
		authCodeCh <- nil
		return
	}

	fmt.Fprintf(w, "Authentication successful. You can close this window")
	authCodeCh <- &code
}

// openBrowser tries to open the browser with a given URL.
func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
