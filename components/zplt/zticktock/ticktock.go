package zticktock

import (
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/ticktock"
	"github.com/hootuu/hyle/hcfg"
)

const (
	defTicktockWorker  = "TICKTOCK_WORKER_D"
	defTicktockPostman = "TICKTOCK_POSTMAN_D"
)

var gDefTicktockWorker *ticktock.Worker

func Ticktock() *ticktock.Worker {
	tickTockRunning := hcfg.GetBool("helix.ticktock.running", true)
	if !tickTockRunning {
		return nil
	}
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
