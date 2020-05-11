package main

import (
	. "contentcache"
	. "service"

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

	// Path for error files
	BasePath = "roofts/www"
)

var cache *ContentCache
var service *Service

func main() {
	log.Printf("Configuring server")

	service = new(Service)
	cache = new(ContentCache)

	service.AddService("ats", "/pandape")
	service.AddService("infojobs", "/infojobs")
	configureCache()

	http.HandleFunc("/", errorHandler())
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Launching server :8080")
	http.ListenAndServe(fmt.Sprintf(":8080"), nil)
	log.Printf("Stopping server :8080")
}

func configureCache() {
	var keyName string
	for i := range service.ServiceCollenction {
		keyName = service.FilePath(BasePath, 404, ".html", service.ServiceCollenction[i].ServiceName)
		cache.AddItem(keyName, keyName)
		keyName = service.FilePath(BasePath, 503, ".html", service.ServiceCollenction[i].ServiceName)
		cache.AddItem(keyName, keyName)
		keyName = service.FilePath(BasePath, 999, ".html", service.ServiceCollenction[i].ServiceName)
		cache.AddItem(keyName, keyName)

	}
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
		var fileDefault string

		start := time.Now()
		showRequestVariables(w, r)

		code, format, file, ext := service.FileDescriptor(BasePath, r.Header.Get(CodeHeader), r.Header.Get(FormatHeader), r.Header.Get(ServiceName))

		w.Header().Set(ContentType, format)
		w.WriteHeader(code)

		fileContent, err := cache.GetItemReader(file)
		if err != nil {
			log.Printf("unexpected error opening file: %v", err)

			genericErrorFile := service.AlternativeErrorMessage(file, code, ext)
			err = cache.AddItem(file, genericErrorFile)
			fileContent, err = cache.GetItemReader(file)
			if err != nil {
				log.Printf("unexpected error opening file: %v", err)

				code, format, fileDefault, ext = service.FileDescriptor(BasePath, "999", r.Header.Get(FormatHeader), r.Header.Get(ServiceName))
				fileContent, err = cache.GetItemReader(fileDefault)
				if err != nil {
					http.NotFound(w, r)
					return
				}
			}
		}
		log.Printf("serving custom error response for code %v and format %v from file %v", code, format, file)
		io.Copy(w, fileContent)

		duration := time.Now().Sub(start).Seconds()
		proto := strconv.Itoa(r.ProtoMajor)
		proto = fmt.Sprintf("%s.%s", proto, strconv.Itoa(r.ProtoMinor))
		requestCount.WithLabelValues(proto).Inc()
		requestDuration.WithLabelValues(proto).Observe(duration)
	}

}
