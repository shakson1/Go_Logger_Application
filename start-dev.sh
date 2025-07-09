#!/bin/bash

echo "🚀 Starting Security Dashboard Development Environment"
echo ""

# Check if backend is running
if ! curl -s http://localhost:8080/api/stats/summary > /dev/null; then
    echo "📦 Starting Backend (Go) server..."
    cd backend
    go run main.go &
    BACKEND_PID=$!
    cd ..
    echo "✅ Backend started on http://localhost:8080"
else
    echo "✅ Backend already running on http://localhost:8080"
fi

# Wait a moment for backend to fully start
sleep 2

# Check if frontend is running
if ! curl -s http://localhost:3000 > /dev/null; then
    echo "📦 Starting Frontend (React) server..."
    cd frontend
    npm start &
    FRONTEND_PID=$!
    cd ..
    echo "✅ Frontend started on http://localhost:3000"
else
    echo "✅ Frontend already running on http://localhost:3000"
fi

echo ""
echo "🎉 Development environment is ready!"
echo "   Frontend: http://localhost:3000"
echo "   Backend API: http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop all services"

# Wait for interrupt signal
trap 'echo ""; echo "🛑 Stopping services..."; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit' INT
wait 