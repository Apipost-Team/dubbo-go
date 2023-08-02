package main

import (
	"encoding/json"
	"fmt"

	"github.com/Apipost-Team/dubbo-go/request"
)

func main() {
	var dubboStruct request.DubboRequest

	// jsonStr := `{
	// 	"target_id":"abc",
	// 	"uuid":"e5f9ac80-d582-4c27-84f0-08e2c6e04cc3",
	// 	"name":"name",
	// 	"debug":"true",
	// 	"dubbo_protocol":"dubbo",
	// 	"api_name":"com.dubbo.service.UserService",
	// 	"function_name":"getById",
	// 	"dubbo_param":[
	// 		{
	// 			"is_checked":1,
	// 			"field_type":"java.lang.Long",
	// 			"key":"",
	// 			"value":"2222222"
	// 		}
	// 	],
	// 	"dubbo_config":{
	// 		"registration_center_name":"nacos",
	// 		"registration_center_address":"172.17.101.199:8848"
	// 	}
	// }`

	jsonStr := `{
		"target_id":"abc",
		"uuid":"e5f9ac80-d582-4c27-84f0-08e2c6e04cc3",
		"name":"name",
		"debug":"true",
		"dubbo_protocol":"dubbo",
		"api_name":"com.dubbo.service.UserService",
		"function_name":"updateByUser",
		"dubbo_param":[
			{
				"is_checked":1,
				"field_type":"com.dubbo.entity.User",
				"key":"",
				"value":"",
				"children":[
					{
						"is_checked":1,
						"field_type":"java.lang.Long",
						"key":"id",
						"value":"333333",
						"children":[]
					},
					{
						"is_checked":1,
						"field_type":"java.lang.String",
						"key":"name",
						"value":"aaa",
						"children":[]
					}
				]
			}
		],
		"dubbo_config":{
			"registration_center_name":"nacos",
			"registration_center_address":"172.17.101.199:8848"
		}
	}`

	// jsonStr := `{
	// 	"target_id":"abc",
	// 	"uuid":"e5f9ac80-d582-4c27-84f0-08e2c6e04cc3",
	// 	"name":"name",
	// 	"debug":"true",
	// 	"dubbo_protocol":"dubbo",
	// 	"api_name":"com.dubbo.service.UserService",
	// 	"function_name":"updateByUser",
	// 	"dubbo_param":[
	// 		{
	// 			"is_checked":1,
	// 			"field_type":"com.dubbo.entity.User",
	// 			"key":"",
	// 			"value":"{ \"id\": 22222, \"name\": \"张三\"}",
	// 			"children":[
	// 			]
	// 		}
	// 	],
	// 	"dubbo_config":{
	// 		"registration_center_name":"nacos",
	// 		"registration_center_address":"172.17.101.199:8848"
	// 	}
	// }`

	json.Unmarshal([]byte(jsonStr), &dubboStruct)

	fmt.Println(dubboStruct)

	resp, err := dubboStruct.Send()
	if err != nil {
		fmt.Println(resp, err)
	} else {
		fmt.Println(resp)
	}

}
