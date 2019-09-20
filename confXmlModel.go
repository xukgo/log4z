package log4z

import "encoding/xml"

type levelDefineXmlModel struct {
	LogPath    string `xml:"path"`
	LogSize    int    `xml:"size,attr"`
	MinLevel   string `xml:"min,attr"`
	MaxLevel   string `xml:"max,attr"`
	MaxDays    int    `xml:"maxDay,attr"`
	MaxBackup  int    `xml:"maxBackup,attr"`
	Encoder    string `xml:"encoder,attr"`
	IsConsole  bool   `xml:"console,attr"`
	LineRecord bool   `xml:"line,attr"`
}

type appenderXmlModel struct {
	Name         string                `xml:"name,attr"`
	LevelDefines []levelDefineXmlModel `xml:"levelDefine"`
}

type loggerXmlModel struct {
	Name         string `xml:"name,attr"`
	AppenderName string `xml:"appender"`
}

type confXmlRoot struct {
	XMLName   xml.Name           `xml:"log4z"`
	Loggers   []loggerXmlModel   `xml:"logger"`
	Appenders []appenderXmlModel `xml:"appender"`
}
