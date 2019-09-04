# log4z
new golang package log config for zap

底层根据zap的log库，超高性能

可以轻松根据配置文件初始哈logger，并且可以快速生成不同的logger，每个logger里面有不同的level定义配置

最后提供了默认console的logger配置，这样在写testcase的时候不用担心logger没有初始化的问题
    
    
        var logCommon *zap.Logger 
        var logWechat *zap.Logger 
        
        func Init() {
        configPath := "./conf/log4z.xml"
        err := log4z.InitConfig(configPath)
        if err != nil {
		    fmt.Printf("warnning: log4z.InitConfig(configPath) configPath=%s; return err=%s\r\n", configPath, err.Error())
            fmt.Println("warnning: now set all logger to default console logger")
            logCommon = log4z.GetConsoleLogger()
            logWechat = log4z.GetConsoleLogger()
        } else {
            logCommon, err = log4z.InitLogger("Common")
            if err != nil {
                fmt.Printf("warnning: log4z.InitLogger(Common) return err=%s\r\n", err.Error())
                fmt.Println("warnning: now set logger Common to default console logger")
                logCommon = log4z.GetConsoleLogger()
            } else {
                fmt.Println("init logCommon success")
            }
            logWechat, err = log4z.InitLogger("Wechat")
            if err != nil {
                fmt.Printf("warnning: log4z.InitLogger(Wechat) return err=%s\r\n", err.Error())
                fmt.Println("warnning: now set logger Wechat to default console logger")
                logWechat = log4z.GetConsoleLogger()
            } else {
                fmt.Println("init logWechat success")
            }
        }
        log4z.UnintConfig()   
    }

根据这样的写法，兼顾了生产代码和测试用例的需求