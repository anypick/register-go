package redisutil

import (
	_ "github.com/anypick/register-go/testx"
	"testing"
)

func TestBaseDao_Add(t *testing.T) {
	data := make(map[string]interface{})
	data["id"] = uint64(10001)
	data["username"] = "Jhone"
	data["password"] = "dqdsaadj kfjfwqe"
	data["height"] = 178.50
	data["phone"] = "18992222393"

	dao := BaseDao{Catalog:"User", Clazz:"user", IdName:"id"}
	dao.SelectFields = make([]string, 2)
	dao.SelectFields[0] = "username"
	dao.SelectFields[1] = "height"
	dao.Add(data, "")

}

func Inter(v interface{}) {

}

func TestBaseDao_Add2(t *testing.T) {

}
