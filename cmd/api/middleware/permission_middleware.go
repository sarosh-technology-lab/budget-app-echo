package middlewares

import (
    "github.com/labstack/echo/v4"
    "budget-backend/internal/models"
	"budget-backend/common"
)

func PermissionMiddleware(requiredPermission string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            roleID, ok := c.Get("role_id").(uint)
            if !ok {
                return common.SendInternalServerErrorResponse(c, "role_id not found")
            }

			db, err := common.Sql()
				if err != nil {
					panic(err)
				}

            // Check if the role has the required permission
            var count int64
           db.Model(&models.RolePermission{}).
                Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
                Where("role_permissions.role_id = ? AND permissions.name = ?", roleID, requiredPermission).
                Count(&count)

            if count == 0 {
                return common.SendForbiddenResponse(c, "permission denied")
            }

            // Proceed to the next handler
            return next(c)
        }
    }
}