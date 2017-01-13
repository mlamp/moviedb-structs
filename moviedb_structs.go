package structs_structs

import (
	"strconv"
	"encoding/json"
	"regexp"
)

type RottenMoviesResponse struct {
	Total  int `json:"total"`
	Movies []RottenMovie `json:"movies"`
}

type RottenTvResponse struct {
	PageCount int `json:"pageCount"`
	TotalCount int `json:"totalCount"`
	TvSeries []RottenTv `json:"tvSeries"`
}

type RottenTv struct {
	Title string `json:"title"`
	EndYear json.Number `json:"endYear,omitempty"`
	StartYear json.Number `json:"startYear"`
	PosterImage string `json:"posterImage,omitempty"`
	MeterClass string `json:"meterClass"`
	Image string `json:"image"`
	URL string `json:"url"`
	MeterValue json.Number `json:"meterValue,omitempty"`
}

type RottenMovie struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Year             json.Number `json:"year"`
	MpaaRating       string `json:"mpaa_rating"`
	Runtime          json.Number `json:"runtime"`
	CriticsConsensus string `json:"critics_consensus,omitempty"`
	ReleaseDates struct {
		Theater string `json:"theater"`
		Dvd     string `json:"dvd"`
	} `json:"release_dates"`
	Ratings struct {
		CriticsRating  string `json:"critics_rating"`
		CriticsScore   json.Number `json:"critics_score"`
		AudienceRating string `json:"audience_rating"`
		AudienceScore  json.Number `json:"audience_score"`
	} `json:"ratings"`
	Synopsis string `json:"synopsis"`
	Posters struct {
		Thumbnail string `json:"thumbnail"`
		Profile   string `json:"profile"`
		Detailed  string `json:"detailed"`
		Original  string `json:"original"`
	} `json:"posters"`
	AbridgedCast []struct {
		Name       string `json:"name"`
		ID         string `json:"id"`
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
	Title        string `json:"Title"`
	Year         string `json:"Year"`
	Rated        string `json:"Rated"`
	Released     string `json:"Released"`
	Runtime      string `json:"Runtime"`
	Genre        string `json:"Genre"`
	Director     string `json:"Director"`
	Writer       string `json:"Writer"`
	Actors       string `json:"Actors"`
	Plot         string `json:"Plot"`
	Language     string `json:"Language"`
	Country      string `json:"Country"`
	Awards       string `json:"Awards"`
	Poster       string `json:"Poster"`
	Metascore    string `json:"Metascore"`
	ImdbRating   string `json:"imdbRating"`
	ImdbVotes    string `json:"imdbVotes"`
	ImdbID       string `json:"imdbID"`
	Type         string `json:"Type"`
	TotalSeasons string `json:"totalSeasons"`
	Response     string `json:"Response"`
}

type MovieRottenScore struct {
	CriticsScore  json.Number `json:"criticsScore"`
	AudienceScore json.Number `json:"audienceScore"`
}

type Movie struct {
	Title        string `json:"title"`
	Year         int `json:"year"`
	RottenLink   string `json:"rottenLink"`
	RottenScores MovieRottenScore `json:"rottenScores"`
	imdbId       string `json:"imdbId"`
}

type Tv struct {
	Title        string `json:"title"`
	YearFrom     int `json:"yearFrom"`
	YearTo       int `json:"yearTo"`
	RottenScores TvRottenScore `json:"rottenScores"`
	RottenLink string `json:"rottenLink"`
	imdbId string `json:"imdbId"`
}

type TvRottenScore struct {
	CriticsScore  int `json:"criticsScore"`
	AudienceScore int `json:"audienceScore"`
}

func MovieFromOmdbSingleResult(omdbMovie OMDbSingleResult) (movie Movie) {
	movie.Title = omdbMovie.Title
	movie.Year, _ = strconv.Atoi(string(omdbMovie.Year))
	movie.imdbId = omdbMovie.ImdbID
	return movie
}

func MovieFromRottenMovie(rottenMovie RottenMovie) (movie Movie) {
	movie.Title = rottenMovie.Title
	movie.Year, _ = strconv.Atoi(string(rottenMovie.Year))
	movie.RottenScores.CriticsScore = rottenMovie.Ratings.CriticsScore
	movie.RottenScores.AudienceScore = rottenMovie.Ratings.AudienceScore
	movie.RottenLink = rottenMovie.Links.Alternate
	return movie
}

func TvFromRottenTv(rottenTv RottenTv) (tv Tv) {
	tv.Title = rottenTv.Title
	tv.YearFrom, _ = strconv.Atoi(string(rottenTv.StartYear))
	tv.YearTo, _ = strconv.Atoi(string(rottenTv.EndYear))
	tv.RottenScores.CriticsScore, _ = strconv.Atoi(string(rottenTv.MeterValue))
	// Cleanup url
	re := regexp.MustCompile("^(/tv/[a-z_-]{1,25})/?")
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
		tv.YearFrom, _ = strconv.Atoi(yearMatches[1]);
		tv.YearTo, _ = strconv.Atoi(yearMatches[2]);
	}
	tv.imdbId = omdbResult.ImdbID
	return tv
}
