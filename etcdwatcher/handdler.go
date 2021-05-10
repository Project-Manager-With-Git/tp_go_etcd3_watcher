package etcdwatcher

import (
	log "github.com/Golang-Tools/loggerhelper"
)

func (w ConsumerProxy) handdler() error {
	for {
		msg, err := w.ReadMessage(-1)
		if err == nil {
			log.Info("get msg", log.Dict{"msg": msg})
		} else {
			// The client will automatically try to recover from all errors.
			log.Error("kafka watcher get error", log.Dict{"err": err, "msg": msg})
			return err
		}
	}
}
