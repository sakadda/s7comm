package s7comm

import (
	_ "embed"
	"fmt"
	"math"
<<<<<<< HEAD
	"sync"
=======
	"os"
>>>>>>> 52ca3ee (clear)
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/processors/dedup"
	"github.com/robinson/gos7"
)

//go:embed sample.conf
var sampleConfig string

type S7Comm struct {
	MetricName    string `toml:"name"`
	Endpoint      string `toml:"plc_ip"`
	Rack          int    `toml:"plc_rack"`
	Slot          int    `toml:"plc_slot"`
	ConnectType   int    `toml:"plc_connect_type" default:"3"`
	DedupInterval int    `toml:"dedup_interval" default:"10"`
	DedupEnable   bool   `toml:"dedup_enable" default:"false"`

	Timeout     config.Duration `toml:"connect_timeout"`
	IdleTimeout config.Duration `toml:"request_timeout"`

	Nodes []NodeSettings  `toml:"nodes"`
	Log   telegraf.Logger `toml:"-"`

	handler *gos7.TCPClientHandler
	client  gos7.Client
	helper  gos7.Helper

	dedup *dedup.Dedup
}

type NodeSettings struct {
	Name        string `toml:"name"`
	FullName    string `toml:"full_name"`
	Address     string `toml:"address"`
	Type        string `toml:"type"`
	EnableDedup bool   `toml:"dedup" default:"false"`
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
		DedupInterval: config.Duration(10 * time.Minute), //TODO use DedupInterval
		FlushTime:     time.Now(),
		Cache:         make(map[uint64]telegraf.Metric),
	}

	return nil
}

func (s *S7Comm) Gather(a telegraf.Accumulator) error {
	results := make(chan map[string]interface{}, len(s.Nodes))
	errs := make(chan error, len(s.Nodes))

<<<<<<< HEAD
	for _, node := range s.Nodes {
		wg.Add(1)
		go func(node NodeSettings) {
			defer wg.Done()
			buf := make([]byte, 8)
<<<<<<< HEAD
			_, err := s.client.Read(node.Address, buf)
			if err != nil {
				errs <- err
=======
=======
<<<<<<< HEAD
	var mu sync.Mutex

=======
>>>>>>> 52ca3ee (clear)
	for _, node := range s.Nodes {
		buf := make([]byte, 8)

		_, err := s.client.Read(node.Address, buf)
		if err != nil {
			errs <- fmt.Errorf("failed to connect for node %s: %v", node.Name, err)
			continue
		}
>>>>>>> 83ff122 (clear)

<<<<<<< HEAD
			mu.Lock()
			defer mu.Unlock()

			for {
				_, err := s.client.Read(node.Address, buf)
				if err != nil {
					s.Log.Error(fmt.Errorf("failed to read from node %s: %v", node.Name, err))
					time.Sleep(time.Duration(s.ReconnectInterval))

					if reconnectErr := s.Connect(); reconnectErr != nil {
						s.Log.Error(fmt.Errorf("failed to reconnect for node %s: %v", node.Name, reconnectErr))

						results <- map[string]interface{}{
							"name":      node.Name,
							"full_name": node.FullName,
							"fields":    nil,
							"dedup":     node.EnableDedup,
						}

						time.Sleep(time.Duration(s.ReconnectInterval))
						continue
					}
					continue
				}

				fields, err := s.readAndConvert(node, buf)
				if err != nil {
					errs <- fmt.Errorf("failed to convert data for node %s: %v", node.Name, err)
					return
				}

				results <- map[string]interface{}{
					"name":      node.Name,
					"full_name": node.FullName,
					"fields":    fields,
					"dedup":     node.EnableDedup,
				}

>>>>>>> 020eb29 (clear)
				return
			}

			fields, err := s.readAndConvert(node, buf)
			if err != nil {
				errs <- err
				return
			}

			results <- map[string]interface{}{
				"name":      node.Name,
				"full_name": node.FullName,
				"fields":    fields,
				"dedup":     node.EnableDedup,
			}
		}(node)
=======
		fields, err := s.readAndConvert(node, buf)
		if err != nil {
			errs <- fmt.Errorf("failed to convert data for node %s: %v", node.Name, err)
			continue
		}

		results <- map[string]interface{}{
			"name":      node.Name,
			"full_name": node.FullName,
			"fields":    fields,
			"dedup":     node.EnableDedup,
		}
>>>>>>> 52ca3ee (clear)
	}

	close(results)
	close(errs)

	for err := range errs {
		s.Log.Error(err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
=======
<<<<<<< HEAD
=======
<<<<<<< HEAD
=======
>>>>>>> 52ca3ee (clear)
>>>>>>> 83ff122 (clear)
	s.processMetrics(a, results)

	return nil
}

func (s *S7Comm) processMetrics(a telegraf.Accumulator, results chan map[string]interface{}) {
>>>>>>> fc2371f (wip)
	var dedupMetrics []telegraf.Metric
	var nonDedupMetrics []telegraf.Metric

<<<<<<< HEAD
=======
	var metrics []telegraf.Metric
>>>>>>> e0e2a2a (wip)
=======
>>>>>>> 52ca3ee (clear)
	for result := range results {
		name := result["name"].(string)
		full_name := result["full_name"].(string)
		fields := result["fields"].(map[string]interface{})
		enableDedup := result["dedup"].(bool)

		metric := metric.New(
			s.MetricName,
			map[string]string{
				"name":      name,
				"full_name": full_name,
			},
			fields,
			time.Now(),
		)
<<<<<<< HEAD
<<<<<<< HEAD

<<<<<<< HEAD
		if s.DedupEnable || enableDedup {
=======
		if s.DedupEnable || result["dedup"].(bool) {
=======
		metrics = append(metrics, metric)
	}

	// s.processMetrics(a, results)
	s.processBatch(a, metrics)

	return nil
}

func (s *S7Comm) processMetrics(a telegraf.Accumulator, metrics []telegraf.Metric) {
	var dedupMetrics []telegraf.Metric
	var nonDedupMetrics []telegraf.Metric
	const batchSize = 100
	var batch []telegraf.Metric

	for _, metric := range metrics {
		if s.DedupEnable || metric.Tags()["dedup"] == "true" {
			dedupMetrics = append(dedupMetrics, metric)
		} else {
			nonDedupMetrics = append(nonDedupMetrics, metric)
		}

		batch = append(batch, metric)
		if len(batch) >= batchSize {
			s.processBatch(a, batch)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		s.processBatch(a, batch)
	}
}

func (s *S7Comm) processBatch(a telegraf.Accumulator, batch []telegraf.Metric) {
	var dedupMetrics []telegraf.Metric
	var nonDedupMetrics []telegraf.Metric

	for _, metric := range batch {
		if s.DedupEnable || metric.Tags()["dedup"] == "true" {
>>>>>>> e0e2a2a (wip)
<<<<<<< HEAD
>>>>>>> fc2371f (wip)
=======
=======

		if s.DedupEnable || result["dedup"].(bool) {
>>>>>>> 52ca3ee (clear)
>>>>>>> 83ff122 (clear)
			dedupMetrics = append(dedupMetrics, metric)
		} else {
			nonDedupMetrics = append(nonDedupMetrics, metric)
		}
	}

	if len(dedupMetrics) > 0 {
		dedupMetrics = s.dedup.Apply(dedupMetrics...)
	}

	for _, metric := range dedupMetrics {
		a.AddMetric(metric)
	}
	for _, metric := range nonDedupMetrics {
		a.AddMetric(metric)
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
