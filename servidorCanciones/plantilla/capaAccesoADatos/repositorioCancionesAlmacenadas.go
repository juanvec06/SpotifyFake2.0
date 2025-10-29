package capaaccesoadatos

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type RepositorioCanciones struct {
	mu sync.Mutex
}

var (
	intancia *RepositorioCanciones
	once     sync.Once
)

// GetRepositorioCanciones aplica patrón Singleton
func GetRepositorioCanciones() *RepositorioCanciones {
	once.Do(func() {
		intancia = &RepositorioCanciones{}
	})
	return intancia
}

func (r *RepositorioCanciones) GuardarCancion(titulo string, genero string, artista string, data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Crear el directorio si no existe
	os.MkdirAll("audios", os.ModePerm)

	// Construir nombre del archivo
	fileName := fmt.Sprintf("%s_%s_%s.mp3", titulo, genero, artista)
	filePath := filepath.Join("audios", fileName)

	// Crear y escribir el archivo
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error guardando la canción: %v", err)
	}

	return nil
}
