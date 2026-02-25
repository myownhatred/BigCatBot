-- 002_schema_improvements.down.sql

BEGIN;

DROP INDEX IF EXISTS idx_time_with_out_chat_id;
DROP INDEX IF EXISTS idx_free_maw_typ;
DROP INDEX IF EXISTS idx_answer_questionid;
DROP INDEX IF EXISTS idx_question_typeid;
DROP INDEX IF EXISTS idx_achlist_uid;
ALTER TABLE achlist DROP CONSTRAINT IF EXISTS achlist_uid_aid_key;
ALTER TABLE metatron DROP CONSTRAINT IF EXISTS metatron_chat_id_key;
DROP INDEX IF EXISTS idx_users_username;

COMMIT;
