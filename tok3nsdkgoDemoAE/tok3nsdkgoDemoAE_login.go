package tok3nsdkgoDemoAE
import (
		"net/http"
		"fmt"
		)

func registerLoginFunctions(){
	http.HandleFunc("/login.do", loginDo)
}

func loginDo(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"try login")
}