/*
Package capaControladores implementa los controladores gRPC para el servidor de streaming de audio.

Este paquete define el struct ControladorServidor que implementa los métodos del servicio gRPC
*/
package capaControladores

import (
	"log"

	"google.golang.org/grpc/metadata"
	comunicacionReproducciones "proyecto.local/servidor-streaming/capaComunicacionReproducciones"
	"proyecto.local/servidor-streaming/capaFachadaServices"
	pb "proyecto.local/servidor-streaming/serviciosCancion"
)

// ControladorServidor es una struct que implementará los métodos de nuestro servicio gRPC.
type ControladorServidor struct {
	pb.UnimplementedStreamingServiceServer
}

// StreamSong es la implementación del método RPC definido en nuestro .proto.
func (s *ControladorServidor) StreamSong(req *pb.StreamSongRequest, stream pb.StreamingService_StreamSongServer) error {
	songTitle := req.GetSongTitle()
	log.Printf("🎵 Controlador: Petición gRPC recibida para '%s'", songTitle)

	// ✅ Extraer usuario ID del metadata de gRPC
	usuarioID := "1" // Valor por defecto si no se proporciona metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		if userIDs := md.Get("user-id"); len(userIDs) > 0 {
			usuarioID = userIDs[0]
			log.Printf("👤 Usuario autenticado desde metadata: %s", usuarioID)
		} else {
			log.Printf("⚠️  No se encontró 'user-id' en metadata, usando usuario por defecto: %s", usuarioID)
		}
	} else {
		log.Printf("⚠️  No hay metadata en la petición, usando usuario por defecto: %s", usuarioID)
	}

	// Registrar reproducción de manera ASÍNCRONA
	// Goroutine asíncrona: no bloquea el streaming de audio
	go func() {
		err := comunicacionReproducciones.RegistrarReproduccion(songTitle, usuarioID)
		if err != nil {
			// Solo loguear el error, no afecta el streaming
			log.Printf("⚠️ Error registrando reproducción (operación asíncrona): %v", err)
		} else {
			log.Printf("✅ Reproducción registrada: '%s' para usuario '%s'", songTitle, usuarioID)
		}
	}()

	// Esta función sabe cómo usar el 'stream' de gRPC para enviar datos.
	sendChunkCallback := func(chunk []byte) error {
		// Creamos el DTO de respuesta (AudioChunk)
		response := &pb.AudioChunk{
			Data: chunk,
		}
		// Usamos el método Send() del stream para enviar el fragmento al cliente.
		return stream.Send(response)
	}

	return capaFachadaServices.StreamAudioFile(songTitle, sendChunkCallback)
}
