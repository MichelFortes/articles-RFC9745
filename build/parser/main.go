package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pb33f/libopenapi"
)

const extensionDeprecatedAt = "x-deprecated-at"
const extensionDeprecatedLink = "x-deprecated-link"
const extensionDeprecatedSunset = "x-deprecated-sunset"

type Deprecated struct {
	At     string `json:"at,omitempty"`
	Link   string `json:"link,omitempty"`
	Sunset string `json:"sunset,omitempty"`
}

type Backend struct {
	Path string `json:"path"`
	Host string `json:"host"`
}

type Operation struct {
	Path       string     `json:"path"`
	Method     string     `json:"method"`
	Deprecated Deprecated `json:"deprecated,omitempty"`
	Backend    Backend    `json:"backend"`
}

func (o Operation) MarshalJSON() ([]byte, error) {

	d := o.Deprecated
	if d.At == "" && d.Link == "" && d.Sunset == "" {
		return json.Marshal(struct {
			Path    string  `json:"path"`
			Method  string  `json:"method"`
			Backend Backend `json:"backend"`
		}{
			Path:    o.Path,
			Method:  o.Method,
			Backend: o.Backend,
		})
	}

	return json.Marshal(struct {
		Path       string     `json:"path"`
		Method     string     `json:"method"`
		Deprecated Deprecated `json:"deprecated"`
		Backend    Backend    `json:"backend"`
	}{
		Path:       o.Path,
		Method:     o.Method,
		Deprecated: d,
		Backend:    o.Backend,
	})
}

type OperationsList struct {
	List []Operation `json:"list"`
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Error: Both 'definition path' and 'template destination' must be provided.")
		fmt.Println("Usage: go run main.go <definition_path> <template_destination>")
		os.Exit(1)
	}

	definitionPath := os.Args[1]
	templateDestination := os.Args[2]

	definition, err := os.ReadFile(definitionPath)
	if err != nil {
		fmt.Printf("Error: Unable to read file '%s': %v\n", definitionPath, err)
		os.Exit(1)
	}

	doc, err := libopenapi.NewDocument(definition)
	if err != nil {
		fmt.Printf("Error: Unable to parse OpenAPI definition: %v\n", err)
		os.Exit(1)
	}

	v3doc, errs := doc.BuildV3Model()
	if len(errs) > 0 {
		fmt.Printf("Error: Unable to build v3 model: %v\n", err)
		os.Exit(1)
	}

	krakendOperations := OperationsList{
		List: []Operation{},
	}

	for pathItemPair := v3doc.Model.Paths.PathItems.First(); pathItemPair != nil; pathItemPair = pathItemPair.Next() {
		path := pathItemPair.Key()
		pathItem := pathItemPair.Value()

		for operationPair := pathItem.GetOperations().First(); operationPair != nil; operationPair = operationPair.Next() {

			method := operationPair.Key()
			operation := operationPair.Value()

			kdOp := Operation{
				Path:       path,
				Method:     method,
				Deprecated: Deprecated{},
				Backend: Backend{
					Host: "http://backend:8888",
					Path: path,
				},
			}

			if nodeDepAt, present := operation.Extensions.Get(extensionDeprecatedAt); present {
				depAt, err := getunixFromTimestampRFC3339(nodeDepAt.Value)
				if err != nil {
					fmt.Printf("Error: Unable to parse deprecated_at timestamp: %v\n", err)
					os.Exit(1)
				}
				kdOp.Deprecated.At = fmt.Sprintf("@%d", depAt)
			}

			if nodeDepLink, present := operation.Extensions.Get(extensionDeprecatedLink); present {
				kdOp.Deprecated.Link = nodeDepLink.Value
			}

			if nodeDepSunset, present := operation.Extensions.Get(extensionDeprecatedSunset); present {
				kdOp.Deprecated.Sunset, err = getRFC1123FromTimestampRFC3339(nodeDepSunset.Value)
				if err != nil {
					fmt.Printf("Error: Unable to parse deprecated_sunset timestamp: %v\n", err)
					os.Exit(1)
				}
			}

			krakendOperations.List = append(krakendOperations.List, kdOp)
		}
	}

	jsonOutput, err := json.MarshalIndent(krakendOperations, "", "  ")
	if err != nil {
		fmt.Printf("Error: Unable to marshal operations to JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(templateDestination, jsonOutput, 0644)
	if err != nil {
		fmt.Printf("Error: Unable to write JSON to file '%s': %v\n", templateDestination, err)
		os.Exit(1)
	}

	fmt.Printf("Operations successfully written to '%s'\n", templateDestination)
}

func getunixFromTimestampRFC3339(timestamp string) (int64, error) {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func getRFC1123FromTimestampRFC3339(timestamp string) (string, error) {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return "", err
	}
	return t.Format(time.RFC1123), nil
}
