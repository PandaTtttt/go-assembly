package zlog

import (
	"github.com/IBM/sarama"
	"github.com/PandaTtttt/go-assembly/errs"
	"sync"
)

type KafkaWriter struct {
	Address []string
	Topic   string
	asyncP  sarama.AsyncProducer
	once    sync.Once
}

func (w *KafkaWriter) callback(logger *Logger) error {
	config := sarama.NewConfig()
	client, err := sarama.NewClient(w.Address, config)
	if err != nil {
		return errs.Internal.Wrap(err)
	}
	w.asyncP, err = sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return errs.Internal.Wrap(err)
	}
	go func() {
		for err := range w.asyncP.Errors() {
			logger.Info(err.Error())
		}
	}()
	return nil
}

func (w *KafkaWriter) Write(b []byte) (n int, err error) {
	w.asyncP.Input() <- &sarama.ProducerMessage{
		Topic: w.Topic,
		Value: sarama.StringEncoder(b),
	}
	return
}
