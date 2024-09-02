package main

import (
	"encoding/json"
	"fmt"
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

func downloadFile(url string) error {	
    // Content-Type: application/x-chess-pgn
    // This content type indicates the type of parser needed to understand the data
    // Content-Disposition: attachment; filename="ChessCom_username_YYYYMM.pgn"
    // This disposition indicates that browser should download, not display, the result, and it suggests a filename based on the source archive

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	contentDisposition := resp.Header.Get("Content-Disposition")
	disposition, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return err
	}
	filename := params["filename"]
	log.Println("getting " + disposition + ": " + filename)

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// TODO: return err
func getMonthlyArchiveUrls(username string) []string {	
	url := "https://api.chess.com/pub/player/" + username + "/games/archives"
	chessClient := http.Client{
		Timeout: time.Second * 60, // Timeout in seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := chessClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	archives := PgnArchives{}
	jsonErr := json.Unmarshal(body, &archives)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return archives.MonthlyUrls
}

func main() {
	
	// TODO: switch for year and month or all
	// helloJson()
	
	// `{[https://api.chess.com/pub/player/barneytron/games/2022/11 https://api.chess.com/pub/player/barneytron/games/2022/12 https://api.chess.com/pub/player/barneytron/games/2023/01 https://api.chess.com/pub/player/barneytron/games/2023/02 https://api.chess.com/pub/player/barneytron/games/2023/03 https://api.chess.com/pub/player/barneytron/games/2023/04 https://api.chess.com/pub/player/barneytron/games/2023/05 https://api.chess.com/pub/player/barneytron/games/2023/06 https://api.chess.com/pub/player/barneytron/games/2023/07 https://api.chess.com/pub/player/barneytron/games/2023/08 https://api.chess.com/pub/player/barneytron/games/2023/09 https://api.chess.com/pub/player/barneytron/games/2023/10 https://api.chess.com/pub/player/barneytron/games/2023/11 https://api.chess.com/pub/player/barneytron/games/2023/12 https://api.chess.com/pub/player/barneytron/games/2024/01 https://api.chess.com/pub/player/barneytron/games/2024/02 https://api.chess.com/pub/player/barneytron/games/2024/03 https://api.chess.com/pub/player/barneytron/games/2024/04 https://api.chess.com/pub/player/barneytron/games/2024/05 https://api.chess.com/pub/player/barneytron/games/2024/06 https://api.chess.com/pub/player/barneytron/games/2024/07 https://api.chess.com/pub/player/barneytron/games/2024/08]}`

	username := "birdmaster3000"
	monthlyUrls := getMonthlyArchiveUrls(username)
	

	for i := 0; i < len(monthlyUrls) - 1; i++ {		
		u := fmt.Sprintf("%s%s", monthlyUrls[i], "/pgn")		
		downloadFile(u)
	}

	// for n, monthlyUrl := range monthlyUrls {
	// 	fmt.Println(n, monthlyUrl)
	// }
}
