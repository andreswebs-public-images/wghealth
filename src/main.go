package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func init() {
	log.SetFlags(0)
}

func main() {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	go handleExitSignal(exitSignal)

	defaultDevice := "wg0"
	defaultPort := 9000

	envDevice, ok := os.LookupEnv("WGHEALTH_DEVICE")
	if ok {
		defaultDevice = envDevice
	}

	envPort, ok := os.LookupEnv("WGHEALTH_PORT")
	if ok {
		if p, err := strconv.Atoi(envPort); err != nil {
			defaultPort = p
		}
	}

	device := flag.String("device", defaultDevice, "device name to check")
	port := flag.Int("port", defaultPort, "port on which to listen")
	isTest := flag.Bool("test", false, "show status and exit")
	flag.Parse()

	if *isTest {
		test(*device)
	} else {
		http.HandleFunc("/", healthCheckHandler(*device))
		addr := fmt.Sprintf(":%d", *port)
		log.Printf("[wghealth] listening on %s\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Printf("[wghealth] error starting server: %v\n", err)
		}
	}
}

func check(device string) bool {
	interfaces, err := net.Interfaces()
	if err != nil {
		return false
	}

	for _, iface := range interfaces {
		if iface.Name == device {
			return iface.Flags&net.FlagUp != 0
		}
	}

	return false
}

func test(device string) {
	if check(device) {
		log.Printf("%s up\n", device)
	} else {
		log.Printf("%s down\n", device)
	}
}

func healthCheckHandler(device string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if check(device) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("OK\n"))
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}
}

func handleExitSignal(exitSignal chan os.Signal) {
	<-exitSignal
	fmt.Println()
	log.Println("[wghealth] received termination signal - shutting down")
	os.Exit(0)
}
