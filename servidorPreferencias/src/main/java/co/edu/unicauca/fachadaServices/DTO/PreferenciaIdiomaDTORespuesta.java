package co.edu.unicauca.fachadaServices.DTO;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class PreferenciaIdiomaDTORespuesta {
    private String nombreIdioma;
    private Integer numeroPreferencias;
}
