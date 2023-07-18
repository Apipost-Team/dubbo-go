package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Apipost-Team/dubbo-go/model"
)

func main() {
	var dubboStruct model.DubboDetail

	jsonStr := `{
		"target_id":"abc",
		"uuid":"e5f9ac80-d582-4c27-84f0-08e2c6e04cc3",
		"name":"name",
		"debug":"true",
		"dubbo_protocol":"dubbo",
		"api_name":"com.dubbo.service.UserService",
		"function_name":"getById",
		"dubbo_param":[
			{
				"is_checked":1,
				"param_type":"java.lang.Long",
				"var":"",
				"val":"2222222"
			}
		],
		"dubbo_config":{
			"registration_center_name":"nacos",
			"registration_center_address":"172.17.101.199:8848"
		}
	}`

	// jsonStr2 := `{
	//     "target_id": "",
	//     "uuid": "00000000-0000-0000-0000-000000000000",
	//     "name": "",
	//     "team_id": "",
	//     "debug": "",
	//     "dubbo_protocol": "dubbo",
	//     "api_name": "com.dubbo.service.UserService",
	//     "function_name": "getById",
	//     "dubbo_param": [
	//         {
	//             "is_checked": 1,
	//             "param_type": "java.lang.Long",
	//             "var": "",
	//             "val": "2222222"
	//         }
	//     ],
	//     "dubbo_assert": [],
	//     "dubbo_regex": [],
	//     "dubbo_config": {
	//         "registration_center_name": "nacos",
	//         "registration_center_address": "172.17.101.199:8848",
	//         "version": ""
	//     },
	//     "configuration": null,
	//     "global_variable": null,
	//     "dubbo_variable": null
	// }`

	json.Unmarshal([]byte(jsonStr), &dubboStruct)

	fmt.Println(dubboStruct)

	var debugMsg map[string]interface{}
	globalVariable := new(sync.Map)
	dubboStruct.Send(debugMsg, globalVariable)

}
