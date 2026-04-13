package api

import (
	"encoding/json"
	"net/http"
	"os"

	"health-monitor/internal/generated"
)

// handleOpenAPISpec serves the OpenAPI specification in JSON format
// Uses the embedded spec from generated code
func (s *Server) handleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	spec, err := generated.GetSwagger()
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get OpenAPI spec")
		http.Error(w, "Failed to load OpenAPI specification", http.StatusInternalServerError)
		return
	}

	// Marshal to JSON
	data, err := json.Marshal(spec)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to marshal OpenAPI spec")
		http.Error(w, "Failed to serialize OpenAPI specification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// handleSwaggerUI serves the Swagger UI HTML page
func (s *Server) handleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Health Monitor API - Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "/openapi.yaml",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// Alternative: serve OpenAPI spec from file system (development mode)
func (s *Server) serveOpenAPISpecFromFile(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("api/openapi.yaml")
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to read OpenAPI spec file")
		http.Error(w, "Failed to load OpenAPI specification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(data)
}
