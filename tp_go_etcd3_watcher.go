package main

import (
	"os"
	"tp_go_etcd3_watcher/kafkasender"
	"tp_go_etcd3_watcher/kafkawatcher"

	log "github.com/Golang-Tools/loggerhelper"
	s "github.com/Golang-Tools/schema-entry-go"
)

func main() {
	root, err := s.New(&s.EntryPointMeta{Name: "tp_go_etcd3_watcher", Usage: "tp_go_etcd3_watcher"})
	if err != nil {
		log.Error("create root node error", log.Dict{"err": err})
		os.Exit(1)
	}
	s.RegistSubNode(root, kafkasender.App)
	s.RegistSubNode(root, kafkawatcher.App)
	root.Parse(os.Args)
}