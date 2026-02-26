package provider

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"watchAlert/pkg/tools"
)

// LogsAdvanced 高级模式的日志结构，支持从查询结果中提取字段值进行告警判断
// 类似 Metrics 的处理方式，每条查询结果可生成独立的告警事件
type LogsAdvanced struct {
	ProviderName string
	Message      []map[string]interface{}
}

// GenerateFingerprint 基于 ruleId 生成指纹（用于日志数量类型的告警）
func (l LogsAdvanced) GenerateFingerprint(ruleId string) string {
	h := md5.New()
	streamString := tools.JsonMarshalToString(map[string]string{
		"ruleId": ruleId,
	})
	h.Write([]byte(streamString))
	fingerprint := hex.EncodeToString(h.Sum(nil))
	return fingerprint
}

// GenerateFingerprintByLabels 基于标签生成指纹（用于从查询结果提取数值的告警，类似 Metrics）
func (l LogsAdvanced) GenerateFingerprintByLabels(labels map[string]interface{}) string {
	if len(labels) == 0 {
		return fmt.Sprintf("%d", tools.HashNew())
	}

	var result uint64
	for labelName, labelValue := range labels {
		sum := tools.HashNew()
		sum = tools.HashAdd(sum, labelName)
		sum = tools.HashAdd(sum, fmt.Sprintf("%v", labelValue))
		result ^= sum
	}

	return fmt.Sprintf("%d", result)
}

// GetValue 从第一条日志消息中提取指定字段的数值
// 用于 ClickHouse 等数据源从查询结果中提取聚合值进行告警判断
func (l LogsAdvanced) GetValue(fieldName string) (float64, bool) {
	if len(l.Message) == 0 {
		return 0, false
	}

	// 从第一条消息中获取指定字段
	value, exists := l.Message[0][fieldName]
	if !exists {
		return 0, false
	}

	// 尝试转换为 float64
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case string:
		// 尝试将字符串转换为数值
		if floatVal, err := fmt.Sscanf(v, "%f", new(float64)); err == nil && floatVal == 1 {
			var result float64
			fmt.Sscanf(v, "%f", &result)
			return result, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// GetMetric 获取第一条消息的所有字段作为标签（类似 Metrics.GetMetric）
func (l LogsAdvanced) GetMetric() map[string]interface{} {
	if len(l.Message) == 0 {
		return make(map[string]interface{})
	}
	return l.Message[0]
}

// GetAnnotations 获取注释信息，处理超长字符串
func (l LogsAdvanced) GetAnnotations() map[string]interface{} {
	msg := make(map[string]interface{})
	if len(l.Message) == 0 {
		return msg
	}

	for k, v := range l.Message[0] {
		if v == nil {
			continue
		}

		switch value := v.(type) {
		case string:
			// 如果是字符串类型，处理长度限制
			content := value
			length := len(content)
			if length > 1000 {
				msg[k] = fmt.Sprintf("%s... 内容过长省略其中 ...%s", content[:500], content[length-500:])
			} else {
				msg[k] = content
			}
		case map[string]interface{}:
			// 如果是嵌套的 map，递归调用处理
			msg[k] = processNestedMapAdvanced(value)
		default:
			// 对于其他类型，直接保留原值
			msg[k] = value
		}
	}
	return msg
}

// processNestedMapAdvanced 辅助函数：递归处理嵌套的 map[string]interface{}
func processNestedMapAdvanced(nestedMap map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range nestedMap {
		if v == nil {
			continue
		}

		switch value := v.(type) {
		case string:
			// 如果是字符串类型，处理长度限制
			content := value
			length := len(content)
			if length > 1000 {
				result[k] = fmt.Sprintf("%s... 内容过长省略其中 ...%s", content[:500], content[length-500:])
			} else {
				result[k] = content
			}
		case map[string]interface{}:
			// 如果是嵌套的 map，继续递归处理
			result[k] = processNestedMapAdvanced(value)
		default:
			// 对于其他类型，直接保留原值
			result[k] = value
		}
	}
	return result
}

