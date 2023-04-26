CREATE TABLE "wallets" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "asset" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "wallet_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_wallet_id" bigint NOT NULL,
  "to_wallet_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
    "id" bigserial PRIMARY KEY,
    "current_attempt_id" TEXT,
    "owner" varchar NOT NULL,
    "is_active" boolean NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "attempts" (
    "id" TEXT NOT NULL,
    "session_id" bigint NOT NULL,
    "target_number" smallint NOT NULL,
    "num_of_dice_throw" smallint NOT NULL,
    "first_dice_throw_value" smallint NOT NULL,
    "second_dice_throw_value" smallint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    PRIMARY KEY (id)
);



CREATE INDEX ON "wallets" ("owner");

CREATE INDEX ON "entries" ("wallet_id");

CREATE INDEX ON "transfers" ("from_wallet_id");

CREATE INDEX ON "transfers" ("to_wallet_id");

CREATE INDEX ON "transfers" ("from_wallet_id", "to_wallet_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "attempts" ADD FOREIGN KEY ("session_id") REFERENCES "sessions" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "sessions" ADD FOREIGN KEY ("current_attempt_id") REFERENCES "attempts" ("id");

ALTER TABLE "sessions" ALTER COLUMN is_active SET DEFAULT TRUE;

ALTER TABLE "attempts" ALTER COLUMN num_of_dice_throw SET DEFAULT 0;

ALTER TABLE "attempts" ALTER COLUMN first_dice_throw_value SET DEFAULT 0;

ALTER TABLE "attempts" ALTER COLUMN second_dice_throw_value SET DEFAULT 0;

ALTER TABLE "wallets" ADD CONSTRAINT "owner_asset_key" UNIQUE ("owner", "asset");

ALTER TABLE "wallets" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

INSERT INTO users (username) VALUES('holding_account');
INSERT INTO wallets (owner, balance, asset) VALUES('holding_account',922337203685477,'sats');
