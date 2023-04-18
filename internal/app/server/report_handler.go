package server

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"screen_recording/internal/util"
)

func reportHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.URL.Query().Get("key"))
	// log.Println(r.URL.Query().Get("channel"))
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(len(data))
	// 解压
	data, err = util.Ugzip(data)
	log.Println(err)

	// 保存图片
	// dec := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(data))
	// bs, _ := io.ReadAll(dec)
	// log.Println(len(data))
	ioutil.WriteFile("./output.jpeg", data, 0777)
}
