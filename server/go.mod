module fileservice

go 1.14

// replace github.com/wup364/pakku => ../../pakku
// replace github.com/wup364/filestorage/opensdk => ../../filestorage/opensdk

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-sql-driver/mysql v1.6.0
	github.com/mattn/go-sqlite3 v1.14.11
	github.com/wup364/filestorage/opensdk v0.0.0-20220626091002-ba8b918aee76
	github.com/wup364/pakku v0.0.1
)
