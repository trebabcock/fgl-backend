package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	handler "fgl-backend/app/handler"
	"fgl-backend/app/model"
	db "fgl-backend/db"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// App holds the database and router
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize connects to the database, migrates the database, and sets up routes
func (a *App) Initialize(dbConfig *db.Config) {
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s connect_timeout=15",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Username,
		dbConfig.Name,
		dbConfig.Password,
		dbConfig.SSLMode,
	)
	fmt.Println("Connecting to database...")
	db, err := gorm.Open(dbConfig.Dialect, dbURI)
	if err != nil {
		log.Println("Could not connect to database")
		log.Fatal(err.Error())
	}

	fmt.Println("Connected to database")
	fmt.Println("Migrating database...")
	a.DB = model.DBMigrate(db)
	fmt.Println("Migrated database")
	a.Router = mux.NewRouter()
	a.setv1Routers()
	a.setStaticServe()
}

func (a *App) setStaticServe() {
	a.Router.HandleFunc("/", a.serveIndex)
	a.Router.HandleFunc("/download", a.serveDownload)
}

func (a *App) setv1Routers() {
	a.get("/api/v1/announcement/{aid}", a.getAnnouncement)
	a.get("/api/v1/announcements", a.getAnnouncements)
	a.post("/api/v1/makeann", a.makeAnnouncement)
	a.put("/api/v1/updateannouncement/{aid}", a.updateAnnouncement)
	a.delete("/api/v1/deleteannouncement/{aid}", a.deleteAnnouncement)
	a.get("/api/v1/discussions/{did}", a.getDiscussionDocument)
	a.get("/api/v1/discussions", a.getDiscussionDocuments)
	a.post("/api/v1/makediscussion", a.makeDiscussionDocument)
	a.put("/api/v1/updatediscussion/{did}", a.updateDiscussionDocument)
	a.delete("/api/v1/deletediscussion/{did}", a.deleteDiscussionDocument)
	a.get("/api/v1/projects/{pid}", a.getProjectSubmission)
	a.get("/api/v1/projects", a.getProjectSubmissions)
	a.post("/api/v1/makeproject", a.makeProjectSubmission)
	a.put("/api/v1/updateproject/{pid}", a.updateProjectSubmission)
	a.delete("/api/v1/deleteproject/{pid}", a.deleteProjectSubmission)
	a.get("/api/v1/users", a.getAllUsers)
	a.post("/api/v1/register", a.registerUser)
	a.post("/api/v1/login", a.userLogin)
	a.get("/api/v1/users/{username}", a.getUser)
	a.put("/api/v1/users/{username}", a.updateUser)
	a.delete("/api/v1/users/{username}", a.deleteUser)

	a.get("/api/v1/messages", a.getMessages)
	a.post("/api/v1/message", a.receivedMessage)

	a.post("/api/v1/recversion", a.receiveVersion)
	a.get("/api/v1/getupdater", a.sendUpdater)

	a.post("/api/v1/auth/{auth_code}", a.authorize)
	a.get("/api/v1/newauth", a.newCode)

	a.get("/api/v1/fglClient", a.downloadClient)

	//a.Router.HandleFunc("/", handler.SendFile)
}

func (a *App) serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}

func (a *App) serveDownload(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]

	if !ok || len(keys[0]) < 1 {
		w.Write([]byte(strconv.Itoa(http.StatusUnauthorized) + " Unauthorized"))
		log.Println("first if")
		return
	}

	if !handler.AuthorizeCode(a.DB, keys[0]) {
		w.Write([]byte(strconv.Itoa(http.StatusUnauthorized) + " Unauthorized"))
		log.Println("second if")
		return
	}

	http.ServeFile(w, r, "public/download.html")
}

func (a *App) downloadClient(w http.ResponseWriter, r *http.Request) {
	code, ok := r.URL.Query()["code"]

	if !ok || len(code[0]) < 1 {
		w.Write([]byte(strconv.Itoa(http.StatusUnauthorized) + " Unauthorized"))
		log.Println("first if")
		return
	}

	if !handler.AuthorizeCode(a.DB, code[0]) {
		w.Write([]byte(strconv.Itoa(http.StatusUnauthorized) + " Unauthorized"))
		log.Println("second if")
		return
	}

	version, ok := r.URL.Query()["version"]

	filename := ""

	if version[1] == "classic" {
		filename = "public/fgl-client.exe"
	} else if version[1] == "gui" {
		filename = "public/fgl-gui.exe"
	}

	Openfile, err := os.Open(filename)
	if err != nil {
		http.Error(w, "File not found.", 404)
		return
	}

	FileHeader := make([]byte, 512)
	Openfile.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)

	FileStat, _ := Openfile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename=fgl-client.exe")
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	Openfile.Seek(0, 0)
	io.Copy(w, Openfile)
	return
}

func (a *App) get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

func (a *App) delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

func (a *App) authorize(w http.ResponseWriter, r *http.Request) {
	handler.AuthorizeDownload(a.DB, w, r)
}

func (a *App) newCode(w http.ResponseWriter, r *http.Request) {
	handler.MakeCode(a.DB, w, r)
}

func (a *App) receiveVersion(w http.ResponseWriter, r *http.Request) {
	handler.ReceiveVersion(w, r)
}

func (a *App) sendUpdater(w http.ResponseWriter, r *http.Request) {
	handler.SendUpdater(w, r)
}

func (a *App) receivedMessage(w http.ResponseWriter, r *http.Request) {
	handler.ReceivedMessage(a.DB, w, r)
}

func (a *App) getMessages(w http.ResponseWriter, r *http.Request) {
	handler.GetMessages(a.DB, w, r)
}

func (a *App) getAnnouncement(w http.ResponseWriter, r *http.Request) {
	handler.GetAnnouncement(a.DB, w, r)
}

func (a *App) getAnnouncements(w http.ResponseWriter, r *http.Request) {
	handler.GetAnnouncements(a.DB, w, r)
}

func (a *App) makeAnnouncement(w http.ResponseWriter, r *http.Request) {
	handler.MakeAnnouncement(a.DB, w, r)
}

func (a *App) updateAnnouncement(w http.ResponseWriter, r *http.Request) {
	handler.UpdateAnnouncement(a.DB, w, r)
}

func (a *App) deleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	handler.DeleteAnnouncement(a.DB, w, r)
}

func (a *App) getDiscussionDocument(w http.ResponseWriter, r *http.Request) {
	handler.GetDiscussionDocument(a.DB, w, r)
}

func (a *App) getDiscussionDocuments(w http.ResponseWriter, r *http.Request) {
	handler.GetDiscussionDocuments(a.DB, w, r)
}

func (a *App) makeDiscussionDocument(w http.ResponseWriter, r *http.Request) {
	handler.MakeDiscussionDocument(a.DB, w, r)
}

func (a *App) updateDiscussionDocument(w http.ResponseWriter, r *http.Request) {
	handler.UpdateDiscussionDocument(a.DB, w, r)
}

func (a *App) deleteDiscussionDocument(w http.ResponseWriter, r *http.Request) {
	handler.DeleteDiscussionDocument(a.DB, w, r)
}

func (a *App) getProjectSubmission(w http.ResponseWriter, r *http.Request) {
	handler.GetProjectSubmission(a.DB, w, r)
}

func (a *App) getProjectSubmissions(w http.ResponseWriter, r *http.Request) {
	handler.GetProjectSubmissions(a.DB, w, r)
}

func (a *App) makeProjectSubmission(w http.ResponseWriter, r *http.Request) {
	handler.MakeProjectSubmission(a.DB, w, r)
}

func (a *App) updateProjectSubmission(w http.ResponseWriter, r *http.Request) {
	handler.UpdateProjectSubmission(a.DB, w, r)
}

func (a *App) deleteProjectSubmission(w http.ResponseWriter, r *http.Request) {
	handler.DeleteProjectSubmission(a.DB, w, r)
}

func (a *App) getAllUsers(w http.ResponseWriter, r *http.Request) {
	handler.GetAllUsers(a.DB, w, r)
}

func (a *App) registerUser(w http.ResponseWriter, r *http.Request) {
	handler.RegisterUser(a.DB, w, r)
}

func (a *App) userLogin(w http.ResponseWriter, r *http.Request) {
	handler.UserLogin(a.DB, w, r)
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	handler.GetUser(a.DB, w, r)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	handler.UpdateUser(a.DB, w, r)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	handler.DeleteUser(a.DB, w, r)
}

// Run starts the server
func (a *App) Run(host string) {
	fmt.Println("Server running at", host)
	log.Fatal(http.ListenAndServe(host, a.Router))
}
