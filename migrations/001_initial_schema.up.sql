-- 001_initial_schema.up.sql
-- Baseline migration: captures the existing PostgreSQL schema as of 2026-02-14.
-- All tables here already exist in production; this migration serves as
-- the single source of truth for the database structure going forward.
--
-- NOTE: If applying to an existing database, mark as already applied:
--   migrate -path migrations -database "$DB_URL" force 1
--
-- Schema reconstructed from production database dump + constraint query.

BEGIN;

-- ============================================================================
-- Users & Identity
-- ============================================================================

-- Telegram users cache.
-- id is Telegram UserID (not auto-generated).
-- chat_role stores per-chat roles as JSON: {"chatID": "role", ...}
-- Column order matters for SELECT * in GetAllUsers:
--   id, first_name, last_name, username, chat_role
CREATE TABLE IF NOT EXISTS users (
    id         BIGINT  PRIMARY KEY,
    first_name VARCHAR,
    last_name  VARCHAR,
    username   VARCHAR,
    chat_role  JSON
);

-- Known Telegram chats the bot operates in.
CREATE TABLE IF NOT EXISTS metatron (
    id        BIGSERIAL PRIMARY KEY,
    chat_id   BIGINT    NOT NULL,
    chat_name VARCHAR   NOT NULL
);

-- ============================================================================
-- Achievement system
-- ============================================================================

-- Achievement categories/groups.
CREATE TABLE IF NOT EXISTS achievegroups (
    id        SERIAL  PRIMARY KEY,
    groupname VARCHAR
);

-- Achievement definitions.
CREATE TABLE IF NOT EXISTS achieves (
    id          SERIAL  PRIMARY KEY,
    groupid     INT     REFERENCES achievegroups(id),
    name        VARCHAR,
    rank        INT,
    description VARCHAR
);

-- User<->Achievement join table (which user earned which achievement).
CREATE TABLE IF NOT EXISTS achlist (
    id     SERIAL    PRIMARY KEY,
    uid    BIGINT    REFERENCES users(id),
    aid    INT       REFERENCES achieves(id),
    date   TIMESTAMP,
    chat   VARCHAR,
    chatid BIGINT
);

-- ============================================================================
-- Anime / Music openings (GORM-managed tables)
-- ============================================================================

-- Anime openings. GORM adds created_at/updated_at/deleted_at automatically.
CREATE TABLE IF NOT EXISTS anime_openings (
    id          SERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_at  TIMESTAMP WITH TIME ZONE,
    deleted_at  TIMESTAMP WITH TIME ZONE,
    description TEXT,
    link        TEXT
);

CREATE INDEX IF NOT EXISTS idx_anime_openings_deleted_at
    ON anime_openings(deleted_at);

-- Grob music table (same structure as anime_openings).
CREATE TABLE IF NOT EXISTS grobs (
    id          SERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_at  TIMESTAMP WITH TIME ZONE,
    deleted_at  TIMESTAMP WITH TIME ZONE,
    description TEXT,
    link        TEXT
);

CREATE INDEX IF NOT EXISTS idx_grobs_deleted_at
    ON grobs(deleted_at);

-- Stunts (same GORM structure).
CREATE TABLE IF NOT EXISTS stunts (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_at  TIMESTAMP WITH TIME ZONE,
    deleted_at  TIMESTAMP WITH TIME ZONE,
    description TEXT,
    link        TEXT
);

CREATE INDEX IF NOT EXISTS idx_stunts_deleted_at
    ON stunts(deleted_at);

-- ============================================================================
-- Free content (memes, media links)
-- ============================================================================

CREATE TABLE IF NOT EXISTS free_maw (
    id          SERIAL PRIMARY KEY,
    typ         VARCHAR,
    description TEXT,
    link        TEXT
);

-- ============================================================================
-- Vector quiz game
-- ============================================================================

-- Quiz categories.
CREATE TABLE IF NOT EXISTS question_types (
    id        SERIAL  PRIMARY KEY,
    type_name VARCHAR,
    protected BOOLEAN
);

-- Quiz questions.
CREATE TABLE IF NOT EXISTS question (
    id              SERIAL  PRIMARY KEY,
    typeid          INT     REFERENCES question_types(id) ON DELETE CASCADE,
    pic_link        VARCHAR,
    question_string VARCHAR,
    userid          BIGINT
);

-- Accepted answers for a question (multiple correct answers possible).
CREATE TABLE IF NOT EXISTS answer (
    id          SERIAL  PRIMARY KEY,
    questionid  INT     REFERENCES question(id) ON DELETE CASCADE,
    answer_text VARCHAR
);

-- Cumulative player scores per quiz type.
-- Primary key doubles as UNIQUE constraint for ON CONFLICT upsert.
CREATE TABLE IF NOT EXISTS vector_scores (
    uid         BIGINT  NOT NULL,
    vector_type INT     NOT NULL,
    score       INT     NOT NULL DEFAULT 0,
    PRIMARY KEY (uid, vector_type)
);

-- ============================================================================
-- Timers ("time without X" tracker)
-- ============================================================================

CREATE TABLE IF NOT EXISTS time_with_out (
    id      SERIAL    PRIMARY KEY,
    name    TEXT      NOT NULL,
    "time"  TIMESTAMP NOT NULL,
    chat_id BIGINT    NOT NULL
);

-- ============================================================================
-- Blog / Posts system (GORM-managed)
-- ============================================================================

CREATE TABLE IF NOT EXISTS posts (
    id          SERIAL PRIMARY KEY,
    content     TEXT,
    author      TEXT   NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_at  TIMESTAMP WITH TIME ZONE,
    deleted_at  TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS comments (
    id          SERIAL  PRIMARY KEY,
    content     TEXT,
    author      VARCHAR,
    post_id     INT,
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_at  TIMESTAMP WITH TIME ZONE,
    deleted_at  TIMESTAMP WITH TIME ZONE
);

-- ============================================================================
-- Sticker tracking (legacy)
-- ============================================================================

CREATE TABLE IF NOT EXISTS sticker_log (
    id        SERIAL  PRIMARY KEY,
    chat_id   BIGINT  NOT NULL DEFAULT 0,
    chat_name VARCHAR NOT NULL DEFAULT '0',
    "user"    VARCHAR NOT NULL DEFAULT '0',
    sticker   VARCHAR NOT NULL DEFAULT '0'
);

CREATE TABLE IF NOT EXISTS sticker_s (
    id          SERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_at  TIMESTAMP WITH TIME ZONE,
    deleted_at  TIMESTAMP WITH TIME ZONE,
    "user"      TEXT,
    sticker     TEXT,
    chat_id     BIGINT,
    chat_name   TEXT
);

CREATE TABLE IF NOT EXISTS stricker_count (
    id      SERIAL  PRIMARY KEY,
    "user"  VARCHAR NOT NULL DEFAULT '0',
    count   INT              DEFAULT 0
);

COMMIT;
