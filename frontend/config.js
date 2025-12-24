const CONFIG = {
    // If frontend is served by backend, use relative URL
    // If frontend is on different port, use full URL
    API_BASE_URL: window.location.origin === 'http://localhost:8080' || window.location.origin === 'http://127.0.0.1:8080'
        ? '/v1'  // Same origin as backend
        : 'http://localhost:8080/v1'  // Different origin
};
