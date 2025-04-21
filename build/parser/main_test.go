// filepath: /Users/michelfortes/Development/articles/articles-RFC9745/build/parser/main_test.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestGetUnixFromTimestampRFC3339(t *testing.T) {
	timestamp := "2023-03-15T14:00:00Z"
	expected := int64(1678888800)

	result, err := getunixFromTimestampRFC3339(timestamp)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestGetRFC1123FromTimestampRFC3339(t *testing.T) {
	timestamp := "2023-03-15T14:00:00Z"
	expected := "Wed, 15 Mar 2023 14:00:00 UTC"

	result, err := getRFC1123FromTimestampRFC3339(timestamp)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestOperationMarshalJSON(t *testing.T) {
	op := Operation{
		Path:   "/test",
		Method: "GET",
		Backend: Backend{
			Path: "/test",
			Host: "http://backend:8888",
		},
	}

	// Test without deprecated fields
	data, err := json.Marshal(op)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := `{"path":"/test","method":"GET","backend":{"path":"/test","host":"http://backend:8888"}}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}

	// Test with deprecated fields
	op.Deprecated = Deprecated{
		At:     "@1678888800",
		Link:   "http://example.com",
		Sunset: "Wed, 15 Mar 2023 14:00:00 UTC",
	}

	data, err = json.Marshal(op)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected = `{"path":"/test","method":"GET","deprecated":{"at":"@1678888800","link":"http://example.com","sunset":"Wed, 15 Mar 2023 14:00:00 UTC"},"backend":{"path":"/test","host":"http://backend:8888"}}`
	if string(data) != expected {
		t.Errorf("Expected: %s got %s", expected, string(data))
	}
}

func TestMainFunction(t *testing.T) {
	// Mock input files
	tempDir := t.TempDir()
	definitionPath := fmt.Sprint(tempDir, "/test_definition.json")
	templateDestination := fmt.Sprint(tempDir, "test_output.json")

	definitionContent := `{
		"openapi": "3.1.1",
		"paths": {
			"/test": {
				"get": {
					"x-deprecated-at": "2023-03-15T14:00:00Z",
					"x-deprecated-link": "http://example.com",
					"x-deprecated-sunset": "2023-03-16T14:00:00Z"
				}
			}
		}
	}`

	expectedOutput := `{
  "list": [
    {
      "path": "/test",
      "method": "get",
      "deprecated": {
        "at": "@1678888800",
        "link": "http://example.com",
        "sunset": "Thu, 16 Mar 2023 14:00:00 UTC"
      },
      "backend": {
        "path": "/test",
        "host": "http://backend:8888"
      }
    }
  ]
}`

	err := os.WriteFile(definitionPath, []byte(definitionContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create mock definition file: %v", err)
	}

	// Run the main function
	os.Args = []string{"main", definitionPath, templateDestination}
	main()

	// Verify output file
	output, err := os.ReadFile(templateDestination)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(output) != expectedOutput {
		t.Errorf("\nExpected:\n%s\nGot\n%s", expectedOutput, string(output))
	}
}
