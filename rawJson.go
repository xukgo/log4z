package log4z

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type rawJSON []byte

// MarshalLogObject 实现 zapcore.ObjectMarshaler 接口
func (r rawJSON) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	// 使用 enc.AddReflected 来避免字符串被转义
	// 注意：这种方式要求字符串是合法的JSON
	return enc.AddReflected("raw", json.RawMessage(r))
}

// RawJSONField 创建一个zap.Field，用来存放不需要转义的JSON字符串
func RawJSONField(key string, val []byte) zap.Field {
	return zap.Field{Key: key, Type: zapcore.ObjectMarshalerType, Interface: rawJSON(val)}
}
