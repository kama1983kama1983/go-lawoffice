package operations 

import(
	"database/sql"
	"lawoffice/config"
	"net/http"
	"strconv"
	"github.com/labstack/echo/v4"
)

type Lawyers struct{
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Notes     string `json:"notes"`
	Specialization   string `json:"specialization"`
}

func LawyerAddTemplate(e echo.Context)error{
	return config.RenderTemplate(e,"lawyers/add-lawyer.html",nil)
}

func LawyerTemplate(e echo.Context)error{
	return config.RenderTemplate(e,"lawyers/lawyers.html",nil)
}

func LawyerEditTemplate(e echo.Context,data interface{})error{
	return config.RenderTemplate(e,"lawyers/edit-lawyer.html",data)
}

func DeleteLawyerHandler(e echo.Context)error{
	ids := e.Param("id")
	id , err := strconv.Atoi(ids)
	q := "DELETE FROM Lawyers where lawyer_id = ?"
	_ , err = db.Query(q,id)
	if err != nil{
		return e.JSON(http.StatusNotFound , map[string]string{"message" : "Id Not found!"})
	}
	return e.Redirect(http.StatusSeeOther,"/lawyers")
}

func GetLawyerOne(e echo.Context)error{
	var x Lawyers
	ids := e.Param("id")
	id , err := strconv.Atoi(ids)
	if err != nil{
		return e.JSON(http.StatusBadRequest,map[string]string{"message":"Id not Found"})
	}
	q := "SELECT first_name , last_name , email , phone , specialization from Lawyers where lawyer_id = ?"
	err =  db.QueryRow(q,id).Scan(&x.FirstName,&x.LastName,&x.Email,&x.Phone,&x.Specialization,)
	if err != nil {
		if err == sql.ErrNoRows{
			return e.JSON(http.StatusInternalServerError,map[string]string{"message" : "Never Found Row!"})
		}
		return e.JSON(http.StatusNotFound , map[string]string{"message" :"Not Found Row!"})
	}
	return LawyerEditTemplate(e,x)
}

func GetLawyers(e echo.Context)error{
	q := "SELECT first_name , last_name , email , phone , specialization from Lawyers"
	rows , err := db.Query(q)
	if err != nil{
		return e.JSON(http.StatusInternalServerError, map[string]string{"message" : "Not Found Results " })
	}
	var xe []Lawyers
	for rows.Next(){
		var xr Lawyers 
		rows.Scan(
			&xr.FirstName,
			&xr.LastName,
			&xr.Email,
			&xr.Phone,
			&xr.Specialization,
		)
		xe = append(xe,xr)
	}
	return e.JSON(http.StatusOK,xe)
}



func CreateLawyer(c echo.Context) error {

	var lawyer Lawyers
	// Get form values

	lawyer.FirstName = c.FormValue("firstname")

	lawyer.LastName = c.FormValue("lastname")

	lawyer.Notes = c.FormValue("notes")

	lawyer.Email = c.FormValue("email")

	lawyer.Phone = c.FormValue("phone") // Ensure to get the phone value

	lawyer.Specialization = c.FormValue("specialization")


	// Prepare SQL query

	q := `

		INSERT INTO Lawyers (first_name, last_name, email, phone, specialization,  notes) 

		VALUES (?, ?, ?, ?, ?, ?, ?)

	`


	// Execute the query

	_, err := db.Exec(q, lawyer.FirstName, lawyer.LastName, lawyer.Email, lawyer.Phone, lawyer.Specialization, lawyer.Notes)

	if err != nil {

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error not inserted to database!"})

	}


	return c.JSON(http.StatusOK, map[string]string{"message": "Insertion is successful!"})

}


func UpdateLawyer(e echo.Context)error{
	firstname := e.FormValue("firstname")
	lastname := e.FormValue("lastname")
	location := e.FormValue("location")
	notes := e.FormValue("notes")
	email := e.FormValue("email")
	specialization := e.FormValue("specialization")	
	qr := "UPDATE FROM Lawyers SET first_name = ? ,last_name = ? ,email = ? ,phone = ? ,specialization = ?) VALUES(?,?,?,?,?)"
	_ , err := db.Exec(qr,firstname ,lastname ,location,notes,email,specialization)
	if err != nil{
		return e.JSON(http.StatusNotFound,map[string]string{"message" : "Insert opertion not compelete!"})
	}
	return e.JSON(http.StatusOK , map[string]string{"message" : "Appoinment is Success Inserted!"})   
 }
