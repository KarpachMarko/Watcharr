package cfgmodel

type ArrSettings struct {
	Name string `json:"name,omitempty"`
	Host string `json:"host,omitempty"`
	Key  string `json:"key,omitempty"`
}

type SonarrSettings struct {
	ArrSettings
	QualityProfile  int  `json:"qualityProfile,omitempty"`
	RootFolder      int  `json:"rootFolder,omitempty"`
	LanguageProfile int  `json:"languageProfile,omitempty"`
	AutomaticSearch bool `json:"automaticSearch"`
	// TODO eventually separate profiles and root for anime
	//  content (i can see diff language profile being useful)
}

type RadarrSettings struct {
	ArrSettings
	QualityProfile  int  `json:"qualityProfile,omitempty"`
	RootFolder      int  `json:"rootFolder,omitempty"`
	AutomaticSearch bool `json:"automaticSearch"`
}
