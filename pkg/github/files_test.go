package github

import (
	"bytes"
	"testing"
)

func TestRenderFileTemplate(t *testing.T) {
	tests := []struct {
		name          string
		content       []byte
		defaultBranch string
		want          []byte
	}{
		{
			name:          "replaces single placeholder",
			content:       []byte(`branches: ["{{DEFAULT_BRANCH}}"]`),
			defaultBranch: "main",
			want:          []byte(`branches: ["main"]`),
		},
		{
			name:          "replaces multiple placeholders",
			content:       []byte(`branch: {{DEFAULT_BRANCH}}, target: {{DEFAULT_BRANCH}}`),
			defaultBranch: "develop",
			want:          []byte(`branch: develop, target: develop`),
		},
		{
			name:          "returns unchanged when no placeholder",
			content:       []byte(`no placeholders here`),
			defaultBranch: "main",
			want:          []byte(`no placeholders here`),
		},
		{
			name:          "handles empty content",
			content:       []byte{},
			defaultBranch: "main",
			want:          []byte{},
		},
		{
			name:          "case sensitive - lowercase not replaced",
			content:       []byte(`{{default_branch}}`),
			defaultBranch: "main",
			want:          []byte(`{{default_branch}}`),
		},
		{
			name:          "handles empty default branch",
			content:       []byte(`branch: {{DEFAULT_BRANCH}}`),
			defaultBranch: "",
			want:          []byte(`branch: `),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderFileTemplate(tt.content, tt.defaultBranch)

			if !bytes.Equal(got, tt.want) {
				t.Errorf("renderFileTemplate() = %q, want %q", got, tt.want)
			}
		})
	}
}
