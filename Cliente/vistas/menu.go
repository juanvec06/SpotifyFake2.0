package vistas

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"proyecto.local/cliente/models"
	"proyecto.local/cliente/services"
)

var reader = bufio.NewReader(os.Stdin)

// Función principal que muestra el menú inicial
func MostrarMenuPrincipal(facade *services.MusicFacade) {
	for {
		fmt.Println("\n===== Spotify =====")
		fmt.Println("1. Ver géneros")
		fmt.Println("2. Salir")
		fmt.Print("Seleccione una opción: ")

		input, _ := reader.ReadString('\n')
		switch strings.TrimSpace(input) {
		case "1":
			mostrarMenuGeneros(facade)
		case "2":
			fmt.Println("¡Hasta luego!")
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

// Muestra la lista de géneros
func mostrarMenuGeneros(facade *services.MusicFacade) {
	genres, err := facade.GetGenres()
	if err != nil {
		fmt.Printf("Error al obtener géneros: %v\n", err)
		return
	}

	for {
		fmt.Println("\n--- Géneros Disponibles ---")
		for i, genre := range genres {
			fmt.Printf("%d. %s\n", i+1, genre.Name)
		}
		fmt.Printf("%d. Atrás\n", len(genres)+1)
		fmt.Print("Seleccione un género: ")

		input, _ := reader.ReadString('\n')
		choice, _ := strconv.Atoi(strings.TrimSpace(input))

		if choice > 0 && choice <= len(genres) {
			genre := genres[choice-1]
			mostrarMenuCanciones(facade, genre)
		} else if choice == len(genres)+1 {
			return
		} else {
			fmt.Println("Opción no válida.")
		}
	}
}

// Muestra la lista de canciones para un género
func mostrarMenuCanciones(facade *services.MusicFacade, genre models.Genre) {
	songs, err := facade.GetSongsByGenre(genre.Name)
	if err != nil {
		fmt.Printf("Error al obtener canciones: %v\n", err)
		return
	}

	for {
		fmt.Printf("\n--- Canciones de %s ---\n", genre.Name)
		for i, song := range songs {
			fmt.Printf("%d. %s - %s\n", i+1, song.Artista, song.Titulo)
		}
		fmt.Printf("%d. Atrás\n", len(songs)+1)
		fmt.Print("Seleccione una canción: ")

		input, _ := reader.ReadString('\n')
		choice, _ := strconv.Atoi(strings.TrimSpace(input))

		if choice > 0 && choice <= len(songs) {
			song := songs[choice-1]
			mostrarMenuDetalles(facade, song)
		} else if choice == len(songs)+1 {
			return
		} else {
			fmt.Println("Opción no válida.")
		}
	}
}

// Muestra los detalles de una canción y la opción de reproducir
func mostrarMenuDetalles(facade *services.MusicFacade, song models.Song) {
	fmt.Printf("\n--- Detalles de la Canción ---\n")
	fmt.Printf("	- Título: %s\n", song.Titulo)
	fmt.Printf("	- Artista: %s\n", song.Artista)
	fmt.Printf("	- Álbum: %s\n", song.Album)
	fmt.Printf("	- Año de lanzamiento: %d\n", song.ReleaseYear)
	fmt.Println("-----------------------------")
	fmt.Println("1. Reproducir")
	fmt.Println("2. Atrás")
	fmt.Print("Seleccione una opción: ")

	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(input) == "1" {
		reproducirCancion(facade, song.FilePath)
	}
}

// Llama a la lógica de streaming usando la fachada
func reproducirCancion(facade *services.MusicFacade, filepath string) {
	stopSignal := make(chan bool)

	// Mostrar menú de control de reproducción
	go mostrarMenuReproduccion(stopSignal)

	// Usar la fachada para reproducir la canción
	err := facade.PlaySong(filepath, stopSignal)
	if err != nil {
		fmt.Printf("Error al reproducir la canción: %v\n", err)
		return
	}
}

// Muestra el menú de control durante la reproducción
func mostrarMenuReproduccion(stopSignal chan bool) {
	fmt.Println("\nReproduciendo Canción")
	fmt.Println("1. Salir")
	fmt.Print("Seleccione una opción: ")

	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(input) == "1" {
		stopSignal <- true
	}
}
