package modules

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/MehraB832/olivia_core/user"
	"github.com/MehraB832/olivia_core/util"

	"github.com/MehraB832/olivia_core/language"
)

var (
	// GenresTag is the intent tag for its module
	GenresTag = "movies genres"
	// MoviesTag is the intent tag for its module
	MoviesTag = "movies search"
	// MoviesAlreadyTag is the intent tag for its module
	MoviesAlreadyTag = "already seen movie"
	// MoviesDataTag is the intent tag for its module
	MoviesDataTag = "movies search from data"
)

// GenresReplacer gets the genre specified in the message and adds it to the user information.
// See modules/modules.go#Module.Replacer() for more details.
func GenresReplacer(locale, entry, response, token string) (string, string) {
	genres := language.FindMoviesGenres(locale, entry)

	// If there is no genres then reply with a message from res/datasets/messages.json
	if len(genres) == 0 {
		responseTag := "no genres"
		return responseTag, util.SelectRandomMessage(locale, responseTag)
	}

	// Change the user information to add the new genres
	user.UpdateUserProfile(token, func(information user.UserProfile) user.UserProfile {
		for _, genre := range genres {
			// Append the genre only is it isn't already in the information
			if util.SliceIncludes(information.GenrePreferences, genre) {
				continue
			}

			information.GenrePreferences = append(information.GenrePreferences, genre)
		}
		return information
	})

	return GenresTag, response
}

// MovieSearchReplacer replaces the patterns contained inside the response by the movie's name
// and rating from the genre specified in the message.
// See modules/modules.go#Module.Replacer() for more details.
func MovieSearchReplacer(locale, entry, response, token string) (string, string) {
	genres := language.FindMoviesGenres(locale, entry)

	// If there is no genres then reply with a message from res/datasets/messages.json
	if len(genres) == 0 {
		responseTag := "no genres"
		return responseTag, util.SelectRandomMessage(locale, responseTag)
	}

	movie := language.SearchMovie(genres[0], token)

	return MoviesTag, fmt.Sprintf(response, movie.Name, movie.Rating)
}

// MovieSearchFromInformationReplacer replaces the patterns contained inside the response by the movie's name
// and rating from the genre in the user's information.
// See modules/modules.go#Module.Replacer() for more details.
func MovieSearchFromInformationReplacer(locale, _, response, token string) (string, string) {
	// If there is no genres then reply with a message from res/datasets/messages.json
	genres := user.RetrieveUserProfile(token).GenrePreferences
	if len(genres) == 0 {
		responseTag := "no genres saved"
		return responseTag, util.SelectRandomMessage(locale, responseTag)
	}

	movie := language.SearchMovie(genres[rand.Intn(len(genres))], token)
	genresJoined := strings.Join(genres, ", ")
	return MoviesDataTag, fmt.Sprintf(response, genresJoined, movie.Name, movie.Rating)
}
