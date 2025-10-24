package arr

import (
	"errors"
	"log/slog"

	"github.com/sbondCo/Watcharr/arr"
	"github.com/sbondCo/Watcharr/config"
	"github.com/sbondCo/Watcharr/config/cfgmodel"
)

type ArrTestParams struct {
	Host string `json:"host,omitempty"`
	Key  string `json:"key,omitempty"`
}

type SonarrTestResponse struct {
	QualityProfiles  []arr.QualityProfile  `json:"qualityProfiles"`
	RootFolders      []arr.RootFolder      `json:"rootFolders"`
	LanguageProfiles []arr.LanguageProfile `json:"languageProfiles"`
}

type RadarrTestResponse struct {
	QualityProfiles  []arr.QualityProfile  `json:"qualityProfiles"`
	RootFolders      []arr.RootFolder      `json:"rootFolders"`
	LanguageProfiles []arr.LanguageProfile `json:"languageProfiles"`
}

// Response given to users with PERM_REQUEST_CONTENT - should never include sensitive info
func testSonarr(p ArrTestParams) (SonarrTestResponse, error) {
	sonarr := arr.New(arr.SONARR, &p.Host, &p.Key)
	qps, err := sonarr.GetQualityProfiles()
	if err != nil {
		slog.Error("testSonarr failed to get quality profiles!", "error", err)
		return SonarrTestResponse{}, errors.New("failed to get quality profiles")
	}
	rfs, err := sonarr.GetRootFolders()
	if err != nil {
		slog.Error("testSonarr failed to get root folders!", "error", err)
		return SonarrTestResponse{}, errors.New("failed to get root folders")
	}
	lps, err := sonarr.GetLangaugeProfiles()
	if err != nil {
		slog.Error("testSonarr failed to get language profiles!", "error", err)
		return SonarrTestResponse{}, errors.New("failed to get language profiles")
	}
	return SonarrTestResponse{QualityProfiles: qps, RootFolders: rfs, LanguageProfiles: lps}, nil
}

// Response given to users with PERM_REQUEST_CONTENT - should never include sensitive info
func testRadarr(p ArrTestParams) (RadarrTestResponse, error) {
	radarr := arr.New(arr.RADARR, &p.Host, &p.Key)
	qps, err := radarr.GetQualityProfiles()
	if err != nil {
		slog.Error("testRadarr failed to get quality profiles!", "error", err)
		return RadarrTestResponse{}, errors.New("failed to get quality profiles")
	}
	rfs, err := radarr.GetRootFolders()
	if err != nil {
		slog.Error("testRadarr failed to get root folders!", "error", err)
		return RadarrTestResponse{}, errors.New("failed to get root folders")
	}
	return RadarrTestResponse{QualityProfiles: qps, RootFolders: rfs}, nil
}

// TODO any way to simplify (deduplicate/reuse) these
// methods (and the whole file tbh) would be very good

// Add sonarr server to config
func addSonarr(cfg *config.ServerConfig, s cfgmodel.SonarrSettings) error {
	for _, v := range cfg.SONARR {
		if v.Name == s.Name {
			// Server exists with this name...
			return errors.New("server with that name already exists")
		}
	}
	cfg.SONARR = append(cfg.SONARR, s)
	cfg.Write()
	return nil
}

// Edit sonarr server in config
func editSonarr(cfg *config.ServerConfig, s cfgmodel.SonarrSettings) error {
	for i, v := range cfg.SONARR {
		if v.Name == s.Name {
			cfg.SONARR[i] = s
			cfg.Write()
			return nil
		}
	}
	return errors.New("can't edit server that does not exist")
}

func rmSonarr(cfg *config.ServerConfig, name string) error {
	for i, v := range cfg.SONARR {
		if v.Name == name {
			cfg.SONARR = append(cfg.SONARR[:i], cfg.SONARR[i+1:]...)
			cfg.Write()
			return nil
		}
	}
	return errors.New("can't remove a server that does not exist")
}

func getSonarr(cfg *config.ServerConfig, name string) (cfgmodel.SonarrSettings, error) {
	for i, v := range cfg.SONARR {
		if v.Name == name {
			return cfg.SONARR[i], nil
		}
	}
	return cfgmodel.SonarrSettings{}, errors.New("server not found")
}

// Get list of sonarr servers without api keys.
// Regular users with access to adding to sonarr will request this.
func getSonarrsSafe(cfg *config.ServerConfig) []cfgmodel.SonarrSettings {
	s := []cfgmodel.SonarrSettings{}
	for _, v := range cfg.SONARR {
		s = append(s, v.Safe())
	}
	return s
}

// Add radarr server to config
func addRadarr(cfg *config.ServerConfig, s cfgmodel.RadarrSettings) error {
	for _, v := range cfg.RADARR {
		if v.Name == s.Name {
			// Server exists with this name...
			return errors.New("server with that name already exists")
		}
	}
	cfg.RADARR = append(cfg.RADARR, s)
	cfg.Write()
	return nil
}

// Edit radarr server in config
func editRadarr(cfg *config.ServerConfig, s cfgmodel.RadarrSettings) error {
	for i, v := range cfg.RADARR {
		if v.Name == s.Name {
			cfg.RADARR[i] = s
			cfg.Write()
			return nil
		}
	}
	return errors.New("can't edit server that does not exist")
}

func rmRadarr(cfg *config.ServerConfig, name string) error {
	for i, v := range cfg.RADARR {
		if v.Name == name {
			cfg.RADARR = append(cfg.RADARR[:i], cfg.RADARR[i+1:]...)
			cfg.Write()
			return nil
		}
	}
	return errors.New("can't remove a server that does not exist")
}

func getRadarr(cfg *config.ServerConfig, name string) (cfgmodel.RadarrSettings, error) {
	for i, v := range cfg.RADARR {
		if v.Name == name {
			return cfg.RADARR[i], nil
		}
	}
	return cfgmodel.RadarrSettings{}, errors.New("server not found")
}

// Get list of radarr servers without api keys.
// Regular users with access to adding to radarr will request this.
func getRadarrsSafe(cfg *config.ServerConfig) []cfgmodel.RadarrSettings {
	s := []cfgmodel.RadarrSettings{}
	for _, v := range cfg.RADARR {
		s = append(s, v.Safe())
	}
	return s
}
