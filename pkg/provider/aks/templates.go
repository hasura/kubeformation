package aks

// azureDeployJSON is a raw Go template string used to render azuredeploy.json
// file for Azure Resource Manager.
var azureDeployJSON = `{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "dnsNamePrefix": {
      "metadata": {
        "description": "Sets the Domain name prefix for the cluster.  The concatenation of the domain name and the regionalized DNS zone make up the fully qualified domain name associated with the public IP address."
      },
      "type": "string"
    },
    "adminUsername": {
      "metadata": {
        "description": "User name for the Linux Virtual Machines."
      },
      "type": "string"
    },
    "sshRSAPublicKey": {
      "metadata": {
        "description": "Configure all linux machines with the SSH RSA public key string.  Your key should include three parts, for example 'ssh-rsa AAAAB...snip...UcyupgH azureuser@linuxvm'"
      },
      "type": "string"
    },
    "servicePrincipalClientId": {
      "metadata": {
        "description": "Client ID (used by cloudprovider)"
      },
      "type": "securestring",
      "defaultValue": "n/a"
    },
    "servicePrincipalClientSecret": {
      "metadata": {
        "description": "The Service Principal Client Secret."
      },
      "type": "securestring",
      "defaultValue": "n/a"
    }
  },
  "variables": {
    "agentsEndpointDNSNamePrefix":"[concat(parameters('dnsNamePrefix'),'agents')]"
  },
  "resources": [
    {
      "apiVersion": "2017-08-31",
      "type": "Microsoft.ContainerService/managedClusters",
      "location": "[resourceGroup().location]",
      "name": "{{ .Name }}",
      "properties": {
        "dnsPrefix": "[parameters('dnsNamePrefix')]",
        "agentPoolProfiles": [
          {{- $nodeLength := sub (len .NodePools) }}
          {{- range $i, $node := .NodePools }}
          {
            "name": "{{ $.Name }}-{{ .Name }}",
            "count": {{ .Count }},
            "dnsPrefix": "[variables('agentsEndpointDNSNamePrefix')]",
            {{- if .VMSize }}
            "vmSize": "{{ .VMSize }}",
            {{- else }}
            "vmSize": "Standard_D2_v2",
            {{- end -}}
            {{- if .OSType }}
            "osType": "{{ .OSType }}"
            {{- else }}
            "osType": "Linux"
            {{- end }}
          {{- if ne $i $nodeLength }}
          },
          {{- else }}
          }
          {{- end }}
          {{- end }}
        ],
        "kubernetesVersion": "{{ .K8SVersion }}",
        "linuxProfile": {
          "adminUsername": "[parameters('adminUsername')]",
          "ssh": {
            "publicKeys": [
              {
                "keyData": "[parameters('sshRSAPublicKey')]"
              }
            ]
          }
        },
        "servicePrincipalProfile": {
          "clientId": "[parameters('servicePrincipalClientId')]",
          "secret": "[parameters('servicePrincipalClientSecret')]"
        }
      }
    }
  {{- if not .Volumes }}
  ]
  {{- else }}
  ],
  "outputs": {
    "azure_disk_command": {
      "type": "string",
      "value": "[concat('Run az group deployment create -n ', '{{ .Name }}-disks ', '-g ', toUpper(concat('MC_', resourceGroup().name, '_', '{{ .Name }}', '_', resourceGroup().location, ' ')), '--template-file azureDisk.json')]"
    }
  }
  {{- end }}
}`

// parametersJSON is a raw Go template string used to render azuredeploy.parameters.json
var parametersJSON = `{
  "$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentParameters.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "dnsNamePrefix": {
      "value": "GEN-UNIQUE"
    },
    "adminUsername": {
      "value": "GEN-UNIQUE"
    },
    "sshRSAPublicKey": {
      "value": "GEN-SSH-PUB-KEY"
    },
    "servicePrincipalClientId": {
      "value": "GEN-UNIQUE"
    },
    "servicePrincipalClientSecret": {
      "value": "GEN-UNIQUE"
    }
  }
}`

// azureDisksJSON is a raw Go template for writing azureDisk.json file which
// contains any disk related info.
var azureDisksJSON = `{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "resources": [
    {{- $volumeLength := sub (len .Volumes) }}
    {{- range $i, $volume := .Volumes }}
    {
      "apiVersion": "2017-03-30",
      "type": "Microsoft.Compute/disks",
      "location": "[resourceGroup().location]",
      "name": "{{ .Name }}",
      "properties": {
        "creationData": {
          "createOption": "Empty"
        },
        "diskSizeGB": {{ .SizeGB }}
      }
    {{- if ne $i $volumeLength }}
    },
    {{- else }}
    }
    {{- end }}
    {{- end }}
  ]
}`

// persistentVolumeJinja is a raw Go template for writing volumes.yaml file which
// contains any PersistentVolume related info.
var persistentVolumeJinja = `{{- $volumeLength := sub (len .Volumes) }}
{{- range $i, $volume := .Volumes -}}
apiVersion: v1
kind: PersistentVolume
metadata:
  name: {{ .Name }}
spec:
  capacity:
    storage: {{ .SizeGB }}G
  accessModes:
    - ReadWriteOnce
  azureDisk:
    kind: Managed
    diskName: {{ .Name }}
    diskURI: /subscriptions/SUBSCRIPTIONID/resourceGroups/GROUPNAME/providers/Microsoft.Compute/disks/{{ .Name }}
{{- if ne $i $volumeLength }}
---
{{ end -}}
{{- end -}}`
