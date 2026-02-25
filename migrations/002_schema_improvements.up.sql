-- 002_schema_improvements.up.sql
-- Adds missing indexes, constraints, and structural improvements.
-- Each section is independent and can be partially applied if needed.

BEGIN;

-- ============================================================================
-- users: index on username (used in WHERE Username = $1)
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- ============================================================================
-- metatron: UNIQUE on chat_id (code does SELECT COUNT WHERE chat_id = $1
-- before insert, but without constraint a race condition can cause dupes)
-- ============================================================================
ALTER TABLE metatron ADD CONSTRAINT metatron_chat_id_key UNIQUE (chat_id);

-- ============================================================================
-- achlist: UNIQUE(uid, aid) prevents duplicate achievements per user.
-- Code checks before insert but has no DB-level protection.
-- ============================================================================
ALTER TABLE achlist ADD CONSTRAINT achlist_uid_aid_key UNIQUE (uid, aid);

-- Index on uid for the common query: SELECT ... FROM achlist WHERE uid = $1
CREATE INDEX IF NOT EXISTS idx_achlist_uid ON achlist(uid);

-- ============================================================================
-- question: index on typeid for WHERE typeid=$1 ORDER BY random() LIMIT 1
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_question_typeid ON question(typeid);

-- ============================================================================
-- answer: index on questionid for WHERE questionid=$1
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_answer_questionid ON answer(questionid);

-- ============================================================================
-- free_maw: index on typ for WHERE typ=$1
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_free_maw_typ ON free_maw(typ);

-- ============================================================================
-- time_with_out: index on chat_id for WHERE chat_id=$1
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_time_with_out_chat_id ON time_with_out(chat_id);

COMMIT;
