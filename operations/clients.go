package operations

import (
	"database/sql"
	"lawoffice/config"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var db *sql.DB

// SetDB sets the database connection for the operations package
func SetDB(database *sql.DB) {
	db = database
}

func ClientsAddTemplate(c echo.Context) error {
	return config.RenderTemplate(c, "clients/add-client.html", nil)
}

func ClientsEditTemplate(c echo.Context) error {
	return config.RenderTemplate(c, "clients/edit-client.html", nil)
}

func GetClientTemplate(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Client ID"})
	}

	var client struct {
		ClientID  int    `json:"client_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Address   string `json:"address"`
	}

	err = db.QueryRow("SELECT client_id, first_name, last_name, email, phone, address FROM Clients WHERE client_id = ?", id).
		Scan(&client.ClientID, &client.FirstName, &client.LastName, &client.Email, &client.Phone, &client.Address)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
		}
		log.Printf("Error querying client: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	return config.RenderTemplate(c, "clients/edit-client.html", client)
}

func ClientsTemplate(c echo.Context) error {
	return config.RenderTemplate(c, "clients/clients.html", nil)
}

func GetClientOne(c echo.Context) error {
	// Extract the client ID from the URL path
	idStr := c.Param("id") // Use Echo's Param method to get the ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Client ID"})
	}

	// Prepare a variable to hold the client details

	var client struct {
		ClientID  int    `json:"client_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Address   string `json:"address"`
	}
	// Query the database for the client
	err = db.QueryRow("SELECT client_id, first_name, last_name, email, phone, address FROM Clients WHERE client_id = ?", id).
		Scan(&client.ClientID, &client.FirstName, &client.LastName, &client.Email, &client.Phone, &client.Address)
	// Handle errors from the query
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
		}
		log.Printf("Error querying client: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	// Return the client details in JSON format
	return c.JSON(http.StatusOK, client)

}

// GetClients retrieves clients from the database
func GetClients(c echo.Context) error {
	rows, err := db.Query("SELECT client_id, first_name, last_name, email, phone, address FROM Clients")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching clients: "+err.Error())
	}
	defer rows.Close()

	var clients []map[string]interface{}
	for rows.Next() {
		var clientID int
		var firstName, lastName, email, phone, address *string // Use pointers for nullable fields
		if err := rows.Scan(&clientID, &firstName, &lastName, &email, &phone, &address); err != nil {
			return c.String(http.StatusInternalServerError, "Error scanning client: "+err.Error())
		}

		clients = append(clients, map[string]interface{}{
			"client_id":  clientID,
			"first_name": firstName, // This will be returned as null if nil
			"last_name":  lastName,  // This will be returned as null if nil
			"email":      email,     // This will be returned as null if nil
			"phone":      phone,     // This will be returned as null if nil
			"address":    address,   // This will be returned as null if nil
		})
	}
	return c.JSON(http.StatusOK, clients)
}

// CreateClient adds a new client to the database
func CreateClient(c echo.Context) error {
	firstName := c.FormValue("first_name")
	lastName := c.FormValue("last_name")
	email := c.FormValue("email")
	phone := c.FormValue("phone")
	address := c.FormValue("address")

	_, err := db.Exec("INSERT INTO Clients (first_name, last_name, email, phone, address) VALUES (?, ?, ?, ?, ?)", firstName, lastName, email, phone, address)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creating client: "+err.Error())
	}
	return c.String(http.StatusOK, "Client created successfully")
}

// EditClient updates an existing client in the database
func UpdateClient(c echo.Context) error {
	id := c.Param("id")
	firstName := c.FormValue("first_name")
	lastName := c.FormValue("last_name")
	email := c.FormValue("email")
	phone := c.FormValue("phone")
	address := c.FormValue("address")

	_, err := db.Exec("UPDATE Clients SET first_name = ?, last_name = ?, email = ?, phone = ?, address = ? WHERE client_id = ?", firstName, lastName, email, phone, address, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating client: "+err.Error())
	}

	return c.String(http.StatusOK, "Client updated successfully")
}

// deleteClientHandler handles DELETE requests to remove a client
func DeleteClientHandler(c echo.Context) error {
	// Extract the client ID from the URL path
	idStr := c.Param("id") // Use Echo's Param method to get the ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Client ID"})
	}
	// Attempt to delete the client from the database
	result, err := db.Exec("DELETE FROM Clients WHERE client_id = ?", id)
	if err != nil {
		log.Printf("Error deleting client: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	// Check how many rows were affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	// If no rows were affected, the client was not found
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
	}
	// If the client was successfully deleted, redirect to the clients page
	return c.Redirect(http.StatusSeeOther, "/clients")
}
