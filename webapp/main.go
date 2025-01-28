package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", rawlog)
	http.ListenAndServe(":8000", nil)
}

type WafService struct{}

func (ws *WafService) Check(r []byte) bool {
	if ws.shouldInlineCheck(r) {
		ws.checkInline(r)
		return true
	}
	return false
}
func (ws *WafService) shouldInlineCheck([]byte) bool {
	return false
}
func (ws *WafService) checkInline([]byte) bool {
	return false
}

func rawlog(w http.ResponseWriter, req *http.Request) {
	//ws := WafService{}

	var buf bytes.Buffer
	req.Write(&buf)
	fmt.Println(req.RemoteAddr)
	fmt.Println(req.Cookies())
	fmt.Println(string(buf.Bytes()))

	return
}
