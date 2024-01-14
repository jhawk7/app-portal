package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jhawk7/app-portal/internal/pkg/loggers"
	"github.com/surrealdb/surrealdb.go"
)

var DB_CONN = os.Getenv("DB_CONN")
var DB_NAME = os.Getenv("DB_NAME")
var DB_USER = os.Getenv("DB_USER")
var DB_PASSWORD = os.Getenv("DB_PASSWORD")
var DB_NS = os.Getenv("DB_NS")
var SELECTOR = strings.ToLower(DB_NAME)

type DBClient struct {
	svc *surrealdb.DB
}

type Portal struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	RedirectUrl string `json:"redirectUrl"`
	Img         string `json:"img"`
	Count       int    `json:"count,omitempty"`
}

type UpdatePortal struct {
	ID      string `json:"id"`
	Changes struct {
		RedirectUrl string `json:"redirectUrl,omitempty"`
		Img         string `json:"img,omitempty"`
		Count       int    `json:"count,omitempty"`
	} `json:"changes"`
}

func InitDB() (dbClient *DBClient, err error) {
	db, connErr := surrealdb.New(DB_CONN)
	if connErr != nil {
		err = fmt.Errorf("failed to connect to db; %v", connErr)
		return
	}

	if _, signinErr := db.Signin(map[string]interface{}{
		"user": DB_USER,
		"pass": DB_PASSWORD,
	}); signinErr != nil {
		err = fmt.Errorf("failed to authenticate db user creds; %v", signinErr)
		return
	}

	if _, nsErr := db.Use(DB_NS, DB_NAME); nsErr != nil {
		err = fmt.Errorf("failed to use namespace and dbname; %v", nsErr)
		return
	}

	dbClient = new(DBClient)
	dbClient.svc = db
	loggers.LogInfo("established connection to DB")
	return
}

func (dbClient *DBClient) GetAllPortals() (portals *[]Portal, notFound bool, err error) {
	data, dbErr := dbClient.svc.Select(DB_NAME)
	if dbErr != nil {
		err = fmt.Errorf("failed to retreive portals from db; %v", dbErr)
		return
	}

	if data == nil {
		err = errors.New("no data")
		notFound = true
		return
	}

	if mErr := surrealdb.Unmarshal(data, &portals); mErr != nil {
		err = fmt.Errorf("failed to marshal db response to struct; [func: GetAllPortals]; %v", mErr)
		return
	}
	loggers.LogInfo("retrieved all portals from db")
	return
}

func (dbClient *DBClient) InsertPortal(portal *Portal) (portalID string, err error) {
	id := fmt.Sprintf("%v:%v", SELECTOR, portal.Name)
	data, dbErr := dbClient.svc.Create(id, portal)
	if dbErr != nil {
		err = fmt.Errorf("failed to create portal %v; %v", portal.Name, dbErr)
		return
	}

	newPortal := new(Portal)
	if mErr := surrealdb.Unmarshal(data, &newPortal); mErr != nil {
		err = fmt.Errorf("failed to marshal db response to struct; [func: InsertPortal]; %v", mErr)
		return
	}

	portalID = newPortal.ID
	loggers.LogInfo(fmt.Sprintf("created portal %v", portal.Name))
	return
}

func (dbClient *DBClient) ModifyPortal(update *UpdatePortal) (portal Portal, notFound bool, err error) {
	changes, convErr := struct2map(update.Changes)
	if convErr != nil {
		err = fmt.Errorf("failed to convert struct to map during update; [func: ModifyPortal - struct2map]; %v", convErr)
		return
	}

	data, dbErr := dbClient.svc.Change(update.ID, changes)
	if dbErr != nil {
		err = fmt.Errorf("failed to retreive portals from db; %v", dbErr)
		if dbErr.Error() == "error no row" {
			notFound = true
		}
		return
	}

	if mErr := surrealdb.Unmarshal(data, &portal); mErr != nil {
		err = fmt.Errorf("failed to marshal db response to struct; [func: ModifyPortal]; %v", mErr)
		return
	}

	loggers.LogInfo(fmt.Sprintf("updated portal %v", update.ID))
	return
}

func (dbClient *DBClient) RemovePortal(portalID string) (err error) {
	id := fmt.Sprintf("%v:%v", SELECTOR, portalID)

	if _, dbErr := dbClient.svc.Delete(id); dbErr != nil {
		err = fmt.Errorf("failed to delete portal %v from db; %v", portalID, dbErr)
		return
	}

	loggers.LogInfo(fmt.Sprintf("deleted portal %v", portalID))
	return
}

func struct2map(in interface{}) (out map[string]string, err error) {
	raw, mErr := json.Marshal(in)
	if mErr != nil {
		err = mErr
		return
	}

	if umErr := json.Unmarshal(raw, &out); umErr != nil {
		err = umErr
		return
	}
	return
}
