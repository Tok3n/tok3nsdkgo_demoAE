package tok3nsdkgoDemoAE
import (
		"net/http"
		"fmt"
		"strings"

		"appengine"
		"appengine/datastore"

		"github.com/gorilla/sessions"
		"github.com/Tok3n/tok3nsdkgo"
		)

func registerTok3nFunctions(){
	http.HandleFunc("/tok3ncallback", tok3ncallback)
	http.HandleFunc("/tok3nreturn", tok3nreturn)
	http.HandleFunc("/login.tok3n", askForTok3n)
	http.HandleFunc("/login.tok3nverify", tok3nverify)
}


func tok3ncallback(w http.ResponseWriter, r *http.Request){
	data := r.FormValue("data")
	userkey := r.FormValue("userkey")
	event := r.FormValue("event")
	secret := r.FormValue("secret")

	c := appengine.NewContext(r)
	if data == "" ||userkey == "" ||event == "" ||secret == "" {
		c.Infof("Error: parameters")
		fmt.Fprintf(w,"Error: parameters")
		return
	}
	if secret != tok3nConfig.SecretKey{
		c.Infof("Error: not valid secret key, sorry")
		fmt.Fprintf(w,"Error: not valid secret key, sorry")
		return
	}

	if event == "userAdded"{
		
	
		q := datastore.NewQuery("User").
			Filter("Username = ",data)// Yes we are not ussing passwords (very bad idea but for demo purposes is OK)

		if count,_ := q.Count(c);count>0 { //exists at least one row 
			c.Infof("User finded: %s",data)
			t := q.Run(c)
			var x User
	        key, err := t.Next(&x)
	        if err != nil {
	        	c.Infof("Error: something happend with the database")
	            fmt.Fprintf(w,"Error: something happend with the database")
	           	return
	        }
	        x.Tok3nKey = userkey
	        datastore.Put(c,key,&x)
	        fmt.Fprintf(w, "OK")
		}else{
			fmt.Fprintf(w,"Error: invalid data")
			return
		}
	}
}

func tok3nreturn(w http.ResponseWriter, r *http.Request){
	data := r.FormValue("callbackdata")
	key := r.FormValue("key")

	if data=="" || key == ""{
		fmt.Fprintf(w, "Error: Parameters Error")
		return
	}

	c := appengine.NewContext(r)
	
	q := datastore.NewQuery("User").
		Filter("Username = ",data).
		Filter("Tok3nKey = ",key)

	if count,_ := q.Count(c);count>0 {
		session, _ := sessionsStore.Get(r, "logindata")
        session.Options = &sessions.Options{ //just valid for half an hour
		    Path:   "/",
		    MaxAge: 60*30,
		}
		session.Values["tok3niced"] = "true"
		session.Save(r, w)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}else{
		fmt.Fprintf(w, "Error: parameters data")
		return
	}
}
func askForTok3n(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w,"hello")
	usr := secureWebAccess(w,r)
	if usr.Tok3nKey == ""{
		fmt.Fprintf(w,"There are session errors reset your cookies or wait an hour")
		return
	}

	c := appengine.NewContext(r)
	tok3nInstance := tok3nsdkgo.GetAppEngineTok3nInstance(c,tok3nConfig)

	s := []string{"<html><form action=\"/login.tok3nverify\"><div id=\"tok3n_placeholder\"></div><script language=\"javascript\" src='", tok3nInstance.GetJsClientUrl__v1_5("Login",usr.Tok3nKey), "' ></script></form></html>"}
	fmt.Fprint(w, strings.Join(s,""))
}

func tok3nverify(w http.ResponseWriter, r *http.Request){
	usr := secureWebAccess(w,r)
	otp := r.FormValue("tok3n_otp_field")
	session := r.FormValue("tok3n_sesion")
	sqr := r.FormValue("tok3n_sqr")

	if session=="" || (otp=="" && sqr==""){
		fmt.Fprintf(w, "Error: invalid parameters")
		return
	}

	c := appengine.NewContext(r)
	tok3nInstance := tok3nsdkgo.GetAppEngineTok3nInstance(c,tok3nConfig)

	if otp != ""{
		response,err := tok3nInstance.ValidateOTP(usr.Tok3nKey, otp, session)
		if err!=nil{
			fmt.Fprintf(w,"%s",err)
			return
		}
		fmt.Fprintf(w, response)
	}else if sqr != ""{
		response,err := tok3nInstance.ValidateSqr(usr.Tok3nKey, sqr, session)
		if err!=nil{
			fmt.Fprintf(w,"%s",err)
			return
		}
		fmt.Fprintf(w, response)
	}

	

	
}

