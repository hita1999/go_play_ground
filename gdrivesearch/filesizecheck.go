package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	//"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var (
	config *oauth2.Config
	//client *http.Client
)

func tokenFromFile(tokenFile string) (*oauth2.Token, error) {
	file, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	token := &oauth2.Token{}
	err = gob.NewDecoder(file).Decode(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func getTokenFromWeb() (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	ctx := context.Background()
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}

	return token, nil
}

func saveToken(tokenFile string, token *oauth2.Token) {
	file, err := os.Create(tokenFile)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer file.Close()

	err = gob.NewEncoder(file).Encode(token)
	if err != nil {
		log.Fatalf("Unable to encode token to file: %v", err)
	}
}

// func getClient(config *oauth2.Config, token *oauth2.Token) *http.Client {
// 	return config.Client(oauth2.NoContext, token)
// }

func main() {
	//Load credentials from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	//Set up Oauth2 config
	ctx := context.Background()
	config = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes: 	 []string{drive.DriveReadonlyScope},
		Endpoint:     google.Endpoint,
		RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
	}

	//Load token from file or request authorization
	token, err := tokenFromFile("token.json")
	if err != nil {
		token, err = getTokenFromWeb()
		if err != nil {
			log.Fatalf("Unable to retrieve token from web: %v", err)
		}
		saveToken("token.json", token)
	}

	//Initialize Google Drive API client
	//client := getClient(config, token)
	srv, err := drive.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	// List files in the user's Google Drive
    r, err := srv.Files.List().Fields("nextPageToken, files(name, size)").Do()
    if err != nil {
        log.Fatalf("Unable to retrieve files: %v", err)
    }
    if len(r.Files) == 0 {
        fmt.Println("No files found.")
    } else {
        fmt.Println("Files:")
        for _, f := range r.Files {
            fmt.Printf("%s (%s bytes)\n", f.Name, strconv.FormatInt(f.Size, 10))
        }
    }
}