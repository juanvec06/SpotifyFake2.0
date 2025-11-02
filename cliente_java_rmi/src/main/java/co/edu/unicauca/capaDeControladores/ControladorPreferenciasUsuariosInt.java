package co.edu.unicauca.capaDeControladores;

import java.rmi.Remote;
import java.rmi.RemoteException;
import co.edu.unicauca.fachadaServices.DTO.PreferenciasDTORespuesta;

//Hereda de la clase Remote, lo cual la convierte en una interfaz remota
public interface ControladorPreferenciasUsuariosInt extends Remote {
    //Definicion del metodo remoto
    public PreferenciasDTORespuesta getReferencias(Integer id) throws RemoteException;
    //Cada definicion de metodo remoto debe lanzar la excepcion RemoteException
}
