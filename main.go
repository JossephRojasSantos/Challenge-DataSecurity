package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var (
	host           = os.Getenv("host")
	port           = os.Getenv("port")
	user           = os.Getenv("user")
	passroot       = os.Getenv("passroot")
	dbname         = os.Getenv("dbname")
	client_secret  = os.Getenv("client_secret")
	openia         = os.Getenv("openia")
	mysqlInfo      = fmt.Sprintf(user + ":" + passroot + "@tcp(" + host + ":" + port + ")/" + dbname)
	filePathKEY    string
	filePathCERT   string
	tmpl           = template.Must(template.ParseGlob("template/*"))
	authURL        string
	idDatoToken    string
	table          = Table{}
	arraytable     []Table
	arraytable2    []Table
	data           Data
	ID             string
	FileID         string
	FileName       string
	Extension      string
	FileOwner      string
	Visibility     int
	ViewFile       string
	Classification sql.NullString
	Version        string
	CreatedTime    string
	ModifiedTime   string
	QuotaBytesUsed string
	srv            *drive.Service
)

type Data struct {
	Suma    int    `json:"suma"`
	File_ID string `json:"otroValor"`
}
type Table struct {
	ID             string
	FileID         string
	FileName       string
	Extension      string
	FileOwner      string
	Visibility     string
	ViewFile       string
	Classification string
	Version        string
	CreatedTime    string
	ModifiedTime   string
	QuotaBytesUsed string
}

func Err(err2 error) {
	if err2 != nil {
		log.Println(err2)
	}
}
func Certificados() {
	currentDir, err := os.Getwd()
	Err(err)
	filePathKEY = filepath.Join(currentDir, "/certificados/server.key")
	filePathCERT = filepath.Join(currentDir, "/certificados/certbundle.pem")

}
func ConexionDB() (db *sql.DB) {
	db, err := sql.Open("mysql", mysqlInfo)
	log.Println("Conectando con Base de Datos")
	Err(err)
	return db
}
func URLGDrive() {
	credentialsJSON, err := ioutil.ReadFile(client_secret)
	credentials, err := google.ConfigFromJSON(credentialsJSON, drive.DriveScope)
	Err(err)
	authURL = credentials.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

}
func redirectToHttps(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "https://localhost:443"+request.RequestURI, http.StatusMovedPermanently)
	log.Println("Redireccion a HTTPS")
}

func main() {

	Certificados()
	CrearTablaInventario()

	r := mux.NewRouter()
	r.HandleFunc("/", Inicio)
	r.HandleFunc("/clasificacion", Clasificacion)
	r.HandleFunc("/listado", Listado)
	r.HandleFunc("/openia", OpenIA)

	go func() {
		_ = http.ListenAndServeTLS(":443", filePathCERT, filePathKEY, r)

	}()
	_ = http.ListenAndServe(":80", http.HandlerFunc(redirectToHttps))

}

