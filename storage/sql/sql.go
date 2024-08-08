package sql

func Page(pageSize, page int) (offset, limit int) {
	offset = (page - 1) * pageSize
	limit = pageSize
	return offset, limit
}
