// Code generated by sqlc. DO NOT EDIT.
// source: entity_types.sql

package db

import (
	"context"
)

const getEntityType_ByName = `-- name: GetEntityType_ByName :one
SELECT id, name, description, language, created_at
FROM entity_types 
WHERE name = $1
`

func (q *Queries) GetEntityType_ByName(ctx context.Context, name string) (EntityType, error) {
	row := q.queryRow(ctx, q.getEntityType_ByNameStmt, getEntityType_ByName, name)
	var i EntityType
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Language,
		&i.CreatedAt,
	)
	return i, err
}
