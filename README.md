# log4z
new golang package log config for zap

底层根据zap的log库，超高性能

可以轻松根据配置文件初始哈logger，并且可以快速生成不同的logger，每个logger里面有不同的level定义配置

最后提供了默认console的logger配置，这样在写testcase的时候不用担心logger没有初始化的问题
    
        
        
    var LoggerCommon *zap.Logger
    var LoggerWechat *zap.Logger
    
    func ExampleInit() {        
        configPath := "./conf/log4z.xml"
        loggerMap := log4z.InitLogger(configPath)
        LoggerCommon = getLoggerOrConsole(loggerMap, "Common")
        LoggerWechat = getLoggerOrConsole(loggerMap, "Wechat")
    
        fmt.Println("LoggerCommon", LoggerCommon)
        fmt.Println("LoggerWechat", LoggerWechat)
    }
    
    func getLoggerOrConsole(dict map[string]*zap.Logger, key string) *zap.Logger {
        logger, ok := dict[key]
        if ok {
            fmt.Printf("info: get logger %s success\r\n", key)
        } else {
            fmt.Printf("warnning: log4z get logger (%s) failed\r\n", key)
            fmt.Printf("warnning: now set logger %s to default console logger\r\n", key)
            logger = log4z.GetConsoleLogger()
        }
        return logger
    }

根据这样的写法，兼顾了生产代码和测试用例的需求