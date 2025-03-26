package client

import (
	"errors"
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
func (d ChesscomPgnDownloader) doWork(id int, urlsChan <-chan string, errorsChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range urlsChan {
		url = fmt.Sprintf("%s/pgn", url)
		log.Println("pool worker", id, "downloading from ", url)
		err := d.chessComClient.DownloadPgn(url)
		if err != nil {
			errorsChan <- errors.New(fmt.Sprintf("worker %d: error processing job %s", id, url))
		}
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
	urlsChan := make(chan string, numWorkers)
	errorChan := make(chan error, numWorkers)
	var wg sync.WaitGroup

	// Create worker pool based on number of CPU cores
	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go d.doWork(w, urlsChan, errorChan, &wg)
	}

	// Send url to the urls channel
	for _, monthlyUrl := range monthlyUrls {
		urlsChan <- monthlyUrl
	}
	close(urlsChan)

	// Wait for all workers to complete
	wg.Wait()
	close(errorChan)
	for err := range errorChan {
		fmt.Println("Error:", err)
	}
	return nil
}
