package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hasura/kubeformation/pkg/provider"
	"github.com/hasura/kubeformation/pkg/spec"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	// Register v1 spec handler
	_ "github.com/hasura/kubeformation/pkg/spec/v1"
)

var ErrInvalidUsage = errors.New("kubeformation: invalid command usage")

var rootCmd = &cobra.Command{
	Use:           "kubeformation",
	Short:         "Kubeformation can bootstrap your cloud provider specific template for Kubernetes",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runKubeformation,
}

var kfm Kubeformation

func init() {
	rootCmd.Flags().StringVarP(&kfm.InputFile, "file", "f", "", "cluster spec file to read (- to read from stdin)")
	rootCmd.MarkFlagRequired("file")
	rootCmd.Flags().StringVarP(&kfm.ProviderValue, "provider", "p", "", "managed kubernetes provider for which template has to be bootstrapped")
	rootCmd.Flags().StringVarP(&kfm.OutputDir, "output", "o", "", "output directory to write templates")
	rootCmd.MarkFlagRequired("output")
}

func Execute() error {
	return rootCmd.Execute()
}

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

func runKubeformation(cmd *cobra.Command, args []string) error {
	err := kfm.ParseInputFlags()
	if err != nil {
		return err
	}
	err = kfm.GetHandler()
	if err != nil {
		return err
	}
	err = kfm.RenderOutputFiles()
	if err != nil {
		return err
	}
	err = kfm.WriteFilesToDir()
	if err != nil {
		return err
	}
	return nil
}
