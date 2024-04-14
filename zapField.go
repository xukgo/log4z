package log4z

import (
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DecimalField 创建一个zap.Field来记录decimal.Decimal类型的数据。
func DecimalField(key string, val decimal.Decimal) zap.Field {
	return zap.Field{
		Key:    key,
		Type:   zapcore.StringType,
		String: val.String(),
	}
}

// DecimalsField 创建一个zap.Field来记录decimal.Decimal数组类型的数据。
func DecimalsField(key string, vals []decimal.Decimal) zap.Field {
	return zap.Array(key, zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, v := range vals {
			// 使用enc来添加每个元素，这里依旧是作为字符串
			enc.AppendString(v.String())
		}
		return nil
	}))
}
