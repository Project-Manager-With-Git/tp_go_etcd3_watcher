package etcdwatcher

import (
	"os"
	"strings"

	log "github.com/Golang-Tools/loggerhelper"
	s "github.com/Golang-Tools/schema-entry-go"

	// "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Application struct {
	App_Version                   string   `json:"app_version" jsonschema:"required,title=v,description=应用版本"`
	App_Name                      string   `json:"app_name" jsonschema:"required,title=n,description=应用名"`
	Log_Level                     string   `json:"log_level" jsonschema:"required,title=l,description=log等级,enum=TRACE,enum=DEBUG,enum=INFO,enum=WARN,enum=ERROR"`
	Watch_Kafka_Topics            []string `json:"watch_kafka_topics" jsonschema:"required,title=t,description=监听的topic列表"`
	Watch_Kafka_URLS              []string `json:"watch_kafka_urls" jsonschema:"required,title=u,description=监听的kafka集群地址"`
	Watch_Kafka_Group_ID          string   `json:"watch_kafka_group_id" jsonschema:"required,title=g,description=监听的kafka集群的消费者群组id"`
	Watch_Kafka_Auto_Offset_Reset string   `json:"watch_kafka_auto_offset_reset" jsonschema:"required,title=a,description=测试列表,enum=earliest,enum=latest"`
	Watch_Kafka_Isolation_Level   string   `json:"watch_kafka_isolation_level" jsonschema:"required,title=i,description=kafka的隔离级别,enum=read_uncommitted,enum=read_committed"`
}

func (app *Application) Main() {
	log.Init(app.Log_Level, map[string]interface{}{"app_name": app.App_Name})
	err := Watcher.InitFromOptions(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(app.Watch_Kafka_URLS, ","),
		"group.id":          app.Watch_Kafka_Group_ID,
		"auto.offset.reset": app.Watch_Kafka_Auto_Offset_Reset,
		"isolation.level":   app.Watch_Kafka_Isolation_Level,
	})
	if err != nil {
		log.Error("Failed to init watcher", log.Dict{"error": err})
		os.Exit(1)
	}
	defer Watcher.Close()
	err = Watcher.SubscribeTopics(app.Watch_Kafka_Topics, nil)
	if err != nil {
		log.Error("Failed to Subscribe Topics", log.Dict{"error": err})
		return
	}
	defer Watcher.Unsubscribe()
	err = Watcher.handdler()
	if err != nil {
		log.Error("handdler get error", log.Dict{"error": err})
	}
}

// root, _ := s.New(&s.EntryPointMeta{Name: "tp_go_etcd3_watcher", Usage: "tp_go_etcd3_watcher"})
// nodeb, _ := s.New(&s.EntryPointMeta{Name: "sender", Usage: "tp_go_etcd3_watcher sender"})
var App, _ = s.New(&s.EntryPointMeta{Name: "watcher", Usage: "tp_go_etcd3_watcher watcher"}, &Application{
	App_Version:                   "0.0.0",
	App_Name:                      "tp_go_etcd3_watcher_watcher",
	Log_Level:                     "DEBUG",
	Watch_Kafka_Auto_Offset_Reset: "latest",
	Watch_Kafka_Isolation_Level:   "read_committed",
})

// s.RegistSubNode(root, nodeb)
// nodeb.RegistSubNode(nodec)
// os.Setenv("FOO_BAR_PAR_A", "123")
// root.Parse([]string{"foo", "bar", "par", "--Field=4", "--Field=5", "--Field=6"})
