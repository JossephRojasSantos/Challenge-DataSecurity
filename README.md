# Challenge Data Security

## Resumen:

Se ha desarrollado una API en lenguaje Go que tiene como función principal la enumeración de los archivos almacenados en Google Drive. Posteriormente, esta API permite acceder y modificar los metadatos de estos archivos en función de la criticidad asignada por el usuario de la cuenta. 

![](https://github.com/JossephRojasSantos/Challenge-DataSecurity/blob/master/png/Diagrama%20-%20Arquitectura.png)


Video o Gif
### Pasos Iniciales:
------
#### Habilitar Entorno Google Cloud

1. Crear proyecto -> [Proyecto](https://console.cloud.google.com/projectcreate?previousPage=%2Fapis%2Fdashboard%3Fhl%3Des-419%26project%3Dchallenge-395013&organizationId=0&hl=es-419)
2. Habilitar API de Google Drive -> [API](https://console.cloud.google.com/apis/library/drive.googleapis.com?hl=es-419&organizationId=0&project=challenge-395013)
3. Configura la pantalla de consentimiento de OAuth y agrégate como usuario de prueba -> [Consentimiento](https://console.cloud.google.com/apis/credentials/consent?hl=es-419)
4. Generar credenciales para aplicación de escritorio de tipo **ID de cliente de OAuth** -> [OAuth](https://console.cloud.google.com/apis/credentials?hl=es-419) 
5. Descargar credencial **[client_secret]**
#### Entorno Windows con Docker

1. Descargar e instalar MySQL ->  [MySQL](https://downloads.mysql.com/archives/get/p/23/file/mysql-8.0.33-winx64-debug-test.zip)
2. Confirmar puerto de servicio **[port]**
3. Crear Base de Datos **[dbname]**
4. Generar usuario de lectura y escritura en la base de datos *[dbname]* creada en el punto 3 **[user][passroot]**
5. Descargar e instalar Docker -> [Docker](https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe)
6. Modificar la sección **bind-address** del archivo **my.ini** de la base de datos, ingresando la IP de Origen definida en el Contenedor.
7. Reiniciar servicio de MySQL. 
8. Ingresar los siguientes datos de variables de entorno **ENV** ubicadas en el archivo **Dockerfile**

* ENV host=[IP del equipo donde se ejecuta el servicio MySQL]
* ENV port=[port]-> Definido en el punto 2
* ENV user=[user]-> Definido en el punto 4
* ENV passworddb=[passroot]-> Definido en el punto 4
* ENV dbname=[dbname]-> Definido en el punto 3
* ENV client_secret=[client_secret]-> Ubicación de la llave en formato JSON, archivo descargado una vez configurado el entorno de Google Cloud.

![](https://github.com/JossephRojasSantos/Challenge-DataSecurity/blob/master/png/ENVDockerFile.png)

8. Descargamos el presente repositorio, nos ubicamos con un CMD en el proyecto e ingresamos los siguientes comandos:
```docker
docker build -t servidor:Challenge .
```
```docker
docker run -p 8080:8080 -p 443:443 -p 80:80 servidor:Challenge
```
9. Ingresamos en Docker Desktop y verificamos que en **Containers** nuestra imagen tenga estado **Running**
10. Ingresamos por medio de un navegador a **https://localhost/**
------
#### Entorno Windows sin Docker

1. Descargar e instalar MySQL ->  [MySQL](https://downloads.mysql.com/archives/get/p/23/file/mysql-8.0.33-winx64-debug-test.zip)
2. Confirmar puerto de servicio **[port]**
3. Crear Base de Datos **[dbname]**
4. Generar usuario de lectura y escritura en la base de datos **[dbname]** creada en el punto 3 **[user][pass]**
5. Descargar e instalar Go -> [GO](https://dl.google.com/go/go1.20.4.windows-amd64.msi)
6. Comprobar instalación de Go (desde un CMD ejecutar "go version", retorno de la consola -> go version go1.20.4 windows/amd64)
7. Crear las siguiente variables de entorno:

* host=localhost
* port=[port]-> Definido en el punto 2
* user=[user]-> Definido en el punto 4
* passworddb=[pass]-> Definido en el punto 4
* dbname=[dbname]-> Definido en el punto 3
* passwordadmin=[passadmin]-> Contraseña en sha512
* jwtkey=[jwtkey]-> Contraseña para la firma de token de sesión
* Changepass=1 -> Cuando se encuentra con valor 0, cambia la contraseña de passwordadmin 


![](https://github.com/JossephRojasSantos/Challenge-DataSecurity/blob/master/png/Variables%20de%20Entorno.png)

8. Descargamos el presente repositorio, nos ubicamos con un CMD en el proyecto e ingresamos el siguiente comando:
```go
go build main.go
```
9. Dentro de la carpeta del proyecto, ubicamos y ejecutamos como administrador el archivo **main.exe**.    
10. Ingresamos por medio de un navegador a **https://localhost/**

**Nota:** Es necesario ejercutar el binario como administrador (Punto 9) para obterner los valores de las variables de entorno configuradas en windows

### Modo de Operación 

------
#### Administrador

Rol con capacidad para crear o eliminar usuarios que accedan o consuman el aplicativo. El presente usuario accede por la Web con las siguientes credenciales:

* Usuario=administrator
* Contraseña= (Definida en las variables de entorno como **passwordadmin**)
* OTP= (Extraer de la base de datos **userdata**, columna **tokenmfa**)

------
#### Usuario

Rol de solo lectura con capacidad para visualizar los datos externos consumidos. El presente usuario accede por la Web con las siguientes credenciales:

* Usuario= (Definido e informado por el administrador)
* Contraseña= (Definido e informado por el administrador)
* OTP= (Informado por el administrador)

------
#### Desarrollador o API

Rol de solo lectura con capacidad para visualizar los datos externos consumidos. Este usuario obtiene los datos en formato JSON bajo la siguiente URL:

```html
https://localhost:8080/json?username=[usuario]&token=[token]
```

* [token]= (Método de autenticación (JWT) informado por el administrador, utiliza el parámetro jwtkey definida en las variable de entorno)
* [usuario]= (Datos de usuario a consultar)
