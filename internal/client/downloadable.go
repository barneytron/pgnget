package client

import (
	"fmt"
	"log"
	"sync"
)

const WorkerCount = 4

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

// TODO: deal with errors
func (d ChesscomPgnDownloader) doWork(id int, urls <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range urls {
		url = fmt.Sprintf("%s/pgn", url)
		log.Println("pool worker", id, "downloading from ", url)
		d.chessComClient.DownloadPgn(url)
	}
}

// TODO: use concurrency
func (d ChesscomPgnDownloader) DownloadAll(username string) error {
	monthlyUrls, err := d.chessComClient.GetAllMonthlyArchiveUrls(username)
	if err != nil {
		return err // TODO: wrap error here
	}

	//numWorkers := runtime.NumCPU()
	numWorkers := WorkerCount
	urls := make(chan string, numWorkers)
	var wg sync.WaitGroup

	// Create worker pool based on number of CPU cores
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go d.doWork(w, urls, &wg)
	}

	// Send jobs to the job channel
	for _, monthlyUrl := range monthlyUrls {
		urls <- monthlyUrl
	}
	close(urls)
	// Wait for all workers to complete
	wg.Wait()

	return nil
}
