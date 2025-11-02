
package co.edu.unicauca.capaDeControladores;

import java.rmi.RemoteException;
import java.rmi.server.RMIClientSocketFactory;
import java.rmi.server.RMIServerSocketFactory;
import java.rmi.server.UnicastRemoteObject;
import co.edu.unicauca.fachadaServices.DTO.PreferenciasDTORespuesta;
import co.edu.unicauca.fachadaServices.services.IPreferenciasService;
import co.edu.unicauca.fachadaServices.services.PreferenciasServiceImpl;


public class ControladorPreferenciasUsuariosIml  extends UnicastRemoteObject implements ControladorPreferenciasUsuariosInt{

    private IPreferenciasService ServicioFachadaPreferencias;

    public ControladorPreferenciasUsuariosIml(IPreferenciasService servicioFachadaPreferencias) throws RemoteException {
        super();
        this.ServicioFachadaPreferencias = new PreferenciasServiceImpl();
    }

    @Override
    public PreferenciasDTORespuesta getReferencias(Integer id) throws RemoteException {
        return this.ServicioFachadaPreferencias.getReferencias(id);
    }
}


