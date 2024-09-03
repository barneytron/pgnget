package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"

	//"strconv"
	"time"
)

type PgnArchives struct {
	MonthlyUrls []string `json:"archives"`
}

func downloadPgn(url string) error {
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
func getAllMonthlyArchiveUrls(username string) []string {
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

var (
	username *string
	year     *string
	month    *string
)

func init() {
	log.Println("init...")
	username = flag.String("username", "", "username")
	year = flag.String("year", "0000", "archive year")
	month = flag.String("month", "00", "archive month")
}

func isYearValid(year string) bool {
	if len(year) != 4 {
		return false
	}
	// TODO: validate all chars digits
	return true
}

func isMonthValid(month string) bool {
	if len(month) != 2 {
		return false
	}
	// TODO: validate all chars digits
	return true
}

func getMultiGamePgn(username string, month string) {

}

func isUsernameValid(username string) bool {
	if len(username) < 1 {
		return false
	}
	// TODO: validate all chars digits
	return true
}

func printUsage() {
	fmt.Println("Usage: pgnget.go --username=user --year=2024 --month=09")
	fmt.Println("Usage: pgnget.go --username=user --year=all --month=all")
}

func main() {
	log.Println("main...")

	flag.Parse()
	// TODO:

	values := flag.Args()

	*year = strings.TrimSpace(*year)
	*month = strings.TrimSpace(*month)
	*username = strings.TrimSpace(*username)

	log.Println(len(values))
	log.Println(*year)

	if *year == "0000" && *month == "00" && *username == "" {
		printUsage()
		// flag.PrintDefaults()
		os.Exit(1)
	}

	if !isMonthValid(*month) || !isYearValid(*year) || !isUsernameValid(*username) {
		printUsage()
		os.Exit(1)
	}

	// log.Println("year is " +  strconv.Itoa(*year));
	// log.Println("month is " +  strconv.Itoa(*month));
	// os.Exit(0)
	// helloJson()

	// `{[https://api.chess.com/pub/player/barneytron/games/2022/11 https://api.chess.com/pub/player/barneytron/games/2022/12 https://api.chess.com/pub/player/barneytron/games/2023/01 https://api.chess.com/pub/player/barneytron/games/2023/02 https://api.chess.com/pub/player/barneytron/games/2023/03 https://api.chess.com/pub/player/barneytron/games/2023/04 https://api.chess.com/pub/player/barneytron/games/2023/05 https://api.chess.com/pub/player/barneytron/games/2023/06 https://api.chess.com/pub/player/barneytron/games/2023/07 https://api.chess.com/pub/player/barneytron/games/2023/08 https://api.chess.com/pub/player/barneytron/games/2023/09 https://api.chess.com/pub/player/barneytron/games/2023/10 https://api.chess.com/pub/player/barneytron/games/2023/11 https://api.chess.com/pub/player/barneytron/games/2023/12 https://api.chess.com/pub/player/barneytron/games/2024/01 https://api.chess.com/pub/player/barneytron/games/2024/02 https://api.chess.com/pub/player/barneytron/games/2024/03 https://api.chess.com/pub/player/barneytron/games/2024/04 https://api.chess.com/pub/player/barneytron/games/2024/05 https://api.chess.com/pub/player/barneytron/games/2024/06 https://api.chess.com/pub/player/barneytron/games/2024/07 https://api.chess.com/pub/player/barneytron/games/2024/08]}`

	//	username := "birdmaster3000"

	if strings.EqualFold(*year, "all") && strings.EqualFold(*month, "all") {

		monthlyUrls := getAllMonthlyArchiveUrls(*username)

		for i := 0; i < len(monthlyUrls)-1; i++ {
			u := fmt.Sprintf("%s%s", monthlyUrls[i], "/pgn")
			downloadPgn(u)
		}
	}

	// TODO: get all archives for a year
	if strings.EqualFold(*month, "all") {
		log.Println("get all monthly archives for specified year")
	}
	
	// https://api.chess.com/pub/player/erik/games/2009/10/pgn
	url := "https://api.chess.com/pub/player/" + *username + "/games/" + *year + "/" + *month + "/pgn"
	log.Println("downloading pgn: " + url)
	err := downloadPgn(url)
	if err != nil {
		log.Fatal(err)
	}

	// for n, monthlyUrl := range monthlyUrls {
	// 	fmt.Println(n, monthlyUrl)
	// }
}
