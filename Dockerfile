# Define la imagen base
FROM golang:latest

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos del proyecto al directorio de trabajo en el contenedor
COPY . .

# Descarga las dependencias del módulo Go
RUN go mod download

# Compila la aplicación
RUN go build -o main

# Expone el puerto en el que la aplicación escucha
EXPOSE 443
EXPOSE 80

# Define el comando de inicio de la aplicación
CMD ["./main"]

ENV host=192.168.1.11
ENV port=3306
ENV user=dbadmin
ENV passroot=XXXXXXXXXXXXXXXXXXXXX
ENV dbname=challenge
ENV client_secret=/app/client_secret.json
ENV openia=XXXXXXXXXXXXX

