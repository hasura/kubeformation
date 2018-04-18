package gke

// clusterJinja is a raw Go template string used to render cluster.jinja file
// for Google Deployment Manager.
var clusterJinja = `resources:
- name: {{ .Name }}
  type: container.v1.cluster
  properties:
    zone: {{"{{"}} properties['zone'] {{"}}"}}
    cluster:
      name: {{ .Name }}
      initialClusterVersion: "{{ .K8SVersion }}"
      nodePools:
      {{- range .NodePools }}
      - name: {{ $.Name }}-{{ .Name }}
        version: "{{ $.K8SVersion }}"
        initialNodeCount: {{ .Size }}
        config:
          oauthScopes:
          - https://www.googleapis.com/auth/compute
          - https://www.googleapis.com/auth/devstorage.read_only
          - https://www.googleapis.com/auth/logging.write
          - https://www.googleapis.com/auth/monitoring
          {{- if .MachineType }}
          machineType: {{ .MachineType }}
          {{- else }}
          machineType: n1-standard-1
          {{- end -}}
          {{- if .ImageType }}
          imageType: {{ .ImageType }}
          {{- else }}
          imageType: cos
          {{- end -}}
          {{- if .Labels }}
          labels:
          {{- range $key, $value := .Labels }}
            {{ $key }}: {{ $value }}
          {{- end -}}
          {{- end -}}
      {{- end -}}`

// clusterYaml is the raw Go template string used to render cluster.yaml file
var clusterYaml = `imports:
- path: cluster.jinja

resources:
- name:  {{ .Name }}
  type: cluster.jinja
  properties:
    zone: ZONE
    name: {{ .Name }}`
