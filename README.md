# Challenge Data Security

## Resumen:

Se ha desarrollado una API en lenguaje Go que tiene como función principal la enumeración de los archivos almacenados en Google Drive. Posteriormente, esta API permite acceder y modificar los metadatos de estos archivos en función de la criticidad asignada por el usuario de la cuenta. 

![](https://github.com/JossephRojasSantos/Challenge-DataSecurity/blob/master/png/Diagrama%20-%20Arquitectura.png)

[<img src="https://github.com/JossephRojasSantos/Challenge-DataSecurity/blob/master/png/Portada.jpg" width="100%">](https://www.youtube.com/watch?v=yYExnrmTcaQ "Challenge Data Security")

### Pasos Iniciales:
------
#### Habilitar Entorno Google Cloud

1. Crear proyecto -> [Proyecto](https://console.cloud.google.com/projectcreate?previousPage=%2Fapis%2Fdashboard%3Fhl%3Des-419%26project%3Dchallenge-395013&organizationId=0&hl=es-419)
2. Habilitar API de Google Drive -> [API](https://console.cloud.google.com/apis/library/drive.googleapis.com?hl=es-419&organizationId=0&project=challenge-395013)
3. Configurar la pantalla de consentimiento de OAuth y agrégar usuario de prueba -> [Consentimiento](https://console.cloud.google.com/apis/credentials/consent?hl=es-419)
4. Generar credenciales para aplicación de escritorio de tipo **ID de cliente de OAuth** -> [OAuth](https://console.cloud.google.com/apis/credentials?hl=es-419) 
5. Descargar credencial **[client_secret]**

#### Entorno Windows con Docker

1. Descargar e instalar MySQL ->  [MySQL](https://downloads.mysql.com/archives/get/p/23/file/mysql-8.0.33-winx64-debug-test.zip)
2. Confirmar puerto de servicio **[port]**
3. Crear Base de Datos **[dbname]**
4. Generar usuario de lectura y escritura en la base de datos *[dbname]* creada en el punto 3 **[user][passroot]**
5. Descargar e instalar Docker -> [Docker](https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe)
6. Modificar y/o incluir la sección **bind-address=0.0.0.0** del archivo **my.ini** de la base de datos, ingresando la IP de Origen definida en el Contenedor.
7. Ejecutar desde Workbench el siguiente código:
```SQL
   GRANT ALL PRIVILEGES ON challenge.* TO 'dbadmin'@'%';
```
8. Reiniciar servicio de MySQL. 
9. Ingresar los siguientes datos de variables de entorno **ENV** ubicadas en el archivo **Dockerfile**

* ENV host=[IP del equipo donde se ejecuta el servicio MySQL]
* ENV port=[port]-> Definido en el punto 2
* ENV user=[user]-> Definido en el punto 4
* ENV passroot=[passroot]-> Definido en el punto 4
* ENV dbname=[dbname]-> Definido en el punto 3
* ENV client_secret=[client_secret]-> Ubicación de la llave en formato JSON, archivo descargado una vez configurado el entorno de Google Cloud.
* ENV openia=[Key API IA] -> Llave de conexión al API IA para su consumo

![](https://github.com/JossephRojasSantos/Challenge-DataSecurity/blob/master/png/ENVDockerFile.png)

8. Descargamos el presente repositorio, nos ubicamos con un CMD en el proyecto e ingresamos los siguientes comandos:
```docker
docker build -t servidor:Challenge .
```
```docker
docker run -p 443:443 -p 80:80 servidor:Challenge
```
9. Ingresamos en Docker Desktop y verificamos que en **Containers** nuestra imagen tenga estado **Running**
10. Ingresamos por medio de un navegador a **https://localhost/**
------
#### Entorno Windows sin Docker

1. Descargar e instalar MySQL ->  [MySQL](https://downloads.mysql.com/archives/get/p/23/file/mysql-8.0.33-winx64-debug-test.zip)
2. Confirmar puerto de servicio **[port]**
3. Crear Base de Datos **[dbname]**
4. Generar usuario de lectura y escritura en la base de datos *[dbname]* creada en el punto 3 **[user][passroot]**
5. Descargar e instalar Go -> [GO](https://dl.google.com/go/go1.20.4.windows-amd64.msi)
6. Comprobar instalación de Go (desde un CMD ejecutar "go version", retorno de la consola -> go version go1.20.4 windows/amd64)
7. Crear las siguiente variables de entorno:

* host=localhost
* port=[port]-> Definido en el punto 2
* user=[user]-> Definido en el punto 4
* passroot=[passroot]-> Definido en el punto 4
* dbname=[dbname]-> Definido en el punto 3
* client_secret=[client_secret]-> Ubicación de la llave en formato JSON, archivo descargado una vez configurado el entorno de Google Cloud.
* openia=[Key API IA] -> Llave de conexión al API IA para su consumo

![](https://github.com/JossephRojasSantos/Challenge-DataSecurity/blob/master/png/Variables%20de%20Entorno.png)

8. Descargamos el presente repositorio, nos ubicamos con un CMD en el proyecto e ingresamos el siguiente comando:
```go
go build main.go
```
9. Dentro de la carpeta del proyecto, ubicamos y ejecutamos como administrador el archivo **main.exe**.    
10. Ingresamos por medio de un navegador a **https://localhost/**

**Nota:** Es necesario ejecutar el binario como administrador (Punto 9) para obtener los valores de las variables de entorno configuradas en windows

### Modo de Operación 

------
1. Ingresamos en el navegador la siguiente URL https://localhost/
2. En la pantalla de Bienvenida, damos clic en el botón **Otorgar Permisos**
3. Seguidamente, se es redireccionado para iniciar sesión en la cuenta de Google. Iniciamos y otorgamos permiso al aplicativo con el objetivo de Ver, modificar, crear y eliminar archivos de Google Drive.
4. Retornamos a la pantalla de Bienvenida una vez finalizado el punto 3. En este momento, podemos iniciar con la clasificación de los documentos ingresando al botón **Iniciar Clasificación**
5. En esta nueva ubicación, se visualizará dentro de una tabla el listado de documentos almacenados en Google Drive. Para modificar su criticidad, damos clic en el botón **Modificar** 
6. En esta sección, podemos visualizar algunos datos del documento. Con el objetivo de clasificar el documento, respondemos las preguntas planteadas y damos clic en el botón **Guardar Clasificación**
7. Repetimos el punto 6 para todo aquel documento denominado como "Sin Clasificar"

