package server

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
)

const DEFAULT_PER_PAGE = 20

type Pagination struct {
	Page       int
	PerPage    int
	TotalPages int
	NextPage   int
	PrevPage   int

	Skip  int
	Total int
}

// create pagination params from request
func NewPagination() *Pagination {
	return &Pagination{
		Page:       -1,
		PerPage:    -1,
		TotalPages: -1,
		NextPage:   -1,
		PrevPage:   -1,
		Skip:       0,
		Total:      -1,
	}
}

func (this *Pagination) fillFromRequest(req *http.Request) error {
	var err error

	params := req.URL.Query()

	page := params.Get("page")
	perPage := params.Get("per_page")

	if page != "" {
		this.Page, err = strconv.Atoi(page)
		if err == nil {
			if perPage == "" {
				this.PerPage = DEFAULT_PER_PAGE
			} else {
				this.PerPage, err = strconv.Atoi(perPage)
			}

			if err == nil {
				if (this.Page < 1) || (this.PerPage < 1) {
					err = errors.New("Invalid pagination parameters")
				} else {
					this.Skip = (this.Page - 1) * this.PerPage
				}
			}
		}
	}

	return err
}

func (this *Pagination) computePages() {
	if (this.Page != -1) && (this.PerPage != -1) && (this.Total != -1) {
		this.TotalPages = this.Total / this.PerPage
		if math.Mod(float64(this.Total), float64(this.PerPage)) != 0 {
			this.TotalPages += 1
		}

		if this.Page < this.TotalPages {
			this.NextPage = this.Page + 1
		}

		if this.Page > 1 {
			this.PrevPage = this.Page - 1
		}
	}
}

// Implements json.MarshalJSON
func (this *Pagination) MarshalJSON() ([]byte, error) {
	hash := map[string]interface{}{}

	if this.Page != -1 {
		hash["page"] = this.Page
		hash["perPage"] = this.PerPage
		hash["skip"] = this.Skip

		if this.Total != -1 {
			hash["total"] = this.Total
		}

		this.computePages()

		if this.TotalPages != -1 {
			hash["pages"] = this.TotalPages
		}

		if this.NextPage != -1 {
			hash["next"] = this.NextPage
		}

		if this.PrevPage != -1 {
			hash["prev"] = this.PrevPage
		}

	}

	return json.Marshal(hash)
}
