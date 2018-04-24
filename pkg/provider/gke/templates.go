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
      {{- end -}}
{{- range .Volumes }}
- name: {{ .Name }}
  type: compute.v1.disk
  properties:
    zone: {{"{{"}} properties['zone'] {{"}}"}}
    sizeGb: {{ .SizeGB }}
    type: https://www.googleapis.com/compute/v1/projects/{{"{{"}} properties['project'] {{"}}"}}/zones/{{"{{"}} properties['zone'] {{"}}"}}/diskTypes/pd-ssd
{{- end -}}`

// clusterYaml is the raw Go template string used to render cluster.yaml file
var clusterYaml = `imports:
- path: gke-cluster.jinja

resources:
- name:  {{ .Name }}
  type: cluster.jinja
  properties:
    name: {{ .Name }}
    project: PROJECT
    zone: ZONE`

// persistentVolumeYaml is a raw Go template for writing volumes.yaml file which
// contains any PersistentVolume related info.
var persistentVolumeYaml = `{{- $volumeLength := sub (len .Volumes) }}
{{- range $i, $volume := .Volumes -}}
apiVersion: v1
kind: PersistentVolume
metadata:
  name: {{ .Name }}
spec:
  capacity:
    storage: {{ .SizeGB }}Gi
  accessModes:
    - ReadWriteOnce
  storageClassName: standard
  gcePersistentDisk:
    pdName: {{ .Name }}
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Name }}-claim
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .SizeGB }}Gi
  volumeName: {{ .Name }}
{{- if ne $i $volumeLength }}
---
{{ end -}}
{{- end -}}`
