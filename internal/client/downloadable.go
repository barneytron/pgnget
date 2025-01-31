package client

import (
	"fmt"
	"log"
)

type Downloadable interface {
	CreatePgnByMonthUrl(username string, year string, month string) string
	DownloadByMonth(url string) error
	DownloadAll(username string) error
}

type ChesscomPgnDownloader struct {
	chessComClient ChessComClient
}

func NewChesscomPgnDownloader(chessComClient ChessComClient) Downloadable {
	return &ChesscomPgnDownloader{
		chessComClient: chessComClient,
	}
}

func (d ChesscomPgnDownloader) CreatePgnByMonthUrl(username string, year string, month string) string {
	// https://api.chess.com/pub/player/erik/games/2009/10/pgn
	return fmt.Sprintf("https://api.chess.com/pub/player/%s/games/%s/%s/pgn", username, year, month)
}

func (d ChesscomPgnDownloader) DownloadByMonth(url string) error {
	log.Println("downloading pgn: " + url)
	err := d.chessComClient.DownloadPgn(url)
	if err != nil {
		return err
	}
	return nil
}

// TODO: use concurrency
func (d ChesscomPgnDownloader) DownloadAll(username string) error {
	monthlyUrls, err := d.chessComClient.GetAllMonthlyArchiveUrls(username)
	if err != nil {
		return err // TODO: wrap error here
	}
	for _, monthlyUrl := range monthlyUrls {
		url := fmt.Sprintf("%s%s", monthlyUrl, "/pgn")
		d.chessComClient.DownloadPgn(url)
	}
	return nil
}
