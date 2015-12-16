part of tok3napi_dart;

class Tok3n_V2{
  
  //static String DOMAIN = "https://{[{Token_PlatformHost}]}";
  static String DOMAIN = "https://secure.tok3n.com";
  
  static Future<String> getTransactionId(publicKey,kind,userKey){
    Completer cmpl = new Completer();
    
    HttpRequest.getString("$DOMAIN/api/v2/transaction.new?publicKey=$publicKey&kind=$kind&userKey=$userKey").then((r){
      new Timer(const Duration(milliseconds:100),()=>cmpl.complete(r));
      
    });
    
    return cmpl.future;
  }
  
  static String getQR_URL(publicKey,transactionId,{forlogin:false}){
    var url = "$DOMAIN/api/v2/getQR?transaction=$transactionId&public_key=$publicKey";
    if (forlogin){
      url = "$url&forlogin=1";
    }
    return url;
  }
  
  static Future<String> QR_valid(){
    Completer cmpl = new Completer();
    
    HttpRequest.getString("$DOMAIN/api/v2/sqrValid", withCredentials:true).then((r)=>cmpl.complete(r));
    
    return cmpl.future;
    
  }
  
  static Future<String> validateOTPfromSMS(publicKey,userKey,otp,transaction){
    Completer cmpl = new Completer();
    
    var url="$DOMAIN/api/v2/push.sms.verify?publicKey=$publicKey&UserKey=$userKey&otp=$otp&transaction=$transaction";
    print(url);
    HttpRequest.getString(url).then((r)=>cmpl.complete(r));
    
    return cmpl.future;
  }
  
  static Future<String> emitPushProtocols(publicKey,userKey,transaction){
    Completer cmpl = new Completer();
    
    var url="$DOMAIN/api/v2/push.emit?publicKey=$publicKey&UserKey=$userKey&transaction=$transaction";
    print(url);
    HttpRequest.getString(url).then((r)=>cmpl.complete(r));
    
    return cmpl.future;
  }
  
  ///
}