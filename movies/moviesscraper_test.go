package movies

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"github.com/ledongthuc/goterators"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	models "jwt-authentication-golang/models"
	"strings"
	"testing"
	"time"
)

func _TestLocateCities(t *testing.T) {
	cities, err := LocateCities("file://./resources/overview_haarlem.html")
	fileBytes, _ := ioutil.ReadFile("./resources/cities.json")
	var expectedCities []models.City
	//Content, err := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(cities)
	//if err == nil {
	//	fmt.Println(Content)
	//}
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(fileBytes, &expectedCities)
	assert.ElementsMatch(t, cities, expectedCities)
	if err != nil {
		t.Fatal(" got an error parsing home.html")
	}
	assert.Equal(t, 109, len(cities))
}

func TestLocateMovies(t *testing.T) {
	city, err := LocatePlays("file://./resources/overview_haarlem.html")
	if err != nil {
		t.Fatal(" got an error parsing haarlem.html")
	}
	assert.Equal(t, 4, len(city.Cinemas))
	assert.Equal(t, "", city.Href)
	schuur, _, _ := goterators.Find(city.Cinemas, func(item *models.Cinema) bool {
		return strings.Contains(item.Name, "Schuur")
	})
	moviedIds := make(map[string]models.Movie, 0)
	for i := 0; i < len(schuur.Plays); i++ {
		movieId := schuur.Plays[i].Movie.Id
		_, exists := moviedIds[movieId]
		if !exists {
			moviedIds[movieId] = schuur.Plays[i].Movie
		}
	}

	if len(moviedIds) != 12 {
		t.Fatal("expected 12 movies but was:", len(moviedIds))
	}
	assert.Equal(t, len(schuur.Plays), 34)
	for i := 0; i < len(schuur.Plays); i++ {
		assert.NotNil(t, schuur.Plays[i].Tickethref)
		assert.NotNil(t, schuur.Plays[i].Start)
		assert.NotNil(t, schuur.Plays[i].Movie)
		assert.NotNil(t, schuur.Plays[i].Movie.Id)
		assert.NotNil(t, schuur.Plays[i].Movie.Title)
		assert.NotNil(t, schuur.Plays[i].Movie.Href)
	}
	pastLives, _, _ := goterators.Find(schuur.Plays, func(item models.Play) bool {
		return strings.Contains(item.Movie.Title, "Past")
	})
	assert.NotNil(t, pastLives)
	assert.Equal(t, "past-lives-2023", pastLives.Movie.Id)
	assert.Equal(t, "https://www.filmladder.nl/film/past-lives-2023/synopsis/haarlem", pastLives.Movie.Href)
	// and read the Movie
	movie, err := LoadMovieContent("file://./resources/past-lives.html")
	checkMovieContent(t, err, movie)
	movie, err = LoadMovie("past-lives-2023")
	checkMovieContent(t, err, movie)

}

func checkMovieContent(t *testing.T, err error, movie models.Movie) {
	if err == nil {
		assert.Equal(t, "Nora en Hae Sung, twee diep verbonden jeugdvrienden, worden uit elkaar gerukt als Nora’s familie vanuit Zuid-Korea naar Canada emigreert. Twaalf jaar later studeert Nora in New York en hervinden de twee elkaar via het internet. Ze fantaseren over een wederzien, maar de afstand doet het contact verwateren. Nog eens twaalf jaar later is Nora inmiddels getrouwd en wordt ze tijdens een allesbepalende week herenigd met haar jeugdliefde, als hij haar opzoekt in New York. ", movie.Content)
		assert.Equal(t, 106, movie.Duration)
	}
}
func TestJsonIterUnMarshal(t *testing.T) {
	extra.SupportPrivateFields()
	type TestObject struct {
		field1 string
	}
	obj := TestObject{}
	jsoniter.UnmarshalFromString(`{"field1":"Hello"}`, &obj)
	assert.Equal(t, "Hello", obj.field1)
}
func TestJsonIterMarshal(t *testing.T) {
	extra.SupportPrivateFields()
	type TestObject struct {
		field1 string
	}
	obj := TestObject{field1: "Hello"}
	output, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalIndent(&obj, "", "")
	assert.Equal(t, "{\"field1\":\"Hello\"}", string(output))
}
func TestParseMovie(t *testing.T) {
	movie, err := LoadMovieContent("file://./resources/oppenheimer.html")
	if err != nil {
		t.Fatal(" got an error", err)
	}
	assert.Equal(t, `De theoretisch natuurkundige J. Robert Oppenheimer hielp bij de ontwikkeling van de eerste kernwapens en om de wereld te redden nam hij het risico dat deze werd verwoest.`, movie.Content)
	assert.Equal(t, 180, movie.Duration)
}

