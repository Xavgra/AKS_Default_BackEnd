/*
Copyright 2017 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// FormatHeader name of the header used to extract the format
	FormatHeader = "X-Format"

	// CodeHeader name of the header used as source of the HTTP status code to return
	CodeHeader = "X-Code"

	// ContentType name of the header that defines the format of the reply
	ContentType = "Content-Type"

	// OriginalURI name of the header with the original URL from NGINX
	OriginalURI = "X-Original-URI"

	// Namespace name of the header that contains information about the Ingress namespace
	Namespace = "X-Namespace"

	// IngressName name of the header that contains the matched Ingress
	IngressName = "X-Ingress-Name"

	// ServiceName name of the header that contains the matched Service in the Ingress
	ServiceName = "X-Service-Name"

	// ServicePort name of the header that contains the matched Service port in the Ingress
	ServicePort = "X-Service-Port"

	// RequestId is a unique ID that identifies the request - same as for backend service
	RequestId = "X-Request-ID"

	// ErrFilesPathVar is the name of the environment variable indicating
	// the location on disk of files served by the handler.
	ErrFilesPathVar = "ERROR_FILES_PATH"

	// Dominio panda pe
	PandaPeDomain = "pandape.com.br"

	//Path
	PandaPePath = "/pandape"

	// Dominio InfoJobs
	InfoJobsDomain = "infojobs.com.br"

	//Path
	InfoJobsPath = "/infojobs"

	// VAriable
	DEFAULT_PAGE = "DEFAULT_PAGE"

	// Env variable
	ServiceNameMode = "ServiceName"

	//Env varaible
	DomainMode = "OriginalURI"

	// Path for error files
	BasePath = "roofts/www"

	// html format
	HtmlFormat = "text/html"
)

func main() {
	log.Printf("Configuring server")
	http.HandleFunc("/", errorHandler())
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Launching server :8080")
	http.ListenAndServe(fmt.Sprintf(":8080"), nil)
	log.Printf("Stopping server :8080")
}

func pathByDomain(uri string) string {

	if strings.Contains(uri, PandaPeDomain) {
		return PandaPePath
	}
	return InfoJobsPath
}
func pathAplication(r *http.Request) string {

	completePath := ""

	if os.Getenv(DEFAULT_PAGE) == ServiceNameMode {
		completePath = r.Header.Get(ServiceName) + "/"
	}

	if os.Getenv(DEFAULT_PAGE) == DomainMode {
		completePath = pathByDomain(r.Header.Get(OriginalURI))
	}
	return completePath
}
func filePath(r *http.Request) (int, string, string, string) {

	errFilesPath := BasePath
	if os.Getenv(ErrFilesPathVar) != "" {
		errFilesPath = os.Getenv(ErrFilesPathVar)
	}

	errCode := r.Header.Get(CodeHeader)
	code, err := strconv.Atoi(errCode)
	if err != nil {
		code = 404
		log.Printf("unexpected error reading return code: %v. Using %v", err, code)
	}

	format := r.Header.Get(FormatHeader)
	if format == "" {
		format = HtmlFormat
		log.Printf("format not specified. Using %v", format)
	}

	ext := "html"
	cext, err := mime.ExtensionsByType(format)
	if err != nil {
		log.Printf("unexpected error reading media type extension: %v. Using %v", err, ext)
		format = HtmlFormat
	} else if len(cext) == 0 {
		log.Printf("couldn't get media type extension. Using %v", ext)
	} else {
		ext = cext[0]
	}

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	return code, ext, fmt.Sprintf("%v%v/%v%v", errFilesPath, pathAplication(r), code, ext), ext
}
func contentErrorMessage(file string, code int, ext string) (*os.File, error) {

	f, err := os.Open(file)

	if err != nil {
		log.Printf("unexpected error opening file: %v", err)
		scode := strconv.Itoa(code)
		file = fmt.Sprintf("%v%cxx%v", strings.ReplaceAll(file, scode+ext, ""), scode[0], ext)
		f, err = os.Open(file)
		if err != nil {
			log.Printf("Error getting custom error response for code %v and from file %v", code, file)
		}
	}

	return f, err
}

func errorHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if os.Getenv("DEBUG") != "" {
			w.Header().Set(FormatHeader, r.Header.Get(FormatHeader))
			w.Header().Set(CodeHeader, r.Header.Get(CodeHeader))
			w.Header().Set(ContentType, r.Header.Get(ContentType))
			w.Header().Set(OriginalURI, r.Header.Get(OriginalURI))
			w.Header().Set(Namespace, r.Header.Get(Namespace))
			w.Header().Set(IngressName, r.Header.Get(IngressName))
			w.Header().Set(ServiceName, r.Header.Get(ServiceName))
			w.Header().Set(ServicePort, r.Header.Get(ServicePort))
			w.Header().Set(RequestId, r.Header.Get(RequestId))
		}

		code, format, file, ext := filePath(r)
		w.Header().Set(ContentType, format)
		w.WriteHeader(code)

		f, err := contentErrorMessage(file, code, ext)

		if err != nil {
			log.Printf("unexpected error opening file: %v", err)
			http.NotFound(w, r)
			return
		}

		defer f.Close()
		log.Printf("serving custom error response for code %v and format %v from file %v", code, format, file)
		io.Copy(w, f)

		duration := time.Now().Sub(start).Seconds()

		proto := strconv.Itoa(r.ProtoMajor)
		proto = fmt.Sprintf("%s.%s", proto, strconv.Itoa(r.ProtoMinor))

		requestCount.WithLabelValues(proto).Inc()
		requestDuration.WithLabelValues(proto).Observe(duration)
	}
}
