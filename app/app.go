package app

import (
	"fmt"
	"log"
	"net/http"

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
	a.get("/api/v1/labreport/{rid}", a.getLabReport)
	a.get("/api/v1/labreports", a.getLabReports)
	a.post("/api/v1/makelabreport", a.makeLabReport)
	a.put("/api/v1/updatelabreport/{rid}", a.updateLabReport)
	a.delete("/api/v1/deletelabreport/{rid}", a.deleteLabReport)
	a.get("/api/v1/gadgetreport/{rid}", a.getGadgetReport)
	a.get("/api/v1/gadgetreports", a.getGadgetReports)
	a.post("/api/v1/makegadgetreport", a.makeGadgetReport)
	a.put("/api/v1/updategadgetreport/{rid}", a.updateGadgetReport)
	a.delete("/api/v1/deletegadgetreport/{rid}", a.deleteGadgetReport)
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

	//a.Router.HandleFunc("/", handler.SendFile)
}

func (a *App) serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}

func (a *App) serveDownload(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["auth_code"]

	if !ok || len(keys[0]) < 1 {
		w.Write([]byte(string(http.StatusUnauthorized) + " Unauthorized"))
		return
	}

	if !handler.AuthorizeCode(a.DB, keys[0]) {
		w.Write([]byte(string(http.StatusUnauthorized) + " Unauthorized"))
		return
	}

	http.ServeFile(w, r, "public/download.html")
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

func (a *App) getLabReport(w http.ResponseWriter, r *http.Request) {
	handler.GetLabReport(a.DB, w, r)
}

func (a *App) getLabReports(w http.ResponseWriter, r *http.Request) {
	handler.GetLabReports(a.DB, w, r)
}

func (a *App) makeLabReport(w http.ResponseWriter, r *http.Request) {
	handler.MakeLabReport(a.DB, w, r)
}

func (a *App) updateLabReport(w http.ResponseWriter, r *http.Request) {
	handler.UpdateLabReport(a.DB, w, r)
}

func (a *App) deleteLabReport(w http.ResponseWriter, r *http.Request) {
	handler.DeleteLabReport(a.DB, w, r)
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

func (a *App) updateGadgetReport(w http.ResponseWriter, r *http.Request) {
	handler.UpdateGadgetReport(a.DB, w, r)
}

func (a *App) deleteGadgetReport(w http.ResponseWriter, r *http.Request) {
	handler.DeleteGadgetReport(a.DB, w, r)
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
