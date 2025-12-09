// Package merge provides strategies for merging organization and repository file contents.
package merge

import (
	"bytes"
	"encoding/json"
	"maps"

	"github.com/cockroachdb/errors"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/tidwall/pretty"
	"go.yaml.in/yaml/v4"

	"github.com/smykla-labs/.github/internal/configtypes"
)

const (
	// jsonPrettyWidth is the max column width for single-line arrays in JSON output.
	jsonPrettyWidth = 80
)

// MergeOptions configures array merge behavior for deep and shallow merges.
type MergeOptions struct {
	// ArrayStrategies maps JSONPath expressions (e.g., "$.packageRules") to merge strategies
	// (append/prepend/replace). Only exact string matches are supported.
	ArrayStrategies map[string]string
	// DeduplicateArrays removes duplicate elements from arrays after merging using deep equality.
	DeduplicateArrays bool
}

// mergeArrays combines two arrays according to the specified strategy.
// Strategies: replace (default), append (base+override), prepend (override+base).
// If deduplicate is true, removes duplicate elements using cmp.Equal for deep equality.
func mergeArrays(
	base, override []any,
	strategy string,
	deduplicate bool,
) []any {
	// Handle nil and empty arrays
	if override == nil {
		if base == nil {
			return []any{}
		}

		if deduplicate {
			return deduplicateArray(base)
		}

		return append([]any{}, base...)
	}

	if base == nil {
		if deduplicate {
			return deduplicateArray(override)
		}

		return append([]any{}, override...)
	}

	// Apply strategy
	var result []any

	switch configtypes.ArrayStrategy(strategy) {
	case configtypes.ArrayStrategyAppend:
		// base + override
		result = make([]any, 0, len(base)+len(override))
		result = append(result, base...)
		result = append(result, override...)
	case configtypes.ArrayStrategyPrepend:
		// override + base
		result = make([]any, 0, len(override)+len(base))
		result = append(result, override...)
		result = append(result, base...)
	case configtypes.ArrayStrategyReplace:
		fallthrough
	default:
		// replace (default for invalid strategy)
		result = append([]any{}, override...)
	}

	if deduplicate {
		return deduplicateArray(result)
	}

	return result
}

// deduplicateArray removes duplicate elements from an array, keeping the first occurrence.
// Uses cmp.Equal for deep equality comparison of objects.
func deduplicateArray(arr []any) []any {
	if len(arr) == 0 {
		return arr
	}

	seen := make([]any, 0, len(arr))

	for _, item := range arr {
		found := false

		for _, seenItem := range seen {
			if cmp.Equal(item, seenItem) {
				found = true

				break
			}
		}

		if !found {
			seen = append(seen, item)
		}
	}

	return seen
}

// collectArrays recursively walks a map and collects all arrays with their JSONPath expressions.
// Returns a map of JSONPath -> array slice.
func collectArrays(obj map[string]any, prefix string) map[string][]any {
	result := make(map[string][]any)

	for key, value := range obj {
		// Build JSONPath for this key
		path := prefix + "." + key

		switch v := value.(type) {
		case []any:
			// Found an array - add to result
			result[path] = v
		case map[string]any:
			// Recurse into nested objects
			nested := collectArrays(v, path)
			maps.Copy(result, nested)
		}
	}

	return result
}

// applyArrayStrategies applies array merge strategies to a merged map.
//
// Walks the merged result, finds arrays, matches them against arrayStrategies,
// and re-merges with the specified strategy.
//
//nolint:unused // Will be wired to DeepMerge/ShallowMerge in Phase 3
func applyArrayStrategies(
	merged, base, override map[string]any,
	opts *MergeOptions,
) error {
	if opts == nil || len(opts.ArrayStrategies) == 0 {
		return nil
	}

	// Collect arrays from base and override
	baseArrays := collectArrays(base, "$")
	overrideArrays := collectArrays(override, "$")

	// Apply strategies to matching arrays
	for path, strategy := range opts.ArrayStrategies {
		baseArray := baseArrays[path]
		overrideArray := overrideArrays[path]

		// Skip if neither base nor override has an array at this path
		if baseArray == nil && overrideArray == nil {
			continue
		}

		// Merge arrays with strategy
		mergedArray := mergeArrays(baseArray, overrideArray, strategy, opts.DeduplicateArrays)

		// Update the merged map at this path
		if err := setValueAtPath(merged, path, mergedArray); err != nil {
			return errors.Wrapf(err, "setting array at path %s", path)
		}
	}

	return nil
}

// setValueAtPath sets a value in a map at the specified JSONPath.
//
// Path format: "$.field.nested.array"
//
//nolint:unused // Will be called by applyArrayStrategies in Phase 3
func setValueAtPath(obj map[string]any, path string, value any) error {
	// Parse path - expect format "$.field.nested.field"
	if len(path) < 2 || path[0] != '$' || path[1] != '.' {
		return errors.Newf("invalid path format: %s", path)
	}

	// Remove "$."]
	pathStr := path[2:]

	// Split by "."
	parts := splitPath(pathStr)
	if len(parts) == 0 {
		return errors.New("empty path")
	}

	// Navigate to parent
	current := obj

	for i := range len(parts) - 1 {
		next, ok := current[parts[i]]
		if !ok {
			// Path doesn't exist - this can happen if the array wasn't in the merged result
			return nil
		}

		nextMap, ok := next.(map[string]any)
		if !ok {
			return errors.Newf("path segment %s is not a map", parts[i])
		}

		current = nextMap
	}

	// Set value at final key
	current[parts[len(parts)-1]] = value

	return nil
}

