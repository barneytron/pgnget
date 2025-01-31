package client

import (
	"encoding/json"
	"errors"
	"log"
	"mime"
	"net/http"
)

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type ChessComClient struct {
	httpClient  HttpClient
	byteCopier  Copyable
	fileCreator Creatable
}

func NewChessClient(httpClient HttpClient, byteCopier Copyable, fileCreator Creatable) *ChessComClient {
	return &ChessComClient{
		httpClient:  httpClient,
		byteCopier:  byteCopier,
		fileCreator: fileCreator,
	}
}

type PgnArchives struct {
	MonthlyUrls []string `json:"archives"`
}

func (c ChessComClient) DownloadPgn(url string) error {
	// Content-Type: application/x-chess-pgn
	// This content type indicates the type of parser needed to understand the data
	// Content-Disposition: attachment; filename="ChessCom_username_YYYYMM.pgn"
	// This disposition indicates that browser should download, not display, the result, and it suggests a filename based on the source archive

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return err // TODO: wrap error
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	disposition, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return err // TODO: wrap error
	}

	var filename string
	var exists bool
	if filename, exists = params["filename"]; !exists {
		return errors.New("filename does not exist in content disposition")
	}

	out, err := c.fileCreator.Create(filename)
	if err != nil {
		return err // TODO: wrap error
	}
	defer out.Close()

	_, err = c.byteCopier.Copy(out, resp.Body)
	if err != nil {
		log.Printf("error occurred copying %s to file system", filename)
		return err // TODO: wrap error
	}

	log.Println("copied " + disposition + ": " + filename)
	return nil // TODO: wrap error
}

func (c ChessComClient) GetAllMonthlyArchiveUrls(username string) ([]string, error) {
	url := "https://api.chess.com/pub/player/" + username + "/games/archives"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err // TODO: wrap error
	}

	var archives PgnArchives
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&archives)

	if err != nil {
		log.Println(err) // TODO: wrap error
		return nil, err
	}

	return archives.MonthlyUrls, nil
}
