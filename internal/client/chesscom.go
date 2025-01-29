package client

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
)

type Copyable interface {
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
}

type Copier struct{}

func NewCopier() Copyable {
	return &Copier{}
}

func (copier Copier) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

type Creatable interface {
	Create(name string) (*os.File, error)
}

type Creator struct{}

func NewCreator() Creatable {
	return &Creator{}
}

func (creator Creator) Create(name string) (*os.File, error) {
	return os.Create(name)
}

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
		return err
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	disposition, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return err
	}

	var filename string
	var exists bool
	if filename, exists = params["filename"]; !exists {
		return errors.New("filename does not exist in content disposition")
	}

	out, err := c.fileCreator.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = c.byteCopier.Copy(out, resp.Body)
	if err != nil {
		log.Printf("error occurred copying %s to file system", filename)
		return err
	}

	log.Println("copied " + disposition + ": " + filename)
	return nil
}

func (c ChessComClient) GetAllMonthlyArchiveUrls(username string) ([]string, error) {
	url := "https://api.chess.com/pub/player/" + username + "/games/archives"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var archives PgnArchives
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&archives)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return archives.MonthlyUrls, nil
}
