package controlador

import (
	capafachada "almacenamiento/capaFachadaServices/fachada"
	"encoding/json"
	"fmt"
	"net/http"
)

type ControladorGeneros struct {
	fachada *capafachada.FachadaGeneros
}

// Constructor del Controlador
func NuevoControladorGeneros() *ControladorGeneros {
	return &ControladorGeneros{
		fachada: capafachada.NuevaFachadaGeneros(),
	}
}

func (thisC *ControladorGeneros) ObtenerGeneros(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Obteniendo géneros...\n")
	generos, err := thisC.fachada.GetGenresService()
	if err != nil {
		http.Error(w, "Error al obtener los géneros", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(generos)
}
