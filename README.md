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

### Database Migrations
Database migration is managed with golang-migrate.
migrate can be downloaded directely and placed in the project root directory.

- https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

$ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
$ echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
$ apt-get update
$ apt-get install -y migrate

./migrate -source "file://db/migrations/" -database "postgres://<user>:<password>@pgsql/inno?sslmode=disable&search_path=famapp_local_testing" up




