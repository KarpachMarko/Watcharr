// auth_proxy contains the logic for Trusted Header Authentication/SSO.
//
// This code is inherently dangerous since we are implicitly trusting
// a header for auth, so this should only be configured if you are
// certain your watcharr instance is only available behind your proxy.

package main

import (
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type TrustedHeaderAuthSetting struct {
	// Required: Should header auth be enabled?
	// This bool exists so header auth can be toggled
	// easily without having to remove configuration.
	// To be actually enabled, HEADER_NAME must also
	// be set.
	Enabled bool `json:"enabled,omitempty"`
	// Required: What is the name of the trusted header
	// that will contain the logged in users username?
	HeaderName string `json:"headerName,omitempty"`
	// Should the frontend attempt auto login if
	// trusted header auth is enabled.
	AutoLogin bool `json:"autoLogin,omitempty"`
	// Where can we redirect the user to logout
	// of the auth service?
	LogoutUrl string `json:"logoutUrl,omitempty"`
}

type TrustedHeaderAuthLogoutDetailsResponse struct {
	LogoutUrl string `json:"logoutUrl,omitempty"`
}

// Is trusted header auth configured on this server?
func trustedHeaderAuthIsEnabled() bool {
	return Config.HEADER_AUTH.Enabled && Config.HEADER_AUTH.HeaderName != ""
}

func setTrustedHeaderAuthSetting(has TrustedHeaderAuthSetting) error {
	slog.Debug("setTrustedHeaderAuthSetting: Attempting to update to new provided value", "new_value", has)
	Config.HEADER_AUTH = has
	err := writeConfig()
	if err != nil {
		slog.Error("setTrustedHeaderAuthSetting: Failed to write updated config!", "error", err)
		return errors.New("failed to write config")
	}
	return nil
}

// Gets proxy logout details.
// Details are accessible to any user for the logout flow.
// If proxy configured should be checked before using this.
func getTrustedHeaderAuthLogoutDetails() *TrustedHeaderAuthLogoutDetailsResponse {
	return &TrustedHeaderAuthLogoutDetailsResponse{
		LogoutUrl: Config.HEADER_AUTH.LogoutUrl,
	}
}

// Login via header sso
func loginTrustedHeaderAuth(user *User, db *gorm.DB) (AuthResponse, error) {
	slog.Debug("loginTrustedHeaderAuth: A user is logging in", "username_from_header", user.Username)
	dbUser := new(User)
	res := db.Where("username = ? AND (type IS NULL OR type = 0 OR type = ?)", user.Username, PROXY_USER).Take(&dbUser)
	if res.Error != nil {
		slog.Debug("loginTrustedHeaderAuth: Creating new User from authentication header", "username_from_header", user.Username)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// Record not found, so we should create the user (if configured to do so)
			// dbUser will be empty, so we can just reuse it for this purpose.
			dbUser.Username = user.Username
			dbUser.Type = PROXY_USER
			dbUser.Country = &Config.DEFAULT_COUNTRY

			res = db.Create(&dbUser)
			if res.Error != nil {
				slog.Error("loginTrustedHeaderAuth: Failed to create new user in db", "error", res.Error)
				return AuthResponse{}, errors.New("failed to create new user")
			}
		} else {
			slog.Error("loginTrustedHeaderAuth: An error occurred when looking up user in db", "error", res.Error)
			return AuthResponse{}, errors.New("error locating user in db")
		}
	}
	token, err := signJWT(dbUser)
	if err != nil {
		slog.Error("loginTrustedHeaderAuth: Failed to sign new jwt", "error", err)
		return AuthResponse{}, errors.New("failed to get auth token")
	}
	return AuthResponse{Token: token}, nil
}
