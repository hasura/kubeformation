package aks

import (
	"reflect"
	"testing"

	"github.com/hasura/kubeformation/pkg/provider"
)

func TestGetType(t *testing.T) {
	s := Spec{}
	got := s.GetType()
	if got != provider.AKS {
		t.Fatalf("expected provider '%v', got '%v'", provider.AKS, got)
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
				"aks-deploy.json": []byte(`{
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
      "name": "aks-cluster",
      "properties": {
        "dnsPrefix": "[parameters('dnsNamePrefix')]",
        "agentPoolProfiles": [
          {
            "name": "np-1",
            "count": 1,
            "dnsPrefix": "[variables('agentsEndpointDNSNamePrefix')]",
            "vmSize": "Standard_D2_v2",
            "osType": "Linux"
          }
        ],
        "kubernetesVersion": "1.8.1",
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
}`),
				"aks-params.json": []byte(`{
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
}`),
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
