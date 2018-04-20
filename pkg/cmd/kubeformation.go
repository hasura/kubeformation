package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hasura/kubeformation/pkg/provider"
	"github.com/hasura/kubeformation/pkg/spec"

	// Register v1 spec handler
	_ "github.com/hasura/kubeformation/pkg/spec/v1"
)

type Kubeformation struct {
	// Input Flags
	ProviderValue string
	InputFile     string
	OutputDir     string

	// Processed Data
	Data     []byte
	Provider provider.ProviderType

	// Rendered Files
	OutputFiles map[string][]byte

	// Handler to generate files
	Handler spec.VersionedSpecHandler
}

func (k *Kubeformation) ParseProvider() {
	switch k.ProviderValue {
	case "gke":
		k.Provider = provider.GKE
	case "aks":
		k.Provider = provider.AKS
	case "eks":
		k.Provider = provider.EKS
	default:
		k.Provider = provider.NOP
	}
}

func (k *Kubeformation) ParseInputFlags() error {
	var err error
	if k.InputFile == "-" {
		// read from stdin
		k.Data, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
	} else {
		// read from file
		k.Data, err = ioutil.ReadFile(k.InputFile)
		if err != nil {
			return err
		}
	}
	k.ParseProvider()
	return nil
}

func (k *Kubeformation) GetHandler() error {
	handler, err := spec.Read(k.Data)
	if err != nil {
		return err
	}
	k.Handler = handler
	return nil
}

func (k *Kubeformation) RenderOutputFiles() error {
	var err error
	k.OutputFiles, err = k.Handler.GenerateProviderTemplate(k.Provider)
	if err != nil {
		return err
	}
	return nil
}

func (k *Kubeformation) WriteFilesToDir() error {
	err := os.MkdirAll(k.OutputDir, os.ModePerm)
	if err != nil {
		return err
	}
	for name, data := range k.OutputFiles {
		if len(data) == 0 {
			continue
		}
		err := ioutil.WriteFile(filepath.Join(k.OutputDir, name), data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Kubeformation) PrintFiles() {
	for k, v := range k.OutputFiles {
		log.Printf("%s\n%s\n\n", k, string(v))
	}
}
