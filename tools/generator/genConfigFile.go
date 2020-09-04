package generator

var configTemplateFile = `
base:
  service_idc: testidc
  service_type: 1001
  service_id: 1
  service_ver: 1
  service_name: {{$.ServiceName}}
  service_port: 8888
  service_widget: 10
  service_funcs:{{range .Rpc}}    
    - {{.Name}}{{end}}
prometheus:
  switch_on: false
  listen_port: 9091
  client_histogram: 100,100,5
  server_histogram: 100,100,5
registry:
  type: etcd
  params:
    addr: 47.92.212.70:2379
    path: /myRpc
    timeout: 1
    report_time: 10
    update_time: 10
log:
  switch_on: false
  level: debug
  chan_size: 10000
  params:
    path: ../logs/
    max_size: 5000000
client_limit:
  switch_on: false
  type: token
  params:
    qps: 50000
    all_water: 10000
server_limit:
  switch_on: false
  type: token
  params:
    qps: 50000
    all_water: 10000
trace:
  switch_on: false
  addr: http://47.92.212.70:6831
  sample_type: const
  sample_rate: 1
balance:
  type: random
hystrix:
  switch_on: false
  timeout: 2000
  max_concurrent_requests: 100
  sleep_window: 1
  error_percent_threshold: 30
  request_volume_threshold: 10
other:
`
