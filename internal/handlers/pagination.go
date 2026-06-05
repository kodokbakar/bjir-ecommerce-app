package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

type PaginationQuery struct {
	Page   int
	Limit  int
	Offset int
}

func GetPaginationQuery(c *gin.Context) PaginationQuery {
	page := parseIntOrDefault(c.Query("page"), DefaultPage)
	if page < 1 {
		page = DefaultPage
	}

	limit := parseIntOrDefault(c.Query("limit"), DefaultLimit)
	if limit < 1 {
		limit = DefaultLimit
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}

	return PaginationQuery{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

func parseIntOrDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsedValue
}
