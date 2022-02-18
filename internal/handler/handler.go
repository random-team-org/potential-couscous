package handler

import (
	"devices/internal/server"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type AllocateDeviceRequest struct {
	DeviceID string `json:"device_id"`
	UserID   string `json:"user_id"`
}

func New(addr, dbopts, memcached string) {
	s := server.NewService(dbopts, memcached)

	r := mux.NewRouter()
	r.HandleFunc("/device", handleGetDevices(s)).Methods("GET")
	r.HandleFunc("/device/alloc", handleAllocateDevice(s)).Methods("POST")
	r.HandleFunc("/device/alloc", handleFreeDevice(s)).Methods("DELETE")

	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func handleGetDevices(s *server.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		devices, err := s.GetDevices()

		if err == nil {
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(devices)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func handleAllocateDevice(s *server.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		allocateRequest := AllocateDeviceRequest{}
		_ = json.NewDecoder(r.Body).Decode(&allocateRequest)

		err := s.AllocateDevice(allocateRequest.DeviceID, allocateRequest.UserID)
		if err == nil {
			rw.WriteHeader(http.StatusOK)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func handleFreeDevice(s *server.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		allocateRequest := AllocateDeviceRequest{}
		_ = json.NewDecoder(r.Body).Decode(&allocateRequest)

		err := s.FreeDevice(allocateRequest.DeviceID)
		if err == nil {
			rw.WriteHeader(http.StatusOK)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}
