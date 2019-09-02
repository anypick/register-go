package esutil

type EsOperator struct {
	Index      string // 索引Index
	IndexAlias string // 索引别名
}

func (e *EsOperator) Search() {
}
