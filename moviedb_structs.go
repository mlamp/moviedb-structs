package moviedb_structs

import (
	"encoding/json"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type RottenMoviesResponse struct {
	Total  int           `json:"total"`
	Movies []RottenMovie `json:"movies"`
}

type RottenTvResponse struct {
	PageCount  int        `json:"pageCount"`
	TotalCount int        `json:"totalCount"`
	TvSeries   []RottenTv `json:"tvSeries"`
}

type RottenTv struct {
	Title       string      `json:"title"`
	EndYear     json.Number `json:"endYear,omitempty"`
	StartYear   json.Number `json:"startYear"`
	PosterImage string      `json:"posterImage,omitempty"`
	MeterClass  string      `json:"meterClass"`
	Image       string      `json:"image"`
	URL         string      `json:"url"`
	MeterValue  json.Number `json:"meterValue,omitempty"`
}

type RottenMovie struct {
	ID               string      `json:"id"`
	Title            string      `json:"title"`
	Year             json.Number `json:"year"`
	MpaaRating       string      `json:"mpaa_rating"`
	Runtime          json.Number `json:"runtime"`
	CriticsConsensus string      `json:"critics_consensus,omitempty"`
	ReleaseDates     struct {
		Theater string `json:"theater"`
		Dvd     string `json:"dvd"`
	} `json:"release_dates"`
	Ratings struct {
		CriticsRating  string      `json:"critics_rating"`
		CriticsScore   json.Number `json:"critics_score"`
		AudienceRating string      `json:"audience_rating"`
		AudienceScore  json.Number `json:"audience_score"`
	} `json:"ratings"`
	Synopsis string `json:"synopsis"`
	Posters  struct {
		Thumbnail string `json:"thumbnail"`
		Profile   string `json:"profile"`
		Detailed  string `json:"detailed"`
		Original  string `json:"original"`
	} `json:"posters"`
	AbridgedCast []struct {
		Name       string   `json:"name"`
		ID         string   `json:"id"`
		Characters []string `json:"characters"`
	} `json:"abridged_cast"`
	Links struct {
		Self      string `json:"self"`
		Alternate string `json:"alternate"`
		Cast      string `json:"cast"`
		Reviews   string `json:"reviews"`
		Similar   string `json:"similar"`
	} `json:"links"`
}

type OMDbSingleResult struct {
	Title        string    `json:"Title"`
	Year         string    `json:"Year"`
	Rated        string    `json:"Rated"`
	Released     string    `json:"Released"`
	Runtime      string    `json:"Runtime"`
	Genre        string    `json:"Genre"`
	Director     string    `json:"Director"`
	Writer       string    `json:"Writer"`
	Actors       string    `json:"Actors"`
	Plot         string    `json:"Plot"`
	Language     string    `json:"Language"`
	Country      string    `json:"Country"`
	Awards       string    `json:"Awards"`
	Poster       string    `json:"Poster"`
	Metascore    string    `json:"Metascore"`
	ImdbRating   string    `json:"imdbRating"`
	ImdbVotes    string    `json:"imdbVotes"`
	ImdbID       string    `json:"imdbID"`
	Type         string    `json:"Type"`
	TotalSeasons string    `json:"totalSeasons"`
	Response     string    `json:"Response"`
	LastUpdated  time.Time `json:"-"`
}

type MovieRottenScore struct {
	CriticsScore  json.Number `json:"criticsScore"`
	AudienceScore json.Number `json:"audienceScore"`
}

type Movie struct {
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	Year         int              `json:"year"`
	RottenLink   string           `json:"rottenLink"`
	RottenScores MovieRottenScore `json:"rottenScores"`
	ImdbId       string           `json:"imdbId"`
	Actors       []Actor          `json:"actors"`
	MatchScore   int              `json:"matchScore" datastore:"-"`
	LastUpdated  time.Time        `json:"-"`
}

type Actor struct {
	Name string `json:"name"`
}

type Tv struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	YearFrom     int           `json:"yearFrom"`
	YearTo       int           `json:"yearTo"`
	RottenScores TvRottenScore `json:"rottenScores"`
	RottenLink   string        `json:"rottenLink"`
	ImdbId       string        `json:"imdbId"`
	Actors       []Actor       `json:"actors"`
	MatchScore   int           `json:"matchScore"`
	LastUpdated  time.Time     `json:"-"`
}

