package zticktock

import (
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/ticktock"
)

const (
	defTicktockWorker  = "TICKTOCK_WORKER_D"
	defTicktockPostman = "TICKTOCK_POSTMAN_D"
)

var gDefTicktockWorker *ticktock.Worker

func Ticktock() *ticktock.Worker {
	if gDefTicktockWorker == nil {
		helix.OnceLoad(defTicktockWorker, func() {
			gDefTicktockWorker = ticktock.NewWorker(defTicktockWorker, zplt.HelixRdsCache())
		})
	}
	return gDefTicktockWorker
}

var gDefTicktockPostman *ticktock.Postman

func Postman() *ticktock.Postman {
	if gDefTicktockPostman == nil {
		helix.OnceLoad(defTicktockPostman, func() {
			gDefTicktockPostman = ticktock.NewPostman(defTicktockPostman, zplt.HelixRdsCache())
		})
	}
	return gDefTicktockPostman
}
