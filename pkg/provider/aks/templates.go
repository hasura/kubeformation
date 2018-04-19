package aks

var azureDeployJinja = `{
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
          {{- range .NodePools }}
          {
            "name": "{{ .Name }}",
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
          }
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
  ]
}`

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
