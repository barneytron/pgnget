package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"pgnget/internal/args"
	"pgnget/internal/client"
)

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

func printUsage() {
	fmt.Println("Usage: pgnget.go --username=user --year=2024 --month=09")
	fmt.Println("Usage: pgnget.go --username=user --year=all --month=all")
}

func validateArgs(username *string, month *string, year *string) {
	if *year == "0000" && *month == "00" && *username == "" {
		printUsage()
		// flag.PrintDefaults()
		os.Exit(1)
	}

	if !args.IsMonthValid(*month) || !args.IsYearValid(*year) || !args.IsUsernameValid(*username) {
		printUsage()
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	*year = strings.TrimSpace(*year)
	*month = strings.TrimSpace(*month)
	*username = strings.TrimSpace(*username)

	validateArgs(username, month, year)

	log.Println(*month)
	log.Println(*year)

	// `{[https://api.chess.com/pub/player/barneytron/games/2022/11 https://api.chess.com/pub/player/barneytron/games/2022/12 https://api.chess.com/pub/player/barneytron/games/2023/01 https://api.chess.com/pub/player/barneytron/games/2023/02 https://api.chess.com/pub/player/barneytron/games/2023/03 https://api.chess.com/pub/player/barneytron/games/2023/04 https://api.chess.com/pub/player/barneytron/games/2023/05 https://api.chess.com/pub/player/barneytron/games/2023/06 https://api.chess.com/pub/player/barneytron/games/2023/07 https://api.chess.com/pub/player/barneytron/games/2023/08 https://api.chess.com/pub/player/barneytron/games/2023/09 https://api.chess.com/pub/player/barneytron/games/2023/10 https://api.chess.com/pub/player/barneytron/games/2023/11 https://api.chess.com/pub/player/barneytron/games/2023/12 https://api.chess.com/pub/player/barneytron/games/2024/01 https://api.chess.com/pub/player/barneytron/games/2024/02 https://api.chess.com/pub/player/barneytron/games/2024/03 https://api.chess.com/pub/player/barneytron/games/2024/04 https://api.chess.com/pub/player/barneytron/games/2024/05 https://api.chess.com/pub/player/barneytron/games/2024/06 https://api.chess.com/pub/player/barneytron/games/2024/07 https://api.chess.com/pub/player/barneytron/games/2024/08]}`

	//	username := "birdmaster3000"

	if strings.EqualFold(*year, "all") && strings.EqualFold(*month, "all") {

		monthlyUrls, err := client.GetAllMonthlyArchiveUrls(*username)
		if err != nil {
			log.Fatal(err)
		}

		for _, monthlyUrl := range monthlyUrls {
			u := fmt.Sprintf("%s%s", monthlyUrl, "/pgn")
			client.DownloadPgn(u)
		}
	} else if strings.EqualFold(*month, "all") {
		// TODO: get all archives for a year
		log.Println("get all monthly archives for specified year")
	} else {
		// https://api.chess.com/pub/player/erik/games/2009/10/pgn
		url := "https://api.chess.com/pub/player/" + *username + "/games/" + *year + "/" + *month + "/pgn"
		log.Println("downloading pgn: " + url)
		err := client.DownloadPgn(url)
		if err != nil {
			log.Fatal(err)
		}
	}
}
