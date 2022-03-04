package models

import (
	"casaattivita/lang"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

// swagger:model Activity
type Activity struct {
	// Id of Activity value
	// in: int64
	Id int64 `json:"id"`
	// Value of Activity
	// in: int
	Value int `json:"valore"`
	// Timestamp of insert
	// in: time
	DateInsert time.Time `json:"data_inserimento"`
}

type Activities []Activity

type ReqAddActivity struct {
	// Value of the Activity
	// in: int
	Value int `json:"valore" validate:"required"`
}

// swagger:parameters addActivity
type ReqActivityBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/ReqAddActivity"
	//  required: true
	Body ReqAddActivity `json:"body"`
}

// ErrHandler returns error message bassed on env debug
func ErrHandler(err error) string {
	var errmessage string
	if os.Getenv("DEBUG") == "true" {
		errmessage = err.Error()
	} else {
		errmessage = lang.Get("something_went_wrong")
	}
	return errmessage
}

func GetActivitiesSqlx(db *sql.DB) *Activities {
	activities := Activities{}
	rows, err := db.Query("SELECT id,0, data_inserimento FROM attivita order by id desc")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Activity
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		activities = append(activities, p)
	}
	return &activities
}
func GetLastActivitySqlx(db *sql.DB) *Activities {
	activities := Activities{}
	rows, err := db.Query("SELECT id, 0, data_inserimento FROM attivita where id = (select max(id) from attivita)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Activity
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		activities = append(activities, p)
	}
	return &activities
}
func GetLastHourSqlx(db *sql.DB) *Activities {
	activities := Activities{}

	tFine := time.Now()
	dataFine := tFine.Format("2006-01-02 15:04:05")

	tInizio := time.Now().Add(time.Duration(-1) * time.Hour)
	dataInizio := tInizio.Format("2006-01-02 15:04:05")

	sqlStatement := fmt.Sprintf("SELECT id,0, data_inserimento FROM attivita where data_inserimento  >= '%s' AND data_inserimento <= '%s'", dataInizio, dataFine)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Activity
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		activities = append(activities, p)
	}

	if len(activities) == 0 {
		elemento := GetLastActivitySqlx(db)
		activities = append(activities, *elemento...)
	}
	return &activities
}

// PostActivitySqlx insert Activity value
func PostActivitySqlx(db *sql.DB, reqrain *ReqAddActivity) (*Activity, string) {

	//value := reqrain.Value

	var activity Activity

	lastInsertId := 0

	//sqlStatement := fmt.Sprintf("insert into 'pioggia' ('valore','data_inserimento') values (%d,CURRENT_TIMESTAMP) RETURNING id", value)
	sqlStatement := fmt.Sprintf("insert into attivita (data_inserimento) values (CURRENT_TIMESTAMP) RETURNING id")

	err := db.QueryRow(sqlStatement).Scan(&lastInsertId)

	if err != nil {
		return &activity, ErrHandler(err)
	}

	sqlStatement1 := fmt.Sprintf("SELECT id,0,data_inserimento FROM attivita where id = %d", lastInsertId)
	rows, err := db.Query(sqlStatement1)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Activity
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		activity = p
	}
	if err != nil {
		return &activity, lang.Get("no_result")
	}
	return &activity, ""
}
