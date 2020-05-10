package main

import (
	. "contentcache"
	. "filehelper"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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



	// Dominio panda pe
	PandaPeDomain = "ats"

	//Path
	PandaPePath = "/pandape"

	// Dominio InfoJobs
	InfoJobsDomain = ""

	//Path
	InfoJobsPath = "/infojobs"

	// VAriable
	DEFAULT_PAGE = "DEFAULT_PAGE"



	//Env varaible
	DomainMode = "OriginalURI"

	// Path for error files
	BasePath = "roofts/www"

	
)

var cache *ContentCache

func main() {
	log.Printf("Configuring server")

	cache = new(ContentCache)
	err := cache.AddItem("ke1", "roofts/www/pandape/503.html")
	if err != nil {
		log.Printf("erro configurando cache")
	}

	http.HandleFunc("/", errorHandler())
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Launching server :8080")
	http.ListenAndServe(fmt.Sprintf(":8080"), nil)
	log.Printf("Stopping server :8080")
}

func showRequestVariables(w http.ResponseWriter, r *http.Request) {
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
}

func errorHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		showRequestVariables(w, r)

		code, format, file, ext := filePath(r)
		w.Header().Set(ContentType, format)
		w.WriteHeader(code)

		fileContent, err := cache.GetItemReader(file)
		if err != nil {

			log.Printf("unexpected error opening file: %v", err)
			http.NotFound(w, r)
			return
		}
		io.Copy(w, fileContent)

		duration := time.Now().Sub(start).Seconds()

		proto := strconv.Itoa(r.ProtoMajor)
		proto = fmt.Sprintf("%s.%s", proto, strconv.Itoa(r.ProtoMinor))

		requestCount.WithLabelValues(proto).Inc()
		requestDuration.WithLabelValues(proto).Observe(duration)

		f, err := contentErrorMessage(file, code, ext)

		if err != nil {
			log.Printf("unexpected error opening file: %v", err)
			http.NotFound(w, r)
			return
		}

		defer f.Close()
		log.Printf("serving custom error response for code %v and format %v from file %v", code, format, file)
		//io.Copy(w, f)

	}

}
