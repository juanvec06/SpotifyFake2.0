package co.edu.unicauca.fachadaServices.DTO;

import java.io.Serializable;
import java.util.List;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class PreferenciasDTORespuesta implements Serializable {
    private Integer idUsuario;
    private List<PreferenciaGeneroDTORespuesta> preferenciasGeneros;
    private List<PreferenciaArtistaDTORespuesta> preferenciasArtistas;
    private List<PreferenciaIdiomaDTORespuesta> preferenciasIdiomas;
}
