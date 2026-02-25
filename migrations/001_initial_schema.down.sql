-- 001_initial_schema.down.sql
-- Drops all tables created in the baseline migration.
-- Order matters: tables with potential FK references dropped first.

BEGIN;

-- Sticker tracking
DROP TABLE IF EXISTS stricker_count;
DROP TABLE IF EXISTS sticker_s;
DROP TABLE IF EXISTS sticker_log;

-- Blog
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;

-- Timers
DROP TABLE IF EXISTS time_with_out;

-- Vector quiz
DROP TABLE IF EXISTS vector_scores;
DROP TABLE IF EXISTS answer;
DROP TABLE IF EXISTS question;
DROP TABLE IF EXISTS question_types;

-- Free content
DROP TABLE IF EXISTS free_maw;

-- Openings / media
DROP TABLE IF EXISTS stunts;
DROP TABLE IF EXISTS grobs;
DROP TABLE IF EXISTS anime_openings;

-- Achievements
DROP TABLE IF EXISTS achlist;
DROP TABLE IF EXISTS achieves;
DROP TABLE IF EXISTS achievegroups;

-- Identity
DROP TABLE IF EXISTS metatron;
DROP TABLE IF EXISTS users;

COMMIT;
