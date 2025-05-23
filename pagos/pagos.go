package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Asi es la estructura de los pagos basado en el json de job_bank del microservicio de simulador
type Pago struct {
	SSN             string  `json:"ssn"` // Numero de seguro social del cliente
	Amount          float64 `json:"amount"` // Monto del pago
	TransactionType string  `json:"transaction_type"` // Tipo de transacción
	Status          string  `json:"status"` // Estado de la transacción
	CreationDate    string  `json:"creation_date"` // Fecha de creación del pago
	TransactionID   string  `json:"transaction_id"` // ID de la transacción
	Source          string  `json:"source"` // Fuente del pago
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", Index)
	router.HandleFunc("/records", GetPayments)
	log.Fatal(http.ListenAndServe(":8003", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func GetPayments(w http.ResponseWriter, r *http.Request) {
	// Se obtiene el archivo en la ruta donde se genera el json segun el job_bank.py
	contenidoArchivo, e := ioutil.ReadFile("/data/payment_records.json")
	if e != nil {
		log.Printf("No se pudo leer el archivo: %v", e)
		http.Error(w, "Error leyendo el archivo json", http.StatusInternalServerError)
		return
	}

	// Como no se tiene forma de json correcta se mapea el json a un mapa
	var datosJSON map[string][]Pago
	if e := json.Unmarshal(contenidoArchivo, &datosJSON); e != nil {
		log.Printf("Error de mapeo de los datos JSON: %v", e)
		http.Error(w, "No se pudo dar formato al archivo de los pagos de clientes", http.StatusInternalServerError)
		return
	}

	// Como se tiene el numero de seguro social aparte, se añade como un atributo al pago
	for ssn, payments := range datosJSON {
		for iteradorPagos := range payments {
			payments[iteradorPagos].SSN = ssn
		}
		datosJSON[ssn] = payments
	}

	// Una vez procesados los datos se envia todo el mapa formateado de nuevo como un archivo con formato JSON con los pagos en orden
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datosJSON)
}
