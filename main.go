package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Metadata  map[string]string `json:"metadata"`
}

// InMemoryDB is a simple thread-safe in-memory log store
type InMemoryDB struct {
	logs []LogEntry
	mu   sync.RWMutex
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		logs: make([]LogEntry, 0),
	}
}

func (db *InMemoryDB) Add(entry LogEntry) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.logs = append(db.logs, entry)
}

func (db *InMemoryDB) GetAll() []LogEntry {
	db.mu.RLock()
	defer db.mu.RUnlock()
	logsCopy := make([]LogEntry, len(db.logs))
	copy(logsCopy, db.logs)
	return logsCopy
}

func (db *InMemoryDB) Filter(level, keyword string, from, to time.Time) []LogEntry {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var filtered []LogEntry
	for _, log := range db.logs {
		if level != "" && log.Level != level {
			continue
		}
		if !from.IsZero() && log.Timestamp.Before(from) {
			continue
		}
		if !to.IsZero() && log.Timestamp.After(to) {
			continue
		}
		if keyword != "" && !strings.Contains(log.Message, keyword) {
			continue
		}
		filtered = append(filtered, log)
	}
	return filtered
}

var (
	startTime = time.Now()
)

var db = NewInMemoryDB()

func logIngestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	var entry LogEntry
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&entry)
	if err != nil {
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
	db.Add(entry)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Log entry stored"))
}

func logsAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	level := r.URL.Query().Get("level")
	keyword := r.URL.Query().Get("keyword")
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var from, to time.Time
	var err error
	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid 'from' timestamp"))
			return
		}
	}
	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid 'to' timestamp"))
			return
		}
	}
	logs := db.Filter(level, keyword, from, to)
	json.NewEncoder(w).Encode(logs)
}

func logsStreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logs := db.GetAll()
	json.NewEncoder(w).Encode(logs)
}

func statsAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logs := db.GetAll()
	// Logs per minute (last hour)
	perMinute := make(map[string]int)
	perHour := make(map[string]int)
	levelCounts := make(map[string]int)
	cutoffHour := time.Now().Add(-1 * time.Hour)
	cutoffDay := time.Now().Add(-24 * time.Hour)
	for _, log := range logs {
		if log.Timestamp.After(cutoffHour) {
			min := log.Timestamp.Format("15:04")
			perMinute[min]++
		}
		if log.Timestamp.After(cutoffDay) {
			hour := log.Timestamp.Format("Jan 2 15:00")
			perHour[hour]++
		}
		levelCounts[log.Level]++
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"perMinute":   perMinute,
		"perHour":     perHour,
		"levelCounts": levelCounts,
	})
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	logs := db.GetAll()
	total := len(logs)
	levelCounts := make(map[string]int)
	for _, log := range logs {
		levelCounts[log.Level]++
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
	w.Write([]byte("# HELP logger_uptime_seconds Uptime in seconds\n"))
	w.Write([]byte("# TYPE logger_uptime_seconds gauge\n"))
	w.Write([]byte("logger_uptime_seconds " + strconv.Itoa(uptime) + "\n"))
}

