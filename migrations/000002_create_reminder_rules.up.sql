CREATE TABLE reminder_rules (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    task_status VARCHAR(30) NOT NULL,
    remind_before_due_minute INT NOT NULL,
    repeat_interval_minute INT NOT NULL DEFAULT 60,
    message_template TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    PRIMARY KEY (id),
    INDEX idx_reminder_rules_task_status (task_status),
    INDEX idx_reminder_rules_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
