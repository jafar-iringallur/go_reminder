package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-reminder/internal/service"
)

type ReminderRuleHandler struct {
	reminderService service.ReminderService
}

func NewReminderRuleHandler(reminderService service.ReminderService) *ReminderRuleHandler {
	return &ReminderRuleHandler{reminderService: reminderService}
}

func (h *ReminderRuleHandler) Create(c *gin.Context) {
	var input service.ReminderRuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.reminderService.CreateRule(input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (h *ReminderRuleHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var input service.ReminderRuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.reminderService.UpdateRule(id, input)
	respondRuleOrError(c, rule, err)
}

func (h *ReminderRuleHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.reminderService.DeleteRule(id); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ReminderRuleHandler) List(c *gin.Context) {
	rules, err := h.reminderService.ListRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rules)
}

func (h *ReminderRuleHandler) Activate(c *gin.Context) {
	h.setStatus(c, true)
}

func (h *ReminderRuleHandler) Deactivate(c *gin.Context) {
	h.setStatus(c, false)
}

func (h *ReminderRuleHandler) setStatus(c *gin.Context, active bool) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	rule, err := h.reminderService.SetRuleStatus(id, active)
	respondRuleOrError(c, rule, err)
}

func parseIDParam(c *gin.Context, key string) (uint, bool) {
	id, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return uint(id), true
}

func respondRuleOrError(c *gin.Context, rule any, err error) {
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, rule)
}

func respondError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, service.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
