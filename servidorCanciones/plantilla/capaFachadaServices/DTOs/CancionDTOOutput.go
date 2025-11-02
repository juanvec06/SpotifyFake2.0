package dtos

type CancionDTOOutput struct {
	Titulo      string `json:"titulo"`
	Artista     string `json:"artista"`
	Album       string `json:"album"`
	Genero      string `json:"genero"`
	ReleaseYear int    `json:"release_year"`
	FilePath    string `json:"file_path"`
	Idioma      string `json:"idioma"`
}
