package authorization

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type AuthorizationDAL struct {
	db *pgx.Conn
}

func NewAuthorizationDAL(connString string) (*AuthorizationDAL, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &AuthorizationDAL{db: conn}, nil
}

func (dal *AuthorizationDAL) Close() {
	dal.db.Close(context.Background())
}

func (dal *AuthorizationDAL) HasRights(req *HasRightsRequest) (*HasRightsResponse, error) {
	query := `SELECT uid, entity, resource, read, write, delete, share, custom, usage_total, usage_used, usage_resets_in_seconds, asset_type, created_at
			   FROM rights
			   WHERE entity = $1 AND resource = ANY($2) AND active = true`

	rows, err := dal.db.Query(context.Background(), query, req.Entity, req.Resources)
	if err != nil {
		return &HasRightsResponse{Valid: false, Error: err.Error()}, err
	}
	defer rows.Close()

	rights := make(map[string]Right)
	for rows.Next() {
		var r Right
		var resourceStr string
		if err := rows.Scan(&r.UID, &r.Entity, &resourceStr, &r.Read, &r.Write, &r.Delete, &r.Share, &r.Custom, &r.UsageTotal, &r.UsageUsed, &r.UsageResetsInSeconds, &r.AssetType, &r.CreatedAt); err != nil {
			return &HasRightsResponse{Valid: false, Error: err.Error()}, err
		}
		rights[resourceStr] = r
	}

	return &HasRightsResponse{
		Entity: req.Entity,
		Rights: rights,
		Valid:  true,
	}, nil
}

func (dal *AuthorizationDAL) GiveRights(req *GiveRightsRequest) (*GiveRightsResponse, error) {
	tx, err := dal.db.Begin(context.Background())
	if err != nil {
		return &GiveRightsResponse{Valid: false, Error: err.Error()}, err
	}
	defer tx.Rollback(context.Background())

	for _, right := range req.Rights {
		query := `INSERT INTO rights (uid, entity, resource, read, write, delete, share, custom, usage_total, usage_used, usage_resets_in_seconds, asset_type)
				   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
				   ON CONFLICT (uid) DO UPDATE SET
				   read = $4, write = $5, delete = $6, share = $7, custom = $8, usage_total = $9, usage_used = $10, usage_resets_in_seconds = $11, asset_type = $12`
		_, err := tx.Exec(context.Background(), query, right.UID, req.Entity, right.Resource, right.Read, right.Write, right.Delete, right.Share, right.Custom, right.UsageTotal, right.UsageUsed, right.UsageResetsInSeconds, right.AssetType)
		if err != nil {
			return &GiveRightsResponse{Valid: false, Error: err.Error()}, err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return &GiveRightsResponse{Valid: false, Error: err.Error()}, err
	}

	return &GiveRightsResponse{
		Entity: req.Entity,
		Valid:  true,
	}, nil
}

func (dal *AuthorizationDAL) RevokeRights(req *RevokeRightsRequest) (*RevokeRightsResponse, error) {
	query := "UPDATE rights SET active = false, deleted_at = current_epoch() WHERE entity = $1 AND uid = ANY($2)"
	_, err := dal.db.Exec(context.Background(), query, req.Entity, req.Rights)
	if err != nil {
		return &RevokeRightsResponse{Valid: false, Error: err.Error()}, err
	}
	return &RevokeRightsResponse{
		Entity: req.Entity,
		Valid:  true,
	}, nil
}

func (dal *AuthorizationDAL) GetRights(req *GetRightsRequest) (*GetRightsResponse, error) {
	query := `SELECT uid, entity, resource, read, write, delete, share, custom, usage_total, usage_used, usage_resets_in_seconds, asset_type, created_at
			   FROM rights
			   WHERE entity = $1 AND active = true`

	rows, err := dal.db.Query(context.Background(), query, req.Entity)
	if err != nil {
		return &GetRightsResponse{Valid: false, Error: err.Error()}, err
	}
	defer rows.Close()

	rights := make(map[string]Right)
	for rows.Next() {
		var r Right
		if err := rows.Scan(&r.UID, &r.Entity, &r.Resource, &r.Read, &r.Write, &r.Delete, &r.Share, &r.Custom, &r.UsageTotal, &r.UsageUsed, &r.UsageResetsInSeconds, &r.AssetType, &r.CreatedAt); err != nil {
			return &GetRightsResponse{Valid: false, Error: err.Error()}, err
		}
		rights[r.UID.String()] = r
	}

	return &GetRightsResponse{
		Entity: req.Entity,
		Rights: rights,
		Valid:  true,
	}, nil
}
