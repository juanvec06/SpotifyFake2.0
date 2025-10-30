package controlador

import (
	dtos "almacenamiento/capaFachadaServices/DTOs"
	capafachada "almacenamiento/capaFachadaServices/fachada"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ControladorAlmacenamientoCanciones struct {
	fachada *capafachada.FachadaAlmacenamiento
}

// Constructor del Controlador
func NuevoControladorAlmacenamientoCanciones() *ControladorAlmacenamientoCanciones {
	return &ControladorAlmacenamientoCanciones{
		fachada: capafachada.NuevaFachadaAlmacenamiento(),
	}
}

func (thisC *ControladorAlmacenamientoCanciones) ObtenerCancionesPorGenero(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Obteniendo canciones por género...\n")
	// gorilla/mux extrae las variables de la URL por nosotros.
	vars := mux.Vars(r)
	genero := vars["genero"] // La clave "genero" debe coincidir con la definida en la ruta.

	fmt.Printf("Consultando canciones para el género: %s...\n", genero)

	canciones, err := thisC.fachada.GetSongsByGenreService(genero)
	if err != nil {
		http.Error(w, "Error al obtener las canciones por género", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(canciones)
}
func (thisC *ControladorAlmacenamientoCanciones) AlmacenarAudioCancion(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Almacenando canción...\n")
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(50 << 20)
	file, _, err := r.FormFile("archivo")
	if err != nil {
		http.Error(w, "Error leyendo el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()
	data, _ := io.ReadAll(file)

	releaseYearStr := r.FormValue("release_year")
	releaseYear, err := strconv.Atoi(releaseYearStr)
	if err != nil {
		http.Error(w, "Release year inválido", http.StatusBadRequest)
		return
	}

	// Leer los campos del DTO
	dto := dtos.CancionAlmacenarDTOInput{
		Titulo:      r.FormValue("titulo"),
		Genero:      r.FormValue("genero"),
		Artista:     r.FormValue("artista"),
		Album:       r.FormValue("album"),
		ReleaseYear: releaseYear,
	}
	thisC.fachada.GuardarCancion(dto, data)
}
