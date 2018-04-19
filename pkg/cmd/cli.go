package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var ErrInvalidUsage = errors.New("kubeformation: invalid command usage")

var rootCmd = &cobra.Command{
	Use:           "kubeformation",
	Short:         "Kubeformation can bootstrap your cloud provider specific template for Kubernetes",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runKubeformation,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the cli version",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(GetVersion())
		return nil
	},
}
var kfm Kubeformation

func init() {
	rootCmd.Flags().StringVarP(&kfm.InputFile, "file", "f", "", "cluster spec file to read (- to read from stdin)")
	rootCmd.MarkFlagRequired("file")
	rootCmd.Flags().StringVarP(&kfm.ProviderValue, "provider", "p", "", "managed kubernetes provider for which template has to be bootstrapped")
	rootCmd.Flags().StringVarP(&kfm.OutputDir, "output", "o", "", "output directory to write templates")
	rootCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(versionCmd)
}

func Execute() error {
	return rootCmd.Execute()
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
