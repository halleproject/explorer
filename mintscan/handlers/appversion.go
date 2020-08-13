package handlers

import (
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/client"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/db"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/models"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/utils"
	"log"
	"net/http"
)

// AppVersion is a AppVersion handler
type AppVersion struct {
	l      *log.Logger
	client *client.Client
	db     *db.Database
}

// NewAppVersion creates a new AppVersion handler with the given params
func NewAppVersion(l *log.Logger, client *client.Client, db *db.Database) *AppVersion {
	return &AppVersion{l, client, db}
}

// GetAppVersion creae new key and save DB, then return AppVersion on the active chain
func (av *AppVersion) SetVersion(rw http.ResponseWriter, r *http.Request) {
	var passwd, version, downloadURL string

	if len(r.URL.Query()["passwd"]) > 0 {
		passwd = r.URL.Query()["passwd"][0]
	} else {
		av.l.Printf("failed to get passwd")
		return
	}
	if len(r.URL.Query()["version"]) > 0 {
		version = r.URL.Query()["version"][0]
	} else {
		av.l.Printf("failed to get version")
		return
	}
	if len(r.URL.Query()["downloadurl"]) > 0 {
		downloadURL = r.URL.Query()["downloadurl"][0]
	} else {
		av.l.Printf("failed to get downloadurl")
		return
	}
	//fmt.Println(passwd, version, downloadURL)

	appVersion, _ := av.db.QueryAppVersion()
	if appVersion.ID == 0 {
		appVersion.Passwd = passwd
		appVersion.Version = version
		appVersion.DownloadURL = downloadURL
		err := av.db.InsertAppVerion(&appVersion)
		if err != nil {
			av.l.Printf("failed to insert AppVersion : %s", err)
			return
		}
		utils.Respond(rw, true)
		return
	} else if appVersion.Passwd != passwd {
		av.l.Printf("failed to set AppVersion : passwd wrong")
		return
	}
	appVersion.Version = version
	appVersion.DownloadURL = downloadURL
	err := av.db.UpdateAppVerion(&appVersion)
	if err != nil {
		av.l.Printf("failed to insert AppVersion : %s", err)
		return
	}

	utils.Respond(rw, true)
	return
}

// GetVersion returns AppVersion information
func (ta *AppVersion) GetVersion(rw http.ResponseWriter, r *http.Request) {
	appVersion, err := ta.db.QueryAppVersion()
	if err != nil {
		ta.l.Printf("failed to query AppVersion : %s", err)
		return
	}
	//fmt.Println(appVersion.ID, appVersion.Version, appVersion.DownloadURL, "hello")
	utils.Respond(rw, models.ResultVersion{appVersion.Version, appVersion.DownloadURL})
	return
}
