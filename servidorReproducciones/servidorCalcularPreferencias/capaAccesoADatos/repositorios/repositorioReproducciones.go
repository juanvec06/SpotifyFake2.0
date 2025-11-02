package repositorios

import (
	"fmt"
	"strconv"
	"sync"
	. "tendencias/capaAccesoADatos/entitys"
	dtos "tendencias/capaFachadaServices/DTOs"
	"time"
)

type RepositorioReproducciones struct {
	mu             sync.Mutex
	reproducciones []ReproduccionEntity
	contadorID     int // Contador para IDs únicos
}

// insancia unica del repositorio (patrón singleton)
var (
	instancia *RepositorioReproducciones
	once      sync.Once
)

// Crear o devolver la unica instancia del repositorio
func GetRepositorio() *RepositorioReproducciones {
	once.Do(func() {
		instancia = &RepositorioReproducciones{
			contadorID: 1,
		}
	})
	return instancia
}

func (r *RepositorioReproducciones) AgregarReproduccion(titulo, usuarioID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Convertir usuarioID string a int
	idUsuario, err := strconv.Atoi(usuarioID)
	if err != nil {
		// Si falla la conversión, intentar asignar un ID por defecto
		fmt.Printf("⚠️ Error al convertir usuarioID '%s' a entero: %v. Usando ID 0.\n", usuarioID, err)
		idUsuario = 0
	}

	// Mock para idCancion: usar longitud del título mod 10 + 1
	// En una implementación real, esto vendría del servidor de canciones
	idCancion := len(titulo)%10 + 1

	reproduccion := ReproduccionEntity{
		ID:        r.contadorID,
		Titulo:    titulo,
		IDUsuario: idUsuario,
		IDCancion: idCancion,
		FechaHora: time.Now().Format("2006-01-02 15:04:05"),
	}
	r.reproducciones = append(r.reproducciones, reproduccion)
	r.contadorID++
	fmt.Printf("✅ Reproducción almacenada: Usuario=%d, Canción=%d, Título=%s\n",
		reproduccion.IDUsuario, reproduccion.IDCancion, reproduccion.Titulo)
	r.mostrarReproducciones()
}

func (r *RepositorioReproducciones) ListarReproducciones() []ReproduccionEntity {
	return r.reproducciones
}

// Método para el servidor de preferencias Java
// Retorna reproducciones en el formato esperado por Java
func (r *RepositorioReproducciones) ListarReproduccionesPorUsuario(idUsuario int) []dtos.ReproduccionDTOOutput {
	r.mu.Lock()
	defer r.mu.Unlock()

	var resultado []dtos.ReproduccionDTOOutput

	for _, repro := range r.reproducciones {
		// Filtrar por IDUsuario
		if repro.IDUsuario == idUsuario {
			resultado = append(resultado, dtos.ReproduccionDTOOutput{
				ID:        repro.ID,
				IDUsuario: repro.IDUsuario,
				Titulo:    repro.Titulo, // Enviar título en lugar de IDCancion
			})
		}
	}

	fmt.Printf("📊 Reproducciones encontradas para usuario %d: %d\n", idUsuario, len(resultado))
	return resultado
}

func (r *RepositorioReproducciones) mostrarReproducciones() {
	fmt.Println("==Reproducciones almacenadas==")
	for i := 0; i < len(r.reproducciones); i++ {
		fmt.Printf(" ID: %d\n", r.reproducciones[i].ID)
		fmt.Printf(" IDUsuario: %d\n", r.reproducciones[i].IDUsuario)
		fmt.Printf(" IDCancion: %d\n", r.reproducciones[i].IDCancion)
		fmt.Printf(" Título: %s\n", r.reproducciones[i].Titulo)
		fmt.Printf(" Fecha y hora: %s\n", r.reproducciones[i].FechaHora)
		fmt.Println("---")
	}
}
