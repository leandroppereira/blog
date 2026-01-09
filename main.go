package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var ready int32 // 0 = not ready, 1 = ready

func main() {
	// Configuráveis por env, mas com defaults bons pra laboratório
	readyAfter := getenvDuration("READY_AFTER", 25*time.Second) // readiness só fica OK depois disso
	dieAfter := getenvDuration("DIE_AFTER", 0)                  // 0 = não morre (opcional)

	start := time.Now()
	log.Printf("starting, ready after=%s, die after=%s", readyAfter, dieAfter)

	// Marca ready depois de um tempo
	go func() {
		time.Sleep(readyAfter)
		atomic.StoreInt32(&ready, 1)
		log.Printf("READY at %s", time.Since(start).Truncate(time.Millisecond))
	}()

	// Opcional: simula travar/morrer (pra testar liveness)
	if dieAfter > 0 {
		go func() {
			time.Sleep(dieAfter)
			log.Printf("Simulating death now (exit 1)")
			os.Exit(1)
		}()
	}

	mux := http.NewServeMux()

	// Endpoint principal: usado pelo readinessProbe (GET /)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&ready) == 1 {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "ok - ready")
			return
		}
		// Enquanto não pronto, devolve 503 (readiness falha)
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "not ready yet")
	})

	// Extra útil: healthz sempre OK (não é exigido pela questão, mas ajuda)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "alive")
	})

	addr := ":8080"
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}

func getenvDuration(name string, def time.Duration) time.Duration {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("invalid %s=%q, using default %s", name, v, def)
		return def
	}
	return d
}

