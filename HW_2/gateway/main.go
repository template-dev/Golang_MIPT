package gateway

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	port := ":8080"
	fmt.Printf("Gateway service started on port %s\n", port)
	fmt.Printf("Test with: curl http://localhost%s/ping\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
