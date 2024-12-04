package operations 

import(
	"database/sql"
	"lawoffice/config"
	"net/http"
	"log"
	"strconv"
	"github.com/labstack/echo/v4"
)

var (

	qrar       string

	today       string

	nextdate    string
)


func SessionsAddTemplate(e echo.Context)error{
  return config.RenderTemplate(e, "sessions/add-sessions.html",nil)
}

func SessionsTemplate(e echo.Context)error{
	return config.RenderTemplate(e, "sessions/sessions.html",nil)
}
  

func SessionsEditTemplate(e echo.Context,data interface{})error{
	return config.RenderTemplate(e, "sessions/edit-sessions.html",data)
}

func GetSessionOne(e echo.Context)error{	
	ids := e.Param("id")
	if ids == ""{
		return e.JSON(http.StatusBadRequest,map[string]string{"message" : "ID is require"})
	}
	id , err := strconv.Atoi(ids)
	if err != nil{
		return e.JSON(http.StatusBadRequest,map[string]string{"message" : "this id is not string"})
	}
	q := `
	   select Case_id,qrar,today,next_date,description from Sessions where id = ?
	`	
	err = db.QueryRow(q,id).Scan(&caseid,&qrar,&today,&nextdate,&description)
	if err != nil{
		if err == sql.ErrNoRows {
            return e.JSON(http.StatusNotFound, map[string]string{"message": "Data not found"})
        }
		return e.JSON(http.StatusInternalServerError,map[string]string{"message" : "Error Not Found this data"})
	}

	data := map[string]interface{}{
		"caseid" : caseid,
		"qrar" : qrar,
		"today" : today,
		"nextdate" : nextdate,
		"description" : description,
	}
	return SessionsEditTemplate(e,data)
}

func GetSessions(e echo.Context) error {
	q := `
	   select Case_id, qrar, today, next_date, description from Sessions
	`
	rows, err := db.Query(q)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"message": "Error: select is not found!"})
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var caseid, qrar, today, nextdate, description string
		err = rows.Scan(&caseid, &qrar, &today, &nextdate, &description)
		if err != nil {
			if err == sql.ErrNoRows {
				return e.JSON(http.StatusNotFound, map[string]string{"message": "Rows Not Found"})
			}
			return e.JSON(http.StatusBadRequest, map[string]string{"message": "Error Not Found"})
		}
		data := map[string]interface{}{
			"caseid":     caseid,
			"qrar":       qrar,
			"today":      today,
			"nextdate":   nextdate,
			"description": description,
		}
		results = append(results, data)
	}

	if len(results) == 0 {
		return e.JSON(http.StatusNotFound, map[string]string{"message": "No Rows Found!"})
	}

	return e.JSON(http.StatusOK, results)
}

func CreateSession(e echo.Context) error {
	caseid := e.FormValue("caseid")
	qrar := e.FormValue("qrars")
	today := e.FormValue("today")
	nextdate := e.FormValue("nextdate")
	description := e.FormValue("description")
	q := `
	  insert into Sessions (Case_id,qrar,today,next_date,description)
	`
	_ , err := db.Exec(q,caseid,qrar,today,nextdate,description)
	if err != nil{
		return e.JSON(http.StatusBadRequest,map[string]string{"message" : "Error insert to database!"})
	}
	return e.String(http.StatusOK ,"Sessions is success")
}
func UpdateSession(e echo.Context) error {
	id := e.Param("id")
	sessionID, err := strconv.Atoi(id)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"message": "Error: id not found"})
	}

	caseid := e.FormValue("caseid")
	qrar := e.FormValue("qrars")
	today := e.FormValue("today")
	nextdate := e.FormValue("nextdate")
	description := e.FormValue("description")

	q := `UPDATE Sessions SET Case_id = ?, qrar = ?, today = ?, next_date = ?, description = ? WHERE session_id = ?`
	result, err := db.Exec(q, caseid, qrar, today, nextdate, description, sessionID)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"message": "Error inserting to database!"})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"message": "Error getting rows affected"})
	}

	if rowsAffected == 0 {
		return e.JSON(http.StatusNotFound, map[string]string{"message": "Session not found"})
	}

	return e.String(http.StatusOK, "Session updated successfully")
}

func CasesCount() (error, int) {
	var casescount int
	d := "select count(*) from Cases"
	err := db.QueryRow(d).Scan(&casescount)
	if err != nil {
		return err, 0
	}
	return nil, casescount
}

func SelectCases(c echo.Context) interface{} {
	// Execute the query to select clients
	rows, err := db.Query("SELECT case_number,mod3le FROM Cases")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Initialize a slice to hold client data
	var cases []map[string]interface{}

	// Iterate through the rows
	for rows.Next() {
		var casenum int
		var mod3le string
		if err := rows.Scan(&casenum,&mod3le); err != nil {
			log.Fatal(err)
		}
		// Create a map for the current client and append it to the clients slice
		cas := map[string]interface{}{
			"casenum":        casenum,
			"mod3le": mod3le,
		}
		cases = append(cases, cas)
	}
	// Check for errors after iterating through rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	// Return the clients as a JSON response
	return cases
}

func DeleteSessionHandler(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return c.JSON(200, map[string]string{"message": "invaild id string!"})
	}
	quer := "delete from Sessions where session_id = ?"
	_, err = db.Exec(quer, id)
	if err != nil {
		return c.JSON(200, map[string]string{"message": "invaild id string!"})
	}
	return c.Redirect(http.StatusSeeOther, "/cases")
}
