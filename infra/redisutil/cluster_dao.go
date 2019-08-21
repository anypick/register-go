package redisutil


type ClusterDao struct {
	// 所属分类，
	Catalog string
	// 表对应的struct
	Clazz string
	// redis cluster hashTag
	HashTag string
	// 数据的Id, 有些表的Id字段命名比较另类
	IdDesc FieldDescriptor
	// 定义需要被查询的字段
	SelectFields []FieldDescriptor
}

