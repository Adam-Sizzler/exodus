ALTER TABLE users
  DROP COLUMN IF EXISTS sub_last_user_agent,
  DROP COLUMN IF EXISTS sub_last_opened_at;
