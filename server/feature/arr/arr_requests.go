package arr

import (
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/sbondCo/Watcharr/arr"
	"github.com/sbondCo/Watcharr/config"
	"github.com/sbondCo/Watcharr/database/entity"
	"github.com/sbondCo/Watcharr/feature/auth/permission"
	"gorm.io/gorm"
)

func deleteArrRequest(db *gorm.DB, id uint) error {
	resp := db.Delete(&entity.ArrRequest{ID: id})
	if resp.Error != nil {
		slog.Error("deleteArrRequest: Failed to remove from db", "error", resp.Error)
		return errors.New("failed when removing request")
	}
	return nil
}

// Gets all requests.
func getArrRequests(db *gorm.DB) ([]entity.ArrRequest, error) {
	var req []entity.ArrRequest
	resp := db.Preload("Content").Preload("User").Find(&req)
	if resp.Error != nil {
		slog.Error("getArrRequests: Failed to search for requests in db", "error", resp.Error)
		return []entity.ArrRequest{}, errors.New("failed to find requests")
	}
	for i := range req {
		req[i].Username = req[i].User.Username
	}
	return req, nil
}

func getArrRequest(db *gorm.DB, requestId uint) (entity.ArrRequest, error) {
	var req entity.ArrRequest
	resp := db.Where("id = ?", requestId).Take(&req)
	if resp.Error != nil {
		slog.Error("getArrRequest: Failed to search for request in db", "error", resp.Error)
		return entity.ArrRequest{}, errors.New("failed to find request")
	}
	return req, nil
}

func getArrRequestByTmdbId(db *gorm.DB, contentType entity.ContentType, tmdbId int) (entity.ArrRequest, error) {
	var req entity.ArrRequest
	resp := db.Joins("JOIN contents ON contents.id = arr_requests.content_id AND contents.tmdb_id = ? AND contents.type = ?", tmdbId, contentType).Find(&req)
	if resp.Error != nil {
		slog.Error("getArrRequestByTmdbId: Failed to search for request in db", "error", resp.Error)
		return entity.ArrRequest{}, errors.New("failed to find request")
	}
	return req, nil
}

func createArrRequest(db *gorm.DB, cp ContentProvider, userId uint, serverName string, contentType entity.ContentType, tmdbId int, reqJson string) (*entity.ArrRequest, error) {
	content, err := cp.GetOrCacheContent(contentType, tmdbId)
	if err != nil {
		slog.Error("createArrRequest: GetOrCacheContent errored.")
		return &entity.ArrRequest{}, err
	}
	req := entity.ArrRequest{UserID: userId, ServerName: serverName, ContentID: &content.ID, RequestJson: reqJson}
	resp := db.Create(&req)
	if resp.Error != nil {
		slog.Error("createArrRequest: Failed when inserting request into db.", "error", err)
		return &entity.ArrRequest{}, errors.New("failed when adding request")
	}
	return &req, nil
}

