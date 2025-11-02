package co.edu.unicauca.fachadaServices.DTO;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class ReproduccionesDTOEntrada {
   private Integer idUsuario;
   private String titulo;  // Título de la canción para matching
}

