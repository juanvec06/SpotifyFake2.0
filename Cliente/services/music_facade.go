/*
Package services implementa la fachada de servicios para el cliente de música.
*/
package services

import (
	"context"
	"io"
	"strings"

	"proyecto.local/cliente/cancionConsumer"
	"proyecto.local/cliente/models"
	"proyecto.local/cliente/streamingConsumer"

	// pbSongs ya no es necesario
	pbStreaming "proyecto.local/servidor-streaming/serviciosCancion"
)

// MusicFacade encapsula todos los servicios de música
type MusicFacade struct {
	cancionConsumer   *cancionConsumer.CancionConsumer
	streamingConsumer *streamingConsumer.StreamingConsumer
	streamingClient   pbStreaming.StreamingServiceClient // El cliente de streaming gRPC se mantiene
}

// NewMusicFacade crea una nueva instancia de la fachada
// Ya no necesita recibir el cliente gRPC de canciones.
func NewMusicFacade(streamingClient pbStreaming.StreamingServiceClient) *MusicFacade {
	return &MusicFacade{
		cancionConsumer:   cancionConsumer.NewCancionConsumer(), // Se crea el nuevo consumidor HTTP
		streamingConsumer: streamingConsumer.NewStreamingConsumer(),
		streamingClient:   streamingClient,
	}
}

// GetGenres obtiene todos los géneros disponibles
func (f *MusicFacade) GetGenres() ([]models.Genre, error) {
	return f.cancionConsumer.GetGenres()
}

// GetSongsByGenre obtiene todas las canciones de un género específico
// La firma del método ahora acepta un string.
func (f *MusicFacade) GetSongsByGenre(genreName string) ([]models.Song, error) {
	return f.cancionConsumer.GetSongsByGenre(genreName)
}

// PlaySong se mantiene igual, ya que sigue usando gRPC para el streaming.
func (f *MusicFacade) PlaySong(filename string, stopSignal chan bool) error {
	// Reemplazamos las barras invertidas por barras normales para asegurar
	// la compatibilidad entre sistemas operativos y URLs.

	normalizedFilename := strings.ReplaceAll(filename, "\\", "/")

	// Quitamos el prefijo "audios/" porque la nueva ruta del servidor de streaming
	// ya lo compone internamente.
	cleanFilename := strings.TrimPrefix(normalizedFilename, "audios/")

	req := &pbStreaming.StreamSongRequest{SongTitle: cleanFilename}
	stream, err := f.streamingClient.StreamSong(context.Background(), req)
	if err != nil {
		return err
	}

	reader, writer := io.Pipe()
	done := make(chan struct{})

	go f.streamingConsumer.RecibirCancion(stream, writer)
	go f.streamingConsumer.DecodificarReproducir(reader, done, stopSignal)

	<-done
	return nil
}
