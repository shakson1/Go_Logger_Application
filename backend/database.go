package main

import (
	"database/sql"
	"strings"
	"time"

	"math/rand"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("sqlite3", "./logs.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func createTables(db *sql.DB) error {
	// Create logs table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			level TEXT NOT NULL,
			rule TEXT NOT NULL,
			source_ip TEXT NOT NULL,
			destination_ip TEXT NOT NULL,
			event TEXT NOT NULL,
			description TEXT NOT NULL,
			urgency INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create indexes for better performance
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_logs_rule ON logs(rule)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_logs_source_ip ON logs(source_ip)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_logs_event ON logs(event)`)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) InsertLog(log LogEntry) error {
	_, err := d.db.Exec(`
		INSERT INTO logs (timestamp, level, rule, source_ip, destination_ip, event, description, urgency)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, log.Timestamp, log.Level, log.Rule, log.SourceIP, log.DestinationIP, log.Event, log.Description, log.Urgency)
	return err
}

func (d *Database) GetLogs(limit int) ([]LogEntry, error) {
	rows, err := d.db.Query(`
		SELECT timestamp, level, rule, source_ip, destination_ip, event, description, urgency
		FROM logs
		ORDER BY timestamp DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		err := rows.Scan(&log.Timestamp, &log.Level, &log.Rule, &log.SourceIP, &log.DestinationIP, &log.Event, &log.Description, &log.Urgency)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (d *Database) SearchLogs(ip, event string, limit int) ([]LogEntry, error) {
	query := `
		SELECT timestamp, level, rule, source_ip, destination_ip, event, description, urgency
		FROM logs
		WHERE 1=1
	`
	args := []interface{}{}

	if ip != "" {
		query += ` AND (source_ip LIKE ? OR destination_ip LIKE ?)`
		args = append(args, "%"+ip+"%", "%"+ip+"%")
	}

	if event != "" {
		query += ` AND event LIKE ?`
		args = append(args, "%"+event+"%")
	}

	query += ` ORDER BY timestamp DESC LIMIT ?`
	args = append(args, limit)

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		err := rows.Scan(&log.Timestamp, &log.Level, &log.Rule, &log.SourceIP, &log.DestinationIP, &log.Event, &log.Description, &log.Urgency)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (d *Database) GetLogsByEvent(event string, limit int) ([]LogEntry, error) {
	rows, err := d.db.Query(`
		SELECT timestamp, level, rule, source_ip, destination_ip, event, description, urgency
		FROM logs
		WHERE event = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`, event, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		err := rows.Scan(&log.Timestamp, &log.Level, &log.Rule, &log.SourceIP, &log.DestinationIP, &log.Event, &log.Description, &log.Urgency)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (d *Database) GetSummaryStats() (SummaryStats, error) {
	var stats SummaryStats

	// Count logs by category (access, network, threat, uba)
	accessCount := 0
	networkCount := 0
	threatCount := 0
	ubaCount := 0

	rows, err := d.db.Query(`
		SELECT rule FROM logs
	`)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var rule string
		err := rows.Scan(&rule)
		if err != nil {
			return stats, err
		}
		// Categorize based on rule name (simplified logic)
		switch {
		case strings.Contains(strings.ToLower(rule), "login") || strings.Contains(strings.ToLower(rule), "access"):
			accessCount++
		case strings.Contains(strings.ToLower(rule), "network") || strings.Contains(strings.ToLower(rule), "traffic"):
			networkCount++
		case strings.Contains(strings.ToLower(rule), "threat") || strings.Contains(strings.ToLower(rule), "malware"):
			threatCount++
		case strings.Contains(strings.ToLower(rule), "behavior") || strings.Contains(strings.ToLower(rule), "uba"):
			ubaCount++
		default:
			// Default to access for unknown rules
			accessCount++
		}
	}

	stats = SummaryStats{
		AccessNotables:  StatTile{Total: accessCount, Delta: 0},
		NetworkNotables: StatTile{Total: networkCount, Delta: 0},
		ThreatNotables:  StatTile{Total: threatCount, Delta: 0},
		UBANotables:     StatTile{Total: ubaCount, Delta: 0},
	}

	return stats, nil
}

func (d *Database) GetUrgencyData() (UrgencyData, error) {
	var data UrgencyData

	rows, err := d.db.Query(`
		SELECT urgency, COUNT(*) as count
		FROM logs
		WHERE timestamp >= datetime('now', '-24 hours')
		GROUP BY urgency
	`)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var urgency int
		var count int
		err := rows.Scan(&urgency, &count)
		if err != nil {
			return data, err
		}
		switch urgency {
		case 4: // critical
			data.Critical = count
		case 3: // high
			data.High = count
		case 2: // medium
			data.Medium = count
		case 1: // low
			data.Low = count
		}
	}

	return data, nil
}

func (d *Database) GetTimelineData() (TimelineData, error) {
	var data TimelineData

	// Generate labels for the last 24 hours
	labels := []string{}
	accessData := []int{}
	networkData := []int{}
	threatData := []int{}

	now := time.Now()
	for i := 23; i >= 0; i-- {
		hour := now.Add(-time.Duration(i) * time.Hour)
		labels = append(labels, hour.Format("15:04"))
		accessData = append(accessData, 0)
		networkData = append(networkData, 0)
		threatData = append(threatData, 0)
	}

	// Get actual data from database
	rows, err := d.db.Query(`
		SELECT 
			strftime('%H:%M', timestamp) as hour,
			rule,
			COUNT(*) as count
		FROM logs
		WHERE timestamp >= datetime('now', '-24 hours')
		GROUP BY strftime('%H:%M', timestamp), rule
		ORDER BY hour
	`)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var hour string
		var rule string
		var count int
		err := rows.Scan(&hour, &rule, &count)
		if err != nil {
			return data, err
		}

		// Find the index for this hour
		for i, label := range labels {
			if label == hour {
				// Categorize based on rule
				switch {
				case strings.Contains(strings.ToLower(rule), "login") || strings.Contains(strings.ToLower(rule), "access"):
					accessData[i] += count
				case strings.Contains(strings.ToLower(rule), "network") || strings.Contains(strings.ToLower(rule), "traffic"):
					networkData[i] += count
				case strings.Contains(strings.ToLower(rule), "threat") || strings.Contains(strings.ToLower(rule), "malware"):
					threatData[i] += count
				default:
					accessData[i] += count
				}
				break
			}
		}
	}

	data = TimelineData{
		Labels: labels,
		Series: []TimelineSeries{
			{Name: "Access", Data: accessData, Color: "#3B82F6"},
			{Name: "Network", Data: networkData, Color: "#10B981"},
			{Name: "Threat", Data: threatData, Color: "#EF4444"},
		},
	}

	return data, nil
}

func (d *Database) GetTopEvents() ([]TopEvent, error) {
	rows, err := d.db.Query(`
		SELECT event, COUNT(*) as count
		FROM logs
		GROUP BY event
		ORDER BY count DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []TopEvent
	for rows.Next() {
		var event TopEvent
		err := rows.Scan(&event.RuleName, &event.Count)
		if err != nil {
			return nil, err
		}
		// Generate mock sparkline data
		sparkline := []int{}
		for i := 0; i < 10; i++ {
			sparkline = append(sparkline, event.Count/10+rand.Intn(5))
		}
		event.Sparkline = sparkline
		event.Urgency = "medium" // Default urgency
		events = append(events, event)
	}
	return events, nil
}

func (d *Database) GetTopSources() ([]TopSource, error) {
	rows, err := d.db.Query(`
		SELECT source_ip, COUNT(*) as count
		FROM logs
		GROUP BY source_ip
		ORDER BY count DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []TopSource
	for rows.Next() {
		var source TopSource
		err := rows.Scan(&source.SourceIP, &source.Count)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}
	return sources, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
