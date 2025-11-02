/*
Package services implementa la fachada de servicios para el cliente de música.
*/
package services

import (
	"context"
	"io"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
	"proyecto.local/cliente/cancionConsumer"
	"proyecto.local/cliente/models"
	"proyecto.local/cliente/preferenciasConsumer"
	"proyecto.local/cliente/streamingConsumer"

	pbStreaming "proyecto.local/servidor-streaming/serviciosCancion"
)

// MusicFacade encapsula todos los servicios de música
type MusicFacade struct {
	cancionConsumer      *cancionConsumer.CancionConsumer
	streamingConsumer    *streamingConsumer.StreamingConsumer
	preferenciasConsumer *preferenciasConsumer.PreferenciasConsumer
	streamingClient      pbStreaming.StreamingServiceClient // El cliente de streaming gRPC se mantiene
}

// NewMusicFacade crea una nueva instancia de la fachada
func NewMusicFacade(streamingClient pbStreaming.StreamingServiceClient) *MusicFacade {
	return &MusicFacade{
		cancionConsumer:      cancionConsumer.NewCancionConsumer(), // Se crea el nuevo consumidor HTTP
		streamingConsumer:    streamingConsumer.NewStreamingConsumer(),
		preferenciasConsumer: preferenciasConsumer.NewPreferenciasConsumer(), // Nuevo consumidor de preferencias
		streamingClient:      streamingClient,
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

// GetAllSongs obtiene todas las canciones disponibles
func (f *MusicFacade) GetAllSongs() ([]models.Song, error) {
	return f.cancionConsumer.GetAllSongs()
}

// PlaySong reproduce una canción usando gRPC para el streaming.
// Ahora acepta el usuarioID para registrar la reproducción correctamente.
func (f *MusicFacade) PlaySong(filename string, usuarioID string, stopSignal chan bool) error {
	// Reemplazamos las barras invertidas por barras normales para asegurar
	// la compatibilidad entre sistemas operativos y URLs.

	normalizedFilename := strings.ReplaceAll(filename, "\\", "/")

	// Quitamos el prefijo "audios/" porque la nueva ruta del servidor de streaming
	cleanFilename := strings.TrimPrefix(normalizedFilename, "audios/")

	// ✅ Crear contexto con metadata que incluye el usuario ID
	ctx := context.Background()
	md := metadata.New(map[string]string{
		"user-id": usuarioID,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &pbStreaming.StreamSongRequest{SongTitle: cleanFilename}
	stream, err := f.streamingClient.StreamSong(ctx, req)
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

// GetPreferenciasByUserID obtiene las preferencias de un usuario
func (f *MusicFacade) GetPreferenciasByUserID(userID int) (*models.Preferencias, error) {
	// Convertir el userID a string para pasarlo al comando Java
	return f.preferenciasConsumer.GetPreferenciasByUserID(strconv.Itoa(userID))
}
