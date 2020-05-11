package service

import (
	"fmt"
	"log"
	"mime"
	"strconv"
	"strings"
)

const (
	HtmlFormat = "txtt/html"
)

type ServiceItem struct {
	ServiceName string
	path        string
}

type Service struct {
	ServiceCollenction []ServiceItem
}

func (service *Service) AddService(serviceName string, pathErrorFiles string) {
	serviceItem := new(ServiceItem)
	serviceItem.path = pathErrorFiles
	serviceItem.ServiceName = serviceName

	service.ServiceCollenction = append(service.ServiceCollenction, *serviceItem)
}

func (service *Service) pathByService(serviceName string) string {

	log.Printf("Generating path for serviceName: %s", serviceName)

	for i := range service.ServiceCollenction {
		if strings.Contains(service.ServiceCollenction[i].ServiceName, serviceName) {
			return service.ServiceCollenction[i].path
		}
	}

	return service.ServiceCollenction[0].path
}

func (service *Service) FilePath(errFilesPath string, code int, ext string, serviceName string) string {
	return fmt.Sprintf("%v%v/%v%v", errFilesPath, service.pathByService(serviceName), code, ext)
}

func (service *Service) FileDescriptor(errFilesPath string, errCode string, format string, serviceName string) (int, string, string, string) {

	code, err := strconv.Atoi(errCode)
	if err != nil {
		code = 404
		log.Printf("unexpected error reading return code: %v. Using %v", err, code)
	}

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

	filePath := service.FilePath(errFilesPath, code, ext, serviceName)

	return code, ext, filePath, ext
}

func (service *Service) AlternativeErrorMessage(file string, code int, ext string) string {
	scode := strconv.Itoa(code)
	file = fmt.Sprintf("%v%cxx%v", strings.ReplaceAll(file, scode+ext, ""), scode[0], ext)
	return file
}
