package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetPaginationQuery_DefaultValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := createPaginationTestContext("/products")

	pagination := GetPaginationQuery(c)

	if pagination.Page != DefaultPage {
		t.Fatalf("expected page %d, got %d", DefaultPage, pagination.Page)
	}

	if pagination.Limit != DefaultLimit {
		t.Fatalf("expected limit %d, got %d", DefaultLimit, pagination.Limit)
	}

	if pagination.Offset != 0 {
		t.Fatalf("expected offset 0, got %d", pagination.Offset)
	}
}

func TestGetPaginationQuery_ValidValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := createPaginationTestContext("/products?page=3&limit=10")

	pagination := GetPaginationQuery(c)

	if pagination.Page != 3 {
		t.Fatalf("expected page 3, got %d", pagination.Page)
	}

	if pagination.Limit != 10 {
		t.Fatalf("expected limit 10, got %d", pagination.Limit)
	}

	if pagination.Offset != 20 {
		t.Fatalf("expected offset 20, got %d", pagination.Offset)
	}
}

func TestGetPaginationQuery_InvalidValuesUseDefaults(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := createPaginationTestContext("/products?page=abc&limit=-1")

	pagination := GetPaginationQuery(c)

	if pagination.Page != DefaultPage {
		t.Fatalf("expected default page %d, got %d", DefaultPage, pagination.Page)
	}

	if pagination.Limit != DefaultLimit {
		t.Fatalf("expected default limit %d, got %d", DefaultLimit, pagination.Limit)
	}

	if pagination.Offset != 0 {
		t.Fatalf("expected offset 0, got %d", pagination.Offset)
	}
}

func TestGetPaginationQuery_PageLessThanOneUsesDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := createPaginationTestContext("/products?page=0&limit=10")

	pagination := GetPaginationQuery(c)

	if pagination.Page != DefaultPage {
		t.Fatalf("expected default page %d, got %d", DefaultPage, pagination.Page)
	}

	if pagination.Limit != 10 {
		t.Fatalf("expected limit 10, got %d", pagination.Limit)
	}

	if pagination.Offset != 0 {
		t.Fatalf("expected offset 0, got %d", pagination.Offset)
	}
}

func TestGetPaginationQuery_LimitGreaterThanMaxIsCapped(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := createPaginationTestContext("/products?page=2&limit=999")

	pagination := GetPaginationQuery(c)

	if pagination.Page != 2 {
		t.Fatalf("expected page 2, got %d", pagination.Page)
	}

	if pagination.Limit != MaxLimit {
		t.Fatalf("expected limit capped at %d, got %d", MaxLimit, pagination.Limit)
	}

	if pagination.Offset != MaxLimit {
		t.Fatalf("expected offset %d, got %d", MaxLimit, pagination.Offset)
	}
}

func createPaginationTestContext(target string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodGet, target, nil)
	c.Request = req

	return c, w
}
