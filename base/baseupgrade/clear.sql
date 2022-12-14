DROP INDEX IF EXISTS _sys_announce_update_time_idx;
DROP INDEX IF EXISTS _sys_announce_type_idx;
DROP INDEX IF EXISTS _sys_announce_title_idx;
DROP INDEX IF EXISTS _sys_announce_status_idx;
DROP INDEX IF EXISTS _sys_announce_marked_idx;
ALTER TABLE IF EXISTS _sys_version_object ALTER COLUMN tid DROP DEFAULT;
ALTER TABLE IF EXISTS _sys_announce ALTER COLUMN tid DROP DEFAULT;
DROP SEQUENCE IF EXISTS _sys_version_object_tid_seq;
DROP TABLE IF EXISTS _sys_version_object;
DROP TABLE IF EXISTS _sys_object;
DROP TABLE IF EXISTS _sys_config;
DROP SEQUENCE IF EXISTS _sys_announce_tid_seq;
DROP TABLE IF EXISTS _sys_announce;
