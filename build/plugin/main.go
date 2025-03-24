package main

import (
	"errors"
	"fmt"
	"io"
)

type ResponseWrapper interface {
	Data() map[string]interface{}
	Io() io.Reader
	IsComplete() bool
	StatusCode() int
	Headers() map[string][]string
}

var errUnknownType = errors.New("unknown request type")

func main() {}

func init() {
	fmt.Println(string(ModifierRegisterer), "loaded!")
}

var ModifierRegisterer = registerer("deprecated-header")

type registerer string

func (r registerer) RegisterModifiers(f func(
	name string,
	factoryFunc func(map[string]interface{}) func(interface{}) (interface{}, error),
	appliesToRequest bool,
	appliesToResponse bool,
)) {
	f(string(r), r.responseHandler, false, true)
}

func (r registerer) responseHandler(
	cfg map[string]interface{},
) func(interface{}) (interface{}, error) {

	pluginConfig, ok := cfg["deprecated-header"]
	if !ok {
		return func(input interface{}) (interface{}, error) {
			return input, errors.New("missing deprecated-header config")
		}
	}

	config, ok := pluginConfig.(map[string]interface{})
	if !ok {
		return func(input interface{}) (interface{}, error) {
			return input, errors.New("invalid deprecated-header config")
		}
	}

	cfgDeprecated, hasDeprecated := config["deprecated"].(string)
	cfgSunset, hasSunset := config["sunset"].(string)
	cfgLink, hasLink := config["link"].(string)

	return func(input interface{}) (interface{}, error) {

		resp, ok := input.(ResponseWrapper)
		if !ok {
			return nil, errUnknownType
		}

		if resp.Headers() != nil {
			if hasDeprecated {
				resp.Headers()["Deprecated"] = []string{cfgDeprecated}
			}
			if hasSunset {
				resp.Headers()["Sunset"] = []string{cfgSunset}
			}
			if hasLink {
				resp.Headers()["Link"] = []string{cfgLink}
			}
		}

		return input, nil
	}
}
