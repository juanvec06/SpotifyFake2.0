package capafachada

import (
	capaaccesoadatos "almacenamiento/capaAccesoADatos"
	"log"
)

type FachadaGeneros struct {
	data *capaaccesoadatos.RepositorioCanciones
}

func NuevaFachadaGeneros() *FachadaGeneros {
	return &FachadaGeneros{
		data: capaaccesoadatos.GetRepositorioCanciones(),
	}
}

func (f *FachadaGeneros) GetGenresService() ([]string, error) {
	genres := f.data.GetGenres()

	if len(genres) == 0 {
		log.Println("Fachada: No se encontraron géneros.")
	}

	return genres, nil
}
