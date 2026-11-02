package dashboard

import (
	"net/http"

	"sorantics.com/tasks/internal/lib"
)

func Overview(registry lib.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registry.RenderWithLayout(w, "dashboard/overview.html", nil, "layout/dashboard.html")
	}
}
