package csv

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Given a stream of CSV records, generate a stream of JSON records, one per line. The headers
// are treated as paths into the resulting JSON object, so a CSV file containing the header
// foo.bar,baz and the data 1, 2 will be converted into a JSON object like {"foo": {"bar": 1}, "baz": 2}
//
// If column values can be successfully unmarshalled as JSON numbers, booleans, objects or arrays then
// the value will be encoded as the corresponding JSON object, otherwise it will be encoded as a string.
// Use --strings to force all column values to be encoded as JSON strings.
//
type CsvToJsonProcess struct {
	BaseObject  string
	StringsOnly bool
}

func (proc *CsvToJsonProcess) writeToMap(m map[string]interface{}, p []string, v interface{}) {
	if len(p) == 1 {
		m[p[0]] = v
	} else if len(p) > 1 {
		var o interface{}
		var mo map[string]interface{}
		var ok bool
		if o, ok = m[p[0]]; !ok {
			mo = map[string]interface{}{}
		} else {
			if mo, ok = o.(map[string]interface{}); !ok {
				mo = map[string]interface{}{}
			}
		}
		m[p[0]] = mo
		proc.writeToMap(mo, p[1:], v)
	}

}

func (p *CsvToJsonProcess) Run(reader Reader, encoder *json.Encoder, errCh chan<- error) {
	errCh <- func() (err error) {
		defer reader.Close()

		baseObject := p.BaseObject
		stringsOnly := p.StringsOnly

		// open the reader
		paths := map[string][]string{}
		for _, k := range reader.Header() {
			paths[k] = strings.Split(k, ".")
		}

		if baseObject != "" {
			if _, ok := paths[baseObject]; !ok {
				return fmt.Errorf("fatal: '%s' is not a valid key in the input stream", baseObject)
			}
		}

		// see: http://stackoverflow.com/questions/13340717/json-numbers-regular-expression
		numberMatcher := regexp.MustCompile("^ *-?(?:0|[1-9]\\d*)(?:\\.\\d+)?(?:[eE][+-]?\\d+)? *$")

		for data := range reader.C() {
			dataMap := data.AsMap()
			objectMap := map[string]interface{}{}

			if baseObject != "" {
				if base, ok := dataMap[baseObject]; ok {
					if err := json.Unmarshal([]byte(base), &objectMap); err != nil {
						fmt.Fprintf(os.Stderr, "warning: failed to parse base object: %s: %s\n", base, err)
					}
				}
			}

			for k, v := range dataMap {
				var f float64
				var ov interface{}

				ov = v

				if baseObject != "" && k == baseObject {
					continue
				} else if v == "" {
					continue
				} else if stringsOnly {
					ov = v
				} else if v == "null" {
					continue
				} else if v == "true" || v == "TRUE" {
					ov = true
				} else if v == "false" || v == "FALSE" {
					ov = false
				} else if v[0] == '{' {
					j := map[string]interface{}{}
					if err := json.Unmarshal([]byte(v), &j); err == nil {
						ov = j
					}
				} else if v[0] == '[' {
					aj := make([]interface{}, 0)
					if err := json.Unmarshal([]byte(v), &aj); err == nil {
						ov = aj
					}
				} else if numberMatcher.MatchString(v) {
					if _, err := fmt.Sscanf(v, "%f", &f); err == nil {
						ov = f
					}
				}
				p.writeToMap(objectMap, paths[k], ov)
			}
			encoder.Encode(objectMap)
		}
		return reader.Error()
	}()
}
