package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/barneytron/pgnget/internal/args"
	"github.com/barneytron/pgnget/internal/client"
)

const (
	zeroYear  = "0000"
	zeroMonth = "00"
)

var (
	username *string
	year     *string
	month    *string
)

func init() {
	log.Println("init...")
	username = flag.String("username", "", "username")
	year = flag.String("year", zeroYear, "archive year")
	month = flag.String("month", zeroMonth, "archive month")
}

// TODO: fix usage
func printUsage() {
	fmt.Println("Usage: pgnget --username=user")
	fmt.Println("Usage: pgnget --username=user --year=2024 --month=09")
}

func validateArgs(username *string, month *string, year *string) {
	if *year == zeroYear && *month == zeroMonth && *username == "" {
		printUsage()
		os.Exit(1)
	}

    if *year == zeroYear && *month == zeroMonth && *username != "" {
        return
    }

	if !args.IsMonthValid(*month) || !args.IsYearValid(*year) || !args.IsUsernameValid(*username) {
		printUsage()
		os.Exit(1)
	}
}

type ByteCopier interface {
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
}

type FileCreator interface {
	Create(name string) (*os.File, error)
}

func main() {
	flag.Parse()

	*year = strings.TrimSpace(*year)
	*month = strings.TrimSpace(*month)
	*username = strings.TrimSpace(*username)

	validateArgs(username, month, year)

	log.Println(*month)
	log.Println(*year)

	byteCopier := client.NewCopier()
	fileCreator := client.NewCreator()

	chessComClient := client.NewChessClient(&http.Client{}, byteCopier, fileCreator)
	chessComPgnDownloader := client.NewChesscomPgnDownloader(*chessComClient)

	var err error    
	if *year == zeroYear && *month == zeroMonth {
		err = chessComPgnDownloader.DownloadAll(*username)
	} else {
		url := chessComPgnDownloader.CreatePgnByMonthUrl(*username, *year, *month)
		err = chessComPgnDownloader.DownloadByMonth(url)
	}
	if err != nil {
		log.Fatal(err)
	}
}
