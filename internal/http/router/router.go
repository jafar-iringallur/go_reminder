package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-reminder/internal/http/handlers"
	"go-reminder/internal/service"
)

func New(reminderService service.ReminderService, taskService service.TaskService, auditService service.AuditService) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	ruleHandler := handlers.NewReminderRuleHandler(reminderService)
	taskHandler := handlers.NewTaskHandler(taskService)
	auditHandler := handlers.NewAuditHandler(auditService)

	api := router.Group("/api/v1")
	{
		rules := api.Group("/reminder-rules")
		{
			rules.POST("", ruleHandler.Create)
			rules.GET("", ruleHandler.List)
			rules.PUT("/:id", ruleHandler.Update)
			rules.DELETE("/:id", ruleHandler.Delete)
			rules.PATCH("/:id/activate", ruleHandler.Activate)
			rules.PATCH("/:id/deactivate", ruleHandler.Deactivate)
		}

		api.GET("/tasks", taskHandler.List)
		api.GET("/audit-logs", auditHandler.List)
	}

	return router
}
