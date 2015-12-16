part of OxxoBanamex;

class OxxoLoader{
  ImageElement img;
  OxxoLoader(){
    var container = querySelector("body");
    var div = new Element.html("""
<div>
  <table>
    <tr>
      <td>
        Escanear el codigo para sincronizar
      </td>
    </tr>
    <tr>
      <td>
        <img id="imagen" />
      </td>
    </tr>
  </table>
</div>
""");
    container.nodes.clear();
    container.nodes.add(div);
    img = div.querySelector("#imagen");
    HttpRequest.getString("/ws/getuser").then(initted);
  }
  void initted(response){
    var data = JSON.decode(response);
    SQRL sqrl = new SQRL("d550667e-a85b-5e95-5a93-cd0645853f6a",data["Tok3nKey"],"auth");
    sqrl.init_done.then((_){
      
      img.width=280;
      img.src = sqrl.getQR_URL(ForLogin:false);
      /*new Timer(const Duration(seconds: 1), (){
        img.src = sqrl.getQR_URL(ForLogin:forlogin);
      });*/
      
      
    });
    sqrl.qr_conected.then(validatedQR);
  }
  
  void validatedQR(response){
    new ShowData();
  }
}

class ShowData{
  ShowData(){
    HttpRequest.getString("/ws/getdata").then(dataLoaded);
  }
  void dataLoaded(String result){
    print(result);
    var data = JSON.decode(result);
    var container = querySelector("body");
    var div = new Element.html("""
<div>
  <table>

    <tr>
      <td>
        Nombre:
      </td>
      <td>
        ${data["Nombre"]}
      </td>
    </tr>

    <tr>
      <td>
        Direccion:
      </td>
      <td>
        ${data["Direccion"]}
      </td>
    </tr>

    <tr>
      <td>
        Telefono:
      </td>
      <td>
        ${data["Telefono"]}
      </td>
    </tr>

    <tr>
      <td>
        <input type="button" value="Siguiente..." id="siguiente" />
      </td>
    </tr>
  </table>
</div>
""");
    container.nodes.clear();
    container.nodes.add(div);
    
    div.querySelector("#siguiente").onClick.listen(siguienteclick);
  }
  void siguienteclick(_){
    new OxxoLoader();
  }
}