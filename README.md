# Go Gin Dynamic Reminder System

REST API and background scheduler for configurable task reminders with an audit trail.

## Features

- Gin-based REST API.
- MySQL persistence through GORM.
- Versioned SQL migrations through `golang-migrate`.
- Seeder command similar to Laravel seeders.
- Structured project layout under `cmd` and `internal`.
- Seeded sample tasks and a default active reminder rule for quick demo.
- CRUD APIs for reminder rules.
- Activate/deactivate APIs for reminder rules.
- Background scheduler that checks active rules and prints simulated reminders.
- Audit logs for rule changes and reminder executions.

## Run

Create a MySQL database first:

```sql
CREATE DATABASE go_reminder CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

Install dependencies:

```bash
go mod tidy
```

Set your MySQL DSN:

```bash
export MYSQL_DSN='root:password@tcp(127.0.0.1:3306)/go_reminder?charset=utf8mb4&parseTime=True&loc=Local'
```

Run migrations:

```bash
go run ./cmd/migrate up
```

Seed demo data:

```bash
go run ./cmd/seed
```

Start the API:

```bash
go run ./cmd/server
```

Default server address: `http://localhost:8080`

Optional environment variables:

```bash
HTTP_ADDRESS=:8080
MYSQL_DSN=root:password@tcp(127.0.0.1:3306)/go_reminder?charset=utf8mb4&parseTime=True&loc=Local
SCHEDULER_INTERVAL_SECONDS=60
```

For demo purposes, set `SCHEDULER_INTERVAL_SECONDS=10` to run the scheduler every 10 seconds.

The server also runs pending migrations on startup. The explicit migrate command is still included so schema changes are managed in a Laravel-like, versioned way.

## Migrations and Seeders

Go does not have one built-in migration system like Laravel. In production Go APIs, it is common to use a dedicated migration tool such as `golang-migrate`, `goose`, or `Atlas`.

This project uses:

- SQL migrations in `migrations/`.
- `go run ./cmd/migrate up` to apply migrations.
- `go run ./cmd/migrate down` to roll back the latest migration.
- `go run ./cmd/seed` to insert demo tasks and a default reminder rule.

## API

### Health

```bash
curl http://localhost:8080/health
```

### Create Reminder Rule

```bash
curl -X POST http://localhost:8080/api/v1/reminder-rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Pending tasks due soon",
    "task_status": "pending",
    "remind_before_due_minutes": 60,
    "repeat_interval_minutes": 15,
    "message_template": "Reminder: {{task_title}} is due at {{due_at}}",
    "is_active": true
  }'
```

### List Reminder Rules

```bash
curl http://localhost:8080/api/v1/reminder-rules
```

### List Seeded Tasks

```bash
curl http://localhost:8080/api/v1/tasks
```

### Update Reminder Rule

```bash
curl -X PUT http://localhost:8080/api/v1/reminder-rules/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "In-progress tasks due soon",
    "task_status": "in_progress",
    "remind_before_due_minutes": 30,
    "repeat_interval_minutes": 10,
    "message_template": "{{rule_name}}: {{task_title}} is due at {{due_at}}",
    "is_active": true
  }'
```

### Activate / Deactivate

```bash
curl -X PATCH http://localhost:8080/api/v1/reminder-rules/1/activate
curl -X PATCH http://localhost:8080/api/v1/reminder-rules/1/deactivate
```

### Delete Reminder Rule

```bash
curl -X DELETE http://localhost:8080/api/v1/reminder-rules/1
```

### Audit Logs

```bash
curl http://localhost:8080/api/v1/audit-logs
curl "http://localhost:8080/api/v1/audit-logs?event_type=reminder.triggered&limit=20"
```

## Reminder Rule Fields

- `task_status`: one of `pending`, `in_progress`, `completed`.
- `remind_before_due_minutes`: finds matching tasks due from now until this many minutes ahead.
- `repeat_interval_minutes`: prevents repeated reminders for the same rule/task pair within the interval.
- `message_template`: supports `{{rule_name}}`, `{{task_title}}`, `{{task_status}}`, and `{{due_at}}`.
