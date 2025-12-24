@echo off
echo Starting local web server for frontend...
echo.
echo Frontend will be available at: http://localhost:8000
echo Backend should be running at: http://localhost:8080
echo.
echo Press Ctrl+C to stop the server
echo.
python -m http.server 8000

