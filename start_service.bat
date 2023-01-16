:: 用于go的logAgent项目启动服务
:: etcd kafka zookeeper influxdb grafana es kibana

color 9
@echo off
chcp 65001

set kafkaPath=E:\Go\kafka_2.12-3.3.1
set etcdPath=E:\Go\etcd\etcd-v3.4.23-windows-amd64
set influxdbPath=E:\Go\influxDB\influxdb-1.7.7-1
set grafanaPath=E:\Go\grafana\grafana-9.3.2
set esPath=E:\Go\es\elasticsearch-7.3.0
set kibanaPath=E:\Go\es\kibana-7.3.0-windows-x86_64

@echo script is running......
IF EXIST %kafkaPath% (echo kafka path exist) else (echo kafka path error)
IF EXIST %etcdPath% (echo etcd path exist) else (echo etcd path error)
IF EXIST %influxdbPath% (echo influxdb path exist) else (echo influxdb path error)
IF EXIST %grafanaPath% (echo grafana path exist) else (echo grafana path error)
IF EXIST %esPath% (echo es path exist) else (echo es path error)
IF EXIST %kibanaPath% (echo kiana path exist) else (echo kiana path error)

@echo start service......
start cmd /k "cd /d %kafkaPath%&&bin\windows\zookeeper-server-start.bat config\zookeeper.properties"
start cmd /k "cd /d %etcdPath%&&etcd.exe"
start cmd /k "cd /d %influxdbPath%&&influxd.exe"
start cmd /k "cd /d %grafanaPath%&&bin\grafana-server.exe"
start cmd /k "cd /d %esPath%&&bin\elasticsearch.bat"

:: 睡眠一会，避免zookeeper节点未完全启动导致 kafka 节点启动失败
timeout /nobreak /t 5
start cmd /k "cd /d %kafkaPath%&&bin\windows\kafka-server-start.bat config\server.properties"

@echo start kafka consumer......
start cmd /k "cd /d %kafkaPath%&&bin\windows\kafka-console-consumer.bat --bootstrap-server 127.0.0.1:9092 --topic web_log --from-beginning"
@echo start influxdb cli......
start cmd /k "cd /d %influxdbPath%&&influx.exe"
@echo start kibana......
start cmd /k "cd /d %kibanaPath%&&bin\kibana.bat"

if errorlevel 0 (
    @echo script run success
) else (
    @echo script run failed
)

pause

