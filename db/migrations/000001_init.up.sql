-- PostgreSQL

BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

drop type if exists app_version_os;
create type app_version_os as enum ('ios', 'android');

drop type if exists user_roles;
create type user_roles as ENUM ('superadmin', 'admin', 'manager', 'member');

drop type if exists upload_file_type;
create type upload_file_type as ENUM ('image','video','pdf','doc','excel','others');


DROP TABLE IF EXISTS app_versions;
CREATE TABLE app_versions (
  "id" serial primary key,
  "os" app_version_os NOT NULL,
  "version" varchar(255)  NOT NULL,
  "title" varchar(255)  NOT NULL,
  "content" jsonb NOT NULL default '{}'::jsonb,
  "link" varchar(255) NOT NULL,
  "download_url" varchar(255) NOT NULL,
  "effective_on" timestamp NOT NULL,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);
create unique index idx_app_versions_os_version on app_versions("os", "version");
create index idx_app_verions_effective_on on app_versions("effective_on");


DROP TABLE IF EXISTS configs;
CREATE TABLE configs (
  "id" serial primary key,
  "config_key" varchar(50) unique NOT NULL,
  "group" varchar(30) NOT NULL DEFAULT 'basic',
  "title" varchar(255) NOT NULL,
  "remark" varchar(255) not null default '',
  "config_type" varchar(255) NOT NULL DEFAULT 'string',
  "config_value" varchar(255) not null default '' ,
  "config_params" jsonb not null default '{}'::jsonb,
  "status" boolean NOT NULL DEFAULT true,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);


DROP TABLE IF EXISTS families;
CREATE TABLE families (
  "id" serial primary key,
  "title" varchar(255)  NOT NULL,
  "status" boolean not null default true,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);

/* Create a default user, mobile = 1234567890, password = 123456 */
INSERT INTO families ("title", "status", "created_at", "updated_at") VALUES
('Sample Family', true, '2020-01-01T00:00:00Z',	'2020-01-01T00:00:00Z');



DROP TABLE IF EXISTS folder_file_uploads;
CREATE TABLE folder_file_uploads (
  "id" serial primary key,
  "user_id" int not null,
  "disk" varchar(20)  NOT NULL,
  "file_name" varchar(255)  NOT NULL,
  "file_type" upload_file_type  NOT NULL,
  "file_path" varchar(511)  NOT NULL,
  "file_md5" varchar(255) not null,
  "taken_on" timestamp NULL DEFAULT NULL,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);

DROP TABLE IF EXISTS folder_files;
CREATE TABLE folder_files (
  "id" serial primary key,
  "peer_id" int  DEFAULT NULL,
  "owner_id" int NOT NULL,
  "folder_id" int NOT NULL,
  "disk" varchar(255)  NOT NULL,
  "file_name" varchar(255)  NOT NULL,
  "file_type" upload_file_type  NOT NULL,
  "original_file_path" varchar(255)  DEFAULT '',
  "file_path" varchar(255)  NOT NULL,
  "thumbnail_path" varchar(255)  NOT NULL,
  "hw_original_file_path" varchar(255)  DEFAULT NULL,
  "google_original_file_path" varchar(255)  DEFAULT NULL,
  "google_file_path" varchar(255)  DEFAULT NULL,
  "google_thumbnail_path" varchar(255)  DEFAULT NULL,
  "google_drive_file_id" varchar(255)  DEFAULT NULL,
  "remark" varchar(255)  NOT NULL DEFAULT '',
  "metadata" jsonb not null default '{}'::jsonb ,
  "has_exif" boolean NOT NULL DEFAULT false,
  "is_private" boolean NOT NULL DEFAULT false,
  "views" int NOT NULL DEFAULT 0,
  "synced_at" timestamp NULL DEFAULT NULL,
  "is_downloading" boolean NOT NULL DEFAULT false,
  "taken_on" timestamp(6) unique NOT NULL,
  "created_at" timestamp(6) NULL DEFAULT NULL,
  "updated_at" timestamp(6) NULL DEFAULT NULL
);


DROP TABLE IF EXISTS folder_invitees;
CREATE TABLE folder_invitees (
  "id" serial primary key,
  "folder_id" int  NOT NULL,
  "invitee_id" int NOT NULL,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);


