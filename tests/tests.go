package tests

func newResponseBody(msg string, data string) string {
	if msg != "" && data != "" {
		return `{"message":"` + msg + `","data":"` + data + `"}`
	}
	if msg != "" && data == "" {
		return `{"message":"` + msg + `"}`
	}
	if msg == "" && data != "" {
		return `{"data":"` + data + `"}`
	}
	return "{}"
}
