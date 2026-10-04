// Name this Go module so its packages have a unique import path.
module course.local/go-mysql-app

// Require Go 1.23 or a compatible newer version.
go 1.23

// Provide the database driver used to connect the app to MySQL.
require github.com/go-sql-driver/mysql v1.8.1

// Support secure authentication required internally by the MySQL driver.
require filippo.io/edwards25519 v1.1.0 // indirect
