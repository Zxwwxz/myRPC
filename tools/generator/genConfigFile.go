package generator

var configTemplateFile = `
service_type: 1001
service_id: 1
service_ver: 1
service_name: {{.ServiceName}}
port: 8888
prometheus:
  switch_on: true
  port: 9091
register:
  switch_on: true
  register_path: /myRpc
  timeout: 1
  heart_beat: 10
  register_name: etcd
  register_addr: 47.92.212.70:2379
log:
  level: debug
  path: ./logs/
  chan_size: 10000
  console_log: true
limit:
  switch_on: true
  qps: 50000
  all_water: 100
trace:
  switch_on: true
  report_addr: http://47.92.212.70:9411/api/v1/spans
  sample_type: const
  sample_rate: 1
`
