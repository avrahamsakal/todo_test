package main

import (
	"github.com/jasonlvhit/gocron"
)

type cron struct {}

func (c *cron) startJobs() {
	go func(){<-gocron.Start()}()
}

func (c *cron) sampleCronJob() {
	
}