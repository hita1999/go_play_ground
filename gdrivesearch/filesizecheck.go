package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Credentials struct {
	Type 	   string `json:"type"`
	ProjectID  string `json:"project_id"`
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey string `json:"private_key"`
	ClientEmail string `json:"client_email"`
	ClientID string `json:client_id"`
	AuthURI string `json:"auth_uri"`
	TokenURI string `json:"token_uri"`
	AuthProviderX509 string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL string `json:"client_x509_cert_url"`
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	//OAtuh2.0の認証をリクエスト
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)
	fmt.Println("Enter the authorization code: ")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	//認証コードを使用してトークンを取得
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token
}

func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tok)
	}
	return config.Client(ctx, tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	//キャッシュからトークンを読み込む
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	return t, err
}

func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func printFiles(srv *drive.Service, query string) {
	//Google Drive APIを初期化
	if srv == nil {
		log.Fatal("Unable to retrieve Drive client.")
	}

	//ファイルの一覧を取得
	r, err := srv.Files.List().Q(query).PageSize(1000).Fields("nextPageToken, files(id, name, size)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		fmt.Println("Files:")
		for _, f := range r.Files {
			fmt.Printf("%s (%s): %d bytes\n", f.Name, f.Id, f.Size)
		}
	}
}

func main() {
	creds := &Credentials{}
	err := json.Unmarshal([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")), creds)
	if err != nil {
		log.Fatalf("failed to unmarshal credentials: %v", err)
	}

	//Google Oauth2.0認証を使用してトークンを取得
	ctx := context.Background()
	config, err := google.ConfigFromJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")), drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("failed to create OAuth2 config: %v", err)
	}
	client := getClient(ctx, config)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	//ファイル一覧を取得
	files, err := srv.Files.List().Fields("nextPageToken, files(id, name, size)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	//ファイル情報を表示
	for _, file := range files.Files {
		fmt.Printf("File name: %s, ID: %s, Size: %d bytes\n", file.Name, file.Id, file.Size)
	}
}
