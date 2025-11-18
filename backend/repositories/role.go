package repositories

import (
	"Inquiro/models"
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type RoleRepository struct {
	DB     *sql.DB
	logger *zap.SugaredLogger
}

func (r *RoleRepository) GetRoleByID(ctx context.Context, id int) (models.Role, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT id, name, level, description FROM role WHERE id = $1", id)
	role := models.Role{}
	err := row.Scan(&role.ID, &role.Name, &role.Level, &role.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Errorw("role id does not exists", "error :", err.Error())
			return models.Role{}, ErrUserNotFound
		}
		return models.Role{}, err
	}
	return role, nil
}
