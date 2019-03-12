package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
	"github.com/go-yaml/yaml"
)

func FuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	// marshal
	f["toJson"] = toJson
	f["toYaml"] = toYaml
	f["toToml"] = toToml
	f["toBool"] = toBool
	// file
	f["fileExists"] = fileExists
	f["fileSize"] = fileSize
	f["fileLastModified"] = fileLastModified
	f["fileGetBytes"] = fileGetBytes
	f["fileGetString"] = fileGetString
	return f
}

// toBool takes a string and converts it to a bool.
// On marshal error will panic if in strict mode, otherwise returns false.
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
//
// This is designed to be called from a template.
func toBool(value string) bool {
	result, err := strconv.ParseBool(value)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return false
	}
	return result
}

// toJson takes an interface, marshals it to json, and returns a string.
// On marshal error will panic if in strict mode, otherwise returns empty string.
//
// This is designed to be called from a template.
func toJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return ""
	}
	return string(data)
}

// toYaml takes an interface, marshals it to yaml, and returns a string.
// On marshal error will panic if in strict mode, otherwise returns empty string.
//
// This is designed to be called from a template.
func toYaml(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return ""
	}
	return string(data)
}

// toToml takes an interface, marshals it to toml, and returns a string.
// On marshal error will panic if in strict mode, otherwise returns empty string.
//
// This is designed to be called from a template.
func toToml(v interface{}) string {
	b := bytes.NewBuffer(nil)
	e := toml.NewEncoder(b)
	err := e.Encode(v)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return ""
	}
	return b.String()
}

func fileExists(file string) bool {
	_, err := os.Stat(file)

	return err == nil
}

func fileSize(file string) int64 {
	info, err := os.Stat(file)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return 0
	}
	return info.Size()
}

func fileLastModified(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return time.Unix(0, 0)
	}
	return info.ModTime()
}

func fileGetBytes(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return []byte{}
	}
	return data
}

func fileGetString(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return ""
	}
	return string(data)
}
