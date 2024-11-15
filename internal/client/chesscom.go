package client

import (
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"time"
)

type PgnArchives struct {
	MonthlyUrls []string `json:"archives"`
}

func DownloadPgn(url string) error {
	// Content-Type: application/x-chess-pgn
	// This content type indicates the type of parser needed to understand the data
	// Content-Disposition: attachment; filename="ChessCom_username_YYYYMM.pgn"
	// This disposition indicates that browser should download, not display, the result, and it suggests a filename based on the source archive

	resp, err := http.Get(url)
	if err != nil {
		return err
	}	

	contentDisposition := resp.Header.Get("Content-Disposition")
	disposition, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return err
	}
	filename := params["filename"]
	// TODO: validate if filename exists
	log.Println("getting " + disposition + ": " + filename)

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func GetAllMonthlyArchiveUrls(username string) ([]string, error) {
	url := "https://api.chess.com/pub/player/" + username + "/games/archives"
	chessClient := http.Client{
		Timeout: time.Second * 60,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, getErr := chessClient.Do(req)
	if getErr != nil {
		log.Println(getErr)
		return nil, err
	}

	var archives PgnArchives
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&archives)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return archives.MonthlyUrls, nil
}
