package scheduler

import (
	"reflect"
	"runtime"
	"sync"

	"github.com/jasonlvhit/gocron"
	"github.com/rs/xlog"
)

type Scheduler struct {
	Scheduler *gocron.Scheduler
	Threads   int
}

var s *Scheduler
var once sync.Once

// interval in seconds
func (sc *Scheduler) Run(fn func(), interval uint64) {
	sc.Scheduler.Every(interval).Seconds().Do(fn)

	sc.Threads++

	xlog.Infof("Scheduler %#v running...", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
}

func (sc *Scheduler) Stop() {
	if sc.Threads > 0 {
		sc.Threads--
	}
	sc.Scheduler.Clear()
	xlog.Info("Stopped all")
}

func (sc *Scheduler) Start() {
	sc.Scheduler.Start()
}

func SchedulerProvider() *Scheduler {

	once.Do(func() {
		s = &Scheduler{gocron.NewScheduler(), 0}
	})

	return s
}
