BEGIN;

-- Tables
DROP TABLE IF EXISTS `app_versions`;
DROP TABLE IF EXISTS `configs`;
DROP TABLE IF EXISTS `families`;
DROP TABLE IF EXISTS `folder_file_uploads`;
DROP TABLE IF EXISTS `folder_files`;
DROP TABLE IF EXISTS `folder_invitees`;
DROP TABLE IF EXISTS `folders`;
DROP TABLE IF EXISTS `im_group_user`;
DROP TABLE IF EXISTS `im_groups`;
DROP TABLE IF EXISTS `im_messages`;
DROP TABLE IF EXISTS `locker_note_invitees`;
DROP TABLE IF EXISTS `locker_note_versions`;
DROP TABLE IF EXISTS `locker_notes`;
DROP TABLE IF EXISTS `system_messages`;
DROP TABLE IF EXISTS `traffics`;
DROP TABLE IF EXISTS `users`;

-- Types
drop type if exists app_version_os;
drop type if exists user_roles;
drop type if exists upload_file_type;

-- Sequences


COMMIT;