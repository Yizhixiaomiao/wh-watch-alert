package tools

import (
	"testing"
	"time"

	"watchAlert/pkg/tools"
)

func TestParseWorkingHoursConfig(t *testing.T) {
	tests := []struct {
		name          string
		configStr     string
		wantStartHour int
		wantEndHour   int
		wantErr       bool
	}{
		{
			name:      "valid config",
			configStr: `{"workDays":[1,2,3,4,5],"startTime":"09:00","endTime":"18:00","holidays":[]}`,
			wantErr:   false,
		},
		{
			name:      "invalid json",
			configStr: `{invalid}`,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := tools.ParseWorkingHoursConfig(tt.configStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWorkingHoursConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && config != nil {
				if config.StartTime.Hour() != tt.wantStartHour {
					t.Errorf("StartHour = %v, want %v", config.StartTime.Hour(), tt.wantStartHour)
				}
				if config.EndTime.Hour() != tt.wantEndHour {
					t.Errorf("EndHour = %v, want %v", config.EndTime.Hour(), tt.wantEndHour)
				}
			}
		})
	}
}

func TestIsWorkingDay(t *testing.T) {
	configStr := `{"workDays":[1,2,3,4,5],"startTime":"09:00","endTime":"18:00","holidays":[]}`
	config, _ := tools.ParseWorkingHoursConfig(configStr)

	tests := []struct {
		name string
		t    time.Time
		want bool
	}{
		{
			name: "monday",
			t:    time.Date(2026, 2, 3, 10, 0, 0, 0, time.UTC), // Monday
			want: true,
		},
		{
			name: "saturday",
			t:    time.Date(2026, 2, 7, 10, 0, 0, 0, time.UTC), // Saturday
			want: false,
		},
		{
			name: "sunday",
			t:    time.Date(2026, 2, 8, 10, 0, 0, 0, time.UTC), // Sunday
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tools.IsWorkingDay(tt.t, config); got != tt.want {
				t.Errorf("IsWorkingDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsWorkingTime(t *testing.T) {
	configStr := `{"workDays":[1,2,3,4,5],"startTime":"09:00","endTime":"18:00","holidays":[]}`
	config, _ := tools.ParseWorkingHoursConfig(configStr)

	tests := []struct {
		name string
		t    time.Time
		want bool
	}{
		{
			name: "during working hours",
			t:    time.Date(2026, 2, 3, 10, 0, 0, 0, time.UTC), // Monday 10:00
			want: true,
		},
		{
			name: "before working hours",
			t:    time.Date(2026, 2, 3, 8, 0, 0, 0, time.UTC), // Monday 08:00
			want: false,
		},
		{
			name: "after working hours",
			t:    time.Date(2026, 2, 3, 19, 0, 0, 0, time.UTC), // Monday 19:00
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tools.IsWorkingTime(tt.t, config); got != tt.want {
				t.Errorf("IsWorkingTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateWorkingDays(t *testing.T) {
	configStr := `{"workDays":[1,2,3,4,5],"startTime":"09:00","endTime":"18:00","holidays":[]}`
	config, _ := tools.ParseWorkingHoursConfig(configStr)

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		want      int
	}{
		{
			name:      "one week",
			startTime: time.Date(2026, 2, 3, 10, 0, 0, 0, time.UTC), // Monday
			endTime:   time.Date(2026, 2, 9, 10, 0, 0, 0, time.UTC), // Sunday
			want:      5,
		},
		{
			name:      "one day",
			startTime: time.Date(2026, 2, 3, 10, 0, 0, 0, time.UTC), // Monday
			endTime:   time.Date(2026, 2, 4, 10, 0, 0, 0, time.UTC), // Tuesday
			want:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tools.CalculateWorkingDays(tt.startTime, tt.endTime, config); got != tt.want {
				t.Errorf("CalculateWorkingDays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateSLADueTime(t *testing.T) {
	configStr := `{"workDays":[1,2,3,4,5],"startTime":"09:00","endTime":"18:00","holidays":[]}`
	config, _ := tools.ParseWorkingHoursConfig(configStr)

	tests := []struct {
		name         string
		startTime    time.Time
		totalSeconds int64
		wantHours    float64
	}{
		{
			name:         "8 hours",
			startTime:    time.Date(2026, 2, 3, 9, 0, 0, 0, time.UTC), // Monday 09:00
			totalSeconds: 8 * 3600,
			wantHours:    8,
		},
		{
			name:         "24 hours",
			startTime:    time.Date(2026, 2, 3, 9, 0, 0, 0, time.UTC), // Monday 09:00
			totalSeconds: 24 * 3600,
			wantHours:    24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dueTime := tools.CalculateSLADueTime(tt.startTime, tt.totalSeconds, config)
			duration := dueTime.Sub(tt.startTime).Hours()
			if duration < tt.wantHours-0.1 || duration > tt.wantHours+0.1 {
				t.Errorf("CalculateSLADueTime() duration = %v hours, want %v hours", duration, tt.wantHours)
			}
		})
	}
}

func TestIsOverdue(t *testing.T) {
	configStr := `{"workDays":[1,2,3,4,5],"startTime":"09:00","endTime":"18:00","holidays":[]}`
	config, _ := tools.ParseWorkingHoursConfig(configStr)

	tests := []struct {
		name         string
		startTime    time.Time
		currentTime  time.Time
		totalSeconds int64
		want         bool
	}{
		{
			name:         "not overdue",
			startTime:    time.Date(2026, 2, 3, 9, 0, 0, 0, time.UTC),  // Monday 09:00
			currentTime:  time.Date(2026, 2, 3, 10, 0, 0, 0, time.UTC), // Monday 10:00
			totalSeconds: 24 * 3600,
			want:         false,
		},
		{
			name:         "overdue",
			startTime:    time.Date(2026, 2, 3, 9, 0, 0, 0, time.UTC),  // Monday 09:00
			currentTime:  time.Date(2026, 2, 6, 10, 0, 0, 0, time.UTC), // Thursday 10:00
			totalSeconds: 24 * 3600,
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tools.IsOverdue(tt.startTime, tt.currentTime, tt.totalSeconds, config); got != tt.want {
				t.Errorf("IsOverdue() = %v, want %v", got, tt.want)
			}
		})
	}
}
