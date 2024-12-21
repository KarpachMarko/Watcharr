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
	// What is the name of the trusted header that
	// will contain the logged in users username?
	HEADER_NAME string `json:",omitempty"`
	// Should the frontend attempt auto login if
	// trusted header auth is enabled.
	AUTO_LOGIN bool `json:",omitempty"`
	// Where can we redirect the user to logout
	// of the auth service?
	LOGOUT_URL string `json:",omitempty"`
}

// Is trusted header auth configured on this server?
// If yes, then it is enabled.
func trustedHeaderAuthIsEnabled() bool {
	return Config.HEADER_AUTH.HEADER_NAME != ""
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
