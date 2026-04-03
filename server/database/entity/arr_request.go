package entity

import "time"

type ArrRequestStatus string

const (
	// Pending approval from an admin.
	ARR_REQUEST_PENDING ArrRequestStatus = "PENDING"
	// Request has been approved and should be added to sonarr/radarr.
	ARR_REQUEST_APPROVED      ArrRequestStatus = "APPROVED"
	ARR_REQUEST_AUTO_APPROVED ArrRequestStatus = "AUTO_APPROVED"
	// Request has been denied, not adding content.
	ARR_REQUEST_DENIED ArrRequestStatus = "DENIED"
	// Content was found on sonarr/radarr already, nothing needs to be done.
	ARR_REQUEST_FOUND ArrRequestStatus = "FOUND"
)

type ArrRequest struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    uint      `json:"-" gorm:"not null"`
	User      User      `json:"-"`
	// Username of `User`.
	// We don't want to send back the entire user object, just their name.
	// Not stored in DB, only used for our response from api.
	Username  string   `json:"username" gorm:"-"`
	ContentID *int     `json:"-" gorm:"uniqueIndex:sn_to_cid;not null"`
	Content   *Content `json:"content,omitempty"`
	// Server names are used as an identifier
	ServerName string `json:"serverName" gorm:"uniqueIndex:sn_to_cid;not null"`
	// Sonarr/Radarrs seriesId/movieId
	ArrID int `json:"arrId"`
	// Tracked request status
	Status ArrRequestStatus `json:"status" gorm:"default:PENDING"`
	// Full request made by user (arr.SonarrRequest / arr.RadarrRequest)
	// so we know how to fulfil the request if approved.
	RequestJson string `json:"requestJson"`
}
