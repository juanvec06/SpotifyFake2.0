package entitys

type ReproduccionEntity struct {
	ID        int
	Titulo    string
	IDUsuario int
	IDCancion int // ID de la canción (mock por ahora)
	FechaHora string
}
