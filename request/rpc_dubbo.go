package request

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	dubboConfig "dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/service/local"
	"github.com/Apipost-Team/dubbo-go/constant"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/common"
	uuid "github.com/satori/go.uuid"
)

var RpcServerMap = new(sync.Map)

// 转换json
func convertMap(m map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		switch v := v.(type) {
		case map[interface{}]interface{}:
			result[fmt.Sprint(k)] = convertMap(v)
		default:
			result[fmt.Sprint(k)] = v
		}
	}
	return result
}

type DubboConfig struct {
	RegistrationCenterName    string `json:"registration_center_name"`
	RegistrationCenterAddress string `json:"registration_center_address"`
	Version                   string `json:"version"`
}

type DubboParam struct {
	IsChecked int32  `json:"is_checked"`
	ParamType string `json:"param_type"`
	Var       string `json:"var"`
	Val       string `json:"val"`
}

type DubboRequest struct {
	TargetId string `json:"target_id"`
	Name     string `json:"name"`
	Debug    string `json:"debug"`

	DubboProtocol string `json:"dubbo_protocol"`
	ApiName       string `json:"api_name"`
	FunctionName  string `json:"function_name"`

	DubboParam  []DubboParam `json:"dubbo_param"`
	DubboConfig DubboConfig  `json:"dubbo_config"`
}

func (d *DubboRequest) Send() (any, error) {
	parameterTypes, parameterValues := []string{}, []hessian.Object{}
	var err error
	var rpcServer common.RPCService
	d.DubboConfig.RegistrationCenterName = strings.TrimSpace(d.DubboConfig.RegistrationCenterName)
	d.DubboConfig.RegistrationCenterAddress = strings.TrimSpace(d.DubboConfig.RegistrationCenterAddress)
	d.ApiName = strings.TrimSpace(d.ApiName)
	d.DubboConfig.Version = strings.TrimSpace(d.DubboConfig.Version)
	soleKey := fmt.Sprintf("%s://%s/%s", d.DubboProtocol, d.DubboConfig.RegistrationCenterAddress, d.ApiName)

	if s, ok := RpcServerMap.Load(soleKey); ok {
		rpcServer, ok = s.(common.RPCService)
		if !ok {
			return "rpcServer", fmt.Errorf("rpcServer is not common.RPCService")
		}
	} else {
		rpcServer, err = d.init(soleKey)
	}

	for _, parame := range d.DubboParam {
		if parame.IsChecked != constant.Open {
			break
		}
		var val interface{}
		switch parame.ParamType {
		case constant.JavaInteger:
			val, err = strconv.Atoi(parame.Val)
			if err != nil {
				val = parame
				continue
			}
		case constant.JavaString:
			val = parame.Val
		case constant.JavaBoolean:
			switch parame.Val {
			case "true":
				val = true
			case "false":
				val = false
			default:
				val = parame.Val
			}
		case constant.JavaByte:

		case constant.JavaCharacter:
		case constant.JavaDouble:
			val, err = strconv.ParseFloat(parame.Val, 64)
			if err != nil {
				val = parame.Val
				continue
			}
		case constant.JavaFloat:
			val, err = strconv.ParseFloat(parame.Val, 64)
			if err != nil {
				val = parame.Val
				continue
			}
			val = float32(val.(float64))
		case constant.JavaLong:
			val, err = strconv.ParseInt(parame.Val, 10, 64)
			if err != nil {
				val = parame.Val
				continue
			}
		case constant.JavaMap:
		case constant.JavaList:
		default:
			val = parame.Val
		}
		parameterTypes = append(parameterTypes, parame.ParamType)
		parameterValues = append(parameterValues, val)

	}

	if err != nil {
		return "param", err
	}

	fmt.Println(parameterTypes, parameterValues)

	var resp interface{}
	//var response []byte

	resp, err = rpcServer.(*generic.GenericService).Invoke(
		context.TODO(),
		d.FunctionName,
		parameterTypes,
		parameterValues, // 实参
	)

	if err != nil {
		return "end", err
	}

	if resp == nil {
		return "end2", fmt.Errorf("resp is nil")
	}

	resp2, ok := resp.(map[interface{}]interface{})
	if !ok {
		return resp, nil
	}

	return convertMap(resp2), nil
}

func (d *DubboRequest) init(soleKey string) (rpcServer common.RPCService, err error) {
	registryConfig := &dubboConfig.RegistryConfig{
		Protocol: d.DubboConfig.RegistrationCenterName,
		Address:  d.DubboConfig.RegistrationCenterAddress,
	}

	var zk string
	if d.DubboConfig.RegistrationCenterName == "zookeeper" {
		zk = "zk"
	} else {
		zk = d.DubboConfig.RegistrationCenterName
	}

	refConf := &dubboConfig.ReferenceConfig{
		InterfaceName:  d.ApiName, // 服务接口名，如：org.apache.dubbo.sample.UserProvider
		Cluster:        "failover",
		RegistryIDs:    []string{zk},          // 注册中心
		Protocol:       d.DubboProtocol,       // dubbo  或 tri（triple）  使用的协议
		Generic:        "true",                // true: 使用泛化调用；false: 不适用泛化调用
		Version:        d.DubboConfig.Version, // 版本号
		RequestTimeout: "3",
		Serialization:  "hessian2",
	}
	fmt.Println("refConf: ", refConf)

	// 构造 Root 配置，引入注册中心模块
	rootConfig := dubboConfig.NewRootConfigBuilder().AddRegistry(zk, registryConfig).Build()
	if err = dubboConfig.Load(dubboConfig.WithRootConfig(rootConfig)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rootConfig)

	//if err = rootConfig.Init(); err != nil {
	//	return
	//}

	// Reference 配置初始化，因为需要使用注册中心进行服务发现，需要传入经过配置的 rootConfig
	if err = refConf.Init(rootConfig); err != nil {
		return
	}

	myuuid, _ := uuid.NewV4()
	fmt.Println("uuuuuid      ", myuuid.String())
	refConf.GenericLoad(myuuid.String())
	rpcServer2 := refConf.GetRPCService()
	RpcServerMap.Store(soleKey, rpcServer2) //存储
	rpcServer, ok := rpcServer2.(common.RPCService)
	if !ok {
		err = fmt.Errorf("rpcServer is not common.RPCService2")
	}

	fmt.Println("rpcServer: ", rpcServer)
	return
}
