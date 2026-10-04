# Go + MySQL — Run Locally

This app serves a browser page at <http://localhost:8080> and displays student records stored in a local MySQL database.

## Prerequisites

- Go 1.23 or newer
- MySQL Server and the MySQL client

Verify both tools:

```bash
go version
mysql --version
```

## 1. Start MySQL

On Ubuntu:

```bash
sudo systemctl start mysql
sudo systemctl status mysql
```

## 2. Create the database and user

Open MySQL as an administrator:

```bash
sudo mysql
```

Run these SQL statements:

```sql
CREATE DATABASE IF NOT EXISTS course;
CREATE USER IF NOT EXISTS 'courseuser'@'localhost' IDENTIFIED BY 'coursepass';
GRANT ALL PRIVILEGES ON course.* TO 'courseuser'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

The Go app creates the `students` table automatically when it starts.

## 3. Download the Go dependency

From this directory:

```bash
go mod tidy
```

## 4. Run the app

```bash
go run .
```

Open <http://localhost:8080>. The page will initially show that no students exist.

## 5. Add data to MySQL

Open a second terminal:

```bash
mysql -h 127.0.0.1 -u courseuser -p course
```

Enter `coursepass` when prompted, then run:

```sql
INSERT INTO students (name, course)
VALUES ('Alice', 'DevOps'), ('Bob', 'Docker');

SELECT * FROM students;
EXIT;
```

Refresh <http://localhost:8080> to see the records.

## Useful checks

```bash
curl http://localhost:8080/health
```

To use different database settings, set these environment variables before `go run .`:

```bash
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_NAME=course
export DB_USER=courseuser
export DB_PASSWORD=coursepass
export PORT=8080
```

Stop the app with `Ctrl+C`.
# go-local-app
