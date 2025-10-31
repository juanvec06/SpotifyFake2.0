package cancionConsumer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url" // Importado para construir URLs de forma segura

	"proyecto.local/cliente/models"
)

const (
	baseURL = "http://localhost:5000" // Base de la URL del servidor de canciones
)

// CancionConsumer se encarga de la lógica de obtención de canciones y géneros vía HTTP.
type CancionConsumer struct {
	// Ya no necesita un cliente gRPC.
	// Podríamos agregar aquí un http.Client para configuraciones avanzadas,
	// pero para este caso no es necesario.
}

// NewCancionConsumer crea una nueva instancia de CancionConsumer.
func NewCancionConsumer() *CancionConsumer {
	// El constructor ya no recibe el cliente gRPC.
	return &CancionConsumer{}
}

// GetAllSongs obtiene todas las canciones disponibles desde el endpoint REST.
func (cc *CancionConsumer) GetAllSongs() ([]models.Song, error) {
	resp, err := http.Get(baseURL + "/canciones")
	if err != nil {
		return nil, fmt.Errorf("error al realizar la petición a /canciones: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("el servidor respondió con un estado inesperado: %s", resp.Status)
	}

	var songs []models.Song
	if err := json.NewDecoder(resp.Body).Decode(&songs); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta de canciones: %w", err)
	}

	return songs, nil
}

// GetGenres obtiene todos los géneros disponibles desde el endpoint REST.
func (cc *CancionConsumer) GetGenres() ([]models.Genre, error) {
	resp, err := http.Get(baseURL + "/generos")
	if err != nil {
		return nil, fmt.Errorf("error al realizar la petición a /generos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("el servidor respondió con un estado inesperado: %s", resp.Status)
	}

	// La API devuelve un arreglo de strings: ["rock", "pop", ...]
	var genreNames []string
	if err := json.NewDecoder(resp.Body).Decode(&genreNames); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta de géneros: %w", err)
	}

	// Convertimos el arreglo de strings a un arreglo de models.Genre
	genres := make([]models.Genre, len(genreNames))
	for i, name := range genreNames {
		genres[i] = models.Genre{Name: name}
	}
	return genres, nil
}

// GetSongsByGenre obtiene todas las canciones de un género específico.
// Nota que ahora recibe un string (nombre del género) en lugar de un ID.
func (cc *CancionConsumer) GetSongsByGenre(genreName string) ([]models.Song, error) {
	// Construimos la URL de forma segura para manejar caracteres especiales en el nombre del género
	endpoint := fmt.Sprintf("%s/generos/%s/canciones", baseURL, url.PathEscape(genreName))

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error al realizar la petición de canciones por género: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("el servidor respondió con un estado inesperado: %s", resp.Status)
	}

	var songs []models.Song
	if err := json.NewDecoder(resp.Body).Decode(&songs); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta de canciones: %w", err)
	}

	return songs, nil
}
