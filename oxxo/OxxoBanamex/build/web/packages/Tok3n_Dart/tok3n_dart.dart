library tok3n_dart;

import 'dart:html';
import 'dart:async';
import 'dart:js' as js;
import 'package:dart_node_tree_sanitizers/dart_node_tree_sanitizers.dart';

class Tok3nClient{
  String publickey;
  ScriptElement element;
  Element container;
  Completer responseC;
  Future get response => responseC.future;
  

  Tok3nClient(this.publickey,this.container);
  ScriptElement addUser(){
    responseC = new Completer();
    element = new Element.html("""
<script src="//secure.tok3n.com/api_v2_iframe/tok3n.js" 
  data-tok3n-integration
  data-render-tag
  data-tag-class-name="btn btn-success btn-lrg btn-block"
  action="authorize"
  public-key="$publickey">
</script>
""",treeSanitizer:new NullTreeSanitizer());
    element.onLoad.listen((_){
      querySelector("#tok3n-iframe").onLoad.listen((_){
        print("antes de llamar");
        js.context["Tok3n"].callMethod("showIFrame",["authorize",publickey,null]);
        print("despues de llamar");
      });
      querySelector("#tok3n-authenticate").style.display="none";
    });
    element.on["response"].listen((r)=>responseC.complete(r));
    container.nodes.add(element);
    return element;
  }
  
  ScriptElement askTok3n(String userkey){
      responseC = new Completer();
      element = new Element.html("""
<script src="//secure.tok3n.com/api_v2_iframe/tok3n.js" 
  data-tok3n-integration
  data-render-tag
  data-tag-class-name="btn btn-success btn-lrg btn-block"
  action="authenticate"
  public-key="$publickey"
  user-key="$userkey">
</script>
""",treeSanitizer:new NullTreeSanitizer());
      element.onLoad.listen((_){
        querySelector("#tok3n-iframe").onLoad.listen((_){
          print("antes de llamar");
          js.context["Tok3n"].callMethod("showIFrame",["authenticate",publickey,userkey]);
          print("despues de llamar");
        });
        querySelector("#tok3n-authenticate").style.display="none";
      });
      element.on["response"].listen((r)=>responseC.complete(r));
      container.nodes.add(element);
      return element;
    }
  
  void remove(){
    container.nodes.remove(element);
  }
  
  responded(r){
    print(r);
  }
}