DROP TABLE IF EXISTS folders;
CREATE TABLE folders (
  "id" serial primary key,
  "owner_id" int NOT NULL,
  "parent_id" int NOT NULL,
  "title" varchar(255)  NOT NULL,
  "cover" varchar(255)  NOT NULL,
  "type" varchar(255)  DEFAULT NULL,
  "metadata" jsonb not null default '{}'::jsonb,
  "is_default" boolean NOT NULL DEFAULT false,
  "is_private" boolean NOT NULL DEFAULT false,
  "total_files" int NOT NULL DEFAULT -1,
  "earliest_taken_on" timestamp NULL DEFAULT NULL,
  "latest_taken_on" timestamp NULL DEFAULT NULL,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);


DROP TABLE IF EXISTS locker_note_invitees;
CREATE TABLE locker_note_invitees (
  "id" serial primary key,
  "note_id" int NOT NULL,
  "invitee_id" int NOT NULL,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);


DROP TABLE IF EXISTS locker_note_versions;
CREATE TABLE locker_note_versions (
  "id" serial primary key,
  "note_id" int NOT NULL,
  "title" varchar(255)  NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);


DROP TABLE IF EXISTS locker_notes;
CREATE TABLE locker_notes (
  "id" serial primary key,
  "owner_id" int NOT NULL,
  "title" varchar(255)  NOT NULL,
  "content" text  NOT NULL,
  "deleted_at" timestamp NULL DEFAULT NULL,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);


DROP TABLE IF EXISTS system_messages;
CREATE TABLE system_messages (
  "id" serial primary key,
  "user_id" int NOT NULL,
  "title" varchar(255)  NOT NULL,
  "content" varchar(255)  NOT NULL,
  "action" varchar(255)  NOT NULL,
  "seen" boolean not null default false,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);


DROP TABLE IF EXISTS traffics;
CREATE TABLE traffics (
  "id" serial primary key,
  "ip" varchar(255)  NOT NULL,
  "user_id" int NOT NULL,
  "requester" varchar(50)  NOT NULL,
  "request_uri" varchar(255)  NOT NULL,
  "request_body" varchar(511)  NOT NULL,
  "response_code" varchar(10)  NOT NULL,
  "error_message" varchar(255)  NOT NULL,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);


DROP TABLE IF EXISTS uploads;
CREATE TABLE uploads (
  "id" serial primary key,
  "user_id" int NOT NULL,
  "disk" varchar(20)  NOT NULL,
  "file_name" varchar(255)  NOT NULL,
  "file_type" upload_file_type  NOT NULL,
  "file_path" varchar(511)  NOT NULL,
  "taken_on" timestamp NULL DEFAULT NULL,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE IF EXISTS users;
CREATE TABLE users (
  "id" serial primary key,
  "family_id" int default 0,
  "uuid" uuid NOT NULL DEFAULT public.uuid_generate_v4 (),
  "role" user_roles default 'member'  NOT NULL,
  "name" varchar(255)  NOT NULL,
  "email" varchar(255)  NOT NULL,
  "email_verified_at" timestamp NULL DEFAULT NULL,
  "mobile" varchar(255) unique NOT NULL,
  "avatar" varchar(255)  NOT NULL,
  "password" varchar(255)  NOT NULL,
  "locker_passcode" varchar(255)  NOT NULL,
  "refresh_token" varchar(100)  DEFAULT NULL,
  "device_token" varchar(200)  DEFAULT NULL,
  "notifications" jsonb not null default '{"total":0}'::jsonb,
  "status" boolean not null default true,
  "created_at" timestamp NULL DEFAULT NULL,
  "updated_at" timestamp NULL DEFAULT NULL
);

/* Create a default user, mobile = 1234567890, password = 123456 */
INSERT INTO users ("uuid", "role", "name", "email", "email_verified_at", "mobile", "avatar", "password", "locker_passcode", "refresh_token", "device_token", "status", "created_at", "updated_at") VALUES
('414ca89b-1b80-49c3-97db-1b94ae3c9738',	'admin',	'Logan',	'logan.dai@innovay.dev',	NULL,	'1234567890',	'',	'$2y$10$km/tZvmVfilc67p0FmZ1z.0hPYLKmuQCohf45afkMQ9ny5xChLJaq',	'123456',	'',	'', true,	'2020-01-01 00:00:00',	'2020-01-01 00:00:00');

COMMIT;