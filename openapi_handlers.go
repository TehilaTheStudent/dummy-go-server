package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// OpenAPI schemas directories
var (
	yamlDir = "schemas/yaml"
	jsonDir = "schemas/json"
)

// listAvailableSchemas returns a list of all available OpenAPI schemas
func listAvailableSchemas(w http.ResponseWriter, r *http.Request) {
	// Create a response structure
	response := struct {
		YAMLSchemas []string `json:"yamlSchemas"`
		JSONSchemas []string `json:"jsonSchemas"`
	}{}

	// List YAML schemas
	yamlFiles, err := listFilesInDir(yamlDir, ".yaml")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing YAML schemas: %v", err), http.StatusInternalServerError)
		return
	}
	response.YAMLSchemas = yamlFiles

	// List JSON schemas
	jsonFiles, err := listFilesInDir(jsonDir, ".json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing JSON schemas: %v", err), http.StatusInternalServerError)
		return
	}
	response.JSONSchemas = jsonFiles

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// listFilesInDir returns a list of files with the specified extension in the given directory
func listFilesInDir(dir string, ext string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ext) {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// serveOpenAPISchema serves the OpenAPI schema file based on the requested path
func serveOpenAPISchema(w http.ResponseWriter, r *http.Request) {
	// Extract the path components
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/schemas/"), "/")
	
	if len(pathParts) < 2 {
		http.Error(w, "Invalid path format. Use /schemas/{format}/{filename}", http.StatusBadRequest)
		return
	}
	
	format := pathParts[0]   // "yaml" or "json"
	filename := pathParts[1] // The filename
	
	// Determine the directory based on format
	var dir string
	switch format {
	case "yaml":
		dir = yamlDir
		w.Header().Set("Content-Type", "text/yaml")
	case "json":
		dir = jsonDir
		w.Header().Set("Content-Type", "application/json")
	default:
		http.Error(w, "Unsupported format. Use 'yaml' or 'json'", http.StatusBadRequest)
		return
	}

	// Get the file path
	filePath := filepath.Join(dir, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Schema file not found", http.StatusNotFound)
		return
	}

	// Read and serve the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Error reading schema file", http.StatusInternalServerError)
		return
	}

	w.Write(content)
}

// convertYAMLToJSON converts a YAML file to JSON and serves it
func convertYAMLToJSON(w http.ResponseWriter, r *http.Request) {
	// Extract the filename from the URL path
	filename := strings.TrimPrefix(r.URL.Path, "/convert/yaml-to-json/")
	
	if !strings.HasSuffix(filename, ".yaml") && !strings.HasSuffix(filename, ".yml") {
		http.Error(w, "Only YAML files can be converted to JSON", http.StatusBadRequest)
		return
	}

	// Get the file path
	filePath := filepath.Join(yamlDir, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "YAML file not found", http.StatusNotFound)
		return
	}

	// Read the YAML file
	yamlContent, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Error reading YAML file", http.StatusInternalServerError)
		return
	}

	// Parse YAML
	var data interface{}
	err = yaml.Unmarshal(yamlContent, &data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing YAML: %v", err), http.StatusBadRequest)
		return
	}

	// Convert to JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}
}

// convertJSONToYAML converts a JSON file to YAML and serves it
func convertJSONToYAML(w http.ResponseWriter, r *http.Request) {
	// Extract the filename from the URL path
	filename := strings.TrimPrefix(r.URL.Path, "/convert/json-to-yaml/")
	
	if !strings.HasSuffix(filename, ".json") {
		http.Error(w, "Only JSON files can be converted to YAML", http.StatusBadRequest)
		return
	}

	// Get the file path
	filePath := filepath.Join(jsonDir, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "JSON file not found", http.StatusNotFound)
		return
	}

	// Read the JSON file
	jsonContent, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Error reading JSON file", http.StatusInternalServerError)
		return
	}

	// Parse JSON
	var data interface{}
	err = json.Unmarshal(jsonContent, &data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Convert to YAML
	yamlContent, err := yaml.Marshal(data)
	if err != nil {
		http.Error(w, "Error converting to YAML", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/yaml")
	w.Write(yamlContent)
}
