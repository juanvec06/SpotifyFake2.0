package dtos

type ReproduccionDTOInput struct {
	Titulo    string `json:"titulo"`
	Cliente   string `json:"cliente"`    // Alias: cliente
	UsuarioID string `json:"usuario_id"` // Alias: usuario_id (para compatibilidad)
}

// GetUsuarioID retorna el ID del usuario, usando el campo que esté disponible
func (r *ReproduccionDTOInput) GetUsuarioID() string {
	if r.UsuarioID != "" {
		return r.UsuarioID
	}
	return r.Cliente
}

// ReproduccionDTOOutput es lo que se devuelve al servidor de preferencias Java
type ReproduccionDTOOutput struct {
	ID        int    `json:"id"`
	IDUsuario int    `json:"idUsuario"`
	Titulo    string `json:"titulo"` // Título de la canción para relacionar con el catálogo
}
