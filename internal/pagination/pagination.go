package pagination

import (
	"net/http"
	"strconv"
)

type Params struct {
	Limit  int
	Offset int
}

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

// ParseFromRequest extrai limit e offset da query string
func ParseFromRequest(r *http.Request) Params {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	return Params{Limit: limit, Offset: offset}
}
