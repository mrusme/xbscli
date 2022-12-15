package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mrusme/xbscli/lz77"
	"golang.org/x/crypto/pbkdf2"
)

type Bookmark struct {
	ID       int        `json:"id"`
	Title    string     `json:"title"`
	URL      string     `json:"url,omitempty"`
	Children []Bookmark `json:"children,omitempty"`
}

type BookmarksResponse struct {
	Bookmarks   string    `json:"bookmarks"`
	LastUpdated time.Time `json:"lastUpdated"`
	Version     string    `json:"version"`
}

var serverURL string
var syncID string
var password string
var format string

func main() {
	flag.StringVar(&serverURL, "s", "", "server URL (required)")
	flag.StringVar(&syncID, "i", "", "sync ID (required)")
	flag.StringVar(&password, "p", "", "sync password (required)")
	flag.StringVar(&format, "f", "json", "format: json, pretty")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if serverURL == "" || syncID == "" || password == "" {
		flag.Usage()
		os.Exit(1)
	}

	response, err := http.Get(serverURL + "/bookmarks/" + syncID)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var bookmarksResponse BookmarksResponse
	if err := json.Unmarshal(responseBody, &bookmarksResponse); err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	key := pbkdf2.Key([]byte(password), []byte(syncID), 250000, 32, sha256.New)

	ciphertext, err := base64.StdEncoding.DecodeString(bookmarksResponse.Bookmarks)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	iv := ciphertext[:16]
	encdata := ciphertext[16:]

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	aesgcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	encdata, err = aesgcm.Open(nil, iv, encdata, nil)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	bms := lz77.DecompressBlockToString(encdata)

	switch format {
	case "json":
		fmt.Printf("%s\n", bms)
		os.Exit(0)
	case "pretty":
		var bookmarks []Bookmark
		if err := json.Unmarshal([]byte(bms), &bookmarks); err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		printBookmarks(bookmarks, 0)
		os.Exit(0)
	}
}

func printBookmarks(bookmarks []Bookmark, level int) {
	for _, bookmark := range bookmarks {
		fmt.Printf("%s%s\n", strings.Repeat(" ", level), bookmark.Title)
		if bookmark.Children != nil && len(bookmark.Children) > 0 {
			printBookmarks(bookmark.Children, level+2)
		} else {
			fmt.Printf("%s%s\n", strings.Repeat(" ", level), bookmark.URL)
			fmt.Println()
		}
	}
}
