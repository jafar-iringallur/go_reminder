CREATE TABLE audit_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    event_type VARCHAR(80) NOT NULL,
    entity VARCHAR(80) NOT NULL,
    entity_id BIGINT UNSIGNED NOT NULL,
    rule_id BIGINT UNSIGNED NULL,
    task_id BIGINT UNSIGNED NULL,
    details JSON NOT NULL,
    created_at DATETIME(3) NULL,
    PRIMARY KEY (id),
    INDEX idx_audit_logs_event_type (event_type),
    INDEX idx_audit_logs_entity (entity),
    INDEX idx_audit_logs_entity_id (entity_id),
    INDEX idx_audit_logs_rule_id (rule_id),
    INDEX idx_audit_logs_task_id (task_id),
    INDEX idx_audit_logs_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
