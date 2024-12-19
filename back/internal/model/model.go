package model

import (
	"errors"
)

type Wishlist struct {
	ID          int64       `json:"id"`
	OwnerID     int64       `json:"owner_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	DisplayType DisplayType `json:"display_type"`
}

type DisplayType int

const (
	DisplayTypeNone = DisplayType(iota)
	DisplayTypePublic
	DisplayTypeFriendsOnly
	DisplayTypeByLink
)

var ErrorInvalidDisplayType = errors.New("invalid select display type")

func IntToDisplayType(displayType int) (DisplayType, error) {
	switch displayType {
	case 0:
		return DisplayTypeNone, nil
	case 1:
		return DisplayTypePublic, nil
	case 2:
		return DisplayTypeFriendsOnly, nil
	case 3:
		return DisplayTypeByLink, nil
	}
	return DisplayType(-1), ErrorInvalidDisplayType
}
