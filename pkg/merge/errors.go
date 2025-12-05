package merge

import "github.com/cockroachdb/errors"

var (
	// ErrMergeTypeMismatch indicates org and repo file types are incompatible
	ErrMergeTypeMismatch = errors.New(
		"merge type mismatch: org and repo file types incompatible",
	)
	// ErrMergeParseError indicates a failure to parse file for merge
	ErrMergeParseError = errors.New("failed to parse file for merge")
	// ErrMergeUnsupportedFileType indicates merge only supports JSON and YAML files
	ErrMergeUnsupportedFileType = errors.New("merge only supports JSON and YAML files")
	// ErrMergeUnknownStrategy indicates an unknown merge strategy was specified
	ErrMergeUnknownStrategy = errors.New("unknown merge strategy")
)
