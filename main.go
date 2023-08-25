package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webapp/controller"
	"webapp/dao/mysql"
	"webapp/dao/redis"
	"webapp/logger"
	"webapp/pkg/snowflake"
	"webapp/routers"
	"webapp/settings"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

// @title 这里写标题
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 这里写接口服务的host
// @BasePath 这里写base path
func main() {
	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("set up setting fail, %s \n", err)
	}
	//2.初始化logger,记录日志
	if err := logger.Init(settings.Conf.LogConf, settings.Conf.Mode); err != nil {
		fmt.Printf("set up setting fail, %s \n", err)
		return
	}
	zap.L().Debug("logger init success")
	defer zap.L().Sync()
	//3.初始化mysql
	if err := mysql.InitDB(settings.Conf.MysqlConf); err != nil {
		fmt.Printf("mysql init fail, %s \n", err)
		return
	}
	zap.L().Debug("mysql init success")
	defer mysql.Close()
	//4.初始化redis
	if err := redis.Init(settings.Conf.RedisConf); err != nil {
		fmt.Printf("redis init fail, %s \n", err)
		return
	}
	zap.L().Debug("redis init success")
	defer redis.Close()

	if err := snowflake.Init(1); err != nil {

	}
	//初始化gin框架校验器的翻译器
	err := controller.InitTrans("zh")
	if err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	//5.注册路由
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	routers.Router(r)
	//6.启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		fmt.Println("cccccc")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