func _TestRockEnSeinne(t *testing.T) {
	//podiums, err := ParseRockEnSeinne("https://www.rockenseine.com/en/line-up/grille/?day_prog=2022-08-26")
	podiums, err := ParseRockEnSeinne("file://./resources/line_up.html")
	if err != nil {
		t.Fatal(" got an error", err)
	}
	assert.Equal(t, 5, len(podiums))
	podiumNick, _, err := goterators.Find(podiums, func(podium Podium) bool {
		return strings.Contains(podium.name, "grande-scene")
	})
	if err != nil {
		t.Fatal("podiumNick not found")
	}
	extra.SupportPrivateFields()
	content, err := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(podiums)
	if err != nil {
		t.Fatal(err)
	}
	content = "{ \"podiums\":" + content + "}"
	fmt.Println(content)
	assert.NotNilf(t, podiumNick, "", nil)
	assert.Equal(t, 4, len(podiumNick.events))
}

func TestTimeParsing(t *testing.T) {

	tijd, err := time.Parse(time.RFC3339, "2022-05-10T13:10:00+02:00")
	if err == nil {
		//fmt.Printf(" tijd: %v", tijd)
		year, month, day := tijd.Date()
		assert.Equal(t, "May", month.String())
		assert.Equal(t, 2022, year)
		assert.Equal(t, 10, day)
	} else {
		t.Fatal(err)
	}
}
func _TestFlow(t *testing.T) {
	//Start with all cities
	cities, err := LocateCities("file://./resources/overview_haarlem.html")
	haarlem, _, _ := goterators.Find(cities, func(item models.City) bool {
		return strings.EqualFold(item.Name, "haarlem")
	})
	assert.Equal(t, "https://www.filmladder.nl/haarlem", haarlem.Href)

	playsForCity, err := LocatePlays("file://./resources/overview_haarlem.html")
	if err == nil {
		assert.Equal(t, 4, len(playsForCity.Cinemas))
		schuur, _, _ := goterators.Find(playsForCity.Cinemas, func(item *models.Cinema) bool {
			return strings.EqualFold(item.Name, "Schuur")
		})
		assert.Equal(t, 34, len(schuur.Plays))
		numb := goterators.Filter(schuur.Plays, func(item models.Play) bool {
			return strings.EqualFold(item.Movie.Title, "numb")
		})
		assert.Equal(t, 2, len(numb))
		assert.Equal(t, "numb-2023", numb[0].Movie.Id)
		// nope duration opgehaald bij moviecontent
		//assert.Equal(t, 0, numb[0].Movie.Duration)

		// loadMovieContent
		movie, err2 := LoadMovieContent("file://./resources/oppenheimer.html")
		if err2 == nil {
			assert.NotNil(t, movie)
			assert.NotNil(t, movie.Content)
			assert.NotNil(t, movie.Duration)
			assert.Equal(t, 180, movie.Duration)
			assert.Equal(t, "De theoretisch natuurkundige J. Robert Oppenheimer hielp bij de ontwikkeling van de eerste kernwapens en om de wereld te redden nam hij het risico dat deze werd verwoest.", movie.Content)
		}
		//en live

	}
}

func TestPerformance(t *testing.T) {
	start := time.Now()
	for i := 0; i < 100; i++ {
		TestLiveLocateCitiesAndMovieContent(t)
	}
	end := time.Now()
	fmt.Printf("this took %d", (end.UnixMilli() - start.UnixMilli()))
}

func TestLiveLocateCitiesAndMovieContent(t *testing.T) {
	// live should be the sanme
	//cities, err := LocateCities("file://./resources/overview_haarlem.html")
	cities, err := LocateCities("http://www.filmladder.nl")
	if err != nil {
		t.Fatal(" got an error parsing home.html")
	}
	if len(cities) < 100 {
		t.Fatal(" expected more then 100 cities")
	}
	_, index, err := goterators.Find(cities, func(item models.City) bool {
		return strings.Contains(item.Name, "Haarlem")
	})
	if index == 0 {
		t.Fatal("Haarlem not found")
	}
	cityName := "haarlem"
	city, err := LocatePlaysForCity(cityName)
	if err != nil {
		t.Fatal(" got an error loading Plays for city")
	}
	schuur, _, _ := goterators.Find(city.Cinemas, func(item *models.Cinema) bool {
		return strings.Contains(item.Name, "Schuur")
	})
	assert.Greater(t, len(schuur.Plays), 10)
	movie, err2 := LoadMovieContent("https://www.filmladder.nl/film/oppenheimer-2023")
	if err2 == nil {
		assert.NotNil(t, movie)
		assert.NotNil(t, movie.Content)
		assert.NotNil(t, movie.Duration)
		assert.Equal(t, 180, movie.Duration)
		assert.Equal(t, "De theoretisch natuurkundige J. Robert Oppenheimer hielp bij de ontwikkeling van de eerste kernwapens en om de wereld te redden nam hij het risico dat deze werd verwoest.", movie.Content)
	}

}