type TvRottenScore struct {
	CriticsScore  json.Number `json:"criticsScore"`
	AudienceScore json.Number `json:"audienceScore"`
}

func MoviesAreEqual(movie1 Movie, movie2 Movie) (isEqual bool) {
	if movie1.ID != movie2.ID || movie1.RottenScores.AudienceScore != movie2.RottenScores.AudienceScore ||
		movie1.RottenScores.CriticsScore != movie2.RottenScores.CriticsScore || movie1.RottenLink != movie2.RottenLink ||
		movie1.ImdbId != movie2.ImdbId || movie1.Title != movie2.Title || movie1.Year != movie2.Year {
		return false
	} else {
		return true
	}
}

func MovieFromOmdbSingleResult(omdbMovie OMDbSingleResult) (movie Movie) {
	movie.Title = omdbMovie.Title
	movie.Year, _ = strconv.Atoi(string(omdbMovie.Year))
	movie.ImdbId = omdbMovie.ImdbID
	if omdbMovie.Actors != "" {
		actors := strings.Split(omdbMovie.Actors, ",")
		for _, actorName := range actors {
			actorNames := strings.Fields(actorName)
			movie.Actors = append(movie.Actors, Actor{Name: strings.Join(actorNames, " ")})
		}
	}
	return movie
}

func MovieFromRottenMovie(rottenMovie RottenMovie) (movie Movie) {
	movie.Title = rottenMovie.Title
	movie.Year, _ = strconv.Atoi(string(rottenMovie.Year))
	movie.RottenScores.CriticsScore = rottenMovie.Ratings.CriticsScore
	movie.RottenScores.AudienceScore = rottenMovie.Ratings.AudienceScore
	movie.RottenLink = rottenMovie.Links.Alternate
	if len(rottenMovie.AbridgedCast) > 0 {
		for _, actorName := range rottenMovie.AbridgedCast {
			movie.Actors = append(movie.Actors, Actor{Name: actorName.Name})
		}
	}
	return movie
}

func TvFromRottenTv(rottenTv RottenTv) (tv Tv) {
	tv.Title = rottenTv.Title
	tv.YearFrom, _ = strconv.Atoi(string(rottenTv.StartYear))
	tv.YearTo, _ = strconv.Atoi(string(rottenTv.EndYear))
	tv.RottenScores.CriticsScore = json.Number(rottenTv.MeterValue)
	// Cleanup url
	re := regexp.MustCompile("^(/tv/[0-9a-z_-]{1,25})/?")
	matches := re.FindStringSubmatch(rottenTv.URL)
	if len(matches) > 1 {
		tv.RottenLink = "https://www.rottentomatoes.com" + matches[1]
	}
	return tv
}

func TvFromOmdbSingleResult(omdbResult OMDbSingleResult) (tv Tv) {
	tv.Title = omdbResult.Title

	re := regexp.MustCompile(`(\d*)–(\d*)`)
	yearMatches := re.FindStringSubmatch(omdbResult.Year)
	if len(yearMatches) > 1 {
		tv.YearFrom, _ = strconv.Atoi(yearMatches[1])
		tv.YearTo, _ = strconv.Atoi(yearMatches[2])
	}
	tv.ImdbId = omdbResult.ImdbID
	return tv
}

// MatchScoreSorter sorts movies by matchScore.
type MovieMatchScoreSorter []Movie

func (a MovieMatchScoreSorter) Len() int           { return len(a) }
func (a MovieMatchScoreSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MovieMatchScoreSorter) Less(i, j int) bool { return a[i].MatchScore > a[j].MatchScore }

func SortMoviesByMatchScore(movies []Movie) []Movie {
	sort.Sort(MovieMatchScoreSorter(movies))
	return movies
}

// MatchScoreSorter sorts movies by matchScore.
type TvMatchScoreSorter []Tv

func (a TvMatchScoreSorter) Len() int           { return len(a) }
func (a TvMatchScoreSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TvMatchScoreSorter) Less(i, j int) bool { return a[i].MatchScore > a[j].MatchScore }

func SortTvsByMatchScore(tvs []Tv) []Tv {
	sort.Sort(TvMatchScoreSorter(tvs))
	return tvs
}
