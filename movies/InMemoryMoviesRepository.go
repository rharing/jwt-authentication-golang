package movies

import (
	"jwt-authentication-golang/models"
)

type InMemoryMoviesRepository struct {
	// here the key is movieid
	store map[string]models.Movie
	// here the key is userid
	myMovies map[string]models.MyMovies
}

func NewInMemoryMoviesRepository() models.MoviesRepository {
	return InMemoryMoviesRepository{store: make(map[string]models.Movie, 0), myMovies: make(map[string]models.MyMovies, 0)}

}
func (m InMemoryMoviesRepository) LoadMovieContent(movieId string) (models.Movie, error) {
	movie, found := m.store[movieId]
	if !found {
		loadedMovie, err := LoadMovieContent(movieId)
		if err == nil {
			m.store[movieId] = loadedMovie
			return loadedMovie, nil
		} else {
			return models.Movie{}, err
		}
	} else {
		return movie, nil
	}

}

func (m InMemoryMoviesRepository) SeenMovie(movieid string, userId string) {
	myMovies, found := m.myMovies[userId]
	if !found {
		myMovies = models.MyMovies{}
	}
	myMovies.Seen = append(myMovies.Seen, movieid)
	m.myMovies[userId] = myMovies
}
func (m InMemoryMoviesRepository) WantedMovie(movieid string, userId string) {
	myMovies, found := m.myMovies[userId]
	if !found {
		myMovies = models.MyMovies{}
	}
	myMovies.Wanted = append(myMovies.Wanted, movieid)
	m.myMovies[userId] = myMovies
}
func (m InMemoryMoviesRepository) UnwantedMovie(movieid string, userId string) {
	myMovies, found := m.myMovies[userId]
	if !found {
		myMovies = models.MyMovies{}
	}
	myMovies.Unwanted = append(myMovies.Unwanted, movieid)
	m.myMovies[userId] = myMovies
}
func (m InMemoryMoviesRepository) ResetMovie(movieid string, userId string) {
	myMovies, found := m.myMovies[userId]
	if found {
		myMovies = models.MyMovies{}
	}

	m.myMovies[userId] = myMovies
}
func (m InMemoryMoviesRepository) MyMovies(userId string) models.MyMovies {
	myMovies, found := m.myMovies[userId]
	if found {
		return myMovies
	}
	return models.MyMovies{}
}
