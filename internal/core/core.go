package core

const (
	TitleOrder       = "title"
	RatingOrder      = "rating"
	ReleaseDateOrder = "release_date"

	UserRole  = "user"
	AdminRole = "admin"
)

var AllowedSorts = map[string]struct{}{
	TitleOrder:       {},
	RatingOrder:      {},
	ReleaseDateOrder: {},
}
