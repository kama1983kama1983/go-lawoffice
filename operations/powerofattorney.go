package operations

import (
	"database/sql"
	"lawoffice/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

var (
	caseid,clientid, lawyerid                              int
	dateofissue, expirationdate, description, referancecode string
)

func PowerAddTemplate(e echo.Context) error {
	return config.RenderTemplate(e, "PowerOfAttorney/add-power.html", nil)
}

func PowerTemplate(c echo.Context) error {	
	return config.RenderTemplate(c, "PowerOfAttorney/powers.html", nil)
}

func PowerEditTemplate(e echo.Context, data interface{}) error {
	return config.RenderTemplate(e, "PowerOfAttorney/edit-power.html", data)
}

func GetPowerOne(e echo.Context) error {
	ids := e.Param("id")
	if ids == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"message": "ID is require"})
	}
	id, err := strconv.Atoi(ids)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"message": "id its not string"})
	}
	q := `
		select client_id, lawyer_id,case_id,date_of_issue,expiration_date,description,reference_code from PowerOfAttorney where power_of_attorney_id = ?
	 `
	err = db.QueryRow(q, id).Scan(&clientid, &lawyerid, &caseid, &dateofissue, &expirationdate, &description, &referancecode)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.JSON(http.StatusNotFound, map[string]string{"message": "Data not found"})
		}
		return e.JSON(http.StatusInternalServerError, map[string]string{"message": "Error Not Found this data"})
	}

	data := map[string]interface{}{
		"clientid":        clientid,
		"lawyerid":        lawyerid,
		"caseid":          caseid,
		"dateofissue":     dateofissue,
		"expiration_date": expirationdate,
		"description":     description,
		"referancecode":       referancecode,
	}
	return PowerEditTemplate(e, data)
}

func GetPowers(e echo.Context) error {

	q := `

    	select client_id, lawyer_id,case_id,date_of_issue,expiration_date,description,reference_code from PowerOfAttorney 

	`

	rows, err := db.Query(q)

	if err != nil {

		return e.JSON(http.StatusInternalServerError, map[string]string{"message": "Error : select is not found!"})

	}

	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {

		var clientid, lawyerid, caseid int

		var dateofissue, expirationdate, description, referencecode string

		err = rows.Scan(&clientid, &lawyerid, &caseid, &dateofissue, &expirationdate, &description, &referencecode)

		if err != nil {

			if err == sql.ErrNoRows {

				return e.JSON(http.StatusNotFound, map[string]string{"message": "Rows Not Found"})

			}

			return e.JSON(http.StatusBadRequest, map[string]string{"message": "Error Not Found "})

		}

		data := map[string]interface{}{

			"clientid":        clientid,

			"lawyerid":        lawyerid,

			"caseid":          caseid,

			"dateofissue":     dateofissue,

			"expiration_date": expirationdate,

			"description":     description,

			"referancecode":   referencecode,

		}

		results = append(results, data)

	}

	if len(results) == 0 {

		return e.JSON(http.StatusNotFound, map[string]string{"message": "No Rows Found!"})

	}

	return e.JSON(http.StatusOK, results)

}

func CreatePower(e echo.Context) error {

	clientid := e.FormValue("clientid")

	lawyerid := e.FormValue("lawyerid")

	caseid := e.FormValue("caseid")

	dateofissue := e.FormValue("dateofissue")

	expirationdate := e.FormValue("expirationdate")

	description := e.FormValue("description")

	reference := e.FormValue("reference")

	q := `

        INSERT INTO PowerOfAttorney (client_id,lawyer_id, case_id, date_of_issue, expiration_date, description, reference)

        VALUES (?, ?, ?, ?, ?, ?, ?)

    `

	_, err := db.Exec(q, clientid, lawyerid, caseid, dateofissue, expirationdate, description, reference)

	if err != nil {

		return e.JSON(http.StatusBadRequest, map[string]string{"message": "Error inserting into database!"})

	}

	return e.String(http.StatusOK, "PowerOfAttorney created successfully")

}


func DeletePowerHandler(e echo.Context)error{
	ids := e.Param("id")
	id , err := strconv.Atoi(ids)
	q := "DELETE FROM PowerOfAttorney where power_of_attorney_id = ?"
	_ , err = db.Query(q,id)
	if err != nil{
		return e.JSON(http.StatusNotFound , map[string]string{"message" : "Id Not found!"})
	}
	return e.Redirect(http.StatusSeeOther,"/powers")
}


func UpdatePower(e echo.Context) error {

	id := e.Param("id")

	sessionId, err := strconv.Atoi(id)

	if err != nil {

		return e.JSON(http.StatusBadRequest, map[string]string{"message": "Error: ID not found"})

	}

	clientid := e.FormValue("clientid")

	lawyerid := e.FormValue("lawyerid")

	caseid := e.FormValue("caseid")

	dateofissue := e.FormValue("dateofissue")

	expirationdate := e.FormValue("expiration_date")

	description := e.FormValue("description")

	reference := e.FormValue("reference")

	q := `

        UPDATE PowerOfAttorney 

        SET client_id = ?, cawyer_id = ?, case_id = ?, date_of_issue = ?, expiration_date = ?, description = ?, reference = ?

        WHERE power_of_attorney_id = ?

    `

	_, err = db.Exec(q, clientid, lawyerid, caseid, dateofissue, expirationdate, description, reference, sessionId)

	if err != nil {

		return e.JSON(http.StatusBadRequest, map[string]string{"message": "Error updating database!"})

	}

	return e.String(http.StatusOK, "Session updated successfully")

}
