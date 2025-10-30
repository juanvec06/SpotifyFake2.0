package main

import (
	controlador "almacenamiento/capaControladores"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	ctrl := controlador.NuevoControladorAlmacenamientoCanciones()

	router := mux.NewRouter()

	router.HandleFunc("/canciones/almacenamiento", ctrl.AlmacenarAudioCancion)
	router.HandleFunc("/generos", controlador.NuevoControladorAlmacenamientoCanciones().ObtenerGeneros)
	router.HandleFunc("/generos/{genero}/canciones", controlador.NuevoControladorAlmacenamientoCanciones().ObtenerCancionesPorGenero)

	fmt.Println("✅ Servicio de Tendencias escuchando en el puerto 5000...")
	if err := http.ListenAndServe(":5000", router); err != nil {
		fmt.Println("Error iniciando el servidor:", err)
	}
}
