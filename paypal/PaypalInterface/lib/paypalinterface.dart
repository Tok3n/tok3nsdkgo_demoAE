library PaypalInterface;
import "dart:html";
import "dart:async";
import "dart:convert";
import "package:dart_node_tree_sanitizers/dart_node_tree_sanitizers.dart";
import "package:Tok3n_Dart/tok3n_dart.dart";
import "package:intl/intl.dart";
import "package:Tok3nAPI_Dart/tok3napi_dart.dart";

part "paypalcalculator.dart";

class PartialLoader{
  static Future<String> _getLoginPartial(){
    return HttpRequest.getString("partials/loginpartial.html");
  }
  static Future<Element> getLoginPartial(Element ele){
    Completer<Element> cmpl = new Completer<Element>();
    _getLoginPartial().then((r){
      ele.nodes.clear();
      Element e = new Element.html(r,treeSanitizer: new NullTreeSanitizer());
      ele.nodes.add(e);
      cmpl.complete(ele);
    });
    return cmpl.future;
  }
}

class PayPalManager{
  PayPalManager(){
    HttpRequest.getString("/ws/isloggedin").then((String resp){
      if (resp == "YES"){
        Tok3nVinculation();
        print("logged in");
      }else{
        new LoginWindow();
      }
    });
    
  }
  Tok3nVinculation(){
    HttpRequest.getString("/ws/userHasTok3n").then((String resp){
      if (resp == "YES"){
        var tok3n = new ShowTok3nWindow();
        tok3n.done.then(showCalculator);
      }else{
        new MergeTok3nWindow();
      }
    });
  }
  
  void showCalculator(_){
    new Calculator();
  }
  
}

class LoginWindow{
  LoginWindow(){
    print("init loginwindow");
    PartialLoader.getLoginPartial(querySelector("#container")).then((Element ele){
      ele.querySelector("#loginbutton").onClick.listen((_)=>loginAccount(ele));
      ele.querySelector("#createbutton").onClick.listen((_)=>createAccount(ele));
    });
  }
  void createAccount(ele){
    InputElement e = ele.querySelector("#username");
    if (e.value.length>3){
      HttpRequest.getString("/ws/create?user=${e.value}").then((String response){
        if (response.startsWith("ERROR")){
          window.alert(response);
        }else{
          loginAccount(ele);
        }
      });
    }
  }
  
  void loginAccount(ele){
    InputElement e = ele.querySelector("#username");
    if (e.value.length>3){
      HttpRequest.getString("/ws/login?user=${e.value}").then(loginAccount_response);
    }
  }
  void loginAccount_response(response){
    if (response == "OK"){
      print("loggedin");
      window.location.reload();
    }else{
      print("not loggedin");
    }
  }
}

class ShowTok3nWindow{
  Completer cmpl;
  Future get done => cmpl.future;
  Tok3nClient client;
  ScriptElement  script;
  ShowTok3nWindow(){
    print("login with tok3n");
    cmpl = new Completer();
    HttpRequest.getString("/ws/getuser").then(inited);
    
  }
  void inited(response){
    var data = JSON.decode(response);
    
    client = new Tok3nClient("4c186122-ed23-5cf5-45c2-2e1a17639efa",querySelector("body"));
    script = client.askTok3n(data["Tok3nKey"]);
    client.response.then(complete);
  }
  void complete(response){
    var data = JSON.decode(response.detail);
    print(data);
    
    if (data["Valid"]=="YES"){
      HttpRequest.getString("/ws/authenticateuser?q=${Uri.encodeComponent(response.detail)}").then((_){
        client.remove();
        cmpl.complete(response);
      });
    }
    
  }
}

class MergeTok3nWindow{
  Tok3nClient client;
  ScriptElement script;
  MergeTok3nWindow(){
    print("start vinculation of tok3n");
    client = new Tok3nClient("4c186122-ed23-5cf5-45c2-2e1a17639efa",querySelector("body"));
    script = client.addUser();
    client.response.then(responded);
  }
  void responded(CustomEvent response){
    var data = JSON.decode(response.detail);
    print(data);
    
    if (data["Valid"]=="YES"){
      HttpRequest.getString("/ws/authenticateuser?q=${Uri.encodeComponent(response.detail)}").then((_)=>window.location.reload());
    }
  }
}