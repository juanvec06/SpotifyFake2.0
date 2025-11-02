package main

import (
	"fmt"
	"net/http"
	controlador "tendencias/capaControladores"
)

func main() {
	ctrl := controlador.NuevaControladorPreferencias()

	// Endpoints para Servidor de Streaming
	http.HandleFunc("/reproducciones/registrar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		ctrl.RegistrarReproduccionHandler(w, r)
	})

	// Endpoints para Servidor de Preferencias Java (compatibilidad)
	http.HandleFunc("/reproducciones", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		ctrl.ListarReproduccionesPorUsuarioHandler(w, r)
	})

	// Endpoints originales (retrocompatibilidad)
	http.HandleFunc("/tendencias/reproduccion", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		ctrl.RegistrarReproduccionHandler(w, r)
	})

	http.HandleFunc("/tendencias/listarReproducciones", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		ctrl.ListarReproduccionesHandler(w, r)
	})

	fmt.Println("✅ Servidor de Reproducciones escuchando en el puerto 3000...")

	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Printf("❌ Error al iniciar el servidor: %v\n", err)
	}
}
