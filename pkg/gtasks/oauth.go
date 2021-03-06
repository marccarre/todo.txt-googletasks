package gtasks

import (
	"context"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/marccarre/todo.txt-googletasks/pkg/gtasks/credentials"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	tasks "google.golang.org/api/tasks/v1"
)

// newOAuthClientFromCredentials creates a new OAuth client from the provided
// credentials.
func newOAuthClientFromCredentials(credentials *credentials.Credentials) (*http.Client, error) {
	config := &oauth2.Config{
		ClientID:     credentials.ClientID,
		ClientSecret: credentials.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{tasks.TasksScope},
	}
	return newOAuthClient(context.Background(), config)
}

func newOAuthClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	cacheFile, err := tokenCacheFile(config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OAuth token filename: %s", err)
	}
	token, err := tokenFromFile(cacheFile)
	if err != nil {
		log.WithField("error", err).WithField("file", cacheFile).Info("Cached OAuth token could not be found. Now trying to authenticate online...")
		token, err = tokenFromWeb(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("failed to authenticate online: %s", err)
		}
		err = saveToken(cacheFile, token)
		if err != nil {
			return nil, fmt.Errorf("failed to cache OAuth token: %s", err)
		}
	} else {
		log.WithField("token", token).WithField("file", cacheFile).Info("Using cached OAuth token.")
	}
	return config.Client(ctx, token), nil
}

func tokenCacheFile(config *oauth2.Config) (string, error) {
	hash := sha256.New()
	_, err := io.WriteString(hash, config.ClientID)
	if err != nil {
		return "", err
	}
	_, err = io.WriteString(hash, config.ClientSecret)
	if err != nil {
		return "", err
	}
	_, err = io.WriteString(hash, strings.Join(config.Scopes, " "))
	if err != nil {
		return "", err
	}
	suffix := fmt.Sprintf("%x", hash.Sum(nil))
	filename := fmt.Sprintf("todo.txt-googletasks_%s", suffix[:8])
	return filepath.Join(osUserCacheDir(), filename), nil
}

func osUserCacheDir() string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Caches")
	case "linux", "freebsd":
		return filepath.Join(os.Getenv("HOME"), ".cache")
	case "windows":
		return os.Getenv("USERPROFILE")
	}
	log.WithField("os", runtime.GOOS).Info("Defaulting OS user cache to current directory")
	return "."
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	token := new(oauth2.Token)
	err = gob.NewDecoder(file).Decode(token)
	return token, err
}

func tokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/favicon.ico" {
			http.Error(rw, "", 404)
			return
		}
		if req.FormValue("state") != randState {
			log.Errorf("State doesn't match: req = %#v", req)
			http.Error(rw, "", 500)
			return
		}
		if code := req.FormValue("code"); code != "" {
			fmt.Fprintf(rw, "<h1>Success</h1>Authorized.")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		log.Errorf("No code")
		http.Error(rw, "", 500)
	}))
	defer ts.Close()

	config.RedirectURL = ts.URL
	authURL := config.AuthCodeURL(randState)
	go openURL(authURL)
	log.Infof("Please authorise this app at: %s", authURL)
	code := <-ch
	log.Infof("Got code: %s", code)

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("Token exchange error: %v", err)
	}
	return token, nil
}

func openURL(url string) {
	bins := []string{"xdg-open", "google-chrome", "open"}
	for _, bin := range bins {
		err := exec.Command(bin, url).Run() // #nosec: URL is based on constant google.Endpoint
		if err == nil {
			return
		}
	}
	log.WithField("url", url).Errorf("Error opening URL in browser")
}

func saveToken(filepath string, token *oauth2.Token) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	return gob.NewEncoder(file).Encode(token)
}
