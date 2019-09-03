package log4z

import "encoding/xml"

type LevelDefineXmlModel struct {
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

type AppenderXmlModel struct {
	Name         string                `xml:"name,attr"`
	LevelDefines []LevelDefineXmlModel `xml:"levelDefine"`
}

type LoggerXmlModel struct {
	Name         string `xml:"name,attr"`
	AppenderName string `xml:"appender"`
}

type ConfXmlRoot struct {
	XMLName   xml.Name           `xml:"log4z"`
	Loggers   []LoggerXmlModel   `xml:"logger"`
	Appenders []AppenderXmlModel `xml:"appender"`
}
