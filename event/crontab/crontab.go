package crontab

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"lms-be/configs"

	"github.com/robfig/cron/v3"
)

type CrontabService interface {
	Start(config *configs.Config)
}

type CrontabServiceImpl struct {
	// SchedulePOService transaction.SchedulePOService
}

func ProvideCrontabServiceImpl() *CrontabServiceImpl {
	// return &CrontabServiceImpl{SchedulePOService: service}
	return &CrontabServiceImpl{}
}

// Start crontab scheduler
func (c *CrontabServiceImpl) Start(config *configs.Config) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(loc))
	defer scheduler.Stop()

	// scheduler.AddFunc("*/1 * * * *", func() {
	scheduler.AddFunc("0 8 * * *", func() {
		c.ReminderSchedulePo()
	})

	go scheduler.Start()

	// tunggu signal stop
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

// function bebas, tidak perlu receiver
func (c *CrontabServiceImpl) ReminderSchedulePo() {
	// ctx := context.Background()

	// _, err := c.SchedulePOService.DuedateSchedulePo(ctx)
	// if err != nil {
	// 	fmt.Println("error ReminderSchedulePo:", err)
	// 	return
	// }
	return
}
