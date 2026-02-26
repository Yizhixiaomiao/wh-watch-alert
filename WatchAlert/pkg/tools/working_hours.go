package tools

import (
	"fmt"
	"time"
)

// WorkingHoursConfig 工作时间配置
type WorkingHoursConfig struct {
	StartTime string // 开始时间，格式 "09:00"
	EndTime   string // 结束时间，格式 "18:00"
	Weekdays  []int // 工作日，0=周日, 1=周一, ..., 6=周六
	Holidays  []string // 节假日列表，格式 "2006-01-02"
}

// DefaultWorkingHours 默认工作时间配置
var DefaultWorkingHours = WorkingHoursConfig{
	StartTime: "09:00",
	EndTime:   "18:00",
	Weekdays:  []int{1, 2, 3, 4, 5}, // 周一到周五
	Holidays:  []string{},
}

// ParseWorkingHoursConfig 解析工作时间配置字符串
// 格式: "09:00-18:00;1,2,3,4,5;2024-01-01,2024-02-10"
func ParseWorkingHoursConfig(configStr string) (*WorkingHoursConfig, error) {
	if configStr == "" {
		return &DefaultWorkingHours, nil
	}

	// 解析配置字符串 (简化版，实际需要更复杂的解析)
	// 格式: "09:00-18:00;1,2,3,4,5;2024-01-01,2024-02-10"
	// 这里暂时返回默认配置
	return &DefaultWorkingHours, nil
}

// IsWorkingDay 判断是否是工作日
func IsWorkingDay(t time.Time, config *WorkingHoursConfig) bool {
	if config == nil {
		config = &DefaultWorkingHours
	}

	// 检查是否是节假日
	dateStr := t.Format("2006-01-02")
	for _, holiday := range config.Holidays {
		if holiday == dateStr {
			return false
		}
	}

	// 检查是否是工作日
	weekday := int(t.Weekday())
	for _, wd := range config.Weekdays {
		if wd == weekday {
			return true
		}
	}

	return false
}

// IsWorkingTime 判断是否是工作时间
func IsWorkingTime(t time.Time, config *WorkingHoursConfig) bool {
	if !IsWorkingDay(t, config) {
		return false
	}

	if config == nil {
		config = &DefaultWorkingHours
	}

	// 解析开始和结束时间
	startTime, err := time.Parse("15:04", config.StartTime)
	if err != nil {
		return false
	}

	endTime, err := time.Parse("15:04", config.EndTime)
	if err != nil {
		return false
	}

	// 获取当前时间的小时和分钟
	hour, min, _ := t.Clock()
	currentTime := time.Date(0, 0, 0, hour, min, 0, 0, time.UTC)

	return currentTime.After(startTime) || currentTime.Equal(startTime) && 
	       currentTime.Before(endTime)
}

// CalculateWorkingDays 计算工作日天数
func CalculateWorkingDays(startTime, endTime time.Time, config *WorkingHoursConfig) int {
	if config == nil {
		config = &DefaultWorkingHours
	}

	days := 0
	for t := startTime; !t.After(endTime); t = t.Add(24 * time.Hour) {
		if IsWorkingDay(t, config) {
			days++
		}
	}
	return days
}

// AddWorkingDays 添加工作日
func AddWorkingDays(startTime time.Time, days int, config *WorkingHoursConfig) time.Time {
	if config == nil {
		config = &DefaultWorkingHours
	}

	result := startTime
	addedDays := 0

	for addedDays < days {
		result = result.Add(24 * time.Hour)
		if IsWorkingDay(result, config) {
			addedDays++
		}
	}

	return result
}

// CalculateSLADueTime 计算SLA截止时间（考虑工作日和工作时间）
// seconds: SLA时间（秒）
// config: 工作时间配置
func CalculateSLADueTime(startTime time.Time, seconds int64, config *WorkingHoursConfig) time.Time {
	if config == nil {
		config = &DefaultWorkingHours
	}

	// 如果没有配置工作时间，直接计算
	if config.StartTime == "" || config.EndTime == "" {
		return startTime.Add(time.Duration(seconds) * time.Second)
	}

	// 解析开始和结束时间
	startHourMin, err := time.Parse("15:04", config.StartTime)
	if err != nil {
		return startTime.Add(time.Duration(seconds) * time.Second)
	}

	endHourMin, err := time.Parse("15:04", config.EndTime)
	if err != nil {
		return startTime.Add(time.Duration(seconds) * time.Second)
	}

	// 计算每天的工作时长（秒）
	dailyWorkingSeconds := endHourMin.Sub(startHourMin).Seconds()

	// 需要的工作天数
	workingDays := int(seconds) / int(dailyWorkingSeconds)
	remainingSeconds := seconds % int64(dailyWorkingSeconds)

	// 计算截止日期
	dueDate := AddWorkingDays(startTime, workingDays, config)

	// 计算截止时间
	startTimeOfDay := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 
		startHourMin.Hour(), startHourMin.Minute(), 0, 0, dueDate.Location())
	
	dueTime := startTimeOfDay.Add(time.Duration(remainingSeconds) * time.Second)

	// 如果截止时间超出了当天的工作时间，需要顺延到下一个工作日
	if !IsWorkingTime(dueTime, config) {
		dueTime = time.Date(dueTime.Year(), dueTime.Month(), dueTime.Day(), 
			endHourMin.Hour(), endHourMin.Minute(), 0, 0, dueTime.Location())
		dueTime = AddWorkingDays(dueTime, 1, config)
		dueTime = time.Date(dueTime.Year(), dueTime.Month(), dueTime.Day(), 
			startHourMin.Hour(), startHourMin.Minute(), 0, 0, dueTime.Location())
	}

	return dueTime
}

// CalculateSLARemainingTime 计算SLA剩余时间（考虑工作日和工作时间）
func CalculateSLARemainingTime(startTime, currentTime time.Time, totalSeconds int64, config *WorkingHoursConfig) int64 {
	if config == nil {
		config = &DefaultWorkingHours
	}

	dueTime := CalculateSLADueTime(startTime, totalSeconds, config)
	
	if currentTime.After(dueTime) {
		return 0
	}

	return int64(dueTime.Sub(currentTime).Seconds())
}

// IsOverdue 判断是否逾期（考虑工作日和工作时间）
func IsOverdue(startTime, currentTime time.Time, totalSeconds int64, config *WorkingHoursConfig) bool {
	remaining := CalculateSLARemainingTime(startTime, currentTime, totalSeconds, config)
	return remaining <= 0
}

// GetWorkingHoursRange 获取指定日期的工作时间范围
func GetWorkingHoursRange(date time.Time, config *WorkingHoursConfig) (time.Time, time.Time, error) {
	if config == nil {
		config = &DefaultWorkingHours
	}

	startTime, err := time.Parse("15:04", config.StartTime)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start time format: %v", err)
	}

	endTime, err := time.Parse("15:04", config.EndTime)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end time format: %v", err)
	}

	start := time.Date(date.Year(), date.Month(), date.Day(), 
		startTime.Hour(), startTime.Minute(), 0, 0, date.Location())
	end := time.Date(date.Year(), date.Month(), date.Day(), 
		endTime.Hour(), endTime.Minute(), 0, 0, date.Location())

	return start, end, nil
}