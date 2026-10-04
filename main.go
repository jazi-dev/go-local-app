package main

// Import packages for SQL, HTML templates, logging, HTTP, and environment variables.
import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// student represents one row read from the MySQL students table.
type student struct {
	ID     int
	Name   string
	Course string
}

// page is the HTML template rendered in the browser with the student data.
var page = template.Must(template.New("students").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Course Students</title>
  <style>
    body { font-family: system-ui, sans-serif; max-width: 760px; margin: 3rem auto; padding: 0 1rem; color: #172033; }
    table { width: 100%; border-collapse: collapse; margin-top: 1.5rem; }
    th, td { padding: .75rem; border-bottom: 1px solid #d9dfeb; text-align: left; }
    th { background: #f3f6fb; }
    .empty { padding: 1rem; background: #fff7dc; border-radius: .5rem; }
  </style>
</head>
<body>
  <h1>Course Students</h1>
  <p>This page is served by Go and reads its data from MySQL.</p>
  {{if .}}
  <table>
    <thead><tr><th>ID</th><th>Name</th><th>Course</th></tr></thead>
    <tbody>
      {{range .}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Course}}</td></tr>{{end}}
    </tbody>
  </table>
  {{else}}
  <p class="empty">No students yet. Add a row to MySQL and refresh this page.</p>
  {{end}}
</body>
</html>`))

// env reads a configuration value or returns its default when it is not set.
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// main connects to MySQL, configures the routes, and starts the web server.
func main() {
	// Build the MySQL connection string from environment variables.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		env("DB_USER", "courseuser"),
		env("DB_PASSWORD", "coursepass"),
		env("DB_HOST", "127.0.0.1"),
		env("DB_PORT", "3306"),
		env("DB_NAME", "course"),
	)

	// Open and verify the application connection to MySQL.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("cannot connect to MySQL: %v", err)
	}

	// Create the students table automatically if it does not already exist.
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		course VARCHAR(100) NOT NULL
	)`)
	if err != nil {
		log.Fatalf("cannot create students table: %v", err)
	}

	// Serve the home page by reading all students and rendering the HTML template.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, course FROM students ORDER BY id")
		if err != nil {
			http.Error(w, "database query failed", http.StatusInternalServerError)
			log.Printf("query failed: %v", err)
			return
		}
		defer rows.Close()

		var students []student
		for rows.Next() {
			var s student
			if err := rows.Scan(&s.ID, &s.Name, &s.Course); err != nil {
				http.Error(w, "database result failed", http.StatusInternalServerError)
				return
			}
			students = append(students, s)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "database result failed", http.StatusInternalServerError)
			log.Printf("reading rows failed: %v", err)
			return
		}

		if err := page.Execute(w, students); err != nil {
			log.Printf("template failed: %v", err)
		}
	})

	// Report whether the application can still communicate with MySQL.
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, `{"status":"unhealthy"}`, http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// Listen on every network interface so the app also works inside Docker.
	port := env("PORT", "8080")
	log.Printf("server running at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
