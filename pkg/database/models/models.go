package models

type InOut struct {
	URL string `json:"url"`
}

type URLShort struct {
	ID          int    `json:"id"`
	URLAddress  string `json:"url_address"`
	VisitCounts int    `json:"visit_counts"`
}
