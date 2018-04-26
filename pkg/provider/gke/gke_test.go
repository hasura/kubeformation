package gke

import (
	"reflect"
	"testing"

	"github.com/hasura/kubeformation/pkg/provider"
)

func TestGetType(t *testing.T) {
	s := Spec{}
	got := s.GetType()
	if got != provider.GKE {
		t.Fatalf("expected provider '%v', got '%v'", provider.GKE, got)
	}
}

func TestMarshalFiles(t *testing.T) {
	tt := []struct {
		spec *Spec
		name string
		data map[string][]byte
		err  error
	}{
		{
			NewDefaultSpec(),
			"valid default spec",
			map[string][]byte{
				"gke-cluster.jinja": []byte(`resources:
- name: gke-cluster
  type: container.v1.cluster
  properties:
    zone: {{ properties['zone'] }}
    cluster:
      name: gke-cluster
      initialClusterVersion: "1.9"
      nodePools:
      - name: gke-cluster-np-1
        version: "1.9"
        initialNodeCount: 1
        config:
          oauthScopes:
          - https://www.googleapis.com/auth/compute
          - https://www.googleapis.com/auth/devstorage.read_only
          - https://www.googleapis.com/auth/logging.write
          - https://www.googleapis.com/auth/monitoring
          machineType: n1-standard-1
          imageType: cos`),
				"gke-cluster.yaml": []byte(`imports:
- path: gke-cluster.jinja

resources:
- name:  gke-cluster
  type: gke-cluster.jinja
  properties:
    name: gke-cluster
    project: PROJECT
    zone: ZONE`),
			},
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			files, err := tc.spec.MarshalFiles()
			if err != tc.err {
				t.Fatalf("expected error '%v', got '%v'", tc.err, err)
			}
			if !reflect.DeepEqual(files, tc.data) {
				// t.Log("expected:")
				// for k, v := range tc.data {
				// 	t.Logf("%s:\n%s\n", k, string(v))
				// }
				// t.Log("got:")
				// for k, v := range files {
				// 	t.Logf("%s:\n%s\n", k, string(v))
				// }
				// TODO: print a diff
				t.Fatal("data mismatch")
			}
		})
	}
}
