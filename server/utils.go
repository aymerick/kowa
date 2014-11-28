package server

import (
	"errors"
	"net/http"
	"strconv"
)

const DEFAULT_PER_PAGE = 20

// converts pagination parameters 'page' and 'per_page' to db query parameters 'skip' and 'limit'
func paginationParams(req *http.Request) (int, int, error) {
	var page int
	var perPage int

	var err error
	skip := 0
	limit := -1

	params := req.URL.Query()

	pageStr := params.Get("page")
	if pageStr != "" {
		if page, err = strconv.Atoi(pageStr); err == nil {
			perPageStr := params.Get("per_page")
			if perPageStr == "" {
				perPage = DEFAULT_PER_PAGE
			} else {
				perPage, err = strconv.Atoi(perPageStr)
			}

			if err == nil {
				if (page < 1) || (perPage < 1) {
					err = errors.New("Invalid pagination parameters")
				} else {
					skip = (page - 1) * perPage
					limit = perPage
				}
			}
		}
	}

	return skip, limit, err
}
