package movies

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	goquery "github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

// Init a new HTTP client for use when the client doesn't want to use their own.
var (
	defaultClient = &http.Client{}

	debug = false

	// Headers contains all HTTP headers to send
	Headers = make(map[string]string)

	// Cookies contains all HTTP cookies to send
	Cookies = make(map[string]string)
)

// ErrorType defines types of errors that are possible from soup
type ErrorType int

const (
	ErrCreatingGetRequest ErrorType = iota
	ErrInGetRequest
	WrongURL
	ErrReadingResponse
	ValNotFound
	ParsingError
)

// Error allows easier introspection on the type of error returned.
// If you know you have a Error, you can compare the Type to one of the exported types
// from this package to see what kind of error it is, then further inspect the Error() method
// to see if it has more specific details for you, like in the case of a ErrElementNotFound
// type of error.
type Error struct {
	Type ErrorType
	msg  string
}

func (se Error) Error() string {
	return se.msg
}

func newError(t ErrorType, msg string) Error {
	return Error{Type: t, msg: msg}
}

type Movie struct {
	id       string
	title    string
	href     string
	rating   string
	content  string
	image    string
	duration int
	//plays    []Play
}
type Event struct {
	artist string
	start  time.Time
	end    time.Time
}
type Podium struct {
	name   string
	events []*Event
}

const replacement = " "

var replacer = strings.NewReplacer(
	"\r\n", replacement,
	"\r", replacement,
	"\n", replacement,
	"\v", replacement,
	"\f", replacement,
	"  ", replacement,
	"\u0085", replacement,
	"\u2028", replacement,
	"\u2029", replacement,
)

type Play struct {
	movie      Movie
	start      time.Time
	tickethref string
}
type Cinema struct {
	name  string
	href  string
	plays []Play
}
type City struct {
	href    string
	name    string
	cinemas []*Cinema
}
type Talk struct {
	name  string
	start time.Time
	end   time.Time
}

func LocatePlaysForCity(cityName string) (City, error) {
	url := fmt.Sprintf("http://filmladder.nl/%s", cityName)
	return LocatePlays(url)
}

func LocatePlays(url string) (City, error) {
	doc, err := load(url)
	city := City{}
	if err != nil {
		return city, err
	}
	city.name = strings.Trim(strings.TrimPrefix(doc.Find("title").Text(), "Bioscopen in"), " ")
	cinemasWithplays := locateCineasWithPlays(doc)
	city.cinemas = cinemasWithplays
	return city, nil
}
func locateCineasWithPlays(doc *goquery.Document) []*Cinema {
	cinemasWithplays := make([]*Cinema, 0)

	doc.Find("div.cinema").Each(func(i int, s *goquery.Selection) {
		// first div contains the href of the cinema
		cinemasWithplays = append(cinemasWithplays, locatePlays(s))

	})
	return cinemasWithplays

}

func locatePlays(doc *goquery.Selection) *Cinema {
	//first div node should contain the matching href id for the cinems
	cinema := locateCinema(doc)
	//movies := make([]movie, 0)
	doc.Find("div.hall").Each(func(i int, s *goquery.Selection) {
		// first div contains the href of the cinema
		movie, err := locateMovie(s)
		if err == nil {
			plays := locateTimes(s, movie)
			cinema.plays = append(cinema.plays, plays...)
			//fmt.Printf("%#v\n", movie)
		}
	})
	return &cinema
}

func locateCinema(doc *goquery.Selection) Cinema {
	cinemaSelector := doc.Find("a.cinema-link")
	cinema := Cinema{}
	title, err := locateNodeValue(cinemaSelector, "title")
	if err == nil {
		cinema.name = title
	}
	href, err2 := locateNodeValue(cinemaSelector, "href")
	if err2 == nil {
		cinema.href = href
	}
	return cinema

}

func locateTimes(doc *goquery.Selection, movie Movie) []Play {
	result := make([]Play, 0)
	weekContainer := doc.Find("div.week-sheet")
	weekContainer.Find("div.day").Each(func(i int, s *goquery.Selection) {
		s.Find("div[itemprop=\"startDate\"]").Each(func(i int, timecontainer *goquery.Selection) {
			playTime, _ := locateValue(timecontainer.Nodes[0].Attr, "content")
			tijd, _ := time.Parse(time.RFC3339, playTime)
			ticketLinkNode := timecontainer.Find("a.ticket")
			play := Play{start: tijd, movie: movie}
			value, err := locateNodeValue(ticketLinkNode, "href")
			if err == nil {
				play.tickethref = value
			}
			result = append(result, play)
		})
	})
	return result
}
func ParseRockEnSeinne(url string) ([]Podium, error) {
	podiums := make([]Podium, 0)
	doc, err := load(url)

	if err == nil {
		doc.Find("div.col").Each(func(i int, podiumContainer *goquery.Selection) {
			PodiumDiv := podiumContainer.Find("div.col-day")
			if PodiumDiv.Nodes != nil {
				podiumName, _ := locateValue(PodiumDiv.Nodes[0].Attr, "class")
				colWithPodiumName := strings.Split(podiumName, " ")[1][4:]
				podium := Podium{name: colWithPodiumName}
				podium.events = locateEvents(PodiumDiv)
				podiums = append(podiums, podium)
			}
		})
	}
	return podiums, nil
}

