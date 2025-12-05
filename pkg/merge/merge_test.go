package merge_test

import (
	"testing"

	"github.com/smykla-labs/.github/pkg/config"
	"github.com/smykla-labs/.github/pkg/merge"
)

func TestDeepMerge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		base     map[string]any
		override map[string]any
		want     map[string]any
		wantErr  bool
	}{
		{
			name: "basic object merge",
			base: map[string]any{
				"name": "org",
				"age":  30,
			},
			override: map[string]any{
				"age":  35,
				"city": "NYC",
			},
			want: map[string]any{
				"name": "org",
				"age":  float64(35),
				"city": "NYC",
			},
		},
		{
			name: "nested object merge",
			base: map[string]any{
				"config": map[string]any{
					"timeout": 30,
					"retries": 3,
				},
			},
			override: map[string]any{
				"config": map[string]any{
					"timeout": 60,
					"verbose": true,
				},
			},
			want: map[string]any{
				"config": map[string]any{
					"timeout": float64(60),
					"retries": float64(3),
					"verbose": true,
				},
			},
		},
		{
			name: "array replacement not merge",
			base: map[string]any{
				"items": []any{"a", "b", "c"},
			},
			override: map[string]any{
				"items": []any{"x", "y"},
			},
			want: map[string]any{
				"items": []any{"x", "y"},
			},
		},
		{
			name: "null value deletes key",
			base: map[string]any{
				"keep":   "value",
				"remove": "value",
			},
			override: map[string]any{
				"remove": nil,
			},
			want: map[string]any{
				"keep": "value",
			},
		},
		{
			name: "nested null deletion",
			base: map[string]any{
				"config": map[string]any{
					"keep":   "value",
					"remove": "value",
				},
			},
			override: map[string]any{
				"config": map[string]any{
					"remove": nil,
				},
			},
			want: map[string]any{
				"config": map[string]any{
					"keep": "value",
				},
			},
		},
		{
			name:     "nil base map",
			base:     nil,
			override: map[string]any{"key": "value"},
			want:     map[string]any{"key": "value"},
		},
		{
			name:     "nil override map",
			base:     map[string]any{"key": "value"},
			override: nil,
			want:     map[string]any{"key": "value"},
		},
		{
			name:     "both nil maps",
			base:     nil,
			override: nil,
			want:     map[string]any{},
		},
		{
			name: "type override",
			base: map[string]any{
				"value": "string",
			},
			override: map[string]any{
				"value": float64(123),
			},
			want: map[string]any{
				"value": float64(123),
			},
		},
		{
			name: "deep nested merge",
			base: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": map[string]any{
							"keep": "base",
						},
					},
				},
			},
			override: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": map[string]any{
							"add": "override",
						},
					},
				},
			},
			want: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": map[string]any{
							"keep": "base",
							"add":  "override",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := merge.DeepMerge(tt.base, tt.override)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeepMerge() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !mapsEqual(got, tt.want) {
				t.Errorf("DeepMerge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShallowMerge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		base     map[string]any
		override map[string]any
		want     map[string]any
		wantErr  bool
	}{
		{
			name: "top level merge only",
			base: map[string]any{
				"name": "org",
				"age":  30,
			},
			override: map[string]any{
				"age":  35,
				"city": "NYC",
			},
			want: map[string]any{
				"name": "org",
				"age":  35,
				"city": "NYC",
			},
		},
		{
			name: "nested object replaced not merged",
			base: map[string]any{
				"config": map[string]any{
					"timeout": 30,
					"retries": 3,
				},
			},
			override: map[string]any{
				"config": map[string]any{
					"timeout": 60,
				},
			},
			want: map[string]any{
				"config": map[string]any{
					"timeout": 60,
				},
			},
		},
		{
			name: "null value deletes top level key",
			base: map[string]any{
				"keep":   "value",
				"remove": "value",
			},
			override: map[string]any{
				"remove": nil,
			},
			want: map[string]any{
				"keep": "value",
			},
		},
		{
			name:     "nil base map",
			base:     nil,
			override: map[string]any{"key": "value"},
			want:     map[string]any{"key": "value"},
		},
		{
			name:     "nil override map",
			base:     map[string]any{"key": "value"},
			override: nil,
			want:     map[string]any{"key": "value"},
		},
		{
			name: "array replacement",
			base: map[string]any{
				"items": []any{"a", "b"},
			},
			override: map[string]any{
				"items": []any{"x"},
			},
			want: map[string]any{
				"items": []any{"x"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := merge.ShallowMerge(tt.base, tt.override)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShallowMerge() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !mapsEqual(got, tt.want) {
				t.Errorf("ShallowMerge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeJSON(t *testing.T) {
	t.Parallel()

	base := map[string]any{
		"name": "org",
		"config": map[string]any{
			"timeout": 30,
		},
	}

	override := map[string]any{
		"config": map[string]any{
			"retries": 3,
		},
	}

	tests := []struct {
		name     string
		strategy config.MergeStrategy
		want     map[string]any
		wantErr  bool
	}{
		{
			name:     "deep merge strategy",
			strategy: config.MergeStrategyDeep,
			want: map[string]any{
				"name": "org",
				"config": map[string]any{
					"timeout": float64(30),
					"retries": float64(3),
				},
			},
		},
		{
			name:     "overlay strategy alias",
			strategy: config.MergeStrategyOverlay,
			want: map[string]any{
				"name": "org",
				"config": map[string]any{
					"timeout": float64(30),
					"retries": float64(3),
				},
			},
		},
		{
			name:     "shallow merge strategy",
			strategy: config.MergeStrategyShallow,
			want: map[string]any{
				"name": "org",
				"config": map[string]any{
					"retries": 3,
				},
			},
		},
		{
			name:     "unknown strategy",
			strategy: config.MergeStrategy("invalid"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := merge.MergeJSON(base, override, tt.strategy)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeJSON() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && !mapsEqual(got, tt.want) {
				t.Errorf("MergeJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeYAML(t *testing.T) {
	t.Parallel()

	base := map[string]any{
		"name": "org",
	}

	override := map[string]any{
		"age": 30,
	}

	got, err := merge.MergeYAML(base, override, config.MergeStrategyDeep)
	if err != nil {
		t.Fatalf("MergeYAML() error = %v", err)
	}

	want := map[string]any{
		"name": "org",
		"age":  float64(30),
	}

	if !mapsEqual(got, want) {
		t.Errorf("MergeYAML() = %v, want %v", got, want)
	}
}

func TestParseJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "valid json",
			data: `{"name":"test","age":30}`,
			want: map[string]any{
				"name": "test",
				"age":  float64(30),
			},
		},
		{
			name:    "invalid json",
			data:    `{invalid}`,
			wantErr: true,
		},
		{
			name: "empty json object",
			data: `{}`,
			want: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := merge.ParseJSON([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSON() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && !mapsEqual(got, tt.want) {
				t.Errorf("ParseJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "valid yaml",
			data: "name: test\nage: 30\n",
			want: map[string]any{
				"name": "test",
				"age":  30,
			},
		},
		{
			name:    "invalid yaml",
			data:    "name: test\n  invalid: indent\n",
			wantErr: true,
		},
		{
			name: "empty yaml",
			data: "",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := merge.ParseYAML([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseYAML() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && !mapsEqual(got, tt.want) {
				t.Errorf("ParseYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"name": "test",
		"age":  float64(30),
	}

	result, err := merge.MarshalJSON(data)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	// Parse back to verify
	parsed, err := merge.ParseJSON(result)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if !mapsEqual(parsed, data) {
		t.Errorf("Round trip failed: got %v, want %v", parsed, data)
	}
}

func TestMarshalYAML(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"name": "test",
		"age":  30,
	}

	result, err := merge.MarshalYAML(data)
	if err != nil {
		t.Fatalf("MarshalYAML() error = %v", err)
	}

	// Parse back to verify
	parsed, err := merge.ParseYAML(result)
	if err != nil {
		t.Fatalf("ParseYAML() error = %v", err)
	}

	if !mapsEqual(parsed, data) {
		t.Errorf("Round trip failed: got %v, want %v", parsed, data)
	}
}

// mapsEqual compares two maps deeply.
func mapsEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}

	for k, va := range a {
		vb, ok := b[k]
		if !ok {
			return false
		}

		if !valuesEqual(va, vb) {
			return false
		}
	}

	return true
}

// valuesEqual compares two values deeply.
func valuesEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	switch va := a.(type) {
	case map[string]any:
		vb, ok := b.(map[string]any)
		if !ok {
			return false
		}

		return mapsEqual(va, vb)
	case []any:
		vb, ok := b.([]any)
		if !ok {
			return false
		}

		if len(va) != len(vb) {
			return false
		}

		for i := range va {
			if !valuesEqual(va[i], vb[i]) {
				return false
			}
		}

		return true
	default:
		return a == b
	}
}

func TestDeepMerge_ComplexScenarios(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		base     map[string]any
		override map[string]any
		want     map[string]any
		wantErr  bool
	}{
		{
			name: "mixed types coexistence",
			base: map[string]any{
				"string": "value",
				"number": float64(42),
				"bool":   true,
				"array":  []any{"a", "b"},
				"object": map[string]any{"key": "value"},
			},
			override: map[string]any{
				"string": "updated",
				"array":  []any{"x"},
			},
			want: map[string]any{
				"string": "updated",
				"number": float64(42),
				"bool":   true,
				"array":  []any{"x"},
				"object": map[string]any{"key": "value"},
			},
		},
		{
			name: "deeply nested null deletion",
			base: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": map[string]any{
							"keep":   "value",
							"remove": "value",
						},
					},
				},
			},
			override: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": map[string]any{
							"remove": nil,
						},
					},
				},
			},
			want: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": map[string]any{
							"keep": "value",
						},
					},
				},
			},
		},
		{
			name: "partial object override with null",
			base: map[string]any{
				"settings": map[string]any{
					"enabled":  true,
					"timeout":  float64(30),
					"retries":  float64(3),
					"endpoint": "https://api.example.com",
				},
			},
			override: map[string]any{
				"settings": map[string]any{
					"timeout":  float64(60),
					"endpoint": nil,
				},
			},
			want: map[string]any{
				"settings": map[string]any{
					"enabled": true,
					"timeout": float64(60),
					"retries": float64(3),
				},
			},
		},
		{
			name: "empty object override",
			base: map[string]any{
				"config": map[string]any{
					"key": "value",
				},
			},
			override: map[string]any{
				"config": map[string]any{},
			},
			want: map[string]any{
				"config": map[string]any{
					"key": "value",
				},
			},
		},
		{
			name: "empty array override",
			base: map[string]any{
				"items": []any{"a", "b", "c"},
			},
			override: map[string]any{
				"items": []any{},
			},
			want: map[string]any{
				"items": []any{},
			},
		},
		{
			name: "nested array with objects",
			base: map[string]any{
				"users": []any{
					map[string]any{"name": "alice", "age": float64(30)},
					map[string]any{"name": "bob", "age": float64(25)},
				},
			},
			override: map[string]any{
				"users": []any{
					map[string]any{"name": "charlie", "age": float64(35)},
				},
			},
			want: map[string]any{
				"users": []any{
					map[string]any{"name": "charlie", "age": float64(35)},
				},
			},
		},
		{
			name: "RFC 7396 example from spec",
			base: map[string]any{
				"a": "b",
				"c": map[string]any{
					"d": "e",
					"f": "g",
				},
			},
			override: map[string]any{
				"a": "z",
				"c": map[string]any{
					"f": nil,
				},
			},
			want: map[string]any{
				"a": "z",
				"c": map[string]any{
					"d": "e",
				},
			},
		},
		{
			name: "object to array type change",
			base: map[string]any{
				"value": map[string]any{"key": "val"},
			},
			override: map[string]any{
				"value": []any{"item1", "item2"},
			},
			want: map[string]any{
				"value": []any{"item1", "item2"},
			},
		},
		{
			name: "array to object type change",
			base: map[string]any{
				"value": []any{"item1", "item2"},
			},
			override: map[string]any{
				"value": map[string]any{"key": "val"},
			},
			want: map[string]any{
				"value": map[string]any{"key": "val"},
			},
		},
		{
			name: "multiple null deletions at different levels",
			base: map[string]any{
				"top":    "value",
				"remove": "this",
				"nested": map[string]any{
					"keep":   "value",
					"remove": "this",
					"deep": map[string]any{
						"keep":   "value",
						"remove": "this",
					},
				},
			},
			override: map[string]any{
				"remove": nil,
				"nested": map[string]any{
					"remove": nil,
					"deep": map[string]any{
						"remove": nil,
					},
				},
			},
			want: map[string]any{
				"top": "value",
				"nested": map[string]any{
					"keep": "value",
					"deep": map[string]any{
						"keep": "value",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := merge.DeepMerge(tt.base, tt.override)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeepMerge() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && !mapsEqual(got, tt.want) {
				t.Errorf("DeepMerge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShallowMerge_ComplexScenarios(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		base     map[string]any
		override map[string]any
		want     map[string]any
		wantErr  bool
	}{
		{
			name: "replace entire nested structure",
			base: map[string]any{
				"config": map[string]any{
					"database": map[string]any{
						"host": "localhost",
						"port": float64(5432),
					},
					"cache": map[string]any{
						"ttl": float64(300),
					},
				},
			},
			override: map[string]any{
				"config": map[string]any{
					"database": map[string]any{
						"host": "prod.example.com",
					},
				},
			},
			want: map[string]any{
				"config": map[string]any{
					"database": map[string]any{
						"host": "prod.example.com",
					},
				},
			},
		},
		{
			name: "mixed type replacements",
			base: map[string]any{
				"a": map[string]any{"nested": "value"},
				"b": []any{"item"},
				"c": "string",
				"d": float64(42),
			},
			override: map[string]any{
				"a": "now a string",
				"b": map[string]any{"now": "object"},
				"c": []any{"now", "array"},
				"d": true,
			},
			want: map[string]any{
				"a": "now a string",
				"b": map[string]any{"now": "object"},
				"c": []any{"now", "array"},
				"d": true,
			},
		},
		{
			name: "empty object replacement",
			base: map[string]any{
				"data": map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
			},
			override: map[string]any{
				"data": map[string]any{},
			},
			want: map[string]any{
				"data": map[string]any{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := merge.ShallowMerge(tt.base, tt.override)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShallowMerge() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && !mapsEqual(got, tt.want) {
				t.Errorf("ShallowMerge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJSON_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "malformed json missing brace",
			data:    `{"key": "value"`,
			wantErr: true,
		},
		{
			name:    "malformed json trailing comma",
			data:    `{"key": "value",}`,
			wantErr: true,
		},
		{
			name:    "json with comments",
			data:    `{"key": "value" /* comment */}`,
			wantErr: true,
		},
		{
			name:    "json array instead of object",
			data:    `["item1", "item2"]`,
			wantErr: true,
		},
		{
			name:    "json null value",
			data:    `null`,
			wantErr: false,
		},
		{
			name:    "json string value",
			data:    `"string"`,
			wantErr: true,
		},
		{
			name:    "json number value",
			data:    `123`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := merge.ParseJSON([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseYAML_ErrorCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "invalid indentation",
			data:    "key: value\n  nested: bad",
			wantErr: true,
		},
		{
			name:    "invalid yaml syntax",
			data:    "key: : value",
			wantErr: true,
		},
		{
			name:    "tabs in yaml",
			data:    "key:\t\tvalue",
			wantErr: false,
		},
		{
			name:    "yaml array",
			data:    "- item1\n- item2",
			wantErr: true,
		},
		{
			name:    "yaml null",
			data:    "~",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := merge.ParseYAML([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMergeJSON_RealWorldRenovateExample(t *testing.T) {
	t.Parallel()

	// Simulate a real-world renovate.json merge scenario
	orgTemplate := map[string]any{
		"$schema":             "https://docs.renovatebot.com/renovate-schema.json",
		"extends":             []any{"config:base"},
		"rebaseWhen":          "behind-base-branch",
		"dependencyDashboard": true,
		"packageRules": []any{
			map[string]any{
				"matchUpdateTypes": []any{"major"},
				"automerge":        false,
			},
		},
	}

	repoOverride := map[string]any{
		"rebaseWhen": "conflicted",
		"automerge":  true,
	}

	result, err := merge.MergeJSON(orgTemplate, repoOverride, config.MergeStrategyDeep)
	if err != nil {
		t.Fatalf("MergeJSON() error = %v", err)
	}

	// Verify key fields
	if result["rebaseWhen"] != "conflicted" {
		t.Errorf("rebaseWhen = %v, want 'conflicted'", result["rebaseWhen"])
	}

	if result["automerge"] != true {
		t.Errorf("automerge = %v, want true", result["automerge"])
	}

	// Verify org fields are preserved
	if result["$schema"] != "https://docs.renovatebot.com/renovate-schema.json" {
		t.Errorf("$schema not preserved")
	}

	if result["dependencyDashboard"] != true {
		t.Errorf("dependencyDashboard not preserved")
	}
}

func TestMergeYAML_RealWorldGitHubActionsExample(t *testing.T) {
	t.Parallel()

	// Simulate a real-world GitHub Actions workflow merge scenario
	orgTemplate := map[string]any{
		"name": "CI",
		"on": map[string]any{
			"push": map[string]any{
				"branches": []any{"main"},
			},
		},
		"jobs": map[string]any{
			"test": map[string]any{
				"runs-on": "ubuntu-latest",
				"steps": []any{
					map[string]any{"uses": "actions/checkout@v3"},
					map[string]any{"uses": "actions/setup-go@v4"},
				},
			},
		},
	}

	repoOverride := map[string]any{
		"on": map[string]any{
			"push": map[string]any{
				"branches": []any{"main", "develop"},
			},
			"pull_request": map[string]any{},
		},
	}

	result, err := merge.MergeYAML(orgTemplate, repoOverride, config.MergeStrategyDeep)
	if err != nil {
		t.Fatalf("MergeYAML() error = %v", err)
	}

	// Verify name is preserved
	if result["name"] != "CI" {
		t.Errorf("name not preserved")
	}

	// Verify on triggers were merged
	on, ok := result["on"].(map[string]any)
	if !ok {
		t.Fatalf("on field is not a map")
	}

	if _, hasPR := on["pull_request"]; !hasPR {
		t.Errorf("pull_request trigger not added")
	}

	// Verify push branches were replaced (array replacement, not merge)
	push, isPush := on["push"].(map[string]any)
	if !isPush {
		t.Fatalf("push field is not a map")
	}

	branches, isBranches := push["branches"].([]any)
	if !isBranches {
		t.Fatalf("branches field is not an array")
	}

	if len(branches) != 2 {
		t.Errorf("branches length = %d, want 2", len(branches))
	}
}

func TestMergeJSON_NilHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		base     map[string]any
		override map[string]any
		strategy config.MergeStrategy
		wantErr  bool
	}{
		{
			name:     "nil base and nil override",
			base:     nil,
			override: nil,
			strategy: config.MergeStrategyDeep,
			wantErr:  false,
		},
		{
			name:     "nil base with valid override",
			base:     nil,
			override: map[string]any{"key": "value"},
			strategy: config.MergeStrategyShallow,
			wantErr:  false,
		},
		{
			name:     "valid base with nil override",
			base:     map[string]any{"key": "value"},
			override: nil,
			strategy: config.MergeStrategyDeep,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := merge.MergeJSON(tt.base, tt.override, tt.strategy)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarshalJSON_InvalidData(t *testing.T) {
	t.Parallel()

	// Create a map with a channel (which is not JSON-serializable)
	data := map[string]any{
		"channel": make(chan int),
	}

	_, err := merge.MarshalJSON(data)
	if err == nil {
		t.Error("MarshalJSON() expected error for non-serializable data, got nil")
	}
}

func TestMarshalYAML_ValidData(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"string":  "value",
		"number":  42,
		"boolean": true,
		"array":   []any{"a", "b", "c"},
		"object": map[string]any{
			"nested": "value",
		},
	}

	result, err := merge.MarshalYAML(data)
	if err != nil {
		t.Fatalf("MarshalYAML() error = %v", err)
	}

	// Verify we can parse it back
	parsed, err := merge.ParseYAML(result)
	if err != nil {
		t.Fatalf("ParseYAML() error = %v", err)
	}

	if !mapsEqual(parsed, data) {
		t.Errorf("Round trip failed: got %v, want %v", parsed, data)
	}
}

func TestMergeStrategies_Equivalence(t *testing.T) {
	t.Parallel()

	base := map[string]any{
		"a": "value-a",
		"nested": map[string]any{
			"b": "value-b",
		},
	}

	override := map[string]any{
		"nested": map[string]any{
			"c": "value-c",
		},
	}

	// Test that overlay and deep-merge are equivalent
	deepResult, err := merge.MergeJSON(base, override, config.MergeStrategyDeep)
	if err != nil {
		t.Fatalf("DeepMerge error = %v", err)
	}

	overlayResult, err := merge.MergeJSON(base, override, config.MergeStrategyOverlay)
	if err != nil {
		t.Fatalf("Overlay error = %v", err)
	}

	if !mapsEqual(deepResult, overlayResult) {
		t.Errorf("deep-merge and overlay strategies should be equivalent")
	}

	// Test that shallow-merge produces different result
	shallowResult, err := merge.MergeJSON(base, override, config.MergeStrategyShallow)
	if err != nil {
		t.Fatalf("ShallowMerge error = %v", err)
	}

	// In shallow merge, nested object should be completely replaced
	nestedShallow, ok := shallowResult["nested"].(map[string]any)
	if !ok {
		t.Fatal("nested should be a map")
	}

	if _, hasB := nestedShallow["b"]; hasB {
		t.Error("shallow merge should not preserve nested.b")
	}

	if nestedShallow["c"] != "value-c" {
		t.Error("shallow merge should have nested.c")
	}
}
