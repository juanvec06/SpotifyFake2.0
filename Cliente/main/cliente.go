/*
Package main es el punto de entrada del cliente gRPC que interactúa con los servidores de canciones y streaming.

Este cliente se conecta a dos servidores gRPC: uno para gestionar las canciones y otro para el streaming de audio.

El cliente utiliza las vistas definidas en el paquete "vistas" para presentar un menú interactivo al usuario, permitiéndole listar canciones por genero y reproducirlas mediante streaming.
El cliente utiliza las vistas definidas en el paquete "vistas" para presentar un menú interactivo al usuario, permitiéndole listar canciones por genero y reproducirlas mediante streaming.
*/
package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"proyecto.local/cliente/services"
	"proyecto.local/cliente/vistas"
	pbStreaming "proyecto.local/servidor-streaming/serviciosCancion"
)

// songsServerAddr es la dirección del servidor de canciones
const songsServerAddr = "localhost:50052"

// streamingServerAddr es la dirección del servidor de streaming
const streamingServerAddr = "localhost:50051"

func main() {

	// Conexión al servidor de streaming
	connStreaming, err := grpc.Dial(streamingServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor de streaming: %v", err)
	}
	defer connStreaming.Close()
	streamingClient := pbStreaming.NewStreamingServiceClient(connStreaming)

	// Crear la fachada de servicios
	musicFacade := services.NewMusicFacade(streamingClient)

	// Pasar la fachada a las vistas
	vistas.MostrarMenuPrincipal(musicFacade)
}
