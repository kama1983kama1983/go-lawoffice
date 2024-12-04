package operations 

import(
	"lawoffice/config"
	"github.com/labstack/echo/v4"
)

func DashboardTemplate(e echo.Context)error{
	return config.RenderTemplate(e,"admin/index.html",nil)
}