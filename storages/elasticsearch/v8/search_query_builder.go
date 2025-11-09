package v8

type QueryBuilder interface {
	Build() QueryBody
}
type queryBuilder struct{}

func newQueryBuilder() QueryBuilder {
	return &queryBuilder{}
}

func (q *queryBuilder) Build() QueryBody {
	return nil
}
