package capafachada

import (
	capaaccesoadatos "almacenamiento/capaAccesoADatos"
	dtos "almacenamiento/capaFachadaServices/DTOs"
	componnteconexioncola "almacenamiento/componnteConexionCola"
	"fmt"
)

type FachadaAlmacenamiento struct {
	repo         *capaaccesoadatos.RepositorioCanciones
	conexionCola *componnteconexioncola.RabbitPublisher
}

// Constructor de la fachada
func NuevaFachadaAlmacenamiento() *FachadaAlmacenamiento {
	fmt.Println("🔧 Inicializando fachada de almacenamiento...")

	repo := capaaccesoadatos.GetRepositorioCanciones()

	conexionCola, err := componnteconexioncola.NewRabbitPublisher()
	if err != nil {
		fmt.Println(" Error al conectar con RabbitMQ:", err)
		conexionCola = nil
	}

	return &FachadaAlmacenamiento{
		repo:         repo,
		conexionCola: conexionCola,
	}
}

func (thisF *FachadaAlmacenamiento) GetSongsByGenreService(genero string) ([]dtos.CancionDTOOutput, error) {
	// Aquí se llamaría al repositorio para obtener las canciones por género
	// Por simplicidad, devolvemos una lista vacía
	var canciones []dtos.CancionDTOOutput
	return canciones, nil
}
func (thisF *FachadaAlmacenamiento) GuardarCancion(objCancion dtos.CancionAlmacenarDTOInput, data []byte) error {
	thisF.conexionCola.PublicarNotificacion(componnteconexioncola.NotificacionCancion{
		Titulo:  objCancion.Titulo,
		Artista: objCancion.Artista,
		Genero:  objCancion.Genero,
		Mensaje: "Nueva canción almacenada: " + objCancion.Titulo + " de " + objCancion.Artista,
	})

	return thisF.repo.GuardarCancion(objCancion.Titulo, objCancion.Genero, objCancion.Artista, data)
}
