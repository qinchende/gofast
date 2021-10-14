package stat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	counts := []int{1, 5, 10, 100, 1000, 1000}
	for _, count := range counts {
		m := NewMetrics("foo")
		m.SetName("bar")
		for i := 0; i < count; i++ {
			m.AddItem(ReqItem{
				Duration: time.Millisecond * time.Duration(i),
				//Description: strconv.Itoa(i),
			})
		}
		m.AddDrop()
		var writer mockedWriter
		//SetReportWriter(&writer)
		m.executor.Flush()
		assert.Equal(t, "bar", writer.report.Name)
	}
}

type mockedWriter struct {
	report *MetricInfo
}

func (m *mockedWriter) Write(report *MetricInfo) error {
	m.report = report
	return nil
}
