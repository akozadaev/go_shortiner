package adapter

import (
	"context"
	"net/http"
	"strconv"
)

const defaultOffset = 0
const defaultPerPage = 20

const paramPerPage = "per_page"
const paramOffset = "offset"

const maxPerPage = 100

const Pagination contextKey = "pagination"

type PaginationAdapter struct {
	perPage   uint32
	offset    uint32
	timestamp int64
}

func NewPagination(req *http.Request) *PaginationAdapter {
	adapter := PaginationAdapter{
		perPage:   defaultPerPage,
		offset:    defaultOffset,
		timestamp: 0,
	}

	if req == nil {
		return &adapter
	}

	limit, err := strconv.Atoi(req.URL.Query().Get(paramPerPage))
	if err == nil {
		if limit > maxPerPage {
			limit = maxPerPage
		}
		adapter.perPage = uint32(limit)
	}

	offset, err := strconv.Atoi(req.URL.Query().Get(paramOffset))
	if err == nil {
		adapter.offset = uint32(offset)
	}

	return &adapter
}

func GetPagination(ctx context.Context) (*PaginationAdapter, bool) {
	adapter, ok := ctx.Value(Pagination).(*PaginationAdapter)
	return adapter, ok
}

func (a *PaginationAdapter) GetLimit() uint32 {
	return a.perPage
}

func (a *PaginationAdapter) GetOffset() uint32 {
	return a.offset
}
