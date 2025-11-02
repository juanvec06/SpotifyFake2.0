package comunicacionReproducciones

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ReproduccionDTOInput struct {
	Titulo    string `json:"titulo"`
	UsuarioID string `json:"usuario_id"`
}

func RegistrarReproduccion(titulo string, usuarioID string) error {
	url := "http://localhost:3000/reproducciones/registrar"
	body := ReproduccionDTOInput{
		Titulo:    titulo,
		UsuarioID: usuarioID,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error al convertir a JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error al hacer la petición POST a Reproducciones: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta del servidor de Reproducciones: %s", resp.Status)
	}

	fmt.Println("Reproducción registrada exitosamente.")
	return nil
}
