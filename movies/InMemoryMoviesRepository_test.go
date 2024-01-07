package movies

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryMoviesRepository_SeenMovie(t *testing.T) {
	moviesRepository := NewInMemoryMoviesRepository()
	moviesRepository.SeenMovie("movie1", "user1")
	assert.Equal(t, 1, len(moviesRepository.MyMovies("user1").Seen))
	moviesRepository.SeenMovie("movie2", "user1")
	assert.Equal(t, 2, len(moviesRepository.MyMovies("user1").Seen))
}
