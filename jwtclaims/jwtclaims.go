// Package jwtclaims provide the Claims type providing helper functions to work with JWT claims
package jwtclaims

import (
	"strconv"
	"time"
)

const (
	Iss       = "iss"
	Jti       = "jti"
	Scope     = "scope"
	Sub       = "sub"
	Exp       = "exp"
	Iat       = "iat"
	Nbf       = "nbf"
	Username  = "username"
	ClientID  = "client_id"
	OriginJti = "origin_jti"
	EventID   = "event_id"
	TokenUse  = "token_use"
	AuthTime  = "auth_time"
)

// Claims is a map of JWT claim names and values
type Claims map[string]string

// Iss returns the value of the iss claim
func (c Claims) Iss() string {
	return c[Iss]
}

// Jti returns the value of the jti claim
func (c Claims) Jti() string {
	return c[Jti]
}

// Scope returns the value of the scope claim
func (c Claims) Scope() string {
	return c[Scope]
}

// sub returns the value of the sub claim
func (c Claims) Sub() string {
	return c[Sub]
}

// ClientID returns the value of the client_id claim
func (c Claims) ClientID() string {
	return c[ClientID]
}

// OriginJti returns the value of the origin_jti claim
func (c Claims) OriginJti() string {
	return c[OriginJti]
}

// EventID returns the value of the event_id claim
func (c Claims) EventID() string {
	return c[EventID]
}

// TokenUse returns the value of the token_use claim
func (c Claims) TokenUse() string {
	return c[TokenUse]
}

// Username returns the value of the username claim
func (c Claims) Username() string {
	return c[Username]
}

// Iat returns the value of the iat claim, converted to a time.Time.
// If iat is either undefined or not a valid timestamp, the zero value of time.Time is returned
func (c Claims) Iat() time.Time {
	return c.Timestamp(Iat)
}

// Exp returns the value of the exp claim, converted to a time.Time.
// If exp is either undefined or not a valid timestamp, the zero value of time.Time is returned
func (c Claims) Exp() time.Time {
	return c.Timestamp(Exp)
}

// Nbf returns the value of the nbf claim, converted to a time.Time.
// If nbf is either undefined or not a valid timestamp, the zero value of time.Time is returned
func (c Claims) Nbf() time.Time {
	return c.Timestamp(Nbf)
}

// AuthTime returns the value of the auth_time claim, converted to a time.Time.
// If auth_time is either undefined or not a valid timestamp, the zero value of time.Time is returned
func (c Claims) AuthTime() time.Time {
	return c.Timestamp(AuthTime)
}

// Timestamp returns the value of the given claim converted to a time.Time.
// If the value of the claim is either undefined or not a valid timestamp, the zero value of time.Time is returned
func (c Claims) Timestamp(claim string) time.Time {
	v := c[claim]
	if v == "" {
		return time.Time{}
	}

	t, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(t, 0)
}
