package model

type ReportReq struct {
	Key     string `form:key`
	Channel string `form:channelkey`
	Data    string `form:data`
}
