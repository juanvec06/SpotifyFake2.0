package co.edu.unicauca.infoii.correo.DTOs;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
public class CancionAlmacenarDTOInput {
    private String id;
    private String titulo;
    private String artista;
    private String genero;
}
	