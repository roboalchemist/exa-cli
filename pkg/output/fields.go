package output

import (
	"encoding/json"
	"strings"
)

// FilterFields filters JSON output to only include requested fields.
// For objects with a "results" array, applies field filtering to each result.
// Works on both single objects and arrays of objects.
func FilterFields(data interface{}, fields string) interface{} {
	if fields == "" {
		return data
	}

	fieldList := make(map[string]bool)
	for _, f := range strings.Split(fields, ",") {
		fieldList[strings.TrimSpace(f)] = true
	}

	// Convert to generic map via JSON round-trip
	raw, err := json.Marshal(data)
	if err != nil {
		return data
	}

	// Try array first
	var arr []map[string]interface{}
	if err := json.Unmarshal(raw, &arr); err == nil {
		result := make([]map[string]interface{}, len(arr))
		for i, item := range arr {
			result[i] = filterMap(item, fieldList)
		}
		return result
	}

	// Try single object — check for "results" array and filter its items
	var obj map[string]interface{}
	if err := json.Unmarshal(raw, &obj); err == nil {
		if results, ok := obj["results"]; ok {
			if resultsArr, ok := results.([]interface{}); ok {
				filtered := make([]interface{}, len(resultsArr))
				for i, item := range resultsArr {
					if m, ok := item.(map[string]interface{}); ok {
						filtered[i] = filterMap(m, fieldList)
					} else {
						filtered[i] = item
					}
				}
				return filtered
			}
		}
		// Fall back to top-level field filtering
		return filterMap(obj, fieldList)
	}

	return data
}

func filterMap(m map[string]interface{}, fields map[string]bool) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		if fields[k] {
			result[k] = v
		}
	}
	return result
}
