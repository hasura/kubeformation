package gke

var clusterJinja = `
resources:
- name: {{ "{{ properties['name'] }}" }}
  type: container.v1.cluster
  properties:
    zone: {{ "{{ properties['zone'] }}" }}
    cluster:
      name: {{ "{{ properties['name'] }}" }}
      initialClusterVersion: {{ .K8SVersion }}
      nodePools:
      {{- range .NodePools }}
      - name: {{ $.Name }}-{{ .Name }}
        version: {{ $.K8SVersion }}
        initialNodeCount: {{ .Size }}
        config:
          machineType: {{ .MachineType }}
          oauthScopes:
          - https://www.googleapis.com/auth/compute
          - https://www.googleapis.com/auth/devstorage.read_only
          - https://www.googleapis.com/auth/logging.write
          - https://www.googleapis.com/auth/monitoring
          imageType: {{ .ImageType }}
          labels:
          {{- range $key, $value := .Labels }}
            {{ $key }}: {{ $value }}
          {{- end -}}
      {{- end -}}
`

var clusterYaml = `
imports:
- path: cluster.jinja

resources:
- name:  {{ .Name }}
  type: cluster.jinja
  properties:
    zone: ZONE
    name: {{ .Name }}
`
