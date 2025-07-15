package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// NotableEvent represents a security notable event
type NotableEvent struct {
	ID          string    `json:"id"`
	RuleName    string    `json:"ruleName"`
	Urgency     string    `json:"urgency"`  // critical, high, medium, low
	Category    string    `json:"category"` // access, network, threat, uba
	SourceIP    string    `json:"sourceIP"`
	Destination string    `json:"destination"`
	Count       int       `json:"count"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
}

// SummaryStats represents dashboard summary statistics
type SummaryStats struct {
	AccessNotables  StatTile `json:"accessNotables"`
	NetworkNotables StatTile `json:"networkNotables"`
	ThreatNotables  StatTile `json:"threatNotables"`
	UBANotables     StatTile `json:"ubaNotables"`
}

// StatTile represents a dashboard statistic tile
type StatTile struct {
	Total int `json:"total"`
	Delta int `json:"delta"`
}

// UrgencyData represents bar chart data for urgency levels
type UrgencyData struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
}

// TimelineData represents line chart time series data
type TimelineData struct {
	Labels []string         `json:"labels"`
	Series []TimelineSeries `json:"series"`
}

// TimelineSeries represents a data series for timeline chart
type TimelineSeries struct {
	Name  string `json:"name"`
	Data  []int  `json:"data"`
	Color string `json:"color"`
}

// TopEvent represents a top notable event for table display
type TopEvent struct {
	RuleName  string `json:"ruleName"`
	Sparkline []int  `json:"sparkline"`
	Count     int    `json:"count"`
	Urgency   string `json:"urgency"`
}

// TopSource represents a top event source for table display
type TopSource struct {
	SourceIP  string `json:"sourceIP"`
	Sparkline []int  `json:"sparkline"`
	Count     int    `json:"count"`
	Category  string `json:"category"`
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp     time.Time `json:"timestamp"`
	Level         string    `json:"level"`
	Rule          string    `json:"rule"`
	SourceIP      string    `json:"sourceIP"`
	DestinationIP string    `json:"destinationIP"`
	Event         string    `json:"event"`
	Description   string    `json:"description"`
	Urgency       int       `json:"urgency"`
}

// In-memory log store
var (
	logStore = struct {
		logs []LogEntry
		mu   sync.RWMutex
	}{logs: []LogEntry{}}
)

// Mock data
var mockEvents = []NotableEvent{
	{ID: "1", RuleName: "Suspicious Login Attempt", Urgency: "critical", Category: "access", SourceIP: "192.168.1.100", Count: 45, Timestamp: time.Now().Add(-2 * time.Hour)},
	{ID: "2", RuleName: "Data Exfiltration Detected", Urgency: "high", Category: "threat", SourceIP: "10.0.0.50", Count: 23, Timestamp: time.Now().Add(-1 * time.Hour)},
	{ID: "3", RuleName: "Unusual Network Traffic", Urgency: "medium", Category: "network", SourceIP: "172.16.0.25", Count: 67, Timestamp: time.Now().Add(-30 * time.Minute)},
	{ID: "4", RuleName: "Privilege Escalation", Urgency: "critical", Category: "access", SourceIP: "192.168.1.101", Count: 12, Timestamp: time.Now().Add(-15 * time.Minute)},
	{ID: "5", RuleName: "Malware Detection", Urgency: "high", Category: "threat", SourceIP: "10.0.0.51", Count: 34, Timestamp: time.Now().Add(-10 * time.Minute)},
	{ID: "6", RuleName: "Anomalous User Behavior", Urgency: "medium", Category: "uba", SourceIP: "172.16.0.26", Count: 89, Timestamp: time.Now().Add(-5 * time.Minute)},
	{ID: "7", RuleName: "Brute Force Attack", Urgency: "critical", Category: "access", SourceIP: "192.168.1.102", Count: 156, Timestamp: time.Now().Add(-2 * time.Minute)},
	{ID: "8", RuleName: "Data Breach Attempt", Urgency: "high", Category: "threat", SourceIP: "10.0.0.52", Count: 78, Timestamp: time.Now().Add(-1 * time.Minute)},
}

var startTime = time.Now()

// Helper function to convert urgency string to integer
func getUrgencyValue(urgency string) int {
	switch urgency {
	case "critical":
		return 4
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 2
	}
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.WriteHeader(http.StatusOK)
}

func summaryStatsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	// Use real logs if available, else fallback to mockEvents
	logStore.mu.RLock()
	logs := make([]LogEntry, len(logStore.logs))
	copy(logs, logStore.logs)
	logStore.mu.RUnlock()
	var source []LogEntry
	if len(logs) > 0 {
		source = logs
	} else {
		for _, e := range mockEvents {
			source = append(source, LogEntry{
				Timestamp:     e.Timestamp,
				Level:         "INFO",
				Rule:          e.RuleName,
				SourceIP:      e.SourceIP,
				DestinationIP: e.Destination,
				Event:         e.RuleName,
				Description:   e.Description,
				Urgency:       getUrgencyValue(e.Urgency),
			})
		}
	}
	accessCount := 0
	networkCount := 0
	threatCount := 0
	ubaCount := 0
	for _, log := range source {
		// Try to find category from mockEvents if possible
		cat := ""
		for _, me := range mockEvents {
			if me.RuleName == log.Rule {
				cat = me.Category
				break
			}
		}
		switch cat {
		case "access":
			accessCount++
		case "network":
			networkCount++
		case "threat":
			threatCount++
		case "uba":
			ubaCount++
		}
	}
	stats := SummaryStats{
		AccessNotables:  StatTile{Total: accessCount, Delta: 0},
		NetworkNotables: StatTile{Total: networkCount, Delta: 0},
		ThreatNotables:  StatTile{Total: threatCount, Delta: 0},
		UBANotables:     StatTile{Total: ubaCount, Delta: 0},
	}
	json.NewEncoder(w).Encode(stats)
}

func urgencyDataHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	logStore.mu.RLock()
	logs := make([]LogEntry, len(logStore.logs))
	copy(logs, logStore.logs)
	logStore.mu.RUnlock()
	var source []LogEntry
	if len(logs) > 0 {
		source = logs
	} else {
		for _, e := range mockEvents {
			source = append(source, LogEntry{
				Timestamp:     e.Timestamp,
				Level:         "INFO",
				Rule:          e.RuleName,
				SourceIP:      e.SourceIP,
				DestinationIP: e.Destination,
				Event:         e.RuleName,
				Description:   e.Description,
				Urgency:       getUrgencyValue(e.Urgency),
			})
		}
	}
	critical := 0
	high := 0
	medium := 0
	low := 0
	for _, log := range source {
		urgency := "medium"
		for _, me := range mockEvents {
			if me.RuleName == log.Rule {
				urgency = me.Urgency
				break
			}
		}
		switch urgency {
		case "critical":
			critical++
		case "high":
			high++
		case "medium":
			medium++
		case "low":
			low++
		}
	}
	data := UrgencyData{
		Critical: critical,
		High:     high,
		Medium:   medium,
		Low:      low,
	}
	json.NewEncoder(w).Encode(data)
}

func timelineDataHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	logStore.mu.RLock()
	logs := make([]LogEntry, len(logStore.logs))
	copy(logs, logStore.logs)
	logStore.mu.RUnlock()
	var source []LogEntry
	if len(logs) > 0 {
		source = logs
	} else {
		for _, e := range mockEvents {
			source = append(source, LogEntry{
				Timestamp:     e.Timestamp,
				Level:         "INFO",
				Rule:          e.RuleName,
				SourceIP:      e.SourceIP,
				DestinationIP: e.Destination,
				Event:         e.RuleName,
				Description:   e.Description,
				Urgency:       getUrgencyValue(e.Urgency),
			})
		}
	}
	labels := []string{}
	accessData := []int{}
	networkData := []int{}
	threatData := []int{}
	now := time.Now()
	for i := 23; i >= 0; i-- {
		hour := now.Add(-time.Duration(i) * time.Hour)
		labels = append(labels, hour.Format("15:04"))
		// Count events in this hour
		ac, nc, tc := 0, 0, 0
		for _, log := range source {
			if log.Timestamp.Format("15:04") == hour.Format("15:04") {
				cat := ""
				for _, me := range mockEvents {
					if me.RuleName == log.Rule {
						cat = me.Category
						break
					}
				}
				switch cat {
				case "access":
					ac++
				case "network":
					nc++
				case "threat":
					tc++
				}
			}
		}
		accessData = append(accessData, ac)
		networkData = append(networkData, nc)
		threatData = append(threatData, tc)
	}
	data := TimelineData{
		Labels: labels,
		Series: []TimelineSeries{
			{Name: "Access", Data: accessData, Color: "#3B82F6"},
			{Name: "Network", Data: networkData, Color: "#10B981"},
			{Name: "Threat", Data: threatData, Color: "#EF4444"},
		},
	}
	json.NewEncoder(w).Encode(data)
}

func topEventsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	// Use real logs if available, else fallback to mockEvents
	logStore.mu.RLock()
	logs := make([]LogEntry, len(logStore.logs))
	copy(logs, logStore.logs)
	logStore.mu.RUnlock()
	var source []LogEntry
	if len(logs) > 0 {
		source = logs
	} else {
		for _, e := range mockEvents {
			source = append(source, LogEntry{
				Timestamp:     e.Timestamp,
				Level:         "INFO",
				Rule:          e.RuleName,
				SourceIP:      e.SourceIP,
				DestinationIP: e.Destination,
				Event:         e.RuleName,
				Description:   e.Description,
				Urgency:       getUrgencyValue(e.Urgency),
			})
		}
	}

	// Group by rule name
	ruleCounts := make(map[string]int)
	ruleUrgency := make(map[string]string)
	for _, event := range source {
		ruleCounts[event.Rule]++
		if _, exists := ruleUrgency[event.Rule]; !exists {
			// Try to find urgency from mockEvents if possible
			urgency := "medium"
			for _, me := range mockEvents {
				if me.RuleName == event.Rule {
					urgency = me.Urgency
					break
				}
			}
			ruleUrgency[event.Rule] = urgency
		}
	}

	// Convert to TopEvent slice
	var topEvents []TopEvent
	for ruleName, count := range ruleCounts {
		// Generate mock sparkline data
		sparkline := []int{}
		for i := 0; i < 10; i++ {
			sparkline = append(sparkline, count/10+rand.Intn(5))
		}
		topEvents = append(topEvents, TopEvent{
			RuleName:  ruleName,
			Sparkline: sparkline,
			Count:     count,
			Urgency:   ruleUrgency[ruleName],
		})
	}

	json.NewEncoder(w).Encode(topEvents)
}

func topSourcesHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	logStore.mu.RLock()
	logs := make([]LogEntry, len(logStore.logs))
	copy(logs, logStore.logs)
	logStore.mu.RUnlock()
	var source []LogEntry
	if len(logs) > 0 {
		source = logs
	} else {
		for _, e := range mockEvents {
			source = append(source, LogEntry{
				Timestamp:     e.Timestamp,
				Level:         "INFO",
				Rule:          e.RuleName,
				SourceIP:      e.SourceIP,
				DestinationIP: e.Destination,
				Event:         e.RuleName,
				Description:   e.Description,
				Urgency:       getUrgencyValue(e.Urgency),
			})
		}
	}
	sourceCounts := make(map[string]int)
	sourceCategory := make(map[string]string)
	for _, event := range source {
		sourceCounts[event.SourceIP]++
		if _, exists := sourceCategory[event.SourceIP]; !exists {
			cat := ""
			for _, me := range mockEvents {
				if me.RuleName == event.Rule {
					cat = me.Category
					break
				}
			}
			sourceCategory[event.SourceIP] = cat
		}
	}
	var topSources []TopSource
	for sourceIP, count := range sourceCounts {
		sparkline := []int{}
		for i := 0; i < 10; i++ {
			sparkline = append(sparkline, count/10+rand.Intn(5))
		}
		topSources = append(topSources, TopSource{
			SourceIP:  sourceIP,
			Sparkline: sparkline,
			Count:     count,
			Category:  sourceCategory[sourceIP],
		})
	}
	json.NewEncoder(w).Encode(topSources)
}

// POST /api/logs - ingest a log entry
func logIngestHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	var entry LogEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	if entry.Level == "" {
		entry.Level = "INFO"
	}
	logStore.mu.Lock()
	logStore.logs = append(logStore.logs, entry)
	logStore.mu.Unlock()
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Log entry stored"))
}

// GET /api/logs?ip=...&event=... - search logs
func logSearchHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	ip := r.URL.Query().Get("ip")
	event := r.URL.Query().Get("event")
	results := []LogEntry{}
	logStore.mu.RLock()
	for _, log := range logStore.logs {
		if ip != "" && log.SourceIP != ip {
			continue
		}
		if event != "" && !strings.Contains(strings.ToLower(log.Rule), strings.ToLower(event)) {
			continue
		}
		results = append(results, log)
	}
	logStore.mu.RUnlock()
	json.NewEncoder(w).Encode(results)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	logStore.mu.RLock()
	logs := make([]LogEntry, len(logStore.logs))
	copy(logs, logStore.logs)
	logStore.mu.RUnlock()
	total := len(logs)
	levelCounts := make(map[string]int)
	ruleCounts := make(map[string]int)
	for _, log := range logs {
		levelCounts[log.Level]++
		ruleCounts[log.Rule]++
	}
	uptime := int(time.Since(startTime).Seconds())
	w.Write([]byte("# HELP logger_logs_total Total number of logs ingested\n"))
	w.Write([]byte("# TYPE logger_logs_total counter\n"))
	w.Write([]byte("logger_logs_total " + strconv.Itoa(total) + "\n"))
	w.Write([]byte("# HELP logger_logs_by_level Number of logs by level\n"))
	w.Write([]byte("# TYPE logger_logs_by_level counter\n"))
	for level, count := range levelCounts {
		w.Write([]byte("logger_logs_by_level{level=\"" + level + "\"} " + strconv.Itoa(count) + "\n"))
	}
	w.Write([]byte("# HELP logger_logs_by_rule Number of logs by rule name\n"))
	w.Write([]byte("# TYPE logger_logs_by_rule counter\n"))
	for rule, count := range ruleCounts {
		w.Write([]byte("logger_logs_by_rule{rule=\"" + rule + "\"} " + strconv.Itoa(count) + "\n"))
	}
	w.Write([]byte("# HELP logger_uptime_seconds Uptime in seconds\n"))
	w.Write([]byte("# TYPE logger_uptime_seconds gauge\n"))
	w.Write([]byte("logger_uptime_seconds " + strconv.Itoa(uptime) + "\n"))
}

// DB-backed summary stats handler
func summaryStatsHandlerDB(w http.ResponseWriter, r *http.Request, db *Database) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	stats, err := db.GetSummaryStats()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed to fetch summary stats"}`))
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// DB-backed urgency data handler
func urgencyDataHandlerDB(w http.ResponseWriter, r *http.Request, db *Database) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	data, err := db.GetUrgencyData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed to fetch urgency data"}`))
		return
	}
	json.NewEncoder(w).Encode(data)
}

