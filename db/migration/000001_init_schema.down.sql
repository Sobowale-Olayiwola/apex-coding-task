ALTER TABLE IF EXISTS "entries" DROP CONSTRAINT IF EXISTS "entries_wallet_id_fkey";
ALTER TABLE IF EXISTS "transfers" DROP CONSTRAINT IF EXISTS "transfers_from_wallet_id_fkey";
ALTER TABLE IF EXISTS "transfers" DROP CONSTRAINT IF EXISTS "transfers_to_wallet_id_fkey";
ALTER TABLE IF EXISTS "attempts" DROP CONSTRAINT IF EXISTS "attempts_session_id_fkey";
ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS "sessions_owner_fkey";
ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS "sessions_current_attempt_id_fkey";
ALTER TABLE IF EXISTS "sessions" DROP CONSTRAINT IF EXISTS "sessions_current_attempt_id_fkey";
ALTER TABLE IF EXISTS "wallets" DROP CONSTRAINT IF EXISTS "wallets_owner_key";
ALTER TABLE IF EXISTS "wallets" DROP CONSTRAINT IF EXISTS "wallets_owner_fkey";
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS attempts;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS users;