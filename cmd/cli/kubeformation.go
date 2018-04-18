// A CLI tool that can read cluster spec from a file and write the rendered
// templates to a directory.
package main

import (
	"github.com/hasura/kubeformation/pkg/cmd"
	log "github.com/sirupsen/logrus"
)

// main is the entrypoint function
func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
