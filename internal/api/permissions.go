package api

import (
	"fmt"
	"net/http"

	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/enums"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func forbidByRole(required string, actual enums.UserRole) core.Response {
	return core.Err(http.StatusForbidden, fmt.Errorf("insufficient permissions: required %s role, got %s", required, actual.String()))
}

func allowForAdmin(user *models.User) core.Response {
	if user.Role != enums.UserRoleAdmin {
		return forbidByRole("admin", user.Role)
	}
	return nil
}

func allowForModeratorOrAdmin(user *models.User) core.Response {
	if user.Role != enums.UserRoleModerator && user.Role != enums.UserRoleAdmin {
		return forbidByRole("moderator or admin", user.Role)
	}
	return nil
}

func allowForOrderManager(user *models.User) core.Response {
	if user.Role != enums.UserRoleManager && user.Role != enums.UserRoleModerator && user.Role != enums.UserRoleAdmin {
		return forbidByRole("manager, moderator or admin", user.Role)
	}
	return nil
}
