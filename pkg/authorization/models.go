package authorization

import (
	"github.com/google/uuid"
)

type HasRightsRequest struct {
	Entity    uuid.UUID
	Resources []string
}

type HasRightsResponse struct {
	Entity uuid.UUID
	Rights map[string]Right
	Valid  bool
	Error  string
}

type GiveRightsRequest struct {
	Entity uuid.UUID
	Rights map[string]Right
}

type GiveRightsResponse struct {
	Entity uuid.UUID
	Valid  bool
	Error  string
}

type RevokeRightsRequest struct {
	Entity uuid.UUID
	Rights []uuid.UUID
}

type RevokeRightsResponse struct {
	Entity uuid.UUID
	Valid  bool
	Error  string
}

type GetRightsRequest struct {
	Entity uuid.UUID
}

type GetRightsResponse struct {
	Entity uuid.UUID
	Rights map[string]Right
	Valid  bool
	Error  string
}

type Right struct {
	UID      uuid.UUID
	Entity   uuid.UUID
	Resource uuid.UUID

	Read   bool
	Write  bool
	Delete bool
	Share  bool

	Custom map[string]bool

	UsageTotal           int
	UsageUsed            int
	UsageResetsInSeconds int64

	AssetType string

	Active  bool
	Created int64
}

func NewRight() Right {
	return Right{
		UID:    uuid.Must(uuid.NewV7()),
		Custom: make(map[string]bool),
	}
}
