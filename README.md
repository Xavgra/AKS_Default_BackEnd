# AKS_Default_BackEnd

### Descripción y funcionamiento

Aplicación destinada a retornar páginas de error en función de las cabeceras enviadas en una petición desde ingres de K8s cuando un servicon falla o no existe.

Cabeceras disponible:

* FormatHeader = "X-Format"
* CodeHeader name of the header used as source of the HTTP status code to return, CodeHeader = "X-Code"
* ContentType name of the header that defines the format of the reply, ContentType = "Content-Type"
* OriginalURI name of the header with the original URL from NGINX, OriginalURI = "X-Original-URI"
* Namespace name of the header that contains information about the Ingress namespace, Namespace = "X-Namespace"
* IngressName name of the header that contains the matched Ingress, IngressName = "X-Ingress-Name"
* ServiceName name of the header that contains the matched Service in the Ingress, ServiceName = "X-Service-Name"
* ServicePort name of the header that contains the matched Service port in the Ingress, ServicePort = "X-Service-Port"
* RequestId is a unique ID that identifies the request - same as for backend service, RequestId = "X-Request-ID"


### Docker

Generar la imagen Docker y arrancar

```system
docker-compose build
docker-compose up
Starting aks_default_backend_defaultbackend_1 ... done
Attaching to aks_default_backend_defaultbackend_1
defaultbackend_1  | 2020/05/06 14:02:41 Configuring server
defaultbackend_1  | 2020/05/06 14:02:41 Launching server :8080
```

### Testear 

Ejecutar una llamada via Postman o curl a http://localhost:80/index.htm y los headers siguientes:

* "X-Original-URI": "pandape.com.br/lskfjsdlfj sdl/lskjdflskdj f/lakshdklsad ",
* "X-Format": "text/html",
* "X-Code": "504",
* "X-Service-Name": "atswebui",

El log de la aplicación si no está corriendo en background mostrará lo siguiente:

```system
$ docker-compose up
Starting aks_default_backend_defaultbackend_1 ... done
Attaching to aks_default_backend_defaultbackend_1
defaultbackend_1  | 2020/05/06 14:02:41 Configuring server
defaultbackend_1  | 2020/05/06 14:02:41 Launching server :8080
defaultbackend_1  | 2020/05/06 14:08:36 couldn't get media type extension. Using html
defaultbackend_1  | 2020/05/06 14:08:36 serving custom error response for code 503 and format .html from file roofts/www/pandape/503.html
```

### Instalación en azure 

```system
az login
az account set --subscription YYY
az acr login --name XXX
docker tag defautlbackend:latest empbra.azurecr.io/ZZZZZZZ:20200506.1.master
docker push empbra.azurecr.io/ZZZZZZZ:20200506.1.master
```


export GO111MODULE=on