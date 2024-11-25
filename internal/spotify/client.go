package spotify

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/browser"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const redirectURI = "http://localhost:4949/callback"

type Client struct {
	spotify *spotify.Client
	user    *spotify.PrivateUser
}

var (
	ch            = make(chan *spotify.Client)
	state         string
	codeVerifier  string
	codeChallenge string
	auth          = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopeUserTopRead,
		spotifyauth.ScopeUserReadRecentlyPlayed,
		spotifyauth.ScopePlaylistReadPrivate,
		spotifyauth.ScopePlaylistModifyPrivate,
		spotifyauth.ScopePlaylistModifyPublic,
	))
)

func init() {
	var err error
	state, err = generateRandomString(16)
	if err != nil {
		log.Fatalf("Failed to generate state: %v", err)
	}

	codeVerifier, err = generateRandomString(32)
	if err != nil {
		log.Fatalf("Failed to generate code verifier: %v", err)
	}

	hash := sha256.New()
	hash.Write([]byte(codeVerifier))
	codeChallenge = base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

func NewClient() *Client {

	// Check if the .spotify directory exists, if not, create it
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	spotifyDir := filepath.Join(homeDir, ".spotify")
	if _, err := os.Stat(spotifyDir); os.IsNotExist(err) {
		createSpotifyDir()
	}

	// Load the token from the config file
	var client *spotify.Client
	token, err := loadTokenFromConfig()
	if err != nil || token == nil {
		client, err = login()
		if err != nil {
			log.Fatalf("Failed to login: %v", err)
		}

		t, tErr := client.Token()
		if tErr != nil {
			log.Fatalf("Failed to get token: %v", tErr)
		}

		sErr := saveTokenToConfig(t)
		if sErr != nil {
			log.Fatalf("Failed to save token: %v", sErr)
		}

		user, uErr := client.CurrentUser(context.Background())
		if uErr != nil {
			log.Fatalf("Failed to get user: %v", err)
		}

		return &Client{
			spotify: client,
			user:    user,
		}
	}

	if !token.Valid() {
		// Refresh the token
		tokenSource := oauth2.StaticTokenSource(token)
		token, err = tokenSource.Token()
		if err != nil {
			log.Fatalf("Failed to refresh token: %v", err)
		}

	}
	// Initialize and return a new Client instance
	client = spotify.New(auth.Client(context.Background(), token))
	sErr := saveTokenToConfig(token)
	if sErr != nil {
		return nil
	}

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	return &Client{
		spotify: client,
		user:    user,
	}
}

func (c *Client) GetUser() (*spotify.PrivateUser, error) {
	user, err := c.spotify.CurrentUser(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting current user: %w", err)
	}
	return user, nil
}

func login() (*spotify.Client, error) {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":4949", nil)

	url := auth.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("client_id", os.Getenv("SPOTIFY_ID")),
	)
	err := browser.OpenURL(url)
	if err != nil {
		return nil, fmt.Errorf("opening browser: %w", err)
	}

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting current user: %w", err)
	}
	fmt.Println("You are logged in as:", user.ID)

	return client, nil
}

func saveTokenToConfig(token *oauth2.Token) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFile := filepath.Join(homeDir, ".spotify", "config")
	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(token); err != nil {
		return err
	}

	return nil

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}

func loadTokenFromConfig() (*oauth2.Token, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configFile := filepath.Join(homeDir, ".spotify", "config")
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var token oauth2.Token
	if err := json.NewDecoder(file).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func createSpotifyDir() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	spotifyDir := filepath.Join(homeDir, ".spotify")
	if _, err := os.Stat(spotifyDir); os.IsNotExist(err) {
		err := os.Mkdir(spotifyDir, 0755)
		if err != nil {
			log.Fatalf("failed to create .spotify directory: %v", err)
		}
	}
}

func RemoveConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFile := filepath.Join(homeDir, ".spotify", "config")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil
	}

	err = os.Remove(configFile)
	if err != nil {
		log.Fatalf("failed to remove config file: %v", err)
	}

	return nil
}

// generateRandomString generates a random string of the specified length
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
