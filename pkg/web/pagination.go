package web

type (
	PaginationBody struct {
		TotalCount  int `json:"total_count"`
		PageCount   int `json:"page_count"`
		CurrentPage int `json:"current_page"`
		PerPage     int `json:"per_page"`
	}

	PaginationQuery struct {
		Page  int `query:"page"`
		Limit int `query:"limit" validate:"max=500"`
	}
)

func NewPaginationQuery(page, limit int) PaginationQuery {
	return PaginationQuery{
		Page:  page,
		Limit: limit,
	}
}

func (p PaginationQuery) GetPage() int {
	if p.Page <= 0 {
		return 1
	}

	return p.Page
}

func (p PaginationQuery) GetLimit() int {
	if p.Limit <= 0 {
		return 20
	}

	return p.Limit
}

func (p PaginationQuery) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p PaginationQuery) GetPageCount(total int) int {
	pageCount := total / p.GetLimit()

	if pageCount%p.GetLimit() != 0 || total < p.GetLimit() {
		pageCount++
	}

	return pageCount
}

func (p PaginationQuery) PaginationBody(total int) PaginationBody {
	return PaginationBody{
		TotalCount:  total,
		PageCount:   p.GetPageCount(total),
		CurrentPage: p.GetPage(),
		PerPage:     p.GetLimit(),
	}
}