func createSonarrRequest(cfg *config.ServerConfig, db *gorm.DB, cp ContentProvider, userId uint, userPerms int, ur arr.SonarrRequest) (*entity.ArrRequest, error) {
	server, err := getSonarr(cfg, ur.ServerName)
	if err != nil {
		slog.Error("createSonarrRequest: Failed to get server", "error", err)
		return &entity.ArrRequest{}, errors.New("failed to get server")
	}
	reqJson, err := json.Marshal(ur)
	if err != nil {
		slog.Error("createRadarrRequest: Failed when marshalling json request", "error", err)
		return &entity.ArrRequest{}, errors.New("failed when processing request")
	}
	// Since we create the request in the db now, we don't have to check for duplicates, a unique constraint will error us here if hit.
	arrReq, err := createArrRequest(db, cp, userId, ur.ServerName, entity.SHOW, ur.TMDBID, string(reqJson[:]))
	if err != nil {
		slog.Error("createSonarrRequest: Failed when creating arr request", "error", err)
		return &entity.ArrRequest{}, errors.New("failed when creating request")
	}
	sonarr := arr.New(arr.SONARR, &server.Host, &server.Key)
	// 1. Lookup on Sonarr to check if the show has already been added (via method other than watcharr).
	lookupRes, err := sonarr.LookupByTmdbId(ur.TMDBID)
	if err == nil && len(lookupRes) == 1 {
		slog.Debug("createSonarrRequest: Lookup returned results.")
		found := lookupRes[0] // There should only be one result when looking up by id.
		// If it has an ID, then it will have already been added to Sonarr.
		if found.ID != 0 {
			dbResp := db.
				Model(&entity.ArrRequest{}).
				Where("id = ?", arrReq.ID).
				Update("arr_id", found.ID).
				Update("status", entity.ARR_REQUEST_FOUND)
			if dbResp.Error != nil {
				slog.Error("createSonarrRequest: Failed to update request in db", "error", err)
				return &entity.ArrRequest{}, errors.New("content was requested, but we failed to update the db")
			} else {
				slog.Debug("createSonarrRequest: Result from lookup had an ID. Request in database has been updated with it.", "arr_id", found.ID)
				arrReq.ArrID = found.ID
				return arrReq, nil
			}
		}
	}
	// 2. If user has auto approve perms, add movie to sonarr.
	if permission.Has(userPerms, entity.PERM_REQUEST_CONTENT_AUTO_APPROVE) {
		slog.Debug("createSonarrRequest: User has auto approve permission.. sending request to Sonarr.")
		ur.AutomaticSearch = server.AutomaticSearch
		resp, err := sonarr.AddContent(sonarr.BuildAddShowBody(ur))
		if err != nil {
			slog.Error("createSonarrRequest: Failed to add content", "error", err)
			return &entity.ArrRequest{}, errors.New("failed to add content")
		}
		dbResp := db.
			Model(&entity.ArrRequest{}).
			Where("id = ?", arrReq.ID).
			Update("arr_id", resp["id"]).
			Update("status", entity.ARR_REQUEST_AUTO_APPROVED)
		if dbResp.Error != nil {
			slog.Error("createSonarrRequest: Failed to update request in db", "error", err)
			return &entity.ArrRequest{}, errors.New("content was requested, but we failed to update the db")
		}
		arrId, ok := resp["id"].(float64)
		if !ok {
			slog.Error("createSonarrRequest: Failed to cast arr id as an int", "id", resp["id"])
			return &entity.ArrRequest{}, errors.New("failed to get arr id")
		}
		arrReq.ArrID = int(arrId)
		arrReq.Status = entity.ARR_REQUEST_AUTO_APPROVED
	}
	return arrReq, nil
}

func createRadarrRequest(cfg *config.ServerConfig, db *gorm.DB, cp ContentProvider, userId uint, userPerms int, ur arr.RadarrRequest) (*entity.ArrRequest, error) {
	server, err := getRadarr(cfg, ur.ServerName)
	if err != nil {
		slog.Error("createRadarrRequest: Failed to get server", "error", err)
		return &entity.ArrRequest{}, errors.New("failed to get server")
	}
	reqJson, err := json.Marshal(ur)
	if err != nil {
		slog.Error("createRadarrRequest: Failed when marshalling json request", "error", err)
		return &entity.ArrRequest{}, errors.New("failed when processing request")
	}
	// Since we create the request in the db now, we don't have to check for duplicates, a unique constraint will error us here if hit.
	arrReq, err := createArrRequest(db, cp, userId, ur.ServerName, entity.MOVIE, ur.TMDBID, string(reqJson[:]))
	if err != nil {
		slog.Error("createRadarrRequest: Failed when creating arr request", "error", err)
		return &entity.ArrRequest{}, errors.New("failed when creating request")
	}
	radarr := arr.New(arr.RADARR, &server.Host, &server.Key)
	// 1. Lookup on Radarr to check if the movie has already been added (via method other than watcharr).
	lookupRes, err := radarr.LookupByTmdbId(ur.TMDBID)
	if err == nil && len(lookupRes) == 1 {
		slog.Debug("createRadarrRequest: Lookup returned results.")
		found := lookupRes[0] // There should only be one result when looking up by id.
		// If it has an ID, then it will have already been added to Radarr.
		if found.ID != 0 {
			dbResp := db.
				Model(&entity.ArrRequest{}).
				Where("id = ?", arrReq.ID).
				Update("arr_id", found.ID).
				Update("status", entity.ARR_REQUEST_FOUND)
			if dbResp.Error != nil {
				slog.Error("createRadarrRequest: Failed to update request in db", "error", err)
				return &entity.ArrRequest{}, errors.New("content was requested, but we failed to update the db")
			} else {
				slog.Debug("createRadarrRequest: Result from lookup had an ID. Request in database has been updated with it.", "arr_id", found.ID)
				arrReq.ArrID = found.ID
				return arrReq, nil
			}
		}
	}
	// 2. If user has auto approve perms, add movie to radarr.
	if permission.Has(userPerms, entity.PERM_REQUEST_CONTENT_AUTO_APPROVE) {
		slog.Debug("createRadarrRequest: User has auto approve permission.. sending request to Radarr.")
		ur.AutomaticSearch = server.AutomaticSearch
		resp, err := radarr.AddContent(radarr.BuildAddMovieBody(ur))
		if err != nil {
			slog.Error("createRadarrRequest: Failed to add content", "error", err)
			return &entity.ArrRequest{}, errors.New("failed to add content")
		}
		dbResp := db.
			Model(&entity.ArrRequest{}).
			Where("id = ?", arrReq.ID).
			Update("arr_id", resp["id"]).
			Update("status", entity.ARR_REQUEST_AUTO_APPROVED)
		if dbResp.Error != nil {
			slog.Error("createRadarrRequest: Failed to update request in db", "error", err)
			return &entity.ArrRequest{}, errors.New("content was requested, but we failed to update the db")
		}
		arrId, ok := resp["id"].(float64)
		if !ok {
			slog.Error("createRadarrRequest: Failed to cast arr id as an int", "id", resp["id"])
			return &entity.ArrRequest{}, errors.New("failed to get arr id")
		}
		arrReq.ArrID = int(arrId)
		arrReq.Status = entity.ARR_REQUEST_AUTO_APPROVED
	}
	return arrReq, nil
}