func locateEvents(podiumCol *goquery.Selection) []*Event {
	events := make([]*Event, 0)
	podiumCol.Find("li>a").Each(func(i int, artist *goquery.Selection) {
		artistName, _ := locateNodeValue(artist, "title")
		event := Event{artist: artistName}
		when := artist.Parent().Find("span").Text()
		times := strings.Fields(when)
		start := "2022-08-26T" + times[0] + ":00+02:00"
		end := "2022-08-26T" + times[2] + ":00+02:00"

		startTijd, err := time.Parse(time.RFC3339, start)
		if err == nil {
			endTijd, err1 := time.Parse(time.RFC3339, end)
			if err1 == nil {
				event.start = startTijd
				event.end = endTijd
				events = append(events, &event)
			}
		}
	})
	return events
}
func ParseMovieContent(url string) (Movie, error) {
	doc, err := load(url)
	if err != nil {
		return Movie{}, err
	}
	movie := Movie{}
	movie.content = replaceReplacer(doc.Find("p.synopsis").Text())
	duration := doc.Find("p[itemprop=duration]").Text()
	if strings.HasSuffix(duration, "minuten") {
		movie.duration, _ = strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(duration, "minuten", "")))
	}
	return movie, nil
}

func LoadMovie(id string) (Movie, error) {
	url := fmt.Sprintf("http://filmladder.nl/film/%s", id)

	return ParseMovieContent(url)
}

func replaceReplacer(s string) string {
	content := replacer.Replace(s)
	for strings.Contains(content, "  ") {
		content = strings.ReplaceAll(content, "  ", " ")
	}
	return content
}

func locateMovie(doc *goquery.Selection) (Movie, error) {
	movieLink := doc.Find("a")
	if movieLink.Nodes == nil {
		return Movie{}, newError(ParsingError, "no more movies")
	}
	ratingContainer := doc.Find("span.star-rating a")
	rating := "unknown"
	href := ""
	if len(ratingContainer.Nodes) > 0 {
		href, _ = locateValue(ratingContainer.Nodes[0].Attr, "href")
		rating = ratingContainer.Text()
	}
	title, _ := locateValue(movieLink.Nodes[0].Attr, "title")
	imageContainer := movieLink.Find("img")
	imageHref, _ := locateValue(imageContainer.Nodes[0].Attr, "data-src")
	class, _ := locateNodeValue(doc, "class")
	id := strings.Split(class, " ")[1]

	return Movie{id: id, title: title, image: imageHref, href: href, rating: rating}, nil

}

func LocateCities(url string) ([]City, error) {
	if url == "" {
		return LocateCities("http://www.filmladder.nl")
	}
	doc, err := load(url)
	cities := make([]City, 0)
	if err != nil {
		return nil, err
	}
	doc.Find("div.cities-sheet>a").Each(func(i int, s *goquery.Selection) {
		href, _ := locateNodeValue(s, "href")
		var city = City{href: href,
			name: s.Text()}
		cities = append(cities, city)
	})
	return cities, nil
}

func locateNodeValue(selection *goquery.Selection, key string) (string, error) {
	if selection.Nodes != nil && selection.Nodes[0] != nil && selection.Nodes[0].Attr != nil {
		return locateValue(selection.Nodes[0].Attr, key)
	} else {
		return "", newError(ValNotFound, "queryselector does not have nodes or attr")
	}
}

func locateValue(attr []html.Attribute, key string) (string, error) {
	for index := 0; index < len(attr); index++ {
		if attr[index].Key == key {
			return attr[index].Val, nil
		}
	}
	return "", newError(ValNotFound, "attr has no key"+key)

}

func load(url string) (*goquery.Document, error) {
	if strings.HasPrefix(url, "http") {
		content, err := GetWithClient(url, defaultClient)
		if err == nil {
			return goquery.NewDocumentFromReader(strings.NewReader(content))
		} else {
			return nil, newError(ErrCreatingGetRequest, "getWithClient on "+url+" failed with "+err.Error())
		}
	} else {
		if strings.HasPrefix(url, "file://") {
			fileUrl := url[7:]
			return LoadFromFile(fileUrl)
		} else {
			return nil, newError(WrongURL, "url should start with file or http")
		}
	}
	//return nil, nil
}

func LoadFromFile(url string) (*goquery.Document, error) {
	// create from a file
	f, err := os.Open(url)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	return doc, nil
}

// setHeadersAndCookies helps build a request
func setHeadersAndCookies(req *http.Request) {
	// Set headers
	for hName, hValue := range Headers {
		req.Header.Set(hName, hValue)
	}
	// Set cookies
	for cName, cValue := range Cookies {
		req.AddCookie(&http.Cookie{
			Name:  cName,
			Value: cValue,
		})
	}
}

// GetWithClient returns the HTML returned by the url using a provided HTTP client
func GetWithClient(url string, client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		if debug {
			panic("Couldn't create GET request to " + url)
		}
		return "", newError(ErrCreatingGetRequest, "error creating get request to "+url)
	}

	setHeadersAndCookies(req)

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		if debug {
			panic("Couldn't perform GET request to " + url)
		}
		return "", newError(ErrInGetRequest, "couldn't perform GET request to "+url)
	}
	defer resp.Body.Close()
	utf8Body, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(utf8Body)
	if err != nil {
		if debug {
			panic("Unable to read the response body")
		}
		return "", newError(ErrReadingResponse, "unable to read the response body")
	}
	return string(bytes), nil
}
