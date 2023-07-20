 ## 使用方式
  - 同时支持 websocket 和 http请求
 
 ### http请求
 Post http://127.0.0.1:10718/dubbo_send  
 参数：  
 ```javascript
 {
    "target_id": "abc",
    "debug": "true",
    "dubbo_protocol": "dubbo",
    "api_name": "com.dubbo.service.UserService",
    "function_name": "getById",
    "dubbo_param": [
        {
            "is_checked": 1,
            "field_type": "java.lang.Long",
            "key": "",
            "value": "2222222",
            "children":[]
        }
    ],
    "dubbo_config": {
        "registration_center_name": "nacos",
        "registration_center_address": "172.17.101.199:8848"
    }
}
 ```
 返回值  
 ```javascript
 {
	"code": 1000,
	"msg": "success",
	"data": {
		"target_id": "abc",
		"resp": {
			"class": "com.dubbo.entity.User",
			"id": 2222222,
			"name": "zhangsan"
		}
	}
}
 ```

 ### websocket 请求
 ws://127.0.0.1:10718/websocket