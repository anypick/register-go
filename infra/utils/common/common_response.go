package common

const (
	SuccessCode = 2000
	FailCode    = 5000
)

type ResponseData struct {
	Code    int
	Msg     string `json:",omitempty"`
	Success bool
	Total   int
	Rows    interface{}
}

func NewRespSucc() ResponseData {
	return ResponseData{Code: SuccessCode, Success: true}
}

func NewRespSuccWithData(rows interface{},  total int) ResponseData {
	return ResponseData{Code: SuccessCode, Success: true, Total: total, Rows: rows}
}

func NewRespSuccWithMsg(msg string) ResponseData {
	return ResponseData{Code: SuccessCode, Msg: msg, Success: true, Total: 0}
}

func NewRespFail() ResponseData {
	return ResponseData{Code: FailCode, Success: false}
}

func NewRespFailWithMsg(msg string) ResponseData {
	return ResponseData{Code: FailCode, Success: false, Msg: msg}
}
