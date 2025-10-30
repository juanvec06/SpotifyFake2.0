/*
Package capaControladores implementa los controladores gRPC para el servidor de streaming de audio.

Este paquete define el struct ControladorServidor que implementa los métodos del servicio gRPC
*/
package capaControladores

import (
	"log"

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
	log.Printf("Controlador: Petición gRPC recibida para '%s'", songTitle)

	// Aquí está la magia: definimos una función que será nuestro "callback".
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
