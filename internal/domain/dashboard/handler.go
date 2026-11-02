package dashboard

import "net/http"

func Overview(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("overview"))
}
