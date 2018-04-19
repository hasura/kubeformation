package provider

import "testing"

func TestParse(t *testing.T) {
	tt := []struct {
		name   string
		input  string
		output ProviderType
	}{
		{"valid aks", "aks", AKS},
		{"valid gke", "gke", GKE},
		{"valid eks", "eks", EKS},
		{"invalid empty", "", NOP},
		{"invalid random value", "randomvalue", NOP},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := Parse(tc.input)
			if got != tc.output {
				t.Fatalf("expected '%d', got '%d'", tc.output, got)
			}
		})
	}
}
