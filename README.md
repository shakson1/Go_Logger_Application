# Security Dashboard

A full-stack web application that provides a Splunk Enterprise Security-like dashboard for monitoring security notable events. Built with React (frontend) and Go (backend) with SQLite database persistence.

## Features

🔹 **Frontend (React + TypeScript + Tailwind CSS)**:
- Dark-themed dashboard similar to Splunk Enterprise Security
- Real-time statistics tiles showing Access, Network, Threat, and UBA notables
- Interactive charts using Chart.js:
  - Bar chart: "Notable Events by Urgency" (Critical, High, Medium, Low)
  - Line chart: "Notable Events Over Time" (Access/Network/Threat trends)
- Data tables with sparkline charts:
  - Top notable events table with drilldown (click to see logs)
  - Top event sources table
  - Pagination for tables
- **Log Search**: Search logs by IP or event/rule name, results persist across dashboard refreshes
- **Drilldown**: Click a notable event to see all logs for that rule
- **Home Button**: Instantly scroll to top
- **Refresh Interval Selector**: Choose 5/10/15/30s background refresh, does not reset your view
- Responsive design optimized for desktop and tablet

🔹 **Backend (Go + SQLite)**:
- RESTful API endpoints serving JSON data from SQLite database
- **Persistent log storage** - data survives container restarts
- Log ingestion endpoint: `POST /api/logs`
- Log search endpoint: `GET /api/logs?ip=...&event=...`
- All dashboard endpoints aggregate from real logs stored in SQLite
- CORS support for frontend integration
- Containerized with Docker (Debian-based for SQLite compatibility)
- Prometheus metrics endpoint: `/metrics`

## Project Structure

```
logger/
├── backend/                 # Go backend service
│   ├── main.go             # Main Go application with API handlers
│   ├── database.go         # SQLite database operations
│   ├── go.mod              # Go module file
│   └── Dockerfile          # Backend container (Debian-based)
├── frontend/               # React frontend service
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── services/       # API service
│   │   ├── types/          # TypeScript types
│   │   └── ...
│   ├── Dockerfile          # Frontend container
│   └── nginx.conf          # Nginx configuration
├── docker-compose.yml      # Multi-service orchestration
└── README.md              # This file
```

## Quick Start

### Option 1: Using Docker Compose (Recommended)

1. **Clone and navigate to the project**:
   ```bash
   cd logger
   ```

2. **Build and run both services**:
   ```bash
   docker-compose up --build
   ```

3. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

### Option 2: Development Mode

#### Backend Setup

1. **Navigate to backend directory**:
   ```bash
   cd backend
   ```

2. **Install SQLite dependencies** (if needed):
   ```bash
   # On Ubuntu/Debian
   sudo apt-get install gcc libsqlite3-dev
   
   # On macOS
   brew install sqlite3
   ```

3. **Run the Go server**:
   ```bash
   CGO_ENABLED=1 go run .
   ```

4. **Backend will be available at**: http://localhost:8080

#### Frontend Setup

1. **Navigate to frontend directory**:
   ```bash
   cd frontend
   ```

2. **Install dependencies**:
   ```bash
   npm install
   ```

3. **Start the development server**:
   ```bash
   npm start
   ```

4. **Frontend will be available at**: http://localhost:3000

## API Endpoints

### Log Ingestion
```http
POST /api/logs
Content-Type: application/json

{
  "timestamp": "2024-07-09T12:00:00Z",
  "level": "INFO",
  "rule": "Suspicious Login Attempt",
  "sourceIP": "192.168.1.100",
  "destinationIP": "10.0.0.1",
  "event": "Suspicious Login Attempt",
  "description": "Multiple failed login attempts detected",
  "urgency": 4
}
```

### Log Search
```http
GET /api/logs?ip=192.168.1.100&event=Suspicious&limit=100
```
Returns all logs matching the IP and/or event/rule name (max 1000 results).

### Dashboard Endpoints (all aggregate from SQLite database)
- `GET /api/summary` - Dashboard summary statistics
- `GET /api/urgency` - Bar chart data by urgency
- `GET /api/timeline` - Time series data for line chart
- `GET /api/top-events` - Top notable events (clickable for drilldown)
- `GET /api/top-sources` - Top event sources

### Metrics
```http
GET /metrics
```
Returns Prometheus-formatted metrics including log counts, logs by level, logs by rule, and uptime.

## UI Features
- **Home Button**: Instantly scroll to top
- **Refresh Interval Selector**: Choose 5/10/15/30s background refresh, does not reset your view
- **Log Search**: Search logs by IP or event/rule name, results persist across dashboard refreshes
- **Drilldown**: Click a notable event to see all logs for that rule
- **Pagination**: For tables
- **Responsive Design**: Works on desktop and tablet

## Technologies Used

### Frontend
- **React 18** with TypeScript
- **Tailwind CSS** for styling
- **Chart.js** with react-chartjs-2 for charts
- **Lucide React** for icons
- **Responsive design** with mobile-first approach

### Backend
- **Go 1.21** with standard library
- **SQLite** for persistent data storage
- **RESTful API** design
- **CORS** support for cross-origin requests
- **Prometheus metrics** for monitoring

### DevOps
- **Docker** for containerization
- **Docker Compose** for multi-service orchestration
- **Nginx** for frontend serving and API proxying
- **Debian-based images** for SQLite compatibility

## Development

### Adding New Features

1. **Backend**: Add new endpoints in `backend/main.go` and database operations in `backend/database.go`
2. **Frontend**: Create new components in `frontend/src/components/`
3. **Types**: Update `frontend/src/types/index.ts` for new data structures
4. **API**: Add new methods in `frontend/src/services/api.ts`

### Styling

The application uses Tailwind CSS with custom colors matching Splunk's dark theme:
- `splunk-dark`: #1a1a1a
- `splunk-darker`: #0f0f0f
- `splunk-gray`: #2d2d2d
- `splunk-light-gray`: #404040

### Data Structure

The application uses TypeScript interfaces for type safety:
- `SummaryStats`: Dashboard summary statistics
- `UrgencyData`: Bar chart data by urgency
- `TimelineData`: Time series data for line charts
- `TopEvent`: Notable event with sparkline
- `TopSource`: Event source with sparkline
- `LogEntry`: Ingested log entry

## Troubleshooting

### Common Issues

1. **CORS Errors**: Ensure the backend is running and CORS is properly configured
2. **Port Conflicts**: Check if ports 3000 and 8080 are available
3. **Build Errors**: Ensure all dependencies are installed (`npm install` for frontend, `go mod tidy` for backend)
4. **SQLite Errors**: Ensure CGO is enabled (`CGO_ENABLED=1`) when building Go with SQLite

### Docker Issues

1. **Container Build Fails**: Check Docker is running and has sufficient resources
2. **Port Mapping**: Ensure ports 3000 and 8080 are not used by other services
3. **Network Issues**: Use `docker-compose down` and `docker-compose up --build` to rebuild
4. **SQLite Compatibility**: Backend uses Debian-based images for SQLite support

### Frontend Issues

1. **Black Screen**: Check browser console for JavaScript errors
2. **API Errors**: Verify API endpoints match between frontend and backend
3. **Chart Rendering**: Ensure data is properly formatted for Chart.js components

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is for demonstration purposes. Feel free to use and modify as needed. 