package dtos

type CancionAlmacenarDTOInput struct {
	Titulo      string `json:"titulo"`
	Artista     string `json:"artista"`
	Album       string `json:"album"`
	Genero      string `json:"genero"`
	ReleaseYear int    `json:"release_year"`
	Idioma      string `json:"idioma"`
}
