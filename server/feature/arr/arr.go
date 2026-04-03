package arr

import (
	"errors"
	"log/slog"

	"github.com/sbondCo/Watcharr/arr"
	"github.com/sbondCo/Watcharr/config"
	"github.com/sbondCo/Watcharr/config/cfgmodel"
)

type (
	ArrTestParams struct {
		Host string `json:"host,omitempty"`
		Key  string `json:"key,omitempty"`
	}

	SonarrTestResponse struct {
		QualityProfiles  []arr.QualityProfile  `json:"qualityProfiles"`
		RootFolders      []arr.RootFolder      `json:"rootFolders"`
		LanguageProfiles []arr.LanguageProfile `json:"languageProfiles"`
	}

	RadarrTestResponse struct {
		QualityProfiles  []arr.QualityProfile  `json:"qualityProfiles"`
		RootFolders      []arr.RootFolder      `json:"rootFolders"`
		LanguageProfiles []arr.LanguageProfile `json:"languageProfiles"`
	}
)

// Getting arr servers for public consumption!!!!
// Nothing sensitive to be included!
type arrSettingsPublicResponseBase struct {
	Name string `json:"name"`
	Host string `json:"host"`

	QualityProfile  int  `json:"qualityProfile,omitempty"`
	RootFolder      int  `json:"rootFolder,omitempty"`
	AutomaticSearch bool `json:"automaticSearch"`
}

func newArrSettingsPublicResponseBaseSonarr(
	cfg *cfgmodel.SonarrSettings,
) *arrSettingsPublicResponseBase {
	return &arrSettingsPublicResponseBase{
		Name:            cfg.Name,
		Host:            cfg.Host,
		QualityProfile:  cfg.QualityProfile,
		RootFolder:      cfg.RootFolder,
		AutomaticSearch: cfg.AutomaticSearch,
	}
}

func newArrSettingsPublicResponseBaseRadarr(
	cfg *cfgmodel.RadarrSettings,
) *arrSettingsPublicResponseBase {
	return &arrSettingsPublicResponseBase{
		Name:            cfg.Name,
		Host:            cfg.Host,
		QualityProfile:  cfg.QualityProfile,
		RootFolder:      cfg.RootFolder,
		AutomaticSearch: cfg.AutomaticSearch,
	}
}

// Getting sonarr servers response for ALL USERS!
// Nothing sensitive to be included!
type SonarrSettingsPublicResponseResult struct {
	arrSettingsPublicResponseBase
	LanguageProfile int `json:"languageProfile,omitempty"`
}

func NewSonarrSettingsPublicResponse(
	cfg *cfgmodel.SonarrSettings,
) SonarrSettingsPublicResponseResult {
	base := newArrSettingsPublicResponseBaseSonarr(cfg)
	return SonarrSettingsPublicResponseResult{
		arrSettingsPublicResponseBase: *base,
		LanguageProfile:               cfg.LanguageProfile,
	}
}

// Getting radarr servers response for ALL USERS!
// Nothing sensitive to be included!
type RadarrSettingsPublicResponseResult struct {
	arrSettingsPublicResponseBase
}

func NewRadarrSettingsPublicResponse(
	cfg *cfgmodel.RadarrSettings,
) RadarrSettingsPublicResponseResult {
	base := newArrSettingsPublicResponseBaseRadarr(cfg)
	return RadarrSettingsPublicResponseResult{
		arrSettingsPublicResponseBase: *base,
	}
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
func getSonarrsSafe(cfg *config.ServerConfig) []SonarrSettingsPublicResponseResult {
	s := []SonarrSettingsPublicResponseResult{}
	for _, v := range cfg.SONARR {
		s = append(s, NewSonarrSettingsPublicResponse(&v))
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
func getRadarrsSafe(cfg *config.ServerConfig) []RadarrSettingsPublicResponseResult {
	s := []RadarrSettingsPublicResponseResult{}
	for _, v := range cfg.RADARR {
		s = append(s, NewRadarrSettingsPublicResponse(&v))
	}
	return s
}
