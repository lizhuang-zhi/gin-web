package meta

import "booking-app/micro-service/core/option"

const (
	App       = "App"
	Name      = "Name"
	Version   = "Version"
	BuildTime = "BuildTime"
	Useage    = "Useage"
)

var Data = option.NewMetaData(App)
