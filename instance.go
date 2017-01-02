package polarserver

import (
	"fmt"

	"github.com/Jordanzuo/goutil/logUtil"
	_ "github.com/polariseye/polarserver/common"
	"github.com/polariseye/polarserver/config"
	"github.com/polariseye/polarserver/dataBase"
	"github.com/polariseye/polarserver/moduleManage"
	"github.com/polariseye/polarserver/server"
	"github.com/polariseye/polarserver/server/webServer"
	"github.com/polariseye/polarserver/server/webServer/apiHandle"
)

var (
	// 服务管理对象
	serverManagerObj *server.ServerManagerStruct

	// web服务对象
	webServerObj *webServer.WebServerStruct

	// 配置对象
	configObj config.Configer

	// 处理对象
	handler *apiHandle.Handle4UrlStruct
)

// 初始化
func init() {
	serverManagerObj = server.NewServerManager()
}

// 初始化
// configFileName:配置文件名
// errMsg:错误信息
func Init(configFileName string) (errMsg error) {
	configObj, errMsg = config.NewConfig("json", configFileName)
	if errMsg != nil {
		return
	}

	// 初始化日志记录
	logUtil.SetLogPath(configObj.DefaultString("LogPath", "DefaultLogPath/"))

	// 配置初始化
	errMsg = initDataBaseFromConfig(configObj)
	if errMsg != nil {
		return errMsg
	}

	initWebServerFromConfig(configObj)

	return nil
}

// web服务对象
func WebServerObj() *webServer.WebServerStruct {
	return webServerObj
}

// 服务管理对象
func ServerManagerObj() *server.ServerManagerStruct {
	return serverManagerObj
}

// 配置对象
func ConfigObj() config.Configer {
	return configObj
}

// 初始化web服务
func initWebServerFromConfig(config config.Configer) {
	// web服务初始化
	port := config.DefaultInt("WebPort", 2017)
	webServerObj = webServer.NewWebServer(int32(port), "web 服务")

	// 初始化Api处理
	handler := apiHandle.NewHandle4Json(moduleManage.DefaulApiModuleManager)
	webServerObj.AddRouter("/Api", handler.RequestHandle)

	// 注册模块
	ServerManagerObj().Register(webServerObj)
}

// 从配置文件初始化数据库
func initDataBaseFromConfig(config config.Configer) error {
	tmp, errMsg := config.DIY("DbConnection")
	if errMsg != nil {
		return errMsg
	}
	connectionData, isDataOk := tmp.(map[string]interface{})
	if isDataOk == false {
		return fmt.Errorf("数据库配置不合法,节点：DbConnection")
	}

	for key, connectionInfo := range connectionData {
		connectionItem, isOk := connectionInfo.(map[string]interface{})
		if isOk == false {
			return fmt.Errorf("数据库配置不合法,节点：DbConnection.%v", key)
		}

		var driver, connectionString interface{}

		driver, isOk = connectionItem["Driver"]
		if isOk == false {
			return fmt.Errorf("数据库配置不合法,不存在节点：DbConnection.%v.Driver", key)
		}

		connectionString, isOk = connectionItem["ConnectionString"]
		if isOk == false {
			return fmt.Errorf("数据库配置不合法,不存在节点：DbConnection.%v.ConnectionString", key)
		}

		errMsg = dataBase.AddConnection(key, driver.(string), connectionString.(string))
		if errMsg != nil {
			return errMsg
		}
	}

	return nil
}
