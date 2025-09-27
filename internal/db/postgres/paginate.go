package postgres

type Paginate struct {
	Current int         `json:"current"`
	Size    int         `json:"size"`
	Total   int64       `json:"total"`
	List    interface{} `json:"list"`
}

func NewPaginate(current, size int) *Paginate {
	if current <= 0 {
		current = 1
	}
	if size <= 0 {
		size = 10
	}
	return &Paginate{
		Current: current,
		Size:    size,
	}
}

func (paginate *Paginate) Limit() int {
	return paginate.Size
}

func (paginate *Paginate) Offset() int {
	return (paginate.Current - 1) * paginate.Size
}
