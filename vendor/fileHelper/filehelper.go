package filehelper

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	// Path for error files
	BasePath = "roofts/www"

	// html format
	HtmlFormat = "text/html"

	// ErrFilesPathVar is the name of the environment variable indicating
	// the location on disk of files served by the handler.
	ErrFilesPathVar = "ERROR_FILES_PATH"

	// FormatHeader name of the header used to extract the format
	FormatHeader = "X-Format"

	// CodeHeader name of the header used as source of the HTTP status code to return
	CodeHeader = "X-Code"

		// ServiceName name of the header that contains the matched Service in the Ingress
		ServiceName = "X-Service-Name"
)

func pathByDomain(uri string) string {

	log.Printf("Requested URL: %s", uri)

	if strings.Contains(uri, PandaPeDomain) {
		return PandaPePath
	}
	return InfoJobsPath
}

func FilePath(r *http.Request) (int, string, string, string) {

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

	return code, ext, fmt.Sprintf("%v%v/%v%v", errFilesPath, pathByDomain(r.Header.Get(ServiceName)), code, ext), ext
}

func contentErrorMessage(file string, code int, ext string) (string ) {
	scode := strconv.Itoa(code)
	file = fmt.Sprintf("%v%cxx%v", strings.ReplaceAll(file, scode+ext, ""), scode[0], ext)
	return file
}

