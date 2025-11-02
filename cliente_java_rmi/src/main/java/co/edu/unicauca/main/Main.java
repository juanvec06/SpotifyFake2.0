package co.edu.unicauca.main;

import co.edu.unicauca.capaDeControladores.ControladorPreferenciasUsuariosInt;
import co.edu.unicauca.configuracion.lector.LectorPropiedadesConfig;
import co.edu.unicauca.configuracion.servicios.ClienteDeObjetos;
import co.edu.unicauca.fachadaServices.DTO.PreferenciasDTORespuesta;
import co.edu.unicauca.fachadaServices.services.FachadaGestorUsuariosIml;
import co.edu.unicauca.vista.Menu;
import com.google.gson.Gson; // Importa la librería Gson

public class Main {
    public static void main(String[] args) {
        // Si se pasa un argumento, entra en modo "línea de comandos" para ser llamado por Go
        if (args.length > 0) {
            try {
                int idUsuario = Integer.parseInt(args[0]);
                consultarYSerializarPreferencias(idUsuario);
            } catch (NumberFormatException e) {
                System.err.println("Error: El argumento debe ser un número entero (ID de usuario)");
                System.exit(1); // Finaliza con un código de error
            }
            return; // Termina la ejecución
        }

        // Si no hay argumentos, entra en modo interactivo para uso humano
        ejecutarModoInteractivo();
    }

    // Modo para ser llamado desde Go: obtiene los datos, los convierte a JSON y los imprime.
    private static void consultarYSerializarPreferencias(int idUsuario) {
        try {
            ControladorPreferenciasUsuariosInt objRemoto = obtenerObjetoRemoto();
            // Asumo que el método se llama 'getReferencias' como en las diapositivas
            // Si tu método se llama 'consultar', usa ese nombre.
            PreferenciasDTORespuesta preferencias = objRemoto.getReferencias(idUsuario); 
            
            // Usar Gson para convertir el objeto a un string JSON de forma segura
            Gson gson = new Gson();
            String jsonOutput = gson.toJson(preferencias);

            // Imprimir el string JSON resultante a la salida estándar
            System.out.println(jsonOutput);

        } catch (Exception e) {
            // Imprimir errores en la salida de error estándar para no contaminar el JSON
            System.err.println("Error al consultar preferencias: " + e.getMessage());
            System.exit(1);
        }
    }

    // Modo para que un humano use el cliente directamente
    private static void ejecutarModoInteractivo() {
        System.out.println("Iniciando cliente en modo interactivo...");
        ControladorPreferenciasUsuariosInt objRemoto = obtenerObjetoRemoto();
        FachadaGestorUsuariosIml objFachada = new FachadaGestorUsuariosIml(objRemoto);
        Menu objMenu = new Menu(objFachada);
        objMenu.ejecutarMenuPrincipal();
    }

    // Centralizamos la lógica de conexión para no repetirla
    private static ControladorPreferenciasUsuariosInt obtenerObjetoRemoto() {
        int puertoNS = Integer.parseInt(LectorPropiedadesConfig.get("ns.port"));
        String direccionIPNS = LectorPropiedadesConfig.get("ns.host");
        String identificadorObjetoRemoto = "objControladorPreferenciasUsuarios";
        return ClienteDeObjetos.obtenerObjetoRemoto(direccionIPNS, puertoNS, identificadorObjetoRemoto);
    }
}