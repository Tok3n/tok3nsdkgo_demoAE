part of tok3napi_dart;

class SQRL{
  bool ready = false;
  String publicKey,userKey,kind,transactionID;
  
  Future get init_done => init_doneC.future;
  Future get qr_conected => qr_conectedC.future;
  Completer init_doneC,qr_conectedC;
  
  bool cancelQR;
  
  SQRL(this.publicKey,this.userKey,this.kind){
    init_doneC = new Completer();
    qr_conectedC = new Completer();
    Tok3n_V2.getTransactionId(publicKey, kind, userKey).then(init_response);
  }
  void init_response(String response){
    if (response.startsWith("ERRROR")){
      init_doneC.completeError(response);
    }else{
      transactionID=response;
      init_doneC.complete(response);
    }
  }
  
  String getQR_URL({ForLogin:false}){
    print("the QR is been served");
    startQRRoutine();
    print("the QR has been started scanned");
    cancelQR = false;
    return Tok3n_V2.getQR_URL(publicKey, transactionID,forlogin:ForLogin);
  }
  
  
  void startQRRoutine(){
    print("previous timer");
    new Timer(const Duration(seconds: 5), getDataFromQR);
    print("after timer");
  }
  
  void getDataFromQR(){
    print("validatingQR");
    Tok3n_V2.QR_valid().then((response){
      Map d = JSON.decode(response);
      if (!cancelQR ){
        print("not valid yet");
        if (d["Valid"]=="NO"){
          startQRRoutine();
        }else{
          print("VALID");
          cancelQR = true;
          qr_conectedC.complete(response);
        }
      }
    });
  }
}