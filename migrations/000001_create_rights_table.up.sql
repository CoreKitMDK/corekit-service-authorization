CREATE TABLE IF NOT EXISTS rights (
    uid UUID PRIMARY KEY,
    entity UUID NOT NULL,
    resource UUID NOT NULL,
    read BOOLEAN NOT NULL DEFAULT false,
    write BOOLEAN NOT NULL DEFAULT false,
    delete BOOLEAN NOT NULL DEFAULT false,
    share BOOLEAN NOT NULL DEFAULT false,
    custom JSONB,
    usage_total INTEGER NOT NULL DEFAULT 0,
    usage_used INTEGER NOT NULL DEFAULT 0,
    usage_resets_in_seconds BIGINT NOT NULL DEFAULT 0,
    asset_type TEXT,
    active BOOLEAN NOT NULL DEFAULT false,
    created BIGINT NOT NULL
);
