package fachada

import (
	"tendencias/capaAccesoADatos/entitys"
	capaaccesoadatos "tendencias/capaAccesoADatos/repositorios"
	dtos "tendencias/capaFachadaServices/DTOs"
)

type FachadaTendencias struct {
	repo *capaaccesoadatos.RepositorioReproducciones
}

// Constructor de la fachada
func NuevaFachadaTendencias() *FachadaTendencias {
	return &FachadaTendencias{
		repo: capaaccesoadatos.GetRepositorio(),
	}
}

// Método para registrar una reproducción
func (f *FachadaTendencias) RegistrarReproduccion(titulo, cliente string) {
	f.repo.AgregarReproduccion(titulo, cliente)
}

// Método para obtener todas las reproducciones
func (f *FachadaTendencias) ObtenerReproducciones() []entitys.ReproduccionEntity {
	return f.repo.ListarReproducciones()
}

// Método para obtener reproducciones filtradas por usuario (para servidor de preferencias Java)
func (f *FachadaTendencias) ObtenerReproduccionesPorUsuario(idUsuario int) []dtos.ReproduccionDTOOutput {
	return f.repo.ListarReproduccionesPorUsuario(idUsuario)
}
