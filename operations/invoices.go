package operations

// Invoice struct
import(
    "database/sql"
	"lawoffice/config"
	"net/http"
	"github.com/labstack/echo/v4"
)

type Invoice struct {
	InvoiceID int `json:"invoice_id"`

	ClientID int `json:"client_id"`

	CaseID int `json:"case_id"`

	Amount float64 `json:"amount"`

	Status string `json:"status"`

	DueDate string `json:"due_date"`
}

func CreateInvoiceTemplate(c echo.Context) error{
    return config.RenderTemplate(c,"invoices/add-invoice.html",nil)
}

func InvoicesTemplate(c echo.Context) error{
    return config.RenderTemplate(c,"invoices/invoices.html",nil)
}

func EditInvoiceTemplate(c echo.Context,data interface{}) error{    
    return config.RenderTemplate(c,"invoices/edit-invoice.html",data)
}

// Create Invoice
func createInvoice(c echo.Context) error {
	invoice := Invoice{}
	if err := c.Bind(&invoice); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	query := "INSERT INTO Invoices (client_id, case_id, amount, status, due_date) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, invoice.ClientID, invoice.CaseID, invoice.Amount, invoice.Status, invoice.DueDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, invoice)

}

// Get All Invoices

func getInvoices(c echo.Context) error {
	invoices := []Invoice{}
	rows, err := db.Query("SELECT * FROM Invoices")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer rows.Close()
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(&invoice.InvoiceID, &invoice.ClientID, &invoice.CaseID, &invoice.Amount, &invoice.Status, &invoice.DueDate); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		invoices = append(invoices, invoice)
	}
	return c.JSON(http.StatusOK, invoices)
}
// Get Invoice by ID

func getInvoice(c echo.Context) error {
	id := c.Param("id")
	invoice := Invoice{}
	query := "SELECT * FROM Invoices WHERE invoice_id = ?"
	row := db.QueryRow(query, id)
	err := row.Scan(&invoice.InvoiceID, &invoice.ClientID, &invoice.CaseID, &invoice.Amount, &invoice.Status, &invoice.DueDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, "Invoice not found")
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return EditInvoiceTemplate(c,row)
}

// Update Invoice
func updateInvoice(c echo.Context) error {
	id := c.Param("id")
	invoice := Invoice{}
	if err := c.Bind(&invoice); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	query := "UPDATE Invoices SET client_id = ?, case_id = ?, amount = ?, status = ?, due_date = ? WHERE invoice_id = ?"
	_, err := db.Exec(query, invoice.ClientID, invoice.CaseID, invoice.Amount, invoice.Status, invoice.DueDate, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, invoice)

}

// Delete Invoice

func deleteInvoice(c echo.Context) error {
	id := c.Param("id")
	query := "DELETE FROM Invoices WHERE invoice_id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
