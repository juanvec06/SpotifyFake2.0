package capacontroladores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	dtos "tendencias/capaFachadaServices/DTOs"
	capafachada "tendencias/capaFachadaServices/fachada"
)

type ControladorPreferencias struct {
	fachada *capafachada.FachadaTendencias
}

func NuevaControladorPreferencias() *ControladorPreferencias {
	return &ControladorPreferencias{
		fachada: capafachada.NuevaFachadaTendencias(),
	}
}

func (c *ControladorPreferencias) RegistrarReproduccionHandler(w http.ResponseWriter, r *http.Request) {
	var dto dtos.ReproduccionDTOInput

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Usar el método que obtiene el usuario ID de cualquier campo
	usuarioID := dto.GetUsuarioID()
	fmt.Println("UsuarioID obtenido:", usuarioID)
	c.fachada.RegistrarReproduccion(dto.Titulo, usuarioID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Reproducción registrada con éxito")
}

func (c *ControladorPreferencias) ListarReproduccionesHandler(w http.ResponseWriter, r *http.Request) {
	repros := c.fachada.ObtenerReproducciones()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repros)
}

// Nuevo handler para el servidor de preferencias Java
// GET /reproducciones?idUsuario=1
func (c *ControladorPreferencias) ListarReproduccionesPorUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el parámetro idUsuario de la query string
	idUsuarioStr := r.URL.Query().Get("idUsuario")

	if idUsuarioStr == "" {
		http.Error(w, "Parámetro 'idUsuario' requerido", http.StatusBadRequest)
		return
	}

	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		http.Error(w, "Parámetro 'idUsuario' debe ser un número", http.StatusBadRequest)
		return
	}

	// Obtener reproducciones filtradas por usuario
	repros := c.fachada.ObtenerReproduccionesPorUsuario(idUsuario)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repros)
}
