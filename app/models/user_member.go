package models

import (
	"encoding/json"
	"strings"

	"github.com/innovay-software/famapp-main/app/utils"
)

type UserMember struct {
	ID     uint64
	Role   string
	Name   string
	Mobile string
	Avatar string
}

func (UserMember) TableName() string {
	return "users"
}

func (u UserMember) MarshalJSON() ([]byte, error) {
	// Define a temporary struct to hold the marshalled data
	type UserMemberMarshal struct {
		ID     uint64 `json:"id"`
		Role   string `json:"role"`
		Name   string `json:"name"`
		Mobile string `json:"mobile"`
		Avatar string `json:"avatar"`
	}

	// Construct the MarshalUser object
	if u.Avatar != "" && !strings.HasPrefix(u.Avatar, "http") {
		u.Avatar = utils.GetUrlPath("avatars", u.Avatar)
	}

	return json.Marshal(UserMemberMarshal(u))
}
