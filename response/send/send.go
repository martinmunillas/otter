package send

import "github.com/martinmunillas/otter"

type errorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var Html = htmlSender{
	errorComponent: otter.ErrorAlert,
	logInternals:   true,
}
var Json = jsonSender{
	logInternals: true,
}
