/*
Package capaFachadaServices implementa la lógica de negocio para el servicio de streaming de audio.

Este paquete define la función StreamAudioFile que maneja la lectura del archivo de audio
y utiliza un callback para enviar fragmentos de audio al cliente a través del stream gRPC.
*/
package capaFachadaServices

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Definimos el tamaño de cada fragmento de audio que enviaremos. 64KB es un buen punto de partida.
const chunkSize = 64 * 1024

// StreamAudioFile es la función principal que recibe el título de la canción y una función "callback" que se encargará de enviar cada fragmento.
func StreamAudioFile(songTitle string, sendChunk func(chunk []byte) error) error {
	log.Printf("Fachada: Iniciando streaming para la canción '%s'", songTitle)

	// filePath del archivo en el servidor de canciones
	filePath := "../servidorCanciones/plantilla/audios/" + songTitle
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo de audio: %w", err)
	}
	// Nos aseguramos de que el archivo se cierre al final de la función.
	defer file.Close()

	// Creamos un buffer (un espacio de memoria temporal) para leer los fragmentos.
	buffer := make([]byte, chunkSize)

	for {
		// Leemos un fragmento del archivo y lo guardamos en el buffer.
		bytesRead, err := file.Read(buffer)

		// Si llegamos al final del archivo (End Of File), terminamos el bucle.
		if err == io.EOF {
			log.Println("Fachada: Se ha terminado de leer el archivo.")
			break
		}
		// Si hay otro tipo de error, lo devolvemos.
		if err != nil {
			return fmt.Errorf("error leyendo el archivo: %w", err)
		}

		// Usamos la función callback para enviar el fragmento que acabamos de leer.
		if err := sendChunk(buffer[:bytesRead]); err != nil {
			// Si el callback falla (ej: el cliente se desconectó), devolvemos el error.
			return fmt.Errorf("la función de envío de fragmentos falló: %w", err)
		}
	}

	return nil
}
