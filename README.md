#### Assumptions:
1. Number of media files for any second will not exceed 1 million.


#### Design points:
1. FolderFiles timestamp has a microsecond preceision (6 decimal points)
2. 


### Type and scope definitions:
- repository: an agent responsible for interacting with datasources. Datasources can be local DB, remote DB, or other API etc.
- service: an agent reponsible for completing some predefined tasks
- util: an wrapper for some common unit functions, specific for small functions


### Commands
1. To run tests:
- go test github.com/innovay-software/famapp-main/tests/...

2. To run oapi-codegen to generate code:
- go generate ./...
- reference: https://stackoverflow.com/questions/59794748/go-generate-only-scans-main-go


### Command Queries:
```sql
update folder_files 
set google_original_file_path = '', google_file_path='', google_thumbnail_path='', google_drive_file_id = ''
where id >= 7702
```
