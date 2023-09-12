//go:build go1.16

package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/winjeg/go-commons/log"
)

const defaultSize = 8

type MonitorConfig struct {
	Port int      `json:"port" yaml:"port"`
	Path string   `json:"path" yaml:"path"`
	Tags []string `json:"tags" yaml:"tags"`
}

func NewDefaultMonitorCfg() *MonitorConfig {
	return &MonitorConfig{
		Port: 10000,
		Path: "/metrics",
		Tags: nil,
	}
}

// NewIrisMonitor creates a new Monitor if needed
// recommended usage is app.UseGlobal  not the useRouter function
// for that attackers may attack on those 404 not found uris, which UseRouter function
// will count, causing the metrics too heavy
func NewIrisMonitor(cfg *MonitorConfig) iris.Handler {
	m := NewMonitor(cfg)
	m.Start()

	return func(ctx iris.Context) {
		start := time.Now()
		ctx.Next()
		duration := time.Now().Sub(start)
		cost := float64(duration.Microseconds() / 1000.0)
		errMsg := "None"
		if ctx.Err() != nil {
			errMsg = ctx.Err().Error()
		}

		tags := []string{
			"method", ctx.Method(),
			"uri", ctx.Path(),
			"host", ctx.RemoteAddr(),
			"status", strconv.Itoa(ctx.GetStatusCode()),
			"error", errMsg,
		}
		m.Inc("http_request", tags...)
		m.Timer("http_request_cost", cost, tags...)
	}
}

var (
	monitorLock = sync.Mutex{}
	counterLock = sync.Mutex{}
	gaugeLock   = sync.Mutex{}
	timerLock   = sync.Mutex{}

	counterMap = make(map[string]*prometheus.CounterVec, defaultSize)
	gaugeMap   = make(map[string]*prometheus.GaugeVec, defaultSize)
	timerMap   = make(map[string]*prometheus.HistogramVec, defaultSize)
)

type monitor struct {
	Cfg *MonitorConfig
}

var localMonitor *monitor

func NewMonitor(cfg *MonitorConfig) *monitor {
	if localMonitor != nil {
		return localMonitor
	}
	monitorLock.Lock()
	localMonitor = &monitor{cfg}
	monitorLock.Unlock()
	return localMonitor
}

func (m *monitor) Start() {
	http.Handle(m.Cfg.Path, promhttp.Handler())
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", m.Cfg.Port), nil)
		if err != nil {
			log.GetLogger(nil).Errorln("error starting monitor: " + err.Error())
		}
	}()
	log.GetLogger(nil).Infof("prometheus metrics started at: http://localhost:%d%s", m.Cfg.Port, m.Cfg.Path)
}

func (m *monitor) composeKey(name string, tags []string) string {
	return fmt.Sprintf("%s-%s", name, strings.Join(tags, "-"))
}

func (m *monitor) finalTags(tags []string) []string {
	result := make([]string, 0, defaultSize)
	result = append(result, m.Cfg.Tags...)
	result = append(result, tags...)
	return result
}

func (m *monitor) tagParts(tags []string) ([]string, []string) {
	names := make([]string, 0, defaultSize)
	values := make([]string, 0, defaultSize)
	if len(tags)%2 != 0 {
		return names, values
	}
	for i := range tags {
		if i%2 == 0 {
			names = append(names, tags[i])
		} else {
			values = append(values, tags[i])
		}
	}
	return names, values
}

func (m *monitor) Inc(name string, tags ...string) {
	m.Count(name, 1, tags...)
}

func (m *monitor) Dec(name string, tags ...string) {
	m.Count(name, -1, tags...)
}

func (m *monitor) Count(name string, count float64, tags ...string) {
	names, values := m.tagParts(m.finalTags(tags))
	key := m.composeKey(name, names)
	if _, ok := counterMap[key]; !ok {
		counterLock.Lock()
		if _, ok := counterMap[key]; !ok {
			counterMap[key] = promauto.NewCounterVec(prometheus.CounterOpts{
				Name: name,
				Help: fmt.Sprintf("%s counter", name),
			}, names)
		}
		counterLock.Unlock()
	}
	counterMap[key].WithLabelValues(values...).Add(count)
}

func (m *monitor) Timer(name string, v float64, tags ...string) {
	buckets := prometheus.ExponentialBucketsRange(0.1, 30000, 50)
	m.TimerWithBuckets(name, v, buckets, tags...)
}

func (m *monitor) TimerWithBuckets(name string, v float64, buckets []float64, tags ...string) {
	names, values := m.tagParts(m.finalTags(tags))
	key := m.composeKey(name, names)
	if _, ok := timerMap[key]; !ok {
		timerLock.Lock()
		if _, ok := timerMap[key]; !ok {
			timerMap[key] = promauto.NewHistogramVec(prometheus.HistogramOpts{
				Name:    name,
				Help:    fmt.Sprintf("%s histogram", name),
				Buckets: buckets,
			}, names)
		}
		timerLock.Unlock()
	}
	timerMap[key].WithLabelValues(values...).Observe(v)
}

func (m *monitor) Gauge(name string, v float64, tags ...string) {
	names, values := m.tagParts(m.finalTags(tags))
	key := m.composeKey(name, names)
	if _, ok := gaugeMap[key]; !ok {
		gaugeLock.Lock()
		if _, ok := gaugeMap[key]; !ok {
			gaugeMap[key] = promauto.NewGaugeVec(prometheus.GaugeOpts{
				Name: name,
				Help: fmt.Sprintf("%s gauge", name),
			}, names)
		}
		gaugeLock.Unlock()
	}
	gaugeMap[key].WithLabelValues(values...).Set(v)
}
