package api

import "context"

// Most Bot API results are rendered generically from their raw JSON (the API is
// method-oriented, not a handful of typed resources). We keep typed structs only where the
// code itself needs to read fields — chiefly getMe, used by auth/doctor/whoami.

// User is a Telegram user or bot. https://core.telegram.org/bots/api#user
type User struct {
	ID                      ID     `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name,omitempty"`
	Username                string `json:"username,omitempty"`
	LanguageCode            string `json:"language_code,omitempty"`
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`
}

// DisplayName is a human label: "@username" when present, else the first/last name.
func (u User) DisplayName() string {
	if u.Username != "" {
		return "@" + u.Username
	}
	name := u.FirstName
	if u.LastName != "" {
		name += " " + u.LastName
	}
	return name
}

// GetMe returns the authenticated bot's identity. It is idempotent and used to verify a
// token end-to-end (auth login / auth status / doctor).
func (c *Client) GetMe(ctx context.Context) (*User, error) {
	var u User
	if err := c.CallInto(ctx, "getMe", nil, true, &u); err != nil {
		return nil, err
	}
	return &u, nil
}
