package redisutil

import (
	"github.com/anypick/register-go/infra/utils/common"
	_ "github.com/anypick/register-go/testx"
	"testing"
)

func TestBaseDao_Add(t *testing.T) {
	data := make(map[string]interface{})
	data["id"] = uint64(10002)
	data["username"] = "Jhone"
	data["password"] = "dqdsaadj kfjfwqe"
	data["height"] = 178.50
	data["phone"] = "18992222393"

	dao := BaseDao{Catalog:"User", Clazz:"user", IdName:"id"}
	dao.SelectFields = make([]string, 2)
	dao.SelectFields[0] = "username"
	dao.SelectFields[1] = "height"
	dao.Add(data, common.NilString)

}

func TestBaseDao_Get(t *testing.T) {
	dao := BaseDao{Catalog:"User", Clazz:"user", IdName:"id"}
	dao.SelectFields = make([]string, 2)
	dao.SelectFields[0] = "username"
	dao.SelectFields[1] = "height"

	dao.Get(uint64(10001), common.NilString)
}


func TestBaseDao_GetByField(t *testing.T) {
	dao := BaseDao{Catalog:"User", Clazz:"user", IdName:"id"}
	dao.SelectFields = make([]string, 2)
	dao.SelectFields[0] = "username"
	dao.SelectFields[1] = "height"

	dao.GetByField("Jhone", "username", "")

}

func TestBaseDao_AddHash(t *testing.T) {
	data := make(map[string]interface{})
	data["id"] = uint64(10003)
	data["username"] = "HashJhone"
	data["password"] = "dsadasweedsdfasf"
	data["height"] = 181.50
	data["phone"] = "18857638866"

	dao := BaseDao{Catalog:"User", Clazz:"user", IdName:"id"}
	dao.SelectFields = make([]string, 1)
	dao.SelectFields[0] = "username"

	dao.AddHash(data, common.NilString)


}