func getRadarrRequestInfo(cfg *config.ServerConfig, db *gorm.DB, requestId uint) (arr.MovieSerie, error) {
	if requestId == 0 {
		slog.Error("sonarr info: No request id provided")
		return arr.MovieSerie{}, errors.New("no request id provided")
	}
	arrRequest, err := getArrRequest(db, requestId)
	if err != nil {
		slog.Error("radarr info: Failed to get server", "error", err)
		return arr.MovieSerie{}, errors.New("failed to get server")
	}
	server, err := getRadarr(cfg, arrRequest.ServerName)
	if err != nil {
		slog.Error("radarr info: Failed to get server", "error", err)
		return arr.MovieSerie{}, errors.New("failed to get server")
	}
	radarr := arr.New(arr.RADARR, &server.Host, &server.Key)
	resp, respStatusCode, err := radarr.GetContent(arrRequest.ArrID)
	if err != nil {
		slog.Error("radarr info: Failed to get info", "error", err)
		if (arrRequest.Status == entity.ARR_REQUEST_APPROVED || arrRequest.Status == entity.ARR_REQUEST_AUTO_APPROVED) && respStatusCode == 404 {
			slog.Error("radarr info: 404 returned.. content must've been removed.. removing request.")
			err := deleteArrRequest(db, arrRequest.ID)
			if err != nil {
				return arr.MovieSerie{}, errors.New("failed to delete request for removed content")
			} else {
				return arr.MovieSerie{}, errors.New("request deleted")
			}
		}
		return arr.MovieSerie{}, errors.New("failed to get info")
	}
	return resp, nil
}

func getSonarrRequestInfo(cfg *config.ServerConfig, db *gorm.DB, requestId uint) (arr.MovieSerie, error) {
	if requestId == 0 {
		slog.Error("sonarr info: No request id provided")
		return arr.MovieSerie{}, errors.New("no request id provided")
	}
	arrRequest, err := getArrRequest(db, requestId)
	if err != nil {
		slog.Error("sonarr info: Failed to get server", "error", err)
		return arr.MovieSerie{}, errors.New("failed to get server")
	}
	server, err := getSonarr(cfg, arrRequest.ServerName)
	if err != nil {
		slog.Error("sonarr info: Failed to get server", "error", err)
		return arr.MovieSerie{}, errors.New("failed to get server")
	}
	sonarr := arr.New(arr.SONARR, &server.Host, &server.Key)
	resp, respStatusCode, err := sonarr.GetContent(arrRequest.ArrID)
	if err != nil {
		slog.Error("sonarr info: Failed to get info", "error", err)
		if (arrRequest.Status == entity.ARR_REQUEST_APPROVED || arrRequest.Status == entity.ARR_REQUEST_AUTO_APPROVED) && respStatusCode == 404 {
			slog.Error("sonarr info: 404 returned.. content must've been removed.. removing request.")
			err := deleteArrRequest(db, arrRequest.ID)
			if err != nil {
				return arr.MovieSerie{}, errors.New("failed to delete request for removed content")
			} else {
				return arr.MovieSerie{}, errors.New("request deleted")
			}
		}
		return arr.MovieSerie{}, errors.New("failed to get info")
	}
	return resp, nil
}
