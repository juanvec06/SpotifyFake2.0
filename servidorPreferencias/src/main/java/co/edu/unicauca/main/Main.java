package co.edu.unicauca.main;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;

import co.edu.unicauca.capaDeControladores.ControladorPreferenciasUsuariosIml;
import co.edu.unicauca.configuracion.lector.LectorPropiedadesConfig;
import co.edu.unicauca.configuracion.servicios.ServidorDeObjetos;

@SpringBootApplication
@ComponentScan(basePackages = {"co.edu.unicauca"})
public class Main {

    public static void main(String[] args) {
        SpringApplication.run(Main.class, args);
    }

    @Bean
    public ControladorPreferenciasUsuariosIml getControladorRMI() {
        try {
            // Leer configuración con valores por defecto
            int puertoNS = 2020;
            String direccionIPNS = "localhost";
            
            try {
                String puertoStr = LectorPropiedadesConfig.get("ns.port");
                if (puertoStr != null && !puertoStr.isEmpty()) {
                    puertoNS = Integer.parseInt(puertoStr);
                }
                
                String hostStr = LectorPropiedadesConfig.get("ns.host");
                if (hostStr != null && !hostStr.isEmpty()) {
                    direccionIPNS = hostStr;
                }
            } catch (Exception e) {
                System.out.println("⚠️  Usando valores por defecto para RMI: localhost:2020");
            }

            System.out.println("🚀 Iniciando servidor RMI en " + direccionIPNS + ":" + puertoNS);

            // Paso 1: arrancar o crear el ns
            ServidorDeObjetos.arrancarNS(direccionIPNS, puertoNS);
           
            // Paso 2: crear el objeto remoto
            ControladorPreferenciasUsuariosIml objControladorPreferencias = ServidorDeObjetos.crearObjetoRemoto();
            
            // Paso 3: registrar el objeto remoto en el ns
            String identificadorObjetoRemoto = "objControladorPreferenciasUsuarios";
            ServidorDeObjetos.registrarObjetoRemoto(objControladorPreferencias, direccionIPNS, puertoNS, identificadorObjetoRemoto);

            System.out.println("✅ Servidor RMI iniciado correctamente");
            return objControladorPreferencias;
        } catch (Exception e) {
            System.err.println("❌ Error al iniciar servidor RMI: " + e.getMessage());
            e.printStackTrace();
            throw new RuntimeException("No se pudo iniciar el servidor RMI", e);
        }
    }
}
