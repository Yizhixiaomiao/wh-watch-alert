package provider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"watchAlert/internal/models"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/zeromicro/go-zero/core/logc"
)

// ClickHouseAdvancedProvider 高级模式的 ClickHouse 数据提供者
// 保持原始数据类型，支持从查询结果中提取字段值进行告警判断
type ClickHouseAdvancedProvider struct {
	client         *sql.DB
	ExternalLabels map[string]interface{}
}

// LogsAdvancedFactoryProvider 高级模式的日志工厂接口
type LogsAdvancedFactoryProvider interface {
	QueryAdvanced(options LogQueryOptions) (LogsAdvanced, int, error)
	Check() (bool, error)
	GetExternalLabels() map[string]interface{}
}

// NewClickHouseAdvancedClient 创建高级模式的 ClickHouse 客户端
func NewClickHouseAdvancedClient(ctx context.Context, ds models.AlertDataSource) (LogsAdvancedFactoryProvider, error) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{ds.ClickHouseConfig.Addr},
		Auth: clickhouse.Auth{
			Username: ds.Auth.User,
			Password: ds.Auth.Pass,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * time.Duration(ds.ClickHouseConfig.Timeout),
	})
	if conn == nil {
		return nil, errors.New("clickhouse connection failed")
	}

	return ClickHouseAdvancedProvider{
		client:         conn,
		ExternalLabels: ds.Labels,
	}, nil
}

// QueryAdvanced 高级模式查询，保持原始数据类型
func (c ClickHouseAdvancedProvider) QueryAdvanced(options LogQueryOptions) (LogsAdvanced, int, error) {
	rows, err := c.client.Query(options.ClickHouse.Query)
	if err != nil {
		return LogsAdvanced{}, 0, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return LogsAdvanced{}, 0, err
	}

	// 获取列类型信息
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return LogsAdvanced{}, 0, err
	}

	var (
		// 存储所有日志数据
		message []map[string]interface{}
		// 准备 values 数组，用于接收每行数据
		values = make([]interface{}, len(columns))
	)

	for rows.Next() {
		// 每次循环都重新绑定指针，根据列类型创建对应的接收变量
		for i := range columns {
			values[i] = new(interface{})
		}

		// 扫描数据到 values
		if err := rows.Scan(values...); err != nil {
			logc.Error(context.Background(), "clickhouse advanced scan error:", err)
			return LogsAdvanced{}, 0, err
		}

		// 构造 map，保持原始数据类型
		entry := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			// 解包 interface{}
			if val != nil {
				rawValue := *val.(*interface{})
				if rawValue != nil {
					// 根据实际类型进行适当的转换
					entry[col] = convertClickHouseValueAdvanced(rawValue, columnTypes[i].DatabaseTypeName())
				} else {
					entry[col] = nil
				}
			} else {
				entry[col] = nil
			}
		}

		message = append(message, entry)
	}

	if err := rows.Err(); err != nil {
		return LogsAdvanced{}, 0, err
	}

	return LogsAdvanced{
		ProviderName: ClickHouseDsProviderName,
		Message:      message,
	}, len(message), nil
}

// convertClickHouseValueAdvanced 根据 ClickHouse 的数据类型转换值
// 保持数值类型不变，方便后续告警判断使用
func convertClickHouseValueAdvanced(value interface{}, typeName string) interface{} {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// 字节切片转换为字符串
		return string(v)
	case int8, int16, int32, int64, int:
		// 整数类型保持原样
		return v
	case uint8, uint16, uint32, uint64, uint:
		// 无符号整数类型保持原样
		return v
	case float32, float64:
		// 浮点数类型保持原样
		return v
	case string:
		// 字符串类型保持原样
		return v
	case bool:
		// 布尔类型保持原样
		return v
	case time.Time:
		// 时间类型转换为字符串（ISO 8601 格式）
		return v.Format(time.RFC3339)
	default:
		// 其他类型尝试转换为字符串
		return fmt.Sprintf("%v", v)
	}
}

// Check 检查 ClickHouse 连接是否健康
func (c ClickHouseAdvancedProvider) Check() (bool, error) {
	err := c.client.Ping()
	if err != nil {
		return false, errors.New("check clickhouse advanced datasource is unhealthy")
	}

	return true, nil
}

// GetExternalLabels 获取外部标签
func (c ClickHouseAdvancedProvider) GetExternalLabels() map[string]interface{} {
	return c.ExternalLabels
}

