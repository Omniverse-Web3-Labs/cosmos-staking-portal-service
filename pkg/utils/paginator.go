package utils

type PageRequest struct {
	Current int    `json:"page"`
	Size    int    `json:"size"`
	Order   string `json:"order"`
	Sort    string `json:"sort"`
}

type Paginator struct {
	Current int   `json:"current"`
	Size    int   `json:"size"`
	List    any   `json:"list"`
	Total   int64 `json:"total"`
	HasMore bool  `json:"hasMore"`
}

func NewPaginator(current int, size int) *Paginator {
	if current == 0 {
		current = 1
	}
	if size == 0 {
		size = 10
	}
	return &Paginator{Current: current, Size: size}
}

func (p *Paginator) SetList(v any) {
	p.List = v
}

func (p *Paginator) SetTotal(total int64) {
	p.Total = total
}

func (p *Paginator) SetHasMore(hasMore bool) {
	p.HasMore = hasMore
}

func (p *Paginator) Limit() int {
	return p.Size
}

func (p *Paginator) Offset() int {
	skip := (p.Current - 1) * p.Size
	return skip
}

func (p *Paginator) ToMap() map[string]any {
	return map[string]any{
		"current": p.Current,
		"size":    p.Size,
		"list":    p.List,
		"total":   p.Total,
		"hasMore": p.HasMore,
	}
}
