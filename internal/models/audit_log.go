package models

import (
	"database/sql/driver"
	"errors"
	"time"
)

type AuditEventType string

const (
	AuditEventRuleCreated       AuditEventType = "reminder_rule.created"
	AuditEventRuleUpdated       AuditEventType = "reminder_rule.updated"
	AuditEventRuleDeleted       AuditEventType = "reminder_rule.deleted"
	AuditEventRuleStatusChanged AuditEventType = "reminder_rule.status_changed"
	AuditEventReminderTriggered AuditEventType = "reminder.triggered"
)

type AuditLog struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	EventType AuditEventType  `json:"event_type" gorm:"type:varchar(80);index;not null"`
	Entity    string          `json:"entity" gorm:"type:varchar(80);index;not null"`
	EntityID  uint            `json:"entity_id" gorm:"index;not null"`
	RuleID    *uint           `json:"rule_id,omitempty" gorm:"index"`
	TaskID    *uint           `json:"task_id,omitempty" gorm:"index"`
	Details   JSONDetails     `json:"details" gorm:"type:text;not null"`
	CreatedAt time.Time       `json:"created_at" gorm:"index"`
}

type JSONDetails []byte

func (j JSONDetails) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "{}", nil
	}
	return string(j), nil
}

func (j *JSONDetails) Scan(value any) error {
	switch data := value.(type) {
	case nil:
		*j = JSONDetails("{}")
	case []byte:
		*j = append((*j)[0:0], data...)
	case string:
		*j = append((*j)[0:0], data...)
	default:
		return errors.New("unsupported JSONDetails scan type")
	}
	return nil
}

func (j JSONDetails) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("{}"), nil
	}
	return j, nil
}
