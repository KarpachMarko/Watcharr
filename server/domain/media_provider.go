package domain

type MediaProviderType string

const (
	MediaProviderTypeSub  MediaProviderType = "sub"
	MediaProviderTypeFree MediaProviderType = "free"
)

type MediaProvider struct {
	// Name of the provider.
	Name string `json:"name,omitempty"`
	// The type of service provided by the provider.
	Type MediaProviderType `json:"type,omitempty"`
	// Link to watch.
	// We can't get a direct link to content from tmdb, we are told to link
	// to tmdb to support them instead and from there the user can see deep links.
	Link string `json:"link,omitempty"`
}
