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
			mostrarMenuGeneros(facade, Usuario)
		case "2":
			mostrarMenuCanciones(facade, models.Genre{Name: ""}, Usuario)
		case "3":
			mostrarMenuPreferencias(facade, Usuario)
		case "4":
			fmt.Println("¡Hasta luego!")
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

// Muestra la lista de géneros
func mostrarMenuGeneros(facade *services.MusicFacade, usuario string) {
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
			mostrarMenuCanciones(facade, genre, usuario)
		} else if choice == len(genres)+1 {
			return
		} else {
			fmt.Println("Opción no válida.")
		}
	}
}

// Muestra la lista de canciones para un género
func mostrarMenuCanciones(facade *services.MusicFacade, genre models.Genre, usuario string) {
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
			mostrarMenuDetalles(facade, song, usuario)
		} else if choice == len(songs)+1 {
			return
		} else {
			fmt.Println("Opción no válida.")
		}
	}
}

// Muestra los detalles de una canción y la opción de reproducir
func mostrarMenuDetalles(facade *services.MusicFacade, song models.Song, usuario string) {
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
		reproducirCancion(facade, song.FilePath, usuario)
	}
}

// Llama a la lógica de streaming usando la fachada
func reproducirCancion(facade *services.MusicFacade, filepath string, usuario string) {
	// Convertir nombre de usuario a ID numérico
	// usuario1 -> "1", usuario2 -> "2", usuario3 -> "3"
	usuarioID := "1" // Valor por defecto
	if len(usuario) > 0 {
		// Extraer el último carácter del nombre de usuario
		ultimoChar := usuario[len(usuario)-1:]
		if ultimoChar >= "1" && ultimoChar <= "9" {
			usuarioID = ultimoChar
		}
	}

	fmt.Printf("🎵 Reproduciendo para usuario: %s (ID: %s)\n", usuario, usuarioID)

	stopSignal := make(chan bool)

	// Mostrar menú de control de reproducción
	go mostrarMenuReproduccion(stopSignal)

	// Usar la fachada para reproducir la canción con el usuario ID
	err := facade.PlaySong(filepath, usuarioID, stopSignal)
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

// Muestra el menú de preferencias musicales del usuario
func mostrarMenuPreferencias(facade *services.MusicFacade, usuario string) {
	fmt.Println("\n===== Preferencias Musicales =====")

	// Mapeo de usuario a ID (basado en db.json del servidor de preferencias)
	usuariosIDs := map[string]int{
		"usuario1": 1, // Daniel
		"usuario2": 2, // María
		"usuario3": 3, // Carlos
	}

	userID, exists := usuariosIDs[usuario]
	if !exists {
		fmt.Printf("❌ No se encontró ID para el usuario '%s'\n", usuario)
		fmt.Println("\nPresione Enter para continuar...")
		reader.ReadString('\n')
		return
	}

	fmt.Printf("Obteniendo preferencias para %s (ID: %d)...\n\n", usuario, userID)

	preferencias, err := facade.GetPreferenciasByUserID(userID)
	if err != nil {
		fmt.Printf("❌ Error al obtener preferencias: %v\n", err)
		fmt.Println("\n⚠️  Asegúrate de que:")
		fmt.Println("   1. El servidor de Preferencias esté ejecutándose (puerto 8080)")
		fmt.Println("   2. El servidor de Reproducciones esté ejecutándose (puerto 3000)")
		fmt.Println("   3. Hayas reproducido al menos una canción")
		fmt.Println("\nPresione Enter para continuar...")
		reader.ReadString('\n')
		return
	}

	// Mostrar preferencias
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Printf("║  Usuario ID: %-25d ║\n", preferencias.IDUsuario)
	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Println("║  🎵 GÉNEROS FAVORITOS                 ║")
	fmt.Println("╠════════════════════════════════════════╣")

	if len(preferencias.PreferenciasGeneros) == 0 {
		fmt.Println("║  (No hay géneros registrados)          ║")
	} else {
		for i, genero := range preferencias.PreferenciasGeneros {
			info := fmt.Sprintf("%s (%d reproducciones)", genero.NombreGenero, genero.NumeroPreferencias)
			fmt.Printf("║  %d. %-35s ║\n", i+1, info)
		}
	}

	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Println("║  🎤 ARTISTAS FAVORITOS                ║")
	fmt.Println("╠════════════════════════════════════════╣")

	if len(preferencias.PreferenciasArtistas) == 0 {
		fmt.Println("║  (No hay artistas registrados)         ║")
	} else {
		for i, artista := range preferencias.PreferenciasArtistas {
			info := fmt.Sprintf("%s (%d reproducciones)", artista.NombreArtista, artista.NumeroPreferencias)
			fmt.Printf("║  %d. %-35s ║\n", i+1, info)
		}
	}

	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Println("║  🌍 IDIOMAS FAVORITOS                 ║")
	fmt.Println("╠════════════════════════════════════════╣")

	if len(preferencias.PreferenciasIdiomas) == 0 {
		fmt.Println("║  (No hay idiomas registrados)          ║")
	} else {
		for i, idioma := range preferencias.PreferenciasIdiomas {
			fmt.Printf("║  %d. %-25s (%d) ║\n", i+1, idioma.NombreIdioma, idioma.NumeroPreferencias)
		}
	}

	fmt.Println("╚════════════════════════════════════════╝")

	fmt.Println("\nPresione Enter para volver al menú principal...")
	reader.ReadString('\n')
}
