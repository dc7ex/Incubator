package mysql

import (
	"log"
	"os"
	"time"

	"github.com/arthurkiller/rollingwriter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func NewMysql(connection string, host string, port string, database string, username string, password string, charset string, show bool, logPath string) *xorm.Engine {
	// 异常处理
	defer func() {
		if err := recover(); err != nil {
			log.Println("Mysql Err:", err)
			os.Exit(1)
		}
	}()

	var err error
	var conn *xorm.Engine

	// 单例
	conn, err = xorm.NewEngine(connection, username+":"+password+"@tcp("+host+":"+port+")/"+database+"?charset="+charset) // 链接数据库
	if err != nil {
		panic(err)
	}
	conn.SetMaxIdleConns(30)                                    // 设置连接池的空闲数大小
	conn.SetMaxOpenConns(30)                                    // 设置最大打开连接数
	conn.SetConnMaxLifetime(time.Duration(28800) * time.Second) // 设置连接可以重用的最大时间 空闲时间 28800为Mysql默认时间

	if logPath == "" {
		logPath = "./logs"
	}

	config := rollingwriter.Config{
		LogPath:       logPath,          //日志路径
		TimeTagFormat: "20060102150405", //时间格式串
		FileName:      "mysql_exec",     //日志文件名
		MaxRemain:     7,                //配置日志最大存留数
		// 目前有2中滚动策略: 按照时间滚动按照大小滚动
		// - 时间滚动: 配置策略如同 crontable, 例如,每天0:0切分, 则配置 0 0 0 * * *
		// - 大小滚动: 配置单个日志文件(未压缩)的滚动大小门限, 如1G, 500M
		RollingPolicy:      rollingwriter.TimeRolling, // 配置滚动策略 norolling timerolling volumerolling
		RollingTimePattern: "0 0 0 * * *",             // 配置时间滚动策略
		RollingVolumeSize:  "200MB",                   // 配置截断文件下限大小
		// writer 支持4种不同的 mode:
		// 1. none 2. lock
		// 3. async 4. buffer
		// - 无保护的 writer: 不提供并发安全保障
		// - lock 保护的 writer: 提供由 mutex 保护的并发安全保障
		// - 异步 writer: 异步 write, 并发安全. 异步开启后忽略 Lock 选项
		WriterMode:             "none",
		BufferWriterThershould: 256,
		// Compress will compress log file with gzip
		Compress: true,
	}

	writer, err := rollingwriter.NewWriterFromConfig(&config)
	if err != nil {
		panic(err)
	}

	nsl := xorm.NewSimpleLogger(writer)
	//nsl := xorm.NewSimpleLogger(logger.GetLogger().Writer())
	conn.SetLogger(nsl)
	conn.ShowSQL(show)

	err = conn.Ping() // 检查数据库链接
	if err != nil {
		panic(err)
	}

	return conn
}
