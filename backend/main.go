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
// (reuse or redefine as needed)
type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	RuleName  string            `json:"ruleName"`
	SourceIP  string            `json:"sourceIP"`
	Metadata  map[string]string `json:"metadata"`
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
				Timestamp: e.Timestamp,
				Level:     "INFO",
				Message:   e.Description,
				RuleName:  e.RuleName,
				SourceIP:  e.SourceIP,
				Metadata:  map[string]string{},
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
			if me.RuleName == log.RuleName {
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
				Timestamp: e.Timestamp,
				Level:     "INFO",
				Message:   e.Description,
				RuleName:  e.RuleName,
				SourceIP:  e.SourceIP,
				Metadata:  map[string]string{},
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
			if me.RuleName == log.RuleName {
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
				Timestamp: e.Timestamp,
				Level:     "INFO",
				Message:   e.Description,
				RuleName:  e.RuleName,
				SourceIP:  e.SourceIP,
				Metadata:  map[string]string{},
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
					if me.RuleName == log.RuleName {
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
		// Convert mockEvents to LogEntry for compatibility
		for _, e := range mockEvents {
			source = append(source, LogEntry{
				Timestamp: e.Timestamp,
				Level:     "INFO",
				Message:   e.Description,
				RuleName:  e.RuleName,
				SourceIP:  e.SourceIP,
				Metadata:  map[string]string{},
			})
		}
	}

	// Group by rule name
	ruleCounts := make(map[string]int)
	ruleUrgency := make(map[string]string)
	for _, event := range source {
		ruleCounts[event.RuleName]++
		if _, exists := ruleUrgency[event.RuleName]; !exists {
			// Try to find urgency from mockEvents if possible
			urgency := "medium"
			for _, me := range mockEvents {
				if me.RuleName == event.RuleName {
					urgency = me.Urgency
					break
				}
			}
			ruleUrgency[event.RuleName] = urgency
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
				Timestamp: e.Timestamp,
				Level:     "INFO",
				Message:   e.Description,
				RuleName:  e.RuleName,
				SourceIP:  e.SourceIP,
				Metadata:  map[string]string{},
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
				if me.RuleName == event.RuleName {
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
	if entry.Metadata == nil {
		entry.Metadata = make(map[string]string)
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
		if event != "" && !strings.Contains(strings.ToLower(log.RuleName), strings.ToLower(event)) {
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
		ruleCounts[log.RuleName]++
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

func main() {
	http.HandleFunc("/api/stats/summary", summaryStatsHandler)
	http.HandleFunc("/api/events/urgency", urgencyDataHandler)
	http.HandleFunc("/api/events/timeline", timelineDataHandler)
	http.HandleFunc("/api/events/top", topEventsHandler)
	http.HandleFunc("/api/events/sources", topSourcesHandler)
	http.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			logIngestHandler(w, r)
		} else if r.Method == http.MethodGet {
			logSearchHandler(w, r)
		} else {
			enableCORS(w)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/api/", handleOptions)

	log.Println("Starting backend server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
