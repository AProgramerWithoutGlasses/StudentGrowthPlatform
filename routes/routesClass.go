package routes

import (
	"github.com/gin-gonic/gin"
	"studentGrow/controller/class"
)

func routesClass(r *gin.Engine) {
	rc := r.Group("/class")

	rc.POST("/list", class.GetClassListControl)
	rc.POST("/get_class_by_grade", class.GetClassByGradeControl)
}
