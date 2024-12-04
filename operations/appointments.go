package operations 

import(
  "database/sql"
  "lawoffice/config"
  "net/http"
  "strconv"
  "github.com/labstack/echo/v4"
)

type Appointments struct{
	Clientid int  `json:"clientid"`
    Lawyerid int  `json:"clientid"`
	AppointmentDate string  `json:"appdate"`
	Location string `json:"location"`
	Notes string `json:"notes"`
}

func AppointmentAddTemplate(e echo.Context)error{
	lawyers := SelectLawyers(e)
	clients := SelectClients(e)
	_, lcount := LawyerCount()
	_, ccount := ClientsCount()
	data := map[string]interface{}{
		"lawyers": lawyers,
		"clients": clients,
		"lcount":  lcount,
		"ccount":  ccount,
	}
	return config.RenderTemplate(e,"appointments/add-appointment.html",data)
}

func AppointmentEditTemplate(e echo.Context,data interface{})error{
	return config.RenderTemplate(e,"appointments/edit-appointment.html",data)
}

func GetAppointmentOne(e echo.Context)error{
	var appo Appointments
	ids := e.Param("id")
	id,err := strconv.Atoi(ids)
	if err != nil{
		return e.JSON(http.StatusBadRequest,map[string]string{"message" : "Error : id is not string"})
	}
	g := "select client_id,lawyer_id,appointment_date,location,notes from appointments where appointment_id=?"
	err = db.QueryRow(g,id).Scan(
		&appo.Clientid ,
		&appo.Lawyerid,
		&appo.AppointmentDate,
		&appo.Location,
		&appo.Notes,
	)
	if err != nil{
		if err == sql.ErrNoRows{
			return e.JSON(http.StatusNotFound,map[string]string{"message" : "No Rows"})
		}
		return e.JSON(http.StatusNotFound,map[string]string{"message" : "Appoinments Not Relative!!"})
	}
	data := map[string]interface{}{
		"appo" : appo,
		"clients" : SelectClients(e),
		"lawyers" : SelectLawyers(e),
	}
	return AppointmentEditTemplate(e,data)
}

func CreateAppoinment(e echo.Context)error{
   Lawyerid := e.FormValue("lawyerid")
   Clientid := e.FormValue("clientid")
   AppointmentDate := e.FormValue("appointmentdate")
   Location := e.FormValue("location")
   Notes := e.FormValue("notes")
   qr := "insert into appointments (client_id,lawyer_id,appointment_date,location,notes) VALUES(?,?,?,?,?)"
   _ , err := db.Exec(qr,Clientid ,Lawyerid ,AppointmentDate,Location,Notes)
   if err != nil{
	   return e.JSON(http.StatusNotFound,map[string]string{"message" : "Insert opertion not compelete!"})
   }
   return e.JSON(http.StatusOK , map[string]string{"message" : "Appoinment is Success Inserted!"})   
}

func UpdateAppoinment(e echo.Context)error{
	Lawyerid := e.FormValue("lawyerid")
	Clientid := e.FormValue("clientid")
	AppointmentDate := e.FormValue("appointmentdate")
	Location := e.FormValue("location")
	Notes := e.FormValue("notes")
	qr := "UPDATE FROM appointments SET client_id = ? ,lawyer_id = ? ,appointment_date = ? ,location = ? ,notes = ?) VALUES(?,?,?,?,?)"
	_ , err := db.Exec(qr,Clientid ,Lawyerid ,AppointmentDate,Location,Notes)
	if err != nil{
		return e.JSON(http.StatusNotFound,map[string]string{"message" : "Insert opertion not compelete!"})
	}
	return e.JSON(http.StatusOK , map[string]string{"message" : "Appoinment is Success Inserted!"})   
 }

func AppointmentsTemplate(c echo.Context) error {
	return config.RenderTemplate(c, "appointments/appointments.html", nil)
}
func GetAppointments(e echo.Context)error{
	g := "select client_id,lawyer_id,appointment_date,location,notes from appointments"
	rows , err := db.Query(g)
	if err != nil{
		
		return e.JSON(http.StatusBadRequest , map[string]string{"message" : "Now"})
	}
	var appos []Appointments
	for rows.Next(){
		var appo Appointments
		err := rows.Scan(
			&appo.Clientid , 
			&appo.Lawyerid , 
			&appo.AppointmentDate,
			&appo.Location,
			&appo.Notes,              
		 )
		 if err != nil{

			if err == sql.ErrNoRows{
				return e.JSON(http.StatusNotFound , map[string]string{"message" : "No Rows Found!"})
			} 
			return e.JSON(http.StatusNotFound,map[string]string{"message" : "Not Found in Rows!"})
		 }
		 appos = append(appos,appo)
	}
	return e.JSON(http.StatusOK,appos)
}

func DeleteAppointmentHandler(e echo.Context) error {
	idstr := e.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return e.JSON(200, map[string]string{"message": "invaild id string!"})
	}
	quer := "Delete from Appointments where appointment_id = ?"
	_, err = db.Exec(quer, id)
	if err != nil {
		return e.JSON(200, map[string]string{"message": "invaild id string!"})
	}
	return e.Redirect(http.StatusSeeOther, "/appointments")
}