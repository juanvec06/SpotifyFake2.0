package preferenciasConsumer

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"proyecto.local/cliente/models" // Tu import ya es correcto
)

type PreferenciasConsumer struct {
	javaClientPath string
}

func NewPreferenciasConsumer() *PreferenciasConsumer {
	// Ruta robusta que apunta a la herramienta JAR compilada en el directorio 'bin'
	javaClientPath := "./bin/cliente-preferencias.jar"
	return &PreferenciasConsumer{
		javaClientPath: javaClientPath,
	}
}

func (pc *PreferenciasConsumer) GetPreferenciasByUserID(userID string) (*models.Preferencias, error) {
	// Ejecutar el cliente Java RMI pasando el ID del usuario como argumento
	cmd := exec.Command("java", "-jar", pc.javaClientPath, userID)

	// CombinedOutput captura tanto la salida estándar (JSON) como la de error
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Si hay un error, 'output' contendrá el mensaje de error del programa Java
		return nil, fmt.Errorf("error al ejecutar el cliente Java RMI: %w\nSalida del proceso: %s", err, string(output))
	}

	// 'output' ahora contiene el string JSON que imprimió el cliente Java.
	// Vamos a decodificarlo en tu struct 'models.Preferencias'.

	// 👇 LA ÚNICA LÍNEA QUE CAMBIA: Usamos tu struct 'Preferencias'
	var preferencias models.Preferencias

	if err := json.Unmarshal(output, &preferencias); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta JSON del cliente Java: %w\nRespuesta recibida: %s", err, string(output))
	}

	return &preferencias, nil
}
