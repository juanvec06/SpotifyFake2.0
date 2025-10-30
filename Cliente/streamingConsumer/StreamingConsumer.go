package streamingConsumer

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"

	pbStreaming "proyecto.local/servidor-streaming/serviciosCancion"
)

// Flag es una variable global para identificar cuando el usuario quizo parar la reproduccion
var Flag bool = false

// StreamingConsumer se encarga de la lógica de streaming de audio.
type StreamingConsumer struct{}

// NewStreamingConsumer crea una nueva instancia de StreamingConsumer.
func NewStreamingConsumer() *StreamingConsumer {
	return &StreamingConsumer{}
}

// RecibirCancion escucha el stream y escribe en la tubería.
func (sc *StreamingConsumer) RecibirCancion(stream pbStreaming.StreamingService_StreamSongClient, writer *io.PipeWriter) {
	defer writer.Close()

	for {
		fragmento, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Printf("Error al recibir un fragmento: %v", err)
			return
		}

		_, err = writer.Write(fragmento.Data)
		if err != nil && !Flag {
			log.Printf("Error al escribir en la tubería: %v", err)
			return
		}
	}
}

// DecodificarReproducir lee de la tubería y reproduce el audio.
func (sc *StreamingConsumer) DecodificarReproducir(reader *io.PipeReader, done chan struct{}, stopSignal chan bool) {
	streamer, format, err := mp3.Decode(reader)
	if err != nil {
		log.Printf("Error decodificando el MP3: %v", err)
		close(done)
		return
	}
	defer streamer.Close()

	// Intentar inicializar el speaker, pero manejar el caso donde ya está inicializado
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Printf("Warning: Speaker ya inicializado o error en inicialización: %v", err)
	}

	// Crear un canal interno para esperar a que termine la reproducción
	playbackDone := make(chan bool)

	doneCallback := beep.Callback(func() {
		fmt.Println("\nReproducción finalizada.")
		playbackDone <- true
	})

	speaker.Play(beep.Seq(streamer, doneCallback))

	// Esperar a que termine la reproducción o recibir señal de parada
	select {
	case <-playbackDone:
		// Reproducción terminó naturalmente
	case <-stopSignal:
		// Usuario quiere parar la reproducción
		Flag = true
		speaker.Clear()
		fmt.Println("\nReproducción detenida por el usuario.")
	}

	close(done)
}
