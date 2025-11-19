package response

import (
	"github.com/labstack/echo/v4"
)

func SendError(c echo.Context, err any, status int, msg string) error {
	return c.JSON(status, map[string]interface{}{
		"status":  status,
		"error":   err,
		"message": msg,
	})
}

func SendErrorWithData(c echo.Context, err any, status int, msg string, data any) error {
	return c.JSON(status, map[string]interface{}{
		"status":  status,
		"error":   err,
		"message": msg,
		"data":    data,
	})
}
