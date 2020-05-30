package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	cloudcreation "wordCloud/cloudCreation"

	"github.com/bbalet/stopwords"
	lyrics "github.com/rhnvrm/lyric-api-go"
)

// Conf is a struct for storing font data
type Conf struct {
	FontMaxSize     int          `json:"font_max_size"`
	FontMinSize     int          `json:"font_min_size"`
	RandomPlacement bool         `json:"random_placement"`
	FontFile        string       `json:"font_file"`
	Colors          []color.RGBA `json:"colors"`
	Width           int          `json:"width"`
	Height          int          `json:"height"`
	Mask            MaskConf     `json:"mask"`
}

// MaskConf is a struct for the mask properties
type MaskConf struct {
	File  string     `json:"file"`
	Color color.RGBA `json:"color"`
}

var (
	spotifyBearerToken = "BQBLd_vvy0U_zNXx8pq-kcb-84y3oDUc73HMm2oIC6NHOzIi6oSWBAyLToU0tAg8wiYC7G3yP8P_pgFh_rsjgOwZAUW5cJSc6Fw13pGMYVvjfD-IA2cdm6w8whR7z0N-iWt73kDERhO6y1885cPB2I1a"
	geniusAccessToken  = "rKOfLFltqMYTK-3GtwW6v1V08epuy9Pu0uiJmqWdMiQE6K3W8KLEOkZkEZqoPUo_"

	// DefaultColors represents the color scheme used for DefaultConf
	DefaultColors = []color.RGBA{
		{0x1b, 0x1b, 0x1b, 0xff},
		{0x48, 0x48, 0x4B, 0xff},
		{0x59, 0x3a, 0xee, 0xff},
		{0x65, 0xCD, 0xFA, 0xff},
		{0x70, 0xD6, 0xBF, 0xff},
	}

	// DefaultConf creates a configuration struct for the word cloud
	DefaultConf = Conf{
		FontMaxSize:     700,
		FontMinSize:     10,
		RandomPlacement: false,
		FontFile:        "./fonts/roboto/Roboto-Regular.ttf",
		Colors:          DefaultColors,
		Width:           4096,
		Height:          4096,
		Mask: MaskConf{"", color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}},
	}

	output = flag.String("output", "output.png", "path to output image")
)

func main() {
	items := getTopTracks()
	lyricsString := getLyrics(items)
	lyricsMap := lyricsStringToMap(lyricsString)
	generateWordCloud(lyricsMap)
}

func getTopTracks() interface{} {
	request, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks", nil)
	if err != nil {
		log.Fatalln(err)
	}

	var bearer = "Bearer " + spotifyBearerToken

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", bearer)

	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var dataJSON map[string]interface{}
	err = json.Unmarshal(body, &dataJSON)
	if err != nil {
		log.Fatalln(err)
	}

	return dataJSON["items"]
}

func getLyrics(items interface{}) string {
	var lyricsString string

	for _, value := range items.([]interface{}) {
		album := value.(map[string]interface{})
		song := album["name"].(string)

		for _, val := range album["artists"].([]interface{}) {
			artistStruct := val.(map[string]interface{})
			artist := artistStruct["name"].(string)

			l := lyrics.New(lyrics.WithoutProviders(), lyrics.WithGeniusLyrics(geniusAccessToken))
			lyric, err := l.Search(artist, song)
			if err == nil {
				lyricsString = lyricsString + lyric
			}
			break
		}
	}

	return lyricsString
}

func lyricsStringToMap(lyricsString string) map[string]int {
	replacer := strings.NewReplacer(",", "", ".", "", ";", "", ")", "", "Intro", "", "Pre-Chorus", "", "[", "", "?", "", "]", "", "(", "", "Verse", "", "'", "")
	lyricsString = replacer.Replace(lyricsString)
	cleanLyrics := stopwords.CleanString(lyricsString, "en", true)
	words := strings.Fields(cleanLyrics)

	lyricsMap := make(map[string]int)
	for _, val := range words {
		if _, ok := lyricsMap[val]; ok {
			lyricsMap[val]++
		} else {
			lyricsMap[val] = 1
		}
	}

	return lyricsMap
}

func generateWordCloud(lyricsMap map[string]int) {
	conf := DefaultConf
	confJSON, _ := json.Marshal(conf)

	err := json.Unmarshal(confJSON, &conf)
	if err != nil {
		fmt.Println(err)
	}

	var boxes []*cloudcreation.Box
	if conf.Mask.File != "" {
		boxes = cloudcreation.Mask(
			conf.Mask.File,
			conf.Width,
			conf.Height,
			conf.Mask.Color)
	}

	colors := make([]color.Color, 0)
	for _, c := range conf.Colors {
		colors = append(colors, c)
	}

	w := cloudcreation.NewWordcloud(
		lyricsMap,
		cloudcreation.FontFile(conf.FontFile),
		cloudcreation.FontMaxSize(conf.FontMaxSize),
		cloudcreation.FontMinSize(conf.FontMinSize),
		cloudcreation.Colors(colors),
		cloudcreation.MaskBoxes(boxes),
		cloudcreation.Height(conf.Height),
		cloudcreation.Width(conf.Width),
		cloudcreation.RandomPlacement(conf.RandomPlacement),
	)

	img := w.Draw()
	outputFile, err := os.Create(*output)
	defer outputFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	png.Encode(outputFile, img)
}
