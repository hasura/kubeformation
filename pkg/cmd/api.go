package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

func getContextLogger(l *log.Logger, w http.ResponseWriter, r *http.Request) func(e error, m string, c int) {
	// FIXME: requestID is not set properly
	requestID, ok := r.Context().Value(0).(string)
	if !ok {
		requestID = "unknown"
	}
	return func(e error, m string, c int) {
		if e != nil {
			l.Println(requestID, "ERROR", m, e.Error())
			http.Error(w, m, c)
		} else {
			l.Println(requestID, m)
		}
	}
}

func RenderProviderTemplate(l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := getContextLogger(l, w, r)
		switch r.Method {
		case "POST":
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logger(err, "cannot read requests body", http.StatusBadRequest)
				return
			}

			var download bool
			dlQueryParam := r.URL.Query().Get("download")
			if len(dlQueryParam) > 0 {
				download, err = strconv.ParseBool(dlQueryParam)
				if err != nil {
					logger(err, "error reading query param as bool", http.StatusBadRequest)
					return
				}
			}

			out, err := runKubeformationAPI(data)
			if err != nil {
				logger(err, "error parsing body and generating template", http.StatusInternalServerError)
				return
			}

			response, err := json.Marshal(convertByteMapToString(out))
			if err != nil {
				logger(err, "cannot convert output to JSON", http.StatusInternalServerError)
				return
			}

			if download {
				zipFile, err := createZip(out)
				fmt.Println(zipFile)
				zipFileName := filepath.Base(zipFile)

				if err != nil {
					logger(err, "cannot convert output to zip", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Disposition", "attachment; filename="+zipFileName+".zip")
				w.Header().Set("Content-Type", "application/zip")
				http.ServeFile(w, r, zipFile)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)

		default:
			logger(errors.New("method not allowed"), "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func runKubeformationAPI(data []byte) (map[string][]byte, error) {
	var kfm Kubeformation
	var err error

	kfm.Data = data
	err = kfm.GetHandler()
	if err != nil {
		return nil, err
	}

	err = kfm.RenderOutputFiles()
	if err != nil {
		return nil, err
	}
	return kfm.OutputFiles, nil
}

func convertByteMapToString(m map[string][]byte) map[string]string {
	responseData := make(map[string]string)
	for f, d := range m {
		responseData[f] = string(d)
	}
	return responseData
}

func createZip(data map[string][]byte) (string, error) {
	f, err := ioutil.TempFile("", "kubeformation-")
	if err != nil {
		return "", err
	}
	defer f.Close()

	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()

	for name, content := range data {
		w, err := zipWriter.Create(name)
		if err != nil {
			return "", err
		}
		_, err = w.Write(content)
		if err != nil {
			return "", err
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}
