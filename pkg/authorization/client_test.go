package authorization

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	// Global UUIDs for consistent testing
	TestEntityID   = uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	TestResourceID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	TestRightID1   = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	TestRightID2   = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	Namespace      = "testing-dev"
)

func TestClient_GiveRights(t *testing.T) {
	client := NewClient(nil)
	client.namespace = Namespace

	giveRightsReq := GiveRightsRequest{
		Entity: TestEntityID,
		Rights: map[string]Right{
			TestResourceID.String(): {
				UID:       TestRightID1,
				Entity:    TestEntityID,
				Resource:  TestResourceID,
				Read:      true,
				Write:     true,
				Delete:    false,
				Share:     false,
				Custom:    map[string]bool{"admin": true},
				AssetType: "test_asset",
				Active:    true,
				Created:   1678886400,
			},
		},
	}

	resp, err := client.GiveRights(&giveRightsReq)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Valid)
	assert.Equal(t, TestEntityID, resp.Entity)
}

func TestClient_GetRights(t *testing.T) {
	client := NewClient(nil)
	client.namespace = Namespace

	getRightsReq := GetRightsRequest{
		Entity: TestEntityID,
	}

	resp, err := client.GetRights(&getRightsReq)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Valid)
	assert.Equal(t, TestEntityID, resp.Entity)
	assert.IsType(t, map[string]Right{}, resp.Rights)
}

func TestClient_HasRights(t *testing.T) {
	client := NewClient(nil)
	client.namespace = Namespace

	hasRightsReq := HasRightsRequest{
		Entity:    TestEntityID,
		Resources: []string{TestResourceID.String(), TestResourceID.String()},
	}

	resp, err := client.HasRights(&hasRightsReq)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Valid)
	assert.Equal(t, TestEntityID, resp.Entity)
}

func TestClient_RevokeRights(t *testing.T) {
	client := NewClient(nil)
	client.namespace = Namespace

	revokeRightsReq := RevokeRightsRequest{
		Entity: TestEntityID,
		Rights: []uuid.UUID{TestRightID1, TestRightID2},
	}

	resp, err := client.RevokeRights(&revokeRightsReq)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Valid)
	assert.Equal(t, TestEntityID, resp.Entity)
}
