package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

type movieResult struct {
	Title      string
	Year       string
	Rated      string
	Released   string
	Runtime    string
	Genre      string
	Director   string
	Writer     string
	Actors     string
	Plot       string
	Language   string
	Country    string
	Awards     string
	Poster     string
	Metascore  string
	ImdbRating string
	ImdbVotes  string
	ImdbID     string
	Type       string
	DVD        string
	BoxOffice  string
	Production string
	Website    string
	Response   string
	Error      string
}

const omdbURL = "http://www.omdbapi.com?t="

func formatOmdbURL(title string) string {
	fmtTitle := strings.Replace(title, " ", "+", -1)
	url := omdbURL + fmtTitle + "&y=&plot=full&r=json"
	return url
}

func setUnderline(length int) string {
	var underline bytes.Buffer
	for i := 0; i < length; i++ {
		underline.WriteString("-")
	}
	return underline.String()
}

func getMovieJSON(url string) (*movieResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := new(movieResult)
	err = json.NewDecoder(resp.Body).Decode(&r)

	if err != nil {
		return nil, err
	}

	if r.Response == "False" {
		return r, errors.New(r.Error)
	}

	return r, nil
}

func main() {
	var title string
	app := cli.NewApp()
	app.Name = "movie"
	app.Version = "0.0.1"
	app.Usage = "search for a movie"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "name, n",
			Usage:       "type a movie `TITLE`",
			Destination: &title,
		},
	}
	app.Action = func(c *cli.Context) error {
		url := formatOmdbURL(title)
		color.Set(color.FgYellow)
		fmt.Println("Searching OMDB for your movie...")
		fmt.Println()
		movie, err := getMovieJSON(url)

		if err != nil {
			return err
		}

		fmt.Println("The following movie was found:")
		fmt.Println()
		color.Unset()

		color.Set(color.FgGreen)
		fmt.Println(movie.Title)
		underline := setUnderline(len(movie.Title))
		fmt.Println(underline)
		color.Unset()

		color.Set(color.FgCyan)
		fmt.Println("  released:", movie.Released)
		fmt.Println("  rated:", movie.Rated)
		fmt.Println("  directed by:", movie.Director)
		fmt.Println("  starring:", movie.Actors)
		fmt.Println("  summary:", movie.Plot)
		return nil
	}

	app.Run(os.Args)
}
