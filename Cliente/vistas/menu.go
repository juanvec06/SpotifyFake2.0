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

var usuarios = map[string]string{
	"usuario1": "password1",
	"usuario2": "password2",
	"usuario3": "password3",
}

func MostrarIniciarSesion(facade *services.MusicFacade) {
	fmt.Println("===== Iniciar Sesión =====")
	fmt.Print("Usuario: ")
	usuario, _ := reader.ReadString('\n')
	usuario = strings.TrimSpace(usuario)
	fmt.Print("Contraseña: ")
	contraseña, _ := reader.ReadString('\n')
	contraseña = strings.TrimSpace(contraseña)

	if password, ok := usuarios[usuario]; ok && password == contraseña {
		fmt.Println("Inicio de sesión exitoso.")
		MostrarMenuPrincipal(facade, usuario)
	} else {
		fmt.Println("Usuario o contraseña incorrectos.")
		return
	}
}

// Función principal que muestra el menú inicial
func MostrarMenuPrincipal(facade *services.MusicFacade, Usuario string) {
	for {
		fmt.Println("\n===== Spotify =====")
		fmt.Println("1. Ver géneros")
		fmt.Println("2. Ver canciones")
		fmt.Println("3. Ver preferencias")
		fmt.Println("4. Salir")
		fmt.Print("Seleccione una opción: ")

		input, _ := reader.ReadString('\n')
		switch strings.TrimSpace(input) {
		case "1":
			mostrarMenuGeneros(facade)
		case "2":
			mostrarMenuCanciones(facade, models.Genre{Name: ""})
		case "3":
			//TODO
			//mostrarMenuPreferencias(facade, Usuario)
		case "4":
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
	flagAllSongs := false
	if genre.Name == "" {
		flagAllSongs = true
	}
	songs := []models.Song{}
	err := error(nil)
	if flagAllSongs {
		songs, err = facade.GetAllSongs()
		if err != nil {
			fmt.Printf("Error al obtener todas las canciones: %v\n", err)
			return
		}
	} else {
		songs, err = facade.GetSongsByGenre(genre.Name)
		if err != nil {
			fmt.Printf("Error al obtener canciones: %v\n", err)
			return
		}
	}

	for {
		fmt.Printf("\n--- Canciones  %s ---\n", genre.Name)
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
	playbackDone := make(chan struct{})

	// 1. Inicia la reproducción en una goroutine en segundo plano.
	go func() {
		err := facade.PlaySong(filepath, stopSignal)
		if err != nil {
			// Imprime el error en una nueva línea para no interferir con el prompt del usuario
			fmt.Printf("\nError al reproducir la canción: %v\n", err)
		}
		close(playbackDone) // Señala que la reproducción ha terminado por cualquier motivo.
	}()

	// 2. La goroutine principal ahora maneja la entrada del usuario.
	fmt.Println("\n--- Reproduciendo Canción ---")
	fmt.Println("Presione '1' y Enter para detener y volver.")
	fmt.Print("Opción: ")

	// 3. Crea una goroutine solo para leer la entrada y enviarla por un canal.
	userInputChan := make(chan string)
	go func() {
		input, _ := reader.ReadString('\n')
		userInputChan <- strings.TrimSpace(input)
	}()

	// 4. Usa 'select' para esperar el primer evento que ocurra:
	//    - La canción termina por sí sola.
	//    - El usuario ingresa una opción.
	select {
	case <-playbackDone:
		// La goroutine de userInputChan sigue bloqueada esperando una entrada.
		// Al informar al usuario que presione Enter, esa entrada será consumida por la goroutine "huérfana",
		// evitando que interfiera con el siguiente menú.
		fmt.Println("\nLa canción ha terminado. Presiona Enter para volver al menú.")
		<-userInputChan
		return

	case input := <-userInputChan:
		// El usuario ingresó algo.
		if input == "1" {
			stopSignal <- true // Envía la señal para detener.
			<-playbackDone     // Espera a que la goroutine de reproducción confirme que ha terminado.
			fmt.Println("\nReproducción detenida.")
		}
		// Si el usuario ingresa otra cosa, simplemente volvemos al menú.
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
