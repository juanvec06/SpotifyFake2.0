package capaaccesoadatos

import (
	dtos "almacenamiento/capaFachadaServices/DTOs"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Constante para la ruta del archivo de metadatos
const metadataPath = "audios/metadata.json"

type RepositorioCanciones struct {
	mu        sync.Mutex
	canciones []dtos.CancionDTOOutput // Almacenamos los metadatos en memoria para un acceso rápido
}

var (
	instancia *RepositorioCanciones
	once      sync.Once
)

// GetRepositorioCanciones aplica patrón Singleton y carga los datos al iniciar.
func GetRepositorioCanciones() *RepositorioCanciones {
	once.Do(func() {
		instancia = &RepositorioCanciones{}
		// Al crear la instancia por primera vez, intentamos cargar los metadatos existentes.
		if err := instancia.cargarMetadata(); err != nil {
			fmt.Printf("Advertencia: no se pudo cargar metadata.json: %v. Se creará uno nuevo.\n", err)
		}
	})
	return instancia
}
func (r *RepositorioCanciones) GetAllSongs() ([]dtos.CancionDTOOutput, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Devolvemos una copia del slice para evitar modificaciones externas.
	cancionesCopia := make([]dtos.CancionDTOOutput, len(r.canciones))
	copy(cancionesCopia, r.canciones)
	return cancionesCopia, nil
}

// cargarMetadata lee el archivo JSON y lo vuelca en la estructura del repositorio.
func (r *RepositorioCanciones) cargarMetadata() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.Open(metadataPath)
	if err != nil {
		// Si el archivo no existe, está bien, lo crearemos al guardar la primera canción.
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	// Si el archivo está vacío, no hacemos nada.
	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, &r.canciones)
}

// guardarMetadata escribe el estado actual de los metadatos al archivo JSON.
func (r *RepositorioCanciones) guardarMetadata() error {
	// Esta función no se bloquea porque siempre será llamada desde métodos que ya tienen el lock.
	data, err := json.MarshalIndent(r.canciones, "", "  ") // MarshalIndent para que el JSON sea legible
	if err != nil {
		return err
	}
	// os.MkdirAll se asegura de que el directorio 'audios' exista.
	os.MkdirAll(filepath.Dir(metadataPath), os.ModePerm)
	return os.WriteFile(metadataPath, data, 0644)
}

// GetGenres devuelve la lista de géneros únicos disponibles.
func (r *RepositorioCanciones) GetGenres() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Usamos un mapa para encontrar géneros únicos eficientemente.
	generosUnicos := make(map[string]bool)
	for _, cancion := range r.canciones {
		generosUnicos[cancion.Genero] = true
	}

	// Convertimos el mapa de géneros a un slice.
	generos := make([]string, 0, len(generosUnicos))
	for genero := range generosUnicos {
		generos = append(generos, genero)
	}
	return generos
}

// GetSongsByGenre busca en el slice de memoria las canciones que coincidan con el género.
func (r *RepositorioCanciones) GetSongsByGenre(genero string) ([]dtos.CancionDTOOutput, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var cancionesFiltradas []dtos.CancionDTOOutput
	for _, cancion := range r.canciones {
		if cancion.Genero == genero {
			cancionesFiltradas = append(cancionesFiltradas, cancion)
		}
	}
	return cancionesFiltradas, nil
}

// GuardarCancion guarda el archivo de audio y actualiza los metadatos.
func (r *RepositorioCanciones) GuardarCancion(titulo string, genero string, artista string, album string, releaseYear int, data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	NewTitulo := strings.ReplaceAll(titulo, " ", "_")
	NewArtista := strings.ReplaceAll(artista, " ", "_")
	// 1. Construir nombre del archivo y la ruta.
	fileName := fmt.Sprintf("%s_%s.mp3", NewArtista, NewTitulo)
	filePath := filepath.Join("audios", fileName)

	// 2. Comprobar si ya existe una canción con la misma ruta (evitar duplicados).
	for _, cancion := range r.canciones {
		if cancion.FilePath == filePath {
			return fmt.Errorf("la canción '%s' de '%s' ya existe", titulo, artista)
		}
	}

	// 3. Crear el nuevo DTO con los metadatos.
	nuevaCancion := dtos.CancionDTOOutput{
		Titulo:      titulo,
		Artista:     artista,
		Genero:      genero,
		FilePath:    filePath, // Guardamos la ruta para futura referencia
		Album:       album,
		ReleaseYear: releaseYear,
	}

	// 4. Agregar los nuevos metadatos al slice en memoria.
	r.canciones = append(r.canciones, nuevaCancion)

	// 5. Persistir el slice actualizado en el archivo JSON.
	if err := r.guardarMetadata(); err != nil {
		// Si falla al guardar el JSON, no deberíamos guardar el archivo para mantener la consistencia.
		// Revertimos la adición al slice.
		r.canciones = r.canciones[:len(r.canciones)-1]
		return fmt.Errorf("error guardando metadatos: %v", err)
	}

	// 6. Si los metadatos se guardaron correctamente, ahora guardamos el archivo de audio.
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		// En un caso real y más robusto, aquí deberíamos intentar revertir el guardado de metadatos.
		// Por simplicidad, por ahora solo reportamos el error.
		return fmt.Errorf("error guardando el archivo de audio: %v", err)
	}

	fmt.Printf("Canción guardada exitosamente: %s\n", filePath)
	return nil
}
