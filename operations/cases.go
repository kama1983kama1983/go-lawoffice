package operations

import (
	"database/sql"
	"lawoffice/config"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Case struct {
	CaseID      int    `json:"case_id"`
	ClientID    int    `json:"client_id"`
	LawyerID    int    `json:"lawyer_id"`
	CaseNumber  string `json:"case_number"`
	CaseType    string `json:"case_type"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Court       string `json:"court"`
	Dayra       string `json:"dayra"`
	Subject     string `json:"subject"`
	Mod3le      string `json:"mod3le"`
	Ed3aDate    string `json:"ed3a_date"` // Consider using time.Time for actual date handling
	Qa3aNum     string `json:"qa3a_num"`
}

func CasesAddTemplate(c echo.Context) error {
	lawyers := SelectLawyers(c)
	clients := SelectClients(c)
	_, lcount := LawyerCount()
	_, ccount := ClientsCount()
	data := map[string]interface{}{
		"lawyers": lawyers,
		"clients": clients,
		"lcount":  lcount,
		"ccount":  ccount,
	}
	return config.RenderTemplate(c, "cases/add-case.html", data)
}

func CasesEditTemplate(c echo.Context, Case interface{}) error {
	return config.RenderTemplate(c, "cases/edit-case.html", Case)
}

// GetCaseOne retrieves a case by ID

func GetCaseOne(c echo.Context) error {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Case ID"})
	}
	var caseDetails Case
	// Query the database for the case
	err = db.QueryRow("SELECT case_id, client_id, lawyer_id, case_number, case_type, status, description, court, dayra, subject, mod3le, ed3a_date, qa3a_num FROM Cases WHERE case_id = ?", id).Scan(
		&caseDetails.CaseID,
		&caseDetails.ClientID,
		&caseDetails.LawyerID,
		&caseDetails.CaseNumber,
		&caseDetails.CaseType,
		&caseDetails.Status,
		&caseDetails.Description,
		&caseDetails.Court,
		&caseDetails.Dayra,
		&caseDetails.Subject,
		&caseDetails.Mod3le,
		&caseDetails.Ed3aDate,
		&caseDetails.Qa3aNum,
	)
   
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Case not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error retrieving case"})
	}

	data := map[string]interface{}{
		  "caseDetails" : caseDetails,
		  "clients" : SelectClients(c),
		  "lawyers" : SelectLawyers(c),
	}
	// Return the case details as JSON
	return CasesEditTemplate(c, data)
}



func CreateCase(c echo.Context) error {

	// Retrieve form values

	ClientID := c.FormValue("clientID")

	LawyerID := c.FormValue("lawyerID")

	CaseNumber := c.FormValue("caseNumber")

	CaseType := c.FormValue("caseType")

	Status := c.FormValue("Status")

	Description := c.FormValue("Description")

	Court := c.FormValue("Court")

	Dayra := c.FormValue("Dayra")

	Subject := c.FormValue("Subject")

	Mod3le := c.FormValue("Mod3le")

	Ed3aDate := c.FormValue("Ed3aDate")

	Qa3aNum := c.FormValue("Qa3aNum")

	// Prepare the SQL query

	q := `INSERT INTO Cases (

		client_id,

		lawyer_id,

		case_number,

		case_type,

		status,

		description,

		court,

		dayra,

		subject,

		mod3le,

		ed3a_date,

		qa3a_num) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`


	// Execute the query (assuming db is already defined and connected)

	result, err := db.Exec(q,

		ClientID,

		LawyerID,

		CaseNumber,

	 CaseType,

		Status,

		Description,

		Court,

		Dayra,

		Subject,

		Mod3le,

		Ed3aDate,

		Qa3aNum)


	// Check for errors

	if err != nil {

		log.Printf("Error executing query: %v", err)

		return c.JSON(http.StatusInternalServerError, echo.Map{

			"message": "Error creating case",

			"error":   err.Error(),

		})

	}


	// Check if any rows were affected

	rowsAffected, err := result.RowsAffected()

	if err != nil {

		log.Printf("Error getting rows affected: %v", err)

		return c.JSON(http.StatusInternalServerError, echo.Map{

			"message": "Error checking affected rows",

			"error":   err.Error(),

		})

	}


	if rowsAffected == 0 {

		log.Println("No rows were inserted.")

		return c.JSON(http.StatusInternalServerError, echo.Map{

			"message": "No rows were inserted.",

		})

	}
	return c.Redirect(http.StatusSeeOther, "/cases")
}

func CasesTemplate(c echo.Context) error {
	return config.RenderTemplate(c, "cases/cases.html", nil)
}

func GetCases(c echo.Context) error {
	// Define the SQL query
	q := "SELECT case_id, client_id, lawyer_id, case_number, case_type, status, description, court, dayra, subject, mod3le, ed3a_date, qa3a_num FROM Cases"
	// Execute the query
	rows, err := db.Query(q)
	if err != nil {
		log.Println("Error querying cases:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve cases"})
	}
	defer rows.Close()
	// Prepare a slice to hold the cases
	var cases []Case
	// Iterate through the rows
	for rows.Next() {
		var cas Case
		if err := rows.Scan(&cas.CaseID, &cas.ClientID, &cas.LawyerID, &cas.CaseNumber, &cas.CaseType, &cas.Status, &cas.Description, &cas.Court, &cas.Dayra, &cas.Subject, &cas.Mod3le, &cas.Ed3aDate, &cas.Qa3aNum); err != nil {
			log.Println("Error scanning case:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan case"})
		}
		cases = append(cases, cas)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		log.Println("Error during rows iteration:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error during rows iteration"})
	}

	// Return the cases as a JSON response
	return c.JSON(http.StatusOK, cases)
}

func UpdateCase(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return c.JSON(200, map[string]string{"message": "invaild id string!"})
	}
	ClientID := c.FormValue("clientID")
	LawyerID := c.FormValue("lawyerID")
	CaseNumber := c.FormValue("caseNumber")
	CaseType := c.FormValue("caseType")
	Status := c.FormValue("status")
	Description := c.FormValue("description")
	Court := c.FormValue("court")
	Dayra := c.FormValue("dayra")
	Subject := c.FormValue("subject")
	Mod3le := c.FormValue("mod3le")
	Ed3aDate := c.FormValue("ed3aDate")
	Qa3aNum := c.FormValue("qa3aNum")

	// Prepare the SQL update statement

	query := `
		UPDATE Cases
		SET 
			client_id = ?, 
			lawyer_id = ?, 
			case_number = ?, 
			case_type = ?, 
			status = ?, 
			description = ?, 
			court = ?, 
			dayra = ?, 
			subject = ?, 
			mod3le = ?, 
			ed3a_date = ?, 
			qa3a_num = ?
		WHERE case_id = ?`
	// Execute the update statement
	_, err = db.Exec(query, ClientID, LawyerID, CaseNumber, CaseType, Status, Description, Court, Dayra, Subject, Mod3le, Ed3aDate, Qa3aNum, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update case"})
	}
	// Return a success response
	return c.JSON(http.StatusOK, map[string]string{"message": "Case updated successfully"})
}

func DeleteCaseHandler(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return c.JSON(200, map[string]string{"message": "invaild id string!"})
	}
	quer := "delete from Cases where case_id = ?"
	_, err = db.Exec(quer, id)
	if err != nil {
		return c.JSON(200, map[string]string{"message": "invaild id string!"})
	}
	return c.Redirect(http.StatusSeeOther, "/cases")
}

func LawyerCount() (error, int) {
	var lawyercount int
	d := "select count(*) from Lawyers"
	err := db.QueryRow(d).Scan(&lawyercount)
	if err != nil {
		return err, 0
	}
	return nil, lawyercount
}

func ClientsCount() (error, int) {
	var clientscount int
	d := "select count(*) from Clients"
	err := db.QueryRow(d).Scan(&clientscount)
	if err != nil {
		return err, 0
	}
	return nil, clientscount
}

func SelectClients(c echo.Context) interface{} {
	// Execute the query to select clients
	rows, err := db.Query("SELECT client_id, first_name, last_name FROM Clients")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Initialize a slice to hold client data
	var clients []map[string]interface{}

	// Iterate through the rows
	for rows.Next() {
		var id int
		var firstname, lastname string
		if err := rows.Scan(&id, &firstname, &lastname); err != nil {
			log.Fatal(err)
		}

		// Create a map for the current client and append it to the clients slice
		client := map[string]interface{}{
			"id":        id,
			"firstname": firstname,
			"lastname":  lastname,
		}
		clients = append(clients, client)
	}

	// Check for errors after iterating through rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Return the clients as a JSON response
	return clients
}

func SelectLawyers(c echo.Context) interface{} {
	// Execute the query to select clients
	rows, err := db.Query("SELECT lawyer_id, first_name, last_name FROM Lawyers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Initialize a slice to hold client data
	var lawyers []map[string]interface{}

	// Iterate through the rows
	for rows.Next() {
		var id int
		var firstname, lastname string
		if err := rows.Scan(&id, &firstname, &lastname); err != nil {
			log.Fatal(err)
		}
		// Create a map for the current client and append it to the clients slice
		lawyer := map[string]interface{}{
			"id":        id,
			"firstname": firstname,
			"lastname":  lastname,
		}
		lawyers = append(lawyers, lawyer)
	}
	// Check for errors after iterating through rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	// Return the clients as a JSON response
	return lawyers
}
