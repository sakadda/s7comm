package s7comm

import (
	_ "embed"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/processors/dedup"
	"github.com/robinson/gos7"
)

//go:embed sample.conf
var sampleConfig string

// S7Comm
type S7Comm struct {
	MetricName  string `toml:"name"`
	Endpoint    string `toml:"plc_ip"`
	Rack        int    `toml:"plc_rack"`
	Slot        int    `toml:"plc_slot"`
	ConnectType int    `toml:"plc_connect_type" default:"3"`

	Timeout     config.Duration `toml:"connect_timeout"`
	IdleTimeout config.Duration `toml:"request_timeout"`

	Nodes []NodeSettings  `toml:"nodes"`
	Log   telegraf.Logger `toml:"-"`

	// internal values
	handler *gos7.TCPClientHandler
	client  gos7.Client
	helper  gos7.Helper

	// dedup processor
	dedup *dedup.Dedup
}

type NodeSettings struct {
	Name        string `toml:"name"`
	Address     string `toml:"address"`
	Type        string `toml:"type"`
	EnableDedup bool   `toml:"enable_dedup" default:"false"`
}

func (s *S7Comm) SampleConfig() string {
	return sampleConfig
}

func (s *S7Comm) Connect() error {
	s.handler = gos7.NewTCPClientHandlerWithConnectType(s.Endpoint, s.Rack, s.Slot, s.ConnectType)
	s.handler.Timeout = time.Duration(s.Timeout)
	s.handler.IdleTimeout = time.Duration(s.IdleTimeout)

	err := s.handler.Connect()
	if err != nil {
		return err
	}

	// defer s.handler.Close()

	s.client = gos7.NewClient(s.handler)
	s.helper = gos7.Helper{}

	// s.Log.Debug("Connection successfull")

	return nil
}

func (s *S7Comm) Stop() error {
	err := s.handler.Close()

	return err
}

func (s *S7Comm) Init() error {
	err := s.Connect()
	if err != nil {
		return err
	}

	s.dedup = &dedup.Dedup{
		DedupInterval: config.Duration(10 * time.Minute),
		FlushTime:     time.Now(),
		Cache:         make(map[uint64]telegraf.Metric),
	}

	return nil
}

func (s *S7Comm) Gather(a telegraf.Accumulator) error {
	var wg sync.WaitGroup
	results := make(chan map[string]interface{}, len(s.Nodes))
	errs := make(chan error, len(s.Nodes))

	for _, node := range s.Nodes {
		wg.Add(1)
		go func(node NodeSettings) {
			defer wg.Done()
			buf := make([]byte, 8)
			_, err := s.client.Read(node.Address, buf)
			if err != nil {
				errs <- err
				return
			}

			fields, err := s.readAndConvert(node, buf)
			if err != nil {
				errs <- err
				return
			}

			results <- fields
		}(node)
	}

	wg.Wait()
	close(results)
	close(errs)

	for err := range errs {
		s.Log.Error(err)
	}

	for fields := range results {
		// Создание метрики для dedup
		a.AddFields(s.MetricName, fields, nil)
	}

	return nil
}

func (s *S7Comm) readAndConvert(node NodeSettings, buf []byte) (map[string]interface{}, error) {
	fields := make(map[string]interface{}, 1)

	switch node.Type {
	case "bool":
		var res bool
		s.helper.GetValueAt(buf, 0, &res)
		fields[node.Name] = res
	case "byte":
		var res byte
		s.helper.GetValueAt(buf, 0, &res)
		fields[node.Name] = res
	case "word":
		var res uint16
		s.helper.GetValueAt(buf, 0, &res)
		fields[node.Name] = res
	case "dword":
		var res uint32
		s.helper.GetValueAt(buf, 0, &res)
		fields[node.Name] = res
	case "int":
		var res int16
		s.helper.GetValueAt(buf, 0, &res)
		fields[node.Name] = res
	case "dint":
		var res int32
		s.helper.GetValueAt(buf, 0, &res)
		fields[node.Name] = res
	case "uint":
		var res uint16
		s.helper.GetValueAt(buf, 0, &res)
		fields[node.Name] = res
	case "udint":
		var res uint32
		s.helper.GetValueAt(buf, 0, &res)
		fields[node.Name] = res
	case "real":
		var res float32
		s.helper.GetValueAt(buf, 0, &res)
		if math.IsNaN(float64(res)) {
			res = 0
		}
		fields[node.Name] = res
	case "float64":
		var res float64
		s.helper.GetValueAt(buf, 0, &res)
		fields[node.Name] = res
	case "time":
		var res uint32
		s.helper.GetValueAt(buf, 0, &res)
		fields[node.Name] = res
	default:
		return nil, fmt.Errorf("unknown data type: %s", node.Type)
	}

	return fields, nil
}

func (s *S7Comm) Description() string {
	return "Read data from Siemens PLC using S7 protocol with S7Go"
}

// Add this plugin to telegraf
func init() {
	inputs.Add("s7comm", func() telegraf.Input {
		return &S7Comm{
			MetricName:  "s7comm",
			Endpoint:    "192.168.4.159",
			Rack:        0,
			Slot:        2,
			Timeout:     config.Duration(5 * time.Second),
			IdleTimeout: config.Duration(10 * time.Second),
			Nodes:       nil,
		}
	})
}