// splitPath splits a JSONPath string by dots, handling escaped dots.
// For MVP, we only support simple paths without arrays or wildcards.
func splitPath(path string) []string {
	if path == "" {
		return []string{}
	}

	var parts []string

	var current string

	for i := range len(path) {
		if path[i] == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(path[i])
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}

// DeepMerge recursively merges two maps using RFC 7396 JSON Merge Patch semantics.
//   - Nested objects are merged recursively
//   - Arrays are replaced, not merged
//   - Null values in override explicitly remove keys from base
//   - Type mismatches are handled per RFC 7396 (override wins)
func DeepMerge(base, override map[string]any) (map[string]any, error) {
	// Handle nil cases
	if base == nil && override == nil {
		return make(map[string]any), nil
	}

	if base == nil {
		base = make(map[string]any)
	}

	if override == nil {
		result := make(map[string]any, len(base))
		maps.Copy(result, base)

		return result, nil
	}

	// Convert maps to JSON
	baseJSON, err := json.Marshal(base)
	if err != nil {
		return nil, errors.Wrap(ErrMergeParseError, "marshaling base map to JSON")
	}

	overrideJSON, err := json.Marshal(override)
	if err != nil {
		return nil, errors.Wrap(
			ErrMergeParseError,
			"marshaling override map to JSON",
		)
	}

	// Apply RFC 7396 merge patch
	mergedJSON, err := jsonpatch.MergePatch(baseJSON, overrideJSON)
	if err != nil {
		return nil, errors.Wrap(ErrMergeParseError, "applying merge patch")
	}

	// Convert back to map
	var result map[string]any
	if err := json.Unmarshal(mergedJSON, &result); err != nil {
		return nil, errors.Wrap(
			ErrMergeParseError,
			"unmarshaling merged JSON to map",
		)
	}

	return result, nil
}

// ShallowMerge merges two maps at the top level only.
//   - Only top-level keys are merged
//   - Nested objects are replaced if overridden, not merged recursively
//   - Null values in override explicitly remove keys from base
func ShallowMerge(base, override map[string]any) (map[string]any, error) {
	if base == nil {
		base = make(map[string]any)
	}

	if override == nil {
		// Return a copy of base
		result := make(map[string]any, len(base))
		maps.Copy(result, base)

		return result, nil
	}

	// Create result with base values
	result := make(map[string]any, len(base))
	maps.Copy(result, base)

	// Apply override values at top level only
	for key, overrideVal := range override {
		if overrideVal == nil {
			// Explicit nil means delete the key
			delete(result, key)

			continue
		}

		// Replace with override value (no recursion)
		result[key] = overrideVal
	}

	return result, nil
}

// MergeJSON merges two JSON objects using the specified strategy.
func MergeJSON(
	base, override map[string]any,
	strategy configtypes.MergeStrategy,
) (map[string]any, error) {
	// Default to deep-merge if strategy is empty (not specified in config)
	if strategy == "" {
		strategy = configtypes.MergeStrategyDeep
	}

	switch strategy {
	case configtypes.MergeStrategyDeep, configtypes.MergeStrategyOverlay:
		return DeepMerge(base, override)
	case configtypes.MergeStrategyShallow:
		return ShallowMerge(base, override)
	default:
		return nil, errors.Wrapf(
			ErrMergeUnknownStrategy,
			"strategy: %q",
			strategy,
		)
	}
}

// MergeYAML merges two YAML objects using the specified strategy.
// YAML is converted to JSON internally, merged, then converted back.
func MergeYAML(
	base, override map[string]any,
	strategy configtypes.MergeStrategy,
) (map[string]any, error) {
	// YAML and JSON have compatible data models, so we can use the same merge logic
	return MergeJSON(base, override, strategy)
}

// ParseJSON parses JSON bytes into a map.
func ParseJSON(data []byte) (map[string]any, error) {
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, errors.Wrap(ErrMergeParseError, "parsing JSON")
	}

	return result, nil
}

// ParseYAML parses YAML bytes into a map.
func ParseYAML(data []byte) (map[string]any, error) {
	var result map[string]any
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, errors.Wrap(ErrMergeParseError, "parsing YAML")
	}

	return result, nil
}

// MarshalJSON converts a map to indented JSON bytes for readable config files.
// Uses SetEscapeHTML(false) to preserve <, >, & characters in regex patterns.
// Uses tidwall/pretty to keep short arrays on single lines.
func MarshalJSON(data map[string]any) ([]byte, error) {
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(data); err != nil {
		return nil, errors.Wrap(ErrMergeParseError, "marshaling to JSON")
	}

	// Use tidwall/pretty for formatting with compact short arrays
	opts := &pretty.Options{
		Width:    jsonPrettyWidth,
		Indent:   "  ",
		SortKeys: true,
	}

	result := pretty.PrettyOptions(buf.Bytes(), opts)

	// Trim trailing newline for consistency
	return bytes.TrimSuffix(result, []byte("\n")), nil
}

// MarshalYAML converts a map to YAML bytes.
func MarshalYAML(data map[string]any) ([]byte, error) {
	result, err := yaml.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(ErrMergeParseError, "marshaling to YAML")
	}

	return result, nil
}
