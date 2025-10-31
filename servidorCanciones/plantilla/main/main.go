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

	router.HandleFunc("/canciones/almacenamiento", ctrl.AlmacenarAudioCancion).Methods("POST")
	router.HandleFunc("/generos", controlador.NuevoControladorGeneros().ObtenerGeneros).Methods("GET")
	router.HandleFunc("/generos/{genero}/canciones", controlador.NuevoControladorAlmacenamientoCanciones().ObtenerCancionesPorGenero).Methods("GET")
	router.HandleFunc("/canciones", controlador.NuevoControladorAlmacenamientoCanciones().ObtenerTodasLasCanciones).Methods("GET")
	// --- NUEVA RUTA PARA SERVIR ARCHIVOS DE AUDIO ---
	// Usamos un prefijo y un FileServer para servir de forma segura los archivos del directorio 'audios'.
	// http.StripPrefix quita "/audio/" de la URL para que FileServer busque desde la raíz de 'audios'.
	fs := http.FileServer(http.Dir("./audios/"))
	router.PathPrefix("/audio/").Handler(http.StripPrefix("/audio/", fs))

	fmt.Println("✅ Servicio de Tendencias escuchando en el puerto 5000...")
	if err := http.ListenAndServe(":5000", router); err != nil {
		fmt.Println("Error iniciando el servidor:", err)
	}
}
