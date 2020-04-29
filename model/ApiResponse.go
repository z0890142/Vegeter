package model

type ApiResponse struct {
	ResultCode    string
	ResultMessage interface{}
	LogSN         string
}
type ApiResponseWithTime struct {
	ResultCode    string
	ResultMessage string
	ReceiveTime   string
	ResponseTime  string
	LogSN         string
}

type VoiceApiResponse struct {
	Time          string
	RequestId     string
	ResultCode    string
	ResultMessage string
}
