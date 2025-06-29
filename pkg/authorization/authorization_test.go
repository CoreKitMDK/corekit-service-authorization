package authorization

import (
	"context"
	"github.com/CoreKitMDK/corekit-service-core/v2/pkg/core"
	"github.com/CoreKitMDK/corekit-service-logger/v2/pkg/logger"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAuthorizationDAL(t *testing.T) {
	var Core, _ = core.NewCore()
	connStr, err := Core.Configuration.Get("internal-authorization-db")
	if err != nil {
		Core.Logger.Log(logger.FATAL, "failed to get database connection string: "+err.Error())
	}

	dal, err := NewAuthorizationDAL(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dal.Close()

	// Ensure table exists for testing
	_, err = dal.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS rights (
		uid UUID PRIMARY KEY,
		entity UUID NOT NULL,
		resource UUID NOT NULL,
		read BOOLEAN NOT NULL DEFAULT false,
		write BOOLEAN NOT NULL DEFAULT false,
		"delete" BOOLEAN NOT NULL DEFAULT false,
		share BOOLEAN NOT NULL DEFAULT false,
		custom JSONB,
		usage_total INTEGER NOT NULL DEFAULT 0,
		usage_used INTEGER NOT NULL DEFAULT 0,
		usage_resets_in_seconds BIGINT NOT NULL DEFAULT 0,
		asset_type TEXT,
		active BOOLEAN NOT NULL DEFAULT false,
		created BIGINT NOT NULL
	);`)
	if err != nil {
		t.Fatalf("Failed to create rights table: %v", err)
	}

	// Cleanup after test
	defer dal.db.Exec(context.Background(), "DELETE FROM rights WHERE entity = $1", 123)

	testEntity := uuid.New()
	testRightUID := uuid.New()
	testResource := uuid.New()

	// 1. GiveRights (Add)
	giveReq := &GiveRightsRequest{
		Entity: testEntity,
		Rights: map[string]Right{
			"test_resource": {
				UID:      testRightUID,
				Entity:   testEntity,
				Resource: testResource,
				Read:     true,
				Created:  time.Now().Unix(),
				Active:   true,
			},
		},
	}
	_, err = dal.GiveRights(giveReq)
	if err != nil {
		t.Fatalf("GiveRights failed: %v", err)
	}

	// 2. GetRights (List)
	getReq := &GetRightsRequest{Entity: testEntity}
	getResp, err := dal.GetRights(getReq)
	if err != nil || !getResp.Valid || len(getResp.Rights) != 1 || getResp.Rights[testRightUID.String()].UID != testRightUID {
		t.Fatalf("GetRights failed or returned unexpected result: %v, %v", err, getResp)
	}

	// 3. HasRights
	hasReq := &HasRightsRequest{Entity: testEntity, Resources: []string{testResource.String()}}
	hasResp, err := dal.HasRights(hasReq)
	if err != nil || !hasResp.Valid || !hasResp.Rights[testResource.String()].Read {
		t.Fatalf("HasRights failed or returned unexpected result: %v, %v", err, hasResp)
	}

	// 4. RevokeRights (Delete)
	revokeReq := &RevokeRightsRequest{Entity: testEntity, Rights: []uuid.UUID{testRightUID}}
	_, err = dal.RevokeRights(revokeReq)
	if err != nil {
		t.Fatalf("RevokeRights failed: %v", err)
	}

	// Verify deletion
	getResp, err = dal.GetRights(getReq)
	if err != nil || !getResp.Valid || len(getResp.Rights) != 0 {
		t.Fatalf("RevokeRights verification failed: %v, %v", err, getResp)
	}
}
