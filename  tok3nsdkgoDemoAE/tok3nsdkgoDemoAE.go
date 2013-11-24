package  tok3nsdkgoDemoAE

import (
	"net/http"
	"fmt"

	 
	"github.com/gorilla/sessions"
)

var sessionsStore = sessions.NewCookieStore([]byte("What ever you feel secure"))

func secureWebAccess(w http.ResponseWriter, r *http.Request) *UserWeb{
	session, _ := sessionsStore.Get(r, "logindata")
	value := session.Values["id"]
	
	var u UserWeb

    if (value == nil){
    	http.Redirect(w, r, "/login.do", http.StatusTemporaryRedirect)
		return
	}else{
		c := appengine.NewContext(r)
		k := datastore.NewKey(c,"UserWeb","",value.(int64),nil)
		err := datastore.Get(c,k,&u)
		if err!= nil{
			//what ever you have in the session is not in the datastore, these happends when you delete a user directly from the datastore, and the sesion is still active
			http.Redirect(w, r, "/login.do", http.StatusTemporaryRedirect)
			return
		}
    }

    return &u
}

func init() {
	http.HandleFunc("/_ah/warmup",warmup_method) //usefull for initing changes betwen versions
    http.HandleFunc("/", rootWS)

    registerLoginFunctions()
}

/**
Do nothing just print ok
**/
func warmup_method(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w,"OK")
}

func rootWS(w http.ResponseWriter, r *http.Request) {
	secureWebAccess(w,r)
}