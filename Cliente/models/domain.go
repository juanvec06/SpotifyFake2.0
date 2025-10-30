/*
Package models define las estructuras de dominio del cliente.
Estas estructuras representan los datos de negocio de manera limpia.
*/
package models

// Genre ahora solo contiene el nombre, ya que la API no proporciona un ID.
type Genre struct {
	Name string
}

// Song representa los datos de una canción tal como los devuelve la API REST.
type Song struct {
	Titulo      string `json:"titulo"` // Estas etiquetas `json` ayudan a mapear la respuesta
	Artista     string `json:"artista"`
	Genero      string `json:"genero"`
	FilePath    string `json:"file_path"`
	Album       string `json:"album"`
	ReleaseYear int    `json:"release_year"`
}
