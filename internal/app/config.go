package app

import "time"

type Configuration struct {
	/*
		https://en.wikipedia.org/wiki/Deployment_environment
		local -> development/trunk -> integration -> testing/test/qc/internal acceptnace -> staging/stage/model/pre-production/demo -> production/live
		development -> test -> pre-production -> production
	*/
	Env          Env
	TimezoneID   string // e.g. "Asia/Shanghai"
	TimeLocation *time.Location
	TimeFormat   string // e.g. "2006-02-01 15:04:05"
	Mock         bool   // using mock
}
