package co.edu.unicauca.fachadaServices.services.componenteCalculaPreferencias;

import java.util.Comparator;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.stream.Collectors;

import co.edu.unicauca.fachadaServices.DTO.CancionDTOEntrada;
import co.edu.unicauca.fachadaServices.DTO.PreferenciaArtistaDTORespuesta;
import co.edu.unicauca.fachadaServices.DTO.PreferenciaGeneroDTORespuesta;
import co.edu.unicauca.fachadaServices.DTO.PreferenciaIdiomaDTORespuesta;
import co.edu.unicauca.fachadaServices.DTO.PreferenciasDTORespuesta;
import co.edu.unicauca.fachadaServices.DTO.ReproduccionesDTOEntrada;

public class CalculadorPreferencias {

    public PreferenciasDTORespuesta calcular(Integer idUsuario,
            List<CancionDTOEntrada> canciones,
            List<ReproduccionesDTOEntrada> reproducciones) {
        // Crear mapa de canciones usando el título como clave
        Map<String, CancionDTOEntrada> mapaCanciones = canciones.stream()
                .filter(Objects::nonNull)
                .filter(c -> c.getTitulo() != null)
                .collect(Collectors.toMap(CancionDTOEntrada::getTitulo, c -> c, (a, b) -> a));

        System.out.println("Mapa de canciones creado con " + mapaCanciones.size() + " canciones");

        Map<String, Integer> contadorGeneros = new HashMap<>();
        Map<String, Integer> contadorArtistas = new HashMap<>();
        Map<String, Integer> contadorIdiomas = new HashMap<>();

        for (ReproduccionesDTOEntrada r : reproducciones) {
            String tituloCancion = r.getTitulo();
            if (tituloCancion == null) {
                System.out.println("Reproducción sin título, omitida");
                continue;
            }

            CancionDTOEntrada c = mapaCanciones.get(tituloCancion);
            if (c == null) {
                System.out.println("No se encontró canción con título: '" + tituloCancion + "'");
                continue;
            }

            System.out.println("Procesada: " + tituloCancion + " - Género: " + c.getGenero() + ", Artista: "
                    + c.getArtista() + ", Idioma: " + c.getIdioma());

            String genero = c.getGenero() == null ? "Desconocido" : c.getGenero();
            String artista = c.getArtista() == null ? "Desconocido" : c.getArtista();
            String idioma = c.getIdioma() == null ? "Desconocido" : c.getIdioma();

            contadorGeneros.put(genero, contadorGeneros.getOrDefault(genero, 0) + 1);
            contadorArtistas.put(artista, contadorArtistas.getOrDefault(artista, 0) + 1);
            contadorIdiomas.put(idioma, contadorIdiomas.getOrDefault(idioma, 0) + 1);
        }

        System.out.println("Total géneros encontrados: " + contadorGeneros.size());
        System.out.println("Total artistas encontrados: " + contadorArtistas.size());
        System.out.println("Total idiomas encontrados: " + contadorIdiomas.size());

        List<PreferenciaGeneroDTORespuesta> preferenciasGeneros = contadorGeneros.entrySet().stream()
                .map(e -> {
                    PreferenciaGeneroDTORespuesta dto = new PreferenciaGeneroDTORespuesta();
                    dto.setNombreGenero(e.getKey());
                    dto.setNumeroPreferencias(e.getValue());
                    return dto;
                })
                .sorted(Comparator.comparingInt(PreferenciaGeneroDTORespuesta::getNumeroPreferencias).reversed()
                        .thenComparing(PreferenciaGeneroDTORespuesta::getNombreGenero))
                .collect(Collectors.toList());

        List<PreferenciaArtistaDTORespuesta> preferenciasArtistas = contadorArtistas.entrySet().stream()
                .map(e -> {
                    PreferenciaArtistaDTORespuesta dto = new PreferenciaArtistaDTORespuesta();
                    dto.setNombreArtista(e.getKey());
                    dto.setNumeroPreferencias(e.getValue());
                    return dto;
                })
                .sorted(Comparator.comparingInt(PreferenciaArtistaDTORespuesta::getNumeroPreferencias).reversed()
                        .thenComparing(PreferenciaArtistaDTORespuesta::getNombreArtista))
                .collect(Collectors.toList());

        List<PreferenciaIdiomaDTORespuesta> preferenciasIdiomas = contadorIdiomas.entrySet().stream()
                .map(e -> {
                    PreferenciaIdiomaDTORespuesta dto = new PreferenciaIdiomaDTORespuesta();
                    dto.setNombreIdioma(e.getKey());
                    dto.setNumeroPreferencias(e.getValue());
                    return dto;
                })
                .sorted(Comparator.comparingInt(PreferenciaIdiomaDTORespuesta::getNumeroPreferencias).reversed()
                        .thenComparing(PreferenciaIdiomaDTORespuesta::getNombreIdioma))
                .collect(Collectors.toList());

        PreferenciasDTORespuesta respuesta = new PreferenciasDTORespuesta();
        respuesta.setIdUsuario(idUsuario);
        respuesta.setPreferenciasGeneros(preferenciasGeneros);
        respuesta.setPreferenciasArtistas(preferenciasArtistas);
        respuesta.setPreferenciasIdiomas(preferenciasIdiomas);

        return respuesta;
    }
}
