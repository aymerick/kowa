package server

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
)

const defaultPerPage = 20

type pagination struct {
	Page       int
	PerPage    int
	TotalPages int
	NextPage   int
	PrevPage   int

	Skip  int
	Total int
}

// create pagination params from request
func newPagination() *pagination {
	return &pagination{
		Page:       -1,
		PerPage:    -1,
		TotalPages: -1,
		NextPage:   -1,
		PrevPage:   -1,
		Skip:       0,
		Total:      -1,
	}
}

func (pager *pagination) fillFromRequest(req *http.Request) error {
	var err error

	params := req.URL.Query()

	page := params.Get("page")
	perPage := params.Get("perPage")

	if page != "" {
		pager.Page, err = strconv.Atoi(page)
		if err == nil {
			if perPage == "" {
				pager.PerPage = defaultPerPage
			} else {
				pager.PerPage, err = strconv.Atoi(perPage)
			}

			if err == nil {
				if (pager.Page < 1) || (pager.PerPage < 1) {
					err = errors.New("Invalid pagination parameters")
				} else {
					pager.Skip = (pager.Page - 1) * pager.PerPage
				}
			}
		}
	}

	return err
}

func (pager *pagination) computePages() {
	if (pager.Page != -1) && (pager.PerPage != -1) && (pager.Total != -1) {
		pager.TotalPages = pager.Total / pager.PerPage
		if math.Mod(float64(pager.Total), float64(pager.PerPage)) != 0 {
			pager.TotalPages++
		}

		if pager.Page < pager.TotalPages {
			pager.NextPage = pager.Page + 1
		}

		if pager.Page > 1 {
			pager.PrevPage = pager.Page - 1
		}
	}
}

// MarshalJSON implements json.Marshaler
func (pager *pagination) MarshalJSON() ([]byte, error) {
	hash := map[string]interface{}{}

	if pager.Page != -1 {
		hash["page"] = pager.Page
		hash["perPage"] = pager.PerPage
		hash["skip"] = pager.Skip

		if pager.Total != -1 {
			hash["total"] = pager.Total
		}

		pager.computePages()

		if pager.TotalPages != -1 {
			hash["pages"] = pager.TotalPages
		}

		if pager.NextPage == -1 {
			hash["nextPage"] = nil
		} else {
			hash["nextPage"] = pager.NextPage
		}

		if pager.PrevPage == -1 {
			hash["prevPage"] = nil
		} else {
			hash["prevPage"] = pager.PrevPage
		}

	}

	return json.Marshal(hash)
}
