package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserType uint8

var (
	WATCHARR_USER UserType = 0
	JELLYFIN_USER UserType = 1
	PLEX_USER     UserType = 2
	// Registered via trusted header auth
	PROXY_USER UserType = 3
)

// User Perms
// iota auto increments for us so when adding new
// perms, add to bottom as to not change other perm
// values.
const (
	PERM_NONE int = 1 << iota
	PERM_ADMIN
	PERM_REQUEST_CONTENT
	PERM_REQUEST_CONTENT_AUTO_APPROVE
)

// Holds third party service auth tokens for users.
// Each service may use the fields in their own way.
// Unique index applied between service name and clientID
// to ensure no duplicates (no need to apply it against
// user_id, no accounts should share an integration).
//
// Plex:
//   - AuthToken  : Used for requests against plex.tv
//   - AuthToken2 : Used for requests against home plex server.
type UserServices struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// Service/integration name
	Name string `gorm:"uniqueIndex:svc_name_to_cltid;not null;" json:"-"`
	// The users id on the third party service
	ClientID  string `gorm:"uniqueIndex:svc_name_to_cltid;not null;" json:"-"`
	AuthToken string `gorm:"not null;" json:"-"`
	// Second auth token, generic name so future services can use it without extra confusion.
	// Ex: We require a second auth token for use with our local server for Plex.
	AuthToken2 string `json:"-"`
	UserID     uint   `gorm:"not null;" json:"-"`
}

type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func GetPassArgonParams() *ArgonParams {
	return &ArgonParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

type TokenClaims struct {
	UserID   uint     `json:"userId"`
	Username string   `json:"username"`
	Type     UserType `json:"type"`
	jwt.RegisteredClaims
}
