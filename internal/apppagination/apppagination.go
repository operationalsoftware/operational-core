package apppagination

type Pagination struct {
	Page     uint
	PageSize uint
}

func (p *Pagination) Offset() uint {
	if p.Page == 0 { // default if not set
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) Limit(defaultPageSize uint) uint {
	if p.PageSize == 0 {
		return defaultPageSize
	}
	return p.PageSize
}