// DB-backed timeline data handler
func timelineDataHandlerDB(w http.ResponseWriter, r *http.Request, db *Database) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	data, err := db.GetTimelineData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed to fetch timeline data"}`))
		return
	}
	json.NewEncoder(w).Encode(data)
}

// DB-backed top events handler
func topEventsHandlerDB(w http.ResponseWriter, r *http.Request, db *Database) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	events, err := db.GetTopEvents()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed to fetch top events"}`))
		return
	}
	json.NewEncoder(w).Encode(events)
}

// DB-backed top sources handler
func topSourcesHandlerDB(w http.ResponseWriter, r *http.Request, db *Database) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	sources, err := db.GetTopSources()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed to fetch top sources"}`))
		return
	}
	json.NewEncoder(w).Encode(sources)
}

// DB-backed log ingestion handler
func logIngestHandlerDB(w http.ResponseWriter, r *http.Request, db *Database) {
	enableCORS(w)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	var entry LogEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	if entry.Level == "" {
		entry.Level = "INFO"
	}
	if err := db.InsertLog(entry); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to insert log"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("OK"))
}

// DB-backed log search handler
func logSearchHandlerDB(w http.ResponseWriter, r *http.Request, db *Database) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	ip := r.URL.Query().Get("ip")
	event := r.URL.Query().Get("event")
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}
	logs, err := db.SearchLogs(ip, event, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed to search logs"}`))
		return
	}
	json.NewEncoder(w).Encode(logs)
}

func main() {
	db, err := NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/api/summary", func(w http.ResponseWriter, r *http.Request) { summaryStatsHandlerDB(w, r, db) })
	http.HandleFunc("/api/urgency", func(w http.ResponseWriter, r *http.Request) { urgencyDataHandlerDB(w, r, db) })
	http.HandleFunc("/api/timeline", func(w http.ResponseWriter, r *http.Request) { timelineDataHandlerDB(w, r, db) })
	http.HandleFunc("/api/top-events", func(w http.ResponseWriter, r *http.Request) { topEventsHandlerDB(w, r, db) })
	http.HandleFunc("/api/top-sources", func(w http.ResponseWriter, r *http.Request) { topSourcesHandlerDB(w, r, db) })
	http.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			logIngestHandlerDB(w, r, db)
		} else {
			logSearchHandlerDB(w, r, db)
		}
	})
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/", handleOptions)
	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
