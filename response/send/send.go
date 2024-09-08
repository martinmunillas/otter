package send

import "log/slog"

type errorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var Html = htmlSender{
	logger: slog.Default(),
}
var Json = jsonSender{
	logger: slog.Default(),
}
