module fileservice

go 1.14

// http://git.xxxxx.com:xx/golang/pakku.git
replace pakku => ../../pakku

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-sql-driver/mysql v1.6.0
	github.com/mattn/go-sqlite3 v1.14.9
	pakku v0.0.0-00010101000000-000000000000
)
