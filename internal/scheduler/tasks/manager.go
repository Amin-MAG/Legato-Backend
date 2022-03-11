package tasks

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/redisq"
	"legato_server/config"
	"os"
	"os/signal"
	"syscall"
)

const queueName = "legato_scheduler"
const AccessToken = "Tqqc5He3/;cD,-z/RD'q~QH/BMaXCyw}9pT+D?*7"

type LegatoTaskQueue struct {
	QueueFactory     taskq.Factory
	MainQueue        taskq.Queue
	Tasks            map[string]*taskq.Task
	LegatoServerAddr string
}

func NewLegatoTaskQueue(r *redis.Client, cfg config.Config) (LegatoTaskQueue, error) {
	if r == nil {
		return LegatoTaskQueue{}, errors.New("redis client is empty")
	}

	var tk LegatoTaskQueue
	tk.LegatoServerAddr = fmt.Sprintf("http://%s:%s", cfg.Legato.Host, cfg.Legato.ServingPort)
	tk.QueueFactory = redisq.NewFactory()
	tk.MainQueue = tk.QueueFactory.RegisterQueue(&taskq.QueueOptions{
		Name:  queueName,
		Redis: r,
	})
	tk.Tasks = map[string]*taskq.Task{
		StartScenarioTask: taskq.RegisterTask(&taskq.TaskOptions{
			Name:    StartScenarioTask,
			Handler: startScenario,
		}),
	}

	return tk, nil
}

func (ltq *LegatoTaskQueue) WaitSignal() os.Signal {
	ch := make(chan os.Signal, 2)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			return sig
		}
	}
}

func (ltq *LegatoTaskQueue) Listen() {
	c := context.Background()

	err := ltq.QueueFactory.StartConsumers(c)
	if err != nil {
		panic(err)
	}

	_ = ltq.WaitSignal()

	err = ltq.QueueFactory.Close()
	if err != nil {
		panic(err)
	}
}
