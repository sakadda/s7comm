# S7comm

S7comm is a [telegraf](https://github.com/influxdata/telegraf) external input plugin to gather data from Siemens PLC using the [gos7](https://github.com/robinson/gos7) library.

## Installation

Download the repo

```bash
git clone git@github.com:sakadda/s7comm.git
```

Change the poll interval if needed in the cmd/main.go file

```golang
var pollInterval = flag.Duration("poll_interval", 1*time.Second, "how often to send metrics")
```

Build the binary

```bash
go build -o s7comm cmd/main.go
```

Create your plugin.config file

```toml @plugin.conf
[[inputs.s7comm]]
	name = "S7300"
	plc_ip = "192.168.10.57"
	plc_rack = 0
	plc_slot = 1
    	connect_timeout = "10s"
    	request_timeout = "2s"
    	nodes = [{name= "test_int", address= "DB1.DBW0", type = "int"},
        	{name= "test_real", address= "DB1.DBD2",type = "real"},
        	{name= "test_bool", address= "DB1.DBX10.0",type = "bool"},
        	{name= "test_dint", address= "DB1.DBD12",type = "dint"},
        	{name= "test_uint", address= "DB1.DBW16",type = "uint"},
        	{name= "test_udint", address= "DB1.DBD18",type = "udint"}]
```

| Data Type | Column 2 | Column 3 |
| --------- | -------- | -------- |
| bool      | Cell 2   | Cell 3   |
| byte      | Cell 5   | Cell 6   |
| dword     | Cell 8   | Cell 9   |
| int       | Cell 8   | Cell 9   |
| dint      | Cell 8   | Cell 9   |
| uint      | Cell 8   | Cell 9   |
| udint     | Cell 8   | Cell 9   |
| real      | Cell 8   | Cell 9   |
| float64   | Cell 8   | Cell 9   |
| time      | Cell 8   | Cell 9   |

From here, you can already test the plugin with your config file.

```bash
./s7comm -config plugin.conf
./s7comm.exe -config plugin.conf
```

If everything is ok, you should see something like this

```bash
test_int value=8056i 1623227848846884706
test_real value=403.14764404296875 1623227849849296642
```

## Datatypes

You have to specifie the data type in your config file for each node. At the moment, only those types are implemented :

bool, byte, word, dword, int, dint, uint, udint, real, float64, time

## Telegraf configuration

To use the plugin with telegraf, add this configuration to your main telegraf.conf file. S7comm is an external plugin using [shim](https://github.com/influxdata/telegraf/blob/master/plugins/common/shim/README.md) and [execd](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/execd). Go see their doc for more information.

```toml telegraf.conf
# Run executable as long-running input plugin
[[inputs.execd]]
  ## One program to run as daemon.
  ## NOTE: process and each argument should each be their own string
  #command = ["telegraf-smartctl", "-d", "/dev/sda"]
command = ["/home/r/s7comm/s7comm", "-config", "/home/r/s7comm/CT.conf"]
  ## Environment variables
  ## Array of "key=value" pairs to pass as environment variables
  ## e.g. "KEY=value", "USERNAME=John Doe",
  ## "LD_LIBRARY_PATH=/opt/custom/lib64:/usr/local/libs"
  # environment = []

  ## Define how the process is signaled on each collection interval.
  ## Valid values are:
  ##   "none"    : Do not signal anything. (Recommended for service inputs)
  ##               The process must output metrics by itself.
  ##   "STDIN"   : Send a newline on STDIN. (Recommended for gather inputs)
  ##   "SIGHUP"  : Send a HUP signal. Not available on Windows. (not recommended)
  ##   "SIGUSR1" : Send a USR1 signal. Not available on Windows.
  ##   "SIGUSR2" : Send a USR2 signal. Not available on Windows.
  signal = "STDIN"

  ## Delay before the process is restarted after an unexpected termination
  restart_delay = "30s"

  ## Buffer size used to read from the command output stream
  ## Optional parameter. Default is 64 Kib, minimum is 16 bytes
  # buffer_size = "64Kib"

  ## Data format to consume.
  ## Each data format has its own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
  data_format = "influx"
```

## PLC configuration

S7-300 and S7-400 usually use rack 0 and slot 2 and dont require additional configuration.

S7-1200 and S7-1500 usually use rack 0 and slot 1 and you need to enable the PUT/GET operations in the hardware configuration of your PLC and you have to set DBs as non-optimized.

Be aware of security issue. Once S7 Communication is enabled in a CPU, there is no way to block communication with a partner device. This means that any device on the same network can read and write data to the CPU using the S7 Communication protocol. For this reason, I would recommend using the native OPC.UA server for the newer S7-1200 and S7-1500 PLCs. See the [OPC.UA](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/opcua) telegraf plugin.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Metrics
