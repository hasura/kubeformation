// An API server to which cluster spec can be POST'ed as JSON, responds with
// rendered cloud provider templates. Also provides a download as ZIP endpoint.
package main

import (
	"net/http"

	"github.com/hasura/kubeformation/pkg/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	// http.HandleFunc("/", cmd.SayhelloName) // set router
	http.HandleFunc("/render", cmd.RenderProviderTemplate)
	err := http.ListenAndServe(":8081", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Info("listening...")
}
