package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	handler "fgl-backend/app/handler"
	hub "fgl-backend/app/hub"
	model "fgl-backend/app/model"
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
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	flag.Parse()

	a.get("/announcement/{aid}", a.getAnnouncement)
	a.get("/announcements", a.getAnnouncements)
	a.post("/makeann", a.makeAnnouncement)
	a.get("/labreport/{rid}", a.getLabReport)
	a.get("/labreports", a.getLabReports)
	a.post("/makelabreport", a.makeLabReport)
	a.get("/gadgetreport/{rid}", a.getGadgetReport)
	a.get("/gadgetreports", a.getGadgetReports)
	a.post("/makegadgetreport", a.makeGadgetReport)
	a.get("/users", a.getAllUsers)
	a.post("/register", a.registerUser)
	a.post("/login", a.userLogin)
	a.get("/users/{username}", a.getUser)
	a.put("/users/{username}", a.updateUser)
	a.delete("/users/{username}", a.deleteUser)

	a.post("/recversion", a.receiveVersion)
	a.get("/getupdater", a.sendUpdater)

	//a.Router.HandleFunc("/", handler.SendFile)

	h := hub.NewHub()
	a.chat("/chat", hub.WsHandler{H: h})
}

func (a *App) chat(path string, c hub.WsHandler) {
	a.Router.Handle(path, c)
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

func (a *App) receiveVersion(w http.ResponseWriter, r *http.Request) {
	handler.ReceiveVersion(w, r)
}

func (a *App) sendUpdater(w http.ResponseWriter, r *http.Request) {
	handler.SendUpdater(w, r)
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

func (a *App) getLabReport(w http.ResponseWriter, r *http.Request) {
	handler.GetLabReport(a.DB, w, r)
}

func (a *App) getLabReports(w http.ResponseWriter, r *http.Request) {
	handler.GetLabReports(a.DB, w, r)
}

func (a *App) makeLabReport(w http.ResponseWriter, r *http.Request) {
	handler.MakeLabReport(a.DB, w, r)
}

func (a *App) getGadgetReport(w http.ResponseWriter, r *http.Request) {
	handler.GetGadgetReport(a.DB, w, r)
}

func (a *App) getGadgetReports(w http.ResponseWriter, r *http.Request) {
	handler.GetGadgetReports(a.DB, w, r)
}

func (a *App) makeGadgetReport(w http.ResponseWriter, r *http.Request) {
	handler.MakeGadgetReport(a.DB, w, r)
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
	fmt.Println("Server started")
	log.Fatal(http.ListenAndServe(host, a.Router))
}