const htmlPage = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Logger UI</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 2em; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ccc; padding: 8px; text-align: left; }
        th { background: #f4f4f4; }
        input, select { margin: 0 0.5em 1em 0; }
        .charts { display: flex; gap: 2em; margin-bottom: 2em; }
        .chart-container { width: 400px; }
    </style>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
</head>
<body>
    <h1>Log Viewer</h1>
    <div class="charts">
        <div class="chart-container">
            <canvas id="barChart"></canvas>
        </div>
        <div class="chart-container">
            <canvas id="pieChart"></canvas>
        </div>
    </div>
    <div>
        <label>Level:
            <select id="levelFilter">
                <option value="">All</option>
                <option value="INFO">INFO</option>
                <option value="WARN">WARN</option>
                <option value="ERROR">ERROR</option>
                <option value="DEBUG">DEBUG</option>
            </select>
        </label>
        <label>Keyword:
            <input type="text" id="keywordFilter" placeholder="Search message...">
        </label>
        <label>From:
            <input type="datetime-local" id="fromFilter">
        </label>
        <label>To:
            <input type="datetime-local" id="toFilter">
        </label>
        <button onclick="loadLogs()">Search</button>
    </div>
    <table id="logsTable">
        <thead>
            <tr>
                <th>Timestamp</th>
                <th>Level</th>
                <th>Message</th>
                <th>Metadata</th>
            </tr>
        </thead>
        <tbody></tbody>
    </table>
    <script>
        let pollInterval = null;
        function toRFC3339Local(dt) {
            if (!dt) return '';
            return new Date(dt).toISOString();
        }
        async function loadLogs() {
            const level = document.getElementById('levelFilter').value;
            const keyword = document.getElementById('keywordFilter').value;
            const from = document.getElementById('fromFilter').value;
            const to = document.getElementById('toFilter').value;
            let url = '/api/logs';
            const params = [];
            if (level) params.push('level=' + encodeURIComponent(level));
            if (keyword) params.push('keyword=' + encodeURIComponent(keyword));
            if (from) params.push('from=' + encodeURIComponent(toRFC3339Local(from)));
            if (to) params.push('to=' + encodeURIComponent(toRFC3339Local(to)));
            if (params.length) url += '?' + params.join('&');
            const res = await fetch(url);
            const logs = await res.json();
            renderLogs(logs);
        }
        function renderLogs(logs) {
            const tbody = document.querySelector('#logsTable tbody');
            tbody.innerHTML = '';
            logs.forEach(log => {
                const tr = document.createElement('tr');
                tr.innerHTML = '<td>' + new Date(log.timestamp).toLocaleString() + '</td><td>' + log.level + '</td><td>' + log.message + '</td><td>' + JSON.stringify(log.metadata) + '</td>';
                tbody.appendChild(tr);
            });
        }
        async function loadCharts() {
            const res = await fetch('/api/stats');
            const stats = await res.json();
            // Bar chart: logs per minute (last hour)
            const barLabels = Object.keys(stats.perMinute).sort();
            const barData = barLabels.map(l => stats.perMinute[l]);
            if (window.barChart) window.barChart.destroy();
            window.barChart = new Chart(document.getElementById('barChart'), {
                type: 'bar',
                data: {
                    labels: barLabels,
                    datasets: [{ label: 'Logs/min (last hour)', data: barData, backgroundColor: '#4e79a7' }]
                },
                options: { scales: { x: { title: { display: true, text: 'Minute' } }, y: { beginAtZero: true } } }
            });
            // Pie chart: log level distribution
            const pieLabels = Object.keys(stats.levelCounts);
            const pieData = pieLabels.map(l => stats.levelCounts[l]);
            if (window.pieChart) window.pieChart.destroy();
            window.pieChart = new Chart(document.getElementById('pieChart'), {
                type: 'pie',
                data: {
                    labels: pieLabels,
                    datasets: [{ data: pieData, backgroundColor: ['#4e79a7','#f28e2b','#e15759','#76b7b2'] }]
                },
                options: { plugins: { legend: { position: 'bottom' } } }
            });
        }
        function startPolling() {
            if (pollInterval) clearInterval(pollInterval);
            pollInterval = setInterval(async () => {
                await loadLogs();
                await loadCharts();
            }, 2000);
        }
        window.onload = function() {
            loadLogs();
            loadCharts();
            startPolling();
        };
    </script>
</body>
</html>
`

func uiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlPage))
}

func startLogIngestServer() {
	http.HandleFunc("/logs", logIngestHandler)
	log.Println("Log ingestion endpoint listening on :9000")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatalf("Log ingest server failed: %v", err)
	}
}

func startWebUIServer() {
	http.HandleFunc("/", uiHandler)
	http.HandleFunc("/api/logs", logsAPIHandler)
	http.HandleFunc("/api/logs/stream", logsStreamHandler)
	http.HandleFunc("/api/stats", statsAPIHandler)
	http.HandleFunc("/metrics", metricsHandler)
	log.Println("Web UI listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Web UI server failed: %v", err)
	}
}

func main() {
	go startLogIngestServer()
	go startWebUIServer()
	log.Println("Logger application starting...")
	select {} // Block forever
}
