package v1

import "testing"

func TestGetVersion(t *testing.T) {
	s := ClusterSpec{Version: version}
	specVersion := s.GetVersion()
	if specVersion != "v1" {
		t.Fatalf("getting version failed, expected '%s', got '%s'", version, specVersion)
	}
}

func TestRead(t *testing.T) {
	tt := []struct {
		name string
		data []byte
		err  error
	}{
		// TODO: add validation for the struct and test wrong formats
		{"valid v1 spec yaml", []byte(`version: v1`), nil},
		{"valid v1 spec json", []byte(`{"version": "v1"}`), nil},
		{"invalid v1 spec", []byte(`version: v2`), ErrSpecParseFailure},
		{"invalid yaml", []byte(`version: v1:`), ErrSpecParseFailure},
		{"empty data", []byte(``), ErrEmtpyInput},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := ClusterSpec{Version: version}
			_, err := s.Read(tc.data)
			if err != tc.err {
				t.Fatalf("expected error '%v', got '%v'", tc.err, err)
			}
		})
	}
}
