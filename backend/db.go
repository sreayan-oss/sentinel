package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import driver anonymously
)

var db *sql.DB

// InitDB creates the table if it doesn't exist
func InitDB() {
	var err error
	// Open connection to a file named 'metrics.db'
	db, err = sql.Open("sqlite3", "./metrics.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	// PERFORMANCE HACK 1: WAL Mode (Write-Ahead Logging)
	// This allows readers and writers to work simultaneously without blocking.
	// Essential for high-concurrency monitoring.
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		log.Fatalf("Failed to set WAL mode: %v", err)
	}

	// Create Table
	// We index 'timestamp' because we will always query by time range later.
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS cpu_usage (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		agent_id TEXT,
		cpu_percent REAL,
		timestamp INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON cpu_usage(timestamp);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println("Database initialized (metrics.db)")
}

// RecordMetric inserts a new data point
// We use a "Prepared Statement" here implicitly by passing args.
func RecordMetric(agentId string, cpuUsage float64, timestamp int64) {
	stmt := `INSERT INTO cpu_usage(agent_id, cpu_percent, timestamp) VALUES (?, ?, ?)`
	
	_, err := db.Exec(stmt, agentId, cpuUsage, timestamp)
	if err != nil {
		log.Printf("Error inserting metric: %v", err)
	}
}
type Metric struct {
	ID        int64   `json:"id"`
	AgentID   string  `json:"agent_id"`
	CPU       float64 `json:"cpu"`
	Timestamp int64   `json:"timestamp"`
}

func GetRecentMetrics() ([]Metric, error) {
	// Query the last 20 data points, ordered by time
	rows, err := db.Query(`
		SELECT id, agent_id, cpu_percent, timestamp 
		FROM cpu_usage 
		ORDER BY timestamp DESC 
		LIMIT 20
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []Metric
	for rows.Next() {
		var m Metric
		if err := rows.Scan(&m.ID, &m.AgentID, &m.CPU, &m.Timestamp); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}