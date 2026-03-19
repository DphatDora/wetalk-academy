package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func BuildNextURL(c *gin.Context, total int64, page, limit int) string {
	if limit <= 0 || int64(page*limit) >= total {
		return ""
	}

	query := c.Request.URL.Query()
	query.Set("page", strconv.Itoa(page+1))
	query.Set("limit", strconv.Itoa(limit))

	return c.Request.URL.Path + "?" + query.Encode()
}
