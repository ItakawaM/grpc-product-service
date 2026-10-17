CREATE TABLE IF NOT EXISTS products (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price_cents bigint NOT NULL CHECK (price_cents > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    version integer NOT NULL DEFAULT 1
);