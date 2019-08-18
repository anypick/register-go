package redisutil

import "errors"

type DataStructure uint

const (
	Strings DataStructure = iota
	Hashes
	// 列表，可重复
	Lists
	// 集合无序，不重复, 暂时不用
	Sets
	// 有序集合
	ZSets
)

type OperationType uint

const (
	// 更新数据
	Add OperationType = iota
	// 刷新数据
	Ref
	// 删除数据
	Del
)

const (
	Zero float64 = 0
)

/**
定义Redis操作属性
*/
type RedisOperation struct {
	// 操作类型：ref, add, del
	Operation OperationType
	// 操作的数据类型
	DataType DataStructure
	// key字段
	Key string
	// Hashes 类型的field
	RawKey string
	// 插入的值
	Value string
	// 分数，针对ZSets类型
	Score float64
}

const (
	DefaultFieldType = TypeEq
)

type FieldType uint

const (
	// 该字段和id 一一对应, redis中strings类型
	TypeEq FieldType = iota
	// 该字段和id是多对一关系
	TypeMatch
	// 需要排序
	TypeRange
)

// 字段描述，定义查询的字段
type FieldDescriptor struct {
	// 字段名称
	FieldName string
	// 字段类型
	FieldType FieldType
}

const (
	DefaultPage     = 1
	DefaultPageSize = 16
)

// 分页信息
type PageInfo struct {
	Page      int
	PageSize  int
	TotalPage int
	Total     int
	Value     []map[string]interface{}
}

func NewPageInfo(page, pageSize, total int, value []map[string]interface{}) (pageInfo PageInfo, err error) {
	if total == 0 {
		err = errors.New("total cannot is 0")
		return
	}
	if page == 0 {
		page = DefaultPage
	}

	if pageSize == 0 {
		pageSize = DefaultPageSize
	}
	totalPage := total/pageSize + 1
	pageInfo = PageInfo{Page: page, PageSize: pageSize, TotalPage: totalPage, Total: totalPage, Value: value}
	return
}
