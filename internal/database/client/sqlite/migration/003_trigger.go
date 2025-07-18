package migration

import "github.com/bestruirui/bestsub/internal/database/migration"

func Migration003Trigger() string {
	return `
CREATE TRIGGER IF NOT EXISTS delete_sub_tasks
BEFORE DELETE ON subs
FOR EACH ROW
BEGIN
    DELETE FROM tasks
    WHERE id IN (
        SELECT task_id FROM sub_task_relations
        WHERE sub_id = OLD.id
    );
END;

CREATE TRIGGER IF NOT EXISTS delete_sub_save_tasks
BEFORE DELETE ON sub_save
FOR EACH ROW
BEGIN
    DELETE FROM tasks
    WHERE id IN (
        SELECT task_id FROM save_task_relations
        WHERE save_id = OLD.id
    );
END;
`
}

// init 自动注册迁移
func init() {
	migration.Register(ClientName, 202507171102, "dev", "Triggers", Migration003Trigger)
}
