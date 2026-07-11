package routes

import "net/http"

func NewRouter() *http.ServeMux {

	mux := http.NewServeMux()

	RegisterRoutes(mux)

	return mux
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("API está funcionando normalmente!"))
	})
}
