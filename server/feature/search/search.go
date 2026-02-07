// This package is our "master" search providing one API
// for access to all of our search endpoints, massively
// simplifying access for any client (aka our web ui).

package search

import (
	"errors"
	"log/slog"

	"github.com/sbondCo/Watcharr/config"
	"github.com/sbondCo/Watcharr/domain"
	"github.com/sbondCo/Watcharr/media/tmdb"
	"github.com/sbondCo/Watcharr/util"
	"gorm.io/gorm"
)

type ContentProvider interface {
	SearchContent(query string, pageNum int) (tmdb.TMDBSearchMultiResponse, error)
	SearchMovies(query string, pageNum int) (tmdb.TMDBSearchMoviesResponse, error)
	SearchTv(query string, pageNum int) (tmdb.TMDBSearchShowsResponse, error)
	SearchPeople(query string, pageNum int) (tmdb.TMDBSearchPeopleResponse, error)
}

type Service struct {
	db              *gorm.DB
	cfg             *config.ServerConfig
	contentProvider ContentProvider
}

func NewService(
	db *gorm.DB,
	cfg *config.ServerConfig,
	contentProvider ContentProvider,
) *Service {
	return &Service{
		db,
		cfg,
		contentProvider,
	}
}

// `Limit` is not supported.
func (s *Service) Search(
	r domain.SearchRequest,
	pp util.PaginationParams,
) (domain.SearchResponse, error) {
	resp := domain.SearchResponse{}
	if r.Query == "" {
		return resp, errors.New("a query is required")
	}
	switch r.Type {
	case domain.SearchTypeMulti:
		if err := s.searchMulti(r.Query, pp.Page, &resp); err != nil {
			return resp, errors.New("multi search failed")
		}
	case domain.SearchTypeMovie:
		if err := s.searchMovie(r.Query, pp.Page, &resp); err != nil {
			return resp, errors.New("movie search failed")
		}
	case domain.SearchTypeShow:
		if err := s.searchTv(r.Query, pp.Page, &resp); err != nil {
			return resp, errors.New("tv search failed")
		}
	case domain.SearchTypePerson:
		if err := s.searchPeople(r.Query, pp.Page, &resp); err != nil {
			return resp, errors.New("person search failed")
		}
	case domain.SearchTypeGame:
		if err := s.searchGame(r.Query, pp.Page, &resp); err != nil {
			return resp, errors.New("game search failed")
		}
	}
	return resp, nil
}

// TODO if only one of the requests for data fails, we can still return the data?
// TODO but we'd need a way to tell the client that some data failed to get fetched,
// TODO either with a header OR a result added to array of type error
// SearchMulti is TMDB Multi search but with game data added to first page.
func (s *Service) searchMulti(
	query string,
	page int,
	resp *domain.SearchResponse,
) error {
	// TMDB
	tmdbRes, err := s.contentProvider.SearchContent(query, page)
	if err != nil {
		slog.Error("SearchMulti: Failed to search tmdb!", "error", err)
		return errors.New("content request failed")
	}
	for _, v := range tmdbRes.Results {
		resp.Results = append(
			resp.Results,
			v.AsMedia(),
		)
	}
	// IGDB
	igdbRes, err := s.cfg.TWITCH.Search(query)
	if err != nil {
		slog.Error("SearchMulti: Failed to search tmdb!", "error", err)
		return errors.New("content request failed")
	}
	for _, v := range igdbRes {
		resp.Results = append(
			resp.Results,
			v.AsMedia(),
		)
	}
	resp.Page = tmdbRes.Page
	resp.TotalPages = tmdbRes.TotalPages
	resp.TotalResults = int64(tmdbRes.TotalResults)
	return nil
}

func (s *Service) searchMovie(
	query string,
	page int,
	resp *domain.SearchResponse,
) error {
	tmdbRes, err := s.contentProvider.SearchMovies(query, page)
	if err != nil {
		slog.Error("SearchMovie: Failed to search tmdb!", "error", err)
		return errors.New("content request failed")
	}
	for _, v := range tmdbRes.Results {
		resp.Results = append(
			resp.Results,
			v.AsMedia(),
		)
	}
	resp.Page = tmdbRes.Page
	resp.TotalPages = tmdbRes.TotalPages
	resp.TotalResults = int64(tmdbRes.TotalResults)
	return nil
}

func (s *Service) searchTv(
	query string,
	page int,
	resp *domain.SearchResponse,
) error {
	tmdbRes, err := s.contentProvider.SearchTv(query, page)
	if err != nil {
		slog.Error("SearchMovie: Failed to search tmdb!", "error", err)
		return errors.New("content request failed")
	}
	for _, v := range tmdbRes.Results {
		resp.Results = append(
			resp.Results,
			v.AsMedia(),
		)
	}
	resp.Page = tmdbRes.Page
	resp.TotalPages = tmdbRes.TotalPages
	resp.TotalResults = int64(tmdbRes.TotalResults)
	return nil
}

func (s *Service) searchPeople(
	query string,
	page int,
	resp *domain.SearchResponse,
) error {
	tmdbRes, err := s.contentProvider.SearchPeople(query, page)
	if err != nil {
		slog.Error("SearchMovie: Failed to search tmdb!", "error", err)
		return errors.New("content request failed")
	}
	for _, v := range tmdbRes.Results {
		resp.Results = append(
			resp.Results,
			v.AsMedia(),
		)
	}
	resp.Page = tmdbRes.Page
	resp.TotalPages = tmdbRes.TotalPages
	resp.TotalResults = int64(tmdbRes.TotalResults)
	return nil
}

func (s *Service) searchGame(
	query string,
	page int,
	resp *domain.SearchResponse,
) error {
	igdbRes, err := s.cfg.TWITCH.Search(query)
	if err != nil {
		slog.Error("SearchMulti: Failed to search tmdb!", "error", err)
		return errors.New("content request failed")
	}
	for _, v := range igdbRes {
		resp.Results = append(
			resp.Results,
			v.AsMedia(),
		)
	}
	resp.Page = 1
	resp.TotalPages = 1
	resp.TotalResults = int64(len(igdbRes))
	return nil
}
