# telegraf tailcount

Counts lines in growing file and return metrics on STDOUT.

Use with [`execd`](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/execd) plugin.

## Why ?

Mainly because apache ReqPerSec from mod_status is totally [useless](https://stackoverflow.com/questions/4630654/apache2-server-status-reported-value-for-requests-sec-is-wrong-what-am-i-doi).

## Usage

### Invocation

```
Usage: telegraf-tailcount [options] file
Options:
  -i int     Interval in seconds (default 10)
  -p int     Peak interval in seconds (default 1)
  -m string  Measurement name (default 'tailcount')
  -t string  Comma-separated k=p pairs for tags (default 'file=<file>')
```

e.g.:

```
$ telegraf-tailcount -t foo=bar -i 10 -p 1 /tmp/log
tailcount,file=/tmp/log,foo=bar sum=31,max=5,min=1
tailcount,file=/tmp/log,foo=bar sum=35,max=6,min=1
...
```

- _Interval_ is the metric output interval on STDOUT.
- _Peak interval_ is the interval for which telegraf-tailcount will compute min
  and max

So if you're looking at computing requests per second from an http log file,
you'll want _Peak interval_ to be 1.

Note that _Interval_ must be a multiple of _Peak interval_, and both must be
integers.

### Telegraf config

This tool is intended to be used with the 
[`execd` telegraf plugin](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/execd).

```toml
[[inputs.execd]]
  command = ["telegraf-tailcount", "-t", "foo=bar", "-i", "10", "-p", "1", "/tmp/log"]

  # the values below are the defaults and can be omitted
  signal = "none"
  restart_delay = "10s"
  data_format = "influx"
```

### Metrics

`telegraf-tailcount` outputs the following metrics:

- `sum`: number of lines read during `interval`
- `max`: maximum of lines read during a `peakInterval`
- `min`: minimum of lines read during a `peakInterval`