func CrearTablaInventario() {
	db := ConexionDB()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	log.Println("Creando (si no existe) la tabla inventario")
	createTableQuery := `CREATE TABLE IF NOT EXISTS inventario (
                            id INT AUTO_INCREMENT PRIMARY KEY,
                            FileID TEXT,
                            FileName TEXT,    
							Extension TEXT,   
							FileOwner TEXT,
							Visibility INT,
							Classification TEXT,
							ViewFile TEXT,
                            Version TEXT,
                            CreatedTime TEXT,
                            ModifiedTime TEXT,
                            QuotaBytesUsed TEXT
                         );`
	_, err := db.Exec(createTableQuery)
	Err(err)
}
func ConexionGDrive() {
	ctx := context.Background()
	//timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	//defer cancel()
	credentialsJSON, err := ioutil.ReadFile(client_secret)
	credentials, err := google.ConfigFromJSON(credentialsJSON, drive.DriveScope)
	Err(err)

	if idDatoToken != "" {
		token, err := credentials.Exchange(context.Background(), idDatoToken)
		Err(err)

		client := credentials.Client(ctx, token)

		srv, err = drive.New(client)
		Err(err)

		fileList, err := srv.Files.List().Fields("files(id, name, fileExtension, properties)").Do()
		Err(err)
		for _, file := range fileList.Files {
			if file.FileExtension != "" {
				if len(file.Properties) == 0 {
					properties := make(map[string]string)
					properties["clasificacion"] = "sin clasificar"
					updateFile := &drive.File{
						Properties: properties,
					}
					_, err := srv.Files.Update(file.Id, updateFile).Fields("properties").Do()
					Err(err)
				}
			}
		}

		db := ConexionDB()
		defer func(db *sql.DB) {
			_ = db.Close()
		}(db)

		_, err = db.Exec("DELETE FROM inventario")
		Err(err)
		fileList2, err := srv.Files.List().Fields("files(id, name, fileExtension, webViewLink , owners, shared, properties,version,createdTime,modifiedTime,quotaBytesUsed)").Do()
		Err(err)

		for _, file := range fileList2.Files {
			if file.FileExtension != "" {
				for _, owner := range file.Owners {
					insertQuery := "INSERT INTO inventario (FileID,FileName,Extension,FileOwner,Visibility,ViewFile,Version,CreatedTime,ModifiedTime,QuotaBytesUsed) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
					_, err := db.Exec(insertQuery, file.Id, file.Name, file.FileExtension, owner.EmailAddress, file.Shared, file.WebViewLink, file.Version, file.CreatedTime, file.ModifiedTime, file.QuotaBytesUsed)
					Err(err)

				}
			}
		}

		for _, file := range fileList2.Files {
			if file.FileExtension != "" {
				for _, properties := range file.Properties {
					updateQuery := "UPDATE inventario SET Classification = ? WHERE FileID = ?"
					_, err := db.Exec(updateQuery, properties, file.Id)
					Err(err)
				}
			}
		}

	}
}
func Inicio(w http.ResponseWriter, r *http.Request) {
	URLGDrive()
	_ = tmpl.ExecuteTemplate(w, "inicio", authURL)
	idDatoToken = r.URL.Query().Get("code")
	ConexionGDrive()
}
func Listado(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		if err != nil {
			http.Error(w, "Error al decodificar los datos", http.StatusBadRequest)
			return
		}
		//fmt.Print(strconv.Itoa(data.Suma) + "\n")
	}

	Modificacion(data.Suma, data.File_ID)

	db := ConexionDB()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	arraytable = nil
	registros, err := db.Query("Select ID,FileID,FileName,Extension,FileOwner,Visibility,Classification,ViewFile From inventario")
	Err(err)
	log.Println("Pagina de listado")
	for registros.Next() {
		err = registros.Scan(&ID, &FileID, &FileName, &Extension, &FileOwner, &Visibility, &Classification, &ViewFile)
		Err(err)

		table.ID = ID
		table.FileID = FileID
		table.FileName = FileName
		table.Extension = Extension
		table.FileOwner = FileOwner
		if Visibility == 1 {
			table.Visibility = "Público"
		} else {
			table.Visibility = "Privado"
		}
		table.ViewFile = ViewFile
		if Classification.Valid {
			table.Classification = Classification.String
		} else {
			table.Classification = "Permisos Insuficientes"
		}
		arraytable = append(arraytable, table)
	}
	_ = tmpl.ExecuteTemplate(w, "listado", arraytable)
}
func Clasificacion(w http.ResponseWriter, r *http.Request) {

	idDato := r.URL.Query().Get("id")

	db := ConexionDB()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	query := "Select FileID,FileName,FileOwner,Visibility,Classification,ViewFile,Version,CreatedTime,ModifiedTime,QuotaBytesUsed From inventario where id = ?"
	registros, err := db.Query(query, idDato)
	Err(err)
	arraytable2 = nil
	log.Println("Pagina de Informacion")
	for registros.Next() {
		err = registros.Scan(&FileID, &FileName, &FileOwner, &Visibility, &Classification, &ViewFile, &Version, &CreatedTime, &ModifiedTime, &QuotaBytesUsed)
		Err(err)

		table.FileID = FileID
		table.FileName = FileName
		table.FileOwner = FileOwner
		if Visibility == 1 {
			table.Visibility = "Público"
		} else {
			table.Visibility = "Privado"
		}
		table.ViewFile = ViewFile
		if Classification.Valid {
			table.Classification = Classification.String
		} else {
			table.Classification = "Permisos Insuficientes"
		}
		table.Version = Version
		table.CreatedTime = CreatedTime
		table.ModifiedTime = ModifiedTime
		table.QuotaBytesUsed = QuotaBytesUsed

		arraytable2 = append(arraytable2, table)
	}

	_ = tmpl.ExecuteTemplate(w, "clasificacion", arraytable2)

}
func Modificacion(data int, fileid string) {

	if data >= 1 && data < 4 {
		InsertarCriticidad("Bajo", fileid)
		log.Println("Bajo")
	} else if data >= 4 && data < 6 {
		InsertarCriticidad("Medio", fileid)
		log.Println("Medio")
	} else if data >= 6 && data < 9 {
		InsertarCriticidad("Alto", fileid)
		log.Println("Alto")
	} else if data >= 9 {
		InsertarCriticidad("Crítico", fileid)
		log.Println("Critico")
	} else {
		log.Println("No se clasifico")
	}
}
func InsertarCriticidad(criticidad string, idArchivo string) {

	db := ConexionDB()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	updateQuery := "UPDATE inventario SET Classification = ? WHERE FileID = ?"
	_, err := db.Exec(updateQuery, criticidad, idArchivo)
	Err(err)

	if idArchivo != "" {
		properties := make(map[string]string)
		properties["clasificacion"] = criticidad
		updateFile := &drive.File{
			Properties: properties,
		}
		_, err = srv.Files.Update(idArchivo, updateFile).Fields("properties").Do()
		Err(err)
		log.Println("Archivo clasificado")
		if criticidad == "Critico" {
			privado := 0
			updateQuery := "UPDATE inventario SET Visibility = ? WHERE FileID = ?"
			_, err := db.Exec(updateQuery, privado, idArchivo)
			Err(err)

			nuevoShared := false
			updateFile := &drive.File{
				Shared: nuevoShared,
			}
			_, err = srv.Files.Update(idArchivo, updateFile).Fields("shared").Do()
			Err(err)
		}
	}

}

func OpenIA(w http.ResponseWriter, r *http.Request) {

	openia = "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"
	query := r.FormValue("query")
	log.Println(query)
	requestData := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a helpful assistant.",
			},
			{
				"role":    "user",
				"content": query,
			},
		},
		"model": "gpt-3.5-turbo", // Cambia al modelo de chat adecuado
	}

	// Codifica el mapa como JSON
	jsonData, err := json.Marshal(requestData)
	Err(err)

	apiURL := "https://api.openai.com/v1/chat/completions"
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(jsonData))
	Err(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openia)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error from API: "+resp.Status, http.StatusInternalServerError)
		return
	}

	openAIResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(openAIResponse, &result); err != nil {
		http.Error(w, "Error parsing response JSON", http.StatusInternalServerError)
		return
	}

	if len(result.Choices) > 0 {
		fmt.Println("Response Content:", result.Choices[0].Message.Content)
		_ = tmpl.ExecuteTemplate(w, "openia", result.Choices[0].Message.Content)
	} else {
		_ = tmpl.ExecuteTemplate(w, "openia", nil)
	}
}
