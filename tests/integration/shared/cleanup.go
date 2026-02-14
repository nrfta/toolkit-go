package shared

import (
	"context"
	"database/sql"
	"fmt"
)

// TruncateAll removes all data from test tables while preserving schema
func TruncateAll(ctx context.Context, db *sql.DB) error {
	tables := []string{"orders", "tags", "products", "users"}
	for _, table := range tables {
		_, err := db.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}
	return nil
}
