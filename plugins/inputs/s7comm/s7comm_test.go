package s7comm

import (
	"testing"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/robinson/gos7"
)

// MockAccumulator for testing purposes
type MockAccumulator struct {
	metrics []telegraf.Metric
}

func (m *MockAccumulator) AddCounter(name string, fields map[string]interface{}, tags map[string]string, t ...time.Time) {}
func (m *MockAccumulator) AddFields(name string, fields map[string]interface{}, tags map[string]string, t ...time.Time) {}
func (m *MockAccumulator) AddGauge(name string, fields map[string]interface{}, tags map[string]string, t ...time.Time)  {}
func (m *MockAccumulator) AddSummary(name string, fields map[string]interface{}, tags map[string]string, t ...time.Time) {}
func (m *MockAccumulator) AddHistogram(name string, fields map[string]interface{}, tags map[string]string, t ...time.Time) {}
func (m *MockAccumulator) SetPrecision(precision time.Duration) {}
func (m *MockAccumulator) AddError(err error) {}
func (m *MockAccumulator) AddMetric(metric telegraf.Metric) {}
func (m *MockAccumulator) WithTracking(maxTracked int) telegraf.TrackingAccumulator {
	return nil
}

// Benchmark function for Gather
func BenchmarkGather(b *testing.B) {
	s := &S7Comm{
		MetricName:  "s7comm",
		Endpoint:    "192.168.4.159",
		Rack:        0,
		Slot:        2,
		Timeout:     config.Duration(5 * time.Second),
		IdleTimeout: config.Duration(10 * time.Second),
		Nodes: []NodeSettings{
			{Name: "PAB11CT001.Par", Address: "DB1.DBD74", Type: "real"},
			{Name: "PAB20CT001.Par", Address: "DB1.DBD330", Type: "real"},
			{Name: "PCM01CT001.Par", Address: "DB1.DBD138", Type: "real"},
			// Добавьте больше узлов по мере необходимости
		},
	}

	// Инициализация соединения
	if err := s.Connect(); err != nil {
		b.Fatalf("failed to connect: %v", err)
	}
	defer s.Stop()

	acc := &MockAccumulator{}

	// Запуск тестов
	for i := 0; i < b.N; i++ {
		if err := s.Gather(acc); err != nil {
			b.Fatalf("failed to gather: %v", err)
		}
	}
}

func BenchmarkReadAndConvert(b *testing.B) {
	s := &S7Comm{
		helper: gos7.Helper{},
	}

	node := NodeSettings{Name: "DB1.DBD74", Address: "DB1.DBD74", Type: "real"}
	buf := make([]byte, 8)

	for i := 0; i < b.N; i++ {
		if _, err := s.readAndConvert(node, buf); err != nil {
			b.Fatalf("failed to read and convert: %v", err)
		}
	}
}
