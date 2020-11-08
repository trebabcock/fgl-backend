package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type version struct {
	Version string `json:"version"`
}

// ReceiveVersion handles the version number sent from the client
func ReceiveVersion(w http.ResponseWriter, r *http.Request) {
	cVersion := version{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cVersion); err != nil {
		fmt.Println("error decoding version:", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if needsUpdate(cVersion) {
		RespondJSON(w, http.StatusUpgradeRequired, serverVersion())
	} else {
		RespondJSON(w, http.StatusOK, serverVersion())
	}
}

func needsUpdate(current version) bool {
	version := serverVersion()

	cvFloat, _ := strconv.ParseFloat(current.Version, 64)
	svFloat, _ := strconv.ParseFloat(version.Version, 64)

	return svFloat > cvFloat
}

func serverVersion() version {
	version := version{}
	versionFile, err := os.Open("version.json")
	if err != nil {
		fmt.Println("error opening version file:", err)
	}
	byteValue, _ := ioutil.ReadAll(versionFile)
	if err := json.Unmarshal(byteValue, &version); err != nil {
		fmt.Println("error unmarshalling version.json:", err)
	}
	defer versionFile.Close()

	return version
}

// SendUpdater sends the updater executable to the client
func SendUpdater(w http.ResponseWriter, r *http.Request) {
	updater, err := os.Open("fgl-updater.exe")
	defer updater.Close()
	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	uBytes, err := ioutil.ReadAll(updater)
	if err != nil {
		fmt.Println("ReadAll(): ", err)
	}
	RespondJSON(w, http.StatusOK, uBytes)
}

func sendFile(w http.ResponseWriter, r *http.Request) {
	Filename := r.URL.Query().Get("file")
	if Filename == "" {
		http.Error(w, "Get 'file' not specified in url.", 400)
		return
	}
	fmt.Println("Client requests: " + Filename)

	Openfile, err := os.Open(Filename)
	if err != nil {
		http.Error(w, "File not found.", 404)
		return
	}

	FileHeader := make([]byte, 512)
	Openfile.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)

	FileStat, _ := Openfile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	Openfile.Seek(0, 0)
	io.Copy(w, Openfile)
	return
}
