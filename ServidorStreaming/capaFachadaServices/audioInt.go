/*
Package capaFachadaServices implementa la lógica de negocio para el servicio de streaming de audio.
*/
package capaFachadaServices

import (
	"fmt"
	"io"
	"log"
	"net/http" // Importamos el paquete http
	"strings"  // Importamos el paquete strings
)

const chunkSize = 64 * 1024

// URL base del servidor que almacena las canciones
const songServerBaseURL = "http://localhost:5000"

func StreamAudioFile(songPath string, sendChunk func(chunk []byte) error) error {
	log.Printf("Fachada: Solicitando streaming para la ruta '%s'", songPath)

	// Nos aseguramos de que la ruta use slashes (/) para la URL
	urlPath := strings.ReplaceAll(songPath, "\\", "/")

	// Construimos la URL completa para solicitar el archivo de audio al servidor de canciones.
	// Ej: http://localhost:5000/audio/Los enanitos verdes_Lamento Boliviano.mp3
	fullURL := fmt.Sprintf("%s/audio/%s", songServerBaseURL, urlPath)

	log.Printf("Fachada: Realizando petición GET a %s", fullURL)

	// Hacemos la petición HTTP para obtener el archivo de audio.
	resp, err := http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("no se pudo hacer la petición al servidor de canciones: %w", err)
	}
	defer resp.Body.Close()

	// Verificamos si la canción fue encontrada
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("el servidor de canciones respondió con un error: %s", resp.Status)
	}

	// Creamos un buffer (un espacio de memoria temporal) para leer los fragmentos.
	buffer := make([]byte, chunkSize)

	for {
		// Leemos un fragmento del cuerpo de la respuesta HTTP (que es el audio).
		bytesRead, err := resp.Body.Read(buffer)

		if err == io.EOF {
			log.Println("Fachada: Se ha terminado de leer el stream del servidor de canciones.")
			break
		}
		if err != nil {
			return fmt.Errorf("error leyendo el stream de audio: %w", err)
		}

		// Usamos la función callback para enviar el fragmento al cliente final.
		if err := sendChunk(buffer[:bytesRead]); err != nil {
			return fmt.Errorf("la función de envío de fragmentos falló: %w", err)
		}
	}

	return nil
}
