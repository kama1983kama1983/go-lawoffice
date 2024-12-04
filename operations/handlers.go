package operations

import (
	"github.com/labstack/echo/v4"
)

func DashboardRouter(e *echo.Echo){
   e.GET("/",DashboardTemplate)
}

func ClientsRouter(e *echo.Echo) {
	// Clients
	e.GET("/clients", ClientsTemplate)
	e.GET("/clients/api", GetClients)
	e.GET("/client/create", ClientsAddTemplate)
	e.GET("/client/edit/:id", GetClientTemplate)
	e.POST("/client/add", CreateClient)
	e.POST("/client/update/:id", UpdateClient)       // Edit client
	e.GET("/client/delete/:id", DeleteClientHandler) // Delete client
}

func CasesRouter(e *echo.Echo) {
	e.GET("/cases", CasesTemplate)
	e.GET("/case/create", CasesAddTemplate)
	e.GET("/case/edit/:id", GetCaseOne)
	e.GET("/case/delete/:id", DeleteCaseHandler)
	e.POST("/case/update/:id", UpdateCase)
	e.POST("/case/add", CreateCase)
	e.GET("/cases/api", GetCases)
}

func AppointRouter(e *echo.Echo){
	e.GET("/appoinments", AppointmentsTemplate)
	e.GET("/appoinment/create", AppointmentAddTemplate)
	e.GET("/appoinment/edit/:id", GetAppointmentOne)
	e.GET("/appoinment/delete/:id", DeleteAppointmentHandler)
	e.POST("/appoinment/update/:id", UpdateAppoinment)
	e.POST("/appoinment/add", CreateAppoinment)
	e.GET("/appoinment/api", GetAppointments)
}

func LawyersRouter(e *echo.Echo){
	e.GET("/lawyers", LawyerTemplate)
	e.GET("/lawyer/create",LawyerAddTemplate)
	e.GET("/lawyer/edit/:id", GetLawyerOne)
	e.GET("/lawyer/delete/:id", DeleteLawyerHandler)
	e.POST("/lawyer/update/:id", UpdateLawyer)
	e.POST("/lawyer/add", CreateLawyer)
	e.GET("/lawyers/api", GetLawyers)
}

func PowerRouter(e *echo.Echo){
	e.GET("/powers",PowerTemplate)
	e.GET("/power/create", PowerAddTemplate)
	e.GET("/power/edit/:id", GetPowerOne)
	e.GET("/power/delete/:id", DeletePowerHandler)
	e.POST("/power/update/:id", UpdatePower)
	e.POST("/power/add", CreatePower)
	e.GET("/power/api", GetPowers)
}

func SessionsRouter(e *echo.Echo){
	e.GET("/sessions", SessionsTemplate)
	e.GET("/session/create", SessionsAddTemplate)
	e.GET("/session/edit/:id", GetSessionOne)
	e.GET("/session/delete/:id", DeleteSessionHandler)
	e.POST("/session/update/:id", UpdateSession)
	e.POST("/session/add", CreateSession)
	e.GET("/sessions/api", GetSessions)
}


func InovicesRouter(e *echo.Echo){
	e.GET("/inovices", InvoicesTemplate)
	e.GET("/inovice/create", CreateInvoiceTemplate)
	e.GET("/inovice/edit/:id", getInvoice)
	e.GET("/inovice/delete/:id", deleteInvoice)
	e.POST("/inovice/update/:id", updateInvoice)
	e.POST("/inovice/add", createInvoice)
	e.GET("/inovice/api", getInvoices)
}