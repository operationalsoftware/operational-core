package apppagination

type Pagination struct {
	Page     uint
	PageSize uint
}

func (p *Pagination) Offset() uint {
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) Limit() uint {
	return p.PageSize
}
