/*
Package main es el punto de entrada del servidor gRPC.

Este servidor implementa un servicio de streaming de audio definido en el archivo serviciosCancion.proto.
Utiliza un controlador definido en el paquete capaControladores para manejar las solicitudes entrantes.
*/
package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	// Importamos nuestro controlador
	"proyecto.local/servidor-streaming/capaControladores"
	// Importamos el código generado por protoc
	pb "proyecto.local/servidor-streaming/serviciosCancion"
)

const port = ":50051"

func main() {
	// 1. Abrimos un puerto TCP para escuchar peticiones.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Fallo al escuchar en el puerto %s: %v", port, err)
	}

	// 2. Creamos una nueva instancia de un servidor gRPC.
	grpcServer := grpc.NewServer()

	// 3. Registramos nuestro controlador en el servidor gRPC.
	// Esto le dice al servidor que cualquier llamada al 'StreamingService'
	// debe ser manejada por nuestra struct 'ControladorServidor'.
	pb.RegisterStreamingServiceServer(grpcServer, &capaControladores.ControladorServidor{})

	log.Printf("Servidor gRPC escuchando en el puerto %s", port)

	// 4. Iniciamos el servidor para que empiece a aceptar conexiones.
	// El programa se quedará bloqueado en esta línea, sirviendo peticiones.
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Fallo al iniciar el servidor gRPC: %v", err)
	}
}
