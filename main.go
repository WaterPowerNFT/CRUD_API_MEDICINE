package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Date struct {
	Day   uint `json:"day"`
	Month uint `json:"month"`
	Year  uint `json:"year"`
}

type Fabric struct {
	FabricCountry   string `json:"fabric_country"`
	FabricName      string `json:"fabric_name"`
	FabricBuiltDate Date   `json:"fabric_built_date"`
}

type Medicine struct {
	ID             string  `json:"id"`
	ExpirationDate Date    `json:"expiration_date"`
	ProductionDate Date    `json:"production_date"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	FabricData     *Fabric `json:"fabric"`
}

var medicines []Medicine

func GetMeds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medicines)
}

func DeleteMed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range medicines {
		if item.ID == params["id"] {
			medicines = append(medicines[:index], medicines[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(medicines)
}

func GetOneMed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range medicines {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
		}
	}
}

func CreateMed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var med Medicine
	_ = json.NewDecoder(r.Body).Decode(&med)
	med.ID = strconv.Itoa(rand.Intn(100000000))
	medicines = append(medicines, med)
	json.NewEncoder(w).Encode(medicines)
}

func UpdateMed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range medicines {
		if item.ID == params["ID"] {
			medicines = append(medicines[:index], medicines[index+1:]...)
			var med Medicine
			_ = json.NewDecoder(r.Body).Decode(&med)
			med.ID = params["ID"]
			json.NewEncoder(w).Encode(medicines)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	medicines = append(medicines, Medicine{ID: "1", ExpirationDate: Date{Day: 1, Month: 1, Year: 2011},
		ProductionDate: Date{Day: 2, Month: 2, Year: 2010}, Name: "Paracetamol", Description: "nothing",
		FabricData: &Fabric{FabricCountry: "RU", FabricName: "Leningrad zavod", FabricBuiltDate: Date{Day: 12, Month: 1, Year: 1958}}})

	r.HandleFunc("/medicines", GetMeds).Methods("GET")
	r.HandleFunc("/medicine/{id}", GetOneMed).Methods("GET")
	r.HandleFunc("/medicines", CreateMed).Methods("POST")
	r.HandleFunc("medicines/{id}", UpdateMed).Methods("PUT")
	r.HandleFunc("medicines/{id}", DeleteMed).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
