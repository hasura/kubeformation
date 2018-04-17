package cmd

import (
	"io/ioutil"
	"os"

	"github.com/apex/log"
	"github.com/hasura/kubeformation/pkg/kubeformation/spec"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	_ "github.com/hasura/kubeformation/pkg/kubeformation/spec/v1"
)

var (
	file      string
	provider  string
	outputDir string
)

var ErrInvalidUsage = errors.New("kubeformation: invalid command usage")

var rootCmd = &cobra.Command{
	Use:           "kubeformation",
	Short:         "Kubeformation can bootstrap your cloud provider specific template for Kubernetes",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var data []byte
		var err error
		if file == "-" {
			// read from stdin
			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
		} else {
			// read from file
			data, err = ioutil.ReadFile(file)
			if err != nil {
				return err
			}
		}
		handler, err := spec.Read(data)
		if err != nil {
			return err
		}
		log.Info(handler.GetVersion())
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "cluster spec file to read (- to read from stdin)")
	rootCmd.MarkFlagRequired("file")
	rootCmd.Flags().StringVarP(&provider, "provider", "p", "", "managed kubernetes provider for which template has to be bootstrapped")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "", "output directory to write templates")
}

func Execute() error {
	return rootCmd.Execute()
}
