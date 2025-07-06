package authorization

import (
	"github.com/CoreKitMDK/corekit-service-core/v2/pkg/core"
	"github.com/CoreKitMDK/corekit-service-logger/v2/pkg/logger"
	"github.com/google/uuid"
	"testing"
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
