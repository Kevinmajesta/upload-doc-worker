package routes

import (
	"kevinmajesta/karyawan/controllers"
	"kevinmajesta/karyawan/middlewares"
	"kevinmajesta/karyawan/helpers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {


	router := gin.Default()


	router.POST("/api/register", controllers.Register)

	router.POST("/api/login", controllers.Login)

	authorized := router.Group("/api", middlewares.AuthMiddleware())
	{
		authorized.POST("/employee", controllers.CreateEmployee)
		authorized.PUT("/employee/:id", controllers.UpdateEmployee)
		authorized.DELETE("/employee/:id", controllers.DeleteEmployee)
		authorized.GET("/employees", controllers.GetAllEmployees)
		authorized.GET("/employees/:id", controllers.GetEmployeeByID)
		//document
		authorized.POST("/documents", helpers.CreateDocumentAsync)
		authorized.GET("/documents", controllers.GetDocuments)
		authorized.GET("/documents/:id", controllers.GetDocumentByID)
		authorized.PUT("/documents/:id", controllers.UpdateDocument)
		authorized.DELETE("/documents/:id", controllers.DeleteDocument)
	}

	return router
}
