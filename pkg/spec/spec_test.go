package spec

import (
	"testing"
)

func TestReadVersion(t *testing.T) {
	tt := []struct {
		name    string
		data    []byte
		version string
		err     error
	}{
		{"valid yaml with v1", []byte(`version: v1`), "v1", nil},
		{"valid json with v1", []byte(`{"version": "v1"}`), "v1", nil},
		{"empty content", []byte(``), "", ErrInvalidSpecVersion},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			version, err := readVersion(tc.data)
			if err != tc.err {
				t.Fatalf("reading '%s', expected error '%v', got '%v'", string(tc.data), tc.err, err)
			}
			if version != tc.version {
				t.Fatalf("reading %s, expected version %s, got %s", string(tc.data), tc.version, version)
			}
		})
	}
}
