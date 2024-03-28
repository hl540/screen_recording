package server

import (
	"encoding/json"
	"net/http"

	"screen_recording/internal/channel"
)

func getChannelHandler(w http.ResponseWriter, r *http.Request) {
	names := channel.GetAllName()
	jsonBytes, _ := json.Marshal(names)
	w.Write(jsonBytes)
}
