package tok3nsdkgoDemoAE
import (
	"net/http"
	"fmt"
	"time"
	"errors"
	"encoding/json"
	"appengine"
	"appengine/datastore"
	"github.com/gorilla/sessions"
	"github.com/Tok3n/tok3nsdkgo"
)

func registerWS(){
	http.HandleFunc("/ws/login",wsLogin)
	http.HandleFunc("/ws/create",wsCreateNewUser)
	http.HandleFunc("/ws/isloggedin",wsIsLoggedIn)
	http.HandleFunc("/ws/userHasTok3n",wsUserHasTok3n)
	http.HandleFunc("/ws/authenticateuser",wsAuthenticateUser)
	http.HandleFunc("/ws/getuser",wsGetUser)
	http.HandleFunc("/ws/setmonto",wsSetMonto)
	http.HandleFunc("/ws/getmonto",wsGetMonto)
	

	//banamex
	http.HandleFunc("/ws/setdata",wsSetData)
	http.HandleFunc("/ws/getdata",wsGetData)
}

func wsGetData(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Datos").Order("-Creation").Limit(1)
	var ms []Datos
	q.GetAll(c,&ms)
	jsonstring, _ := json.Marshal(ms[0])
	fmt.Fprintf(w,string(jsonstring))
}

func wsSetData(w http.ResponseWriter, r *http.Request){
	name := r.FormValue("name")
	telefono := r.FormValue("telefono")
	direccion := r.FormValue("direccion")

	var m Datos
	m.Nombre = name
	m.Telefono = telefono
	m.Direccion = direccion
	m.Creation = time.Now()
	c := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(c, "Datos", nil)
	datastore.Put(c,key,&m)
	fmt.Fprintf(w,"1")
}

func wsSetMonto(w http.ResponseWriter, r *http.Request){
	s := r.FormValue("monto")
	var m Monto
	m.Monto = s
	m.Creation = time.Now()
	c := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(c, "Monto", nil)
	datastore.Put(c,key,&m)
	fmt.Fprintf(w,"1")
}

func wsGetMonto(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Monto").Order("-Creation").Limit(1)
	var ms []Monto
	q.GetAll(c,&ms)
	fmt.Fprintf(w, ms[0].Monto)
}

func wsAuthenticateUser(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	q := r.FormValue("q")
	if q == ""{
		fmt.Fprintf(w,"ERROR: Error in parameters")
		return
	}
	var cert Cert
	json.Unmarshal([]byte(q),&cert)
	tok3nInstance := tok3nsdkgo.GetAppEngineTok3nInstance(c,tok3nConfig)
	resp, err := tok3nInstance.ValidateAuth(cert.UserKey,q,cert.TransactionId)
	if err != nil {
        fmt.Fprintf(w,"Error: Error in authentication")
       	return
    }

	
	if resp == "VALID USER"{
		user,key,err := getUserAndKey(r)
		if err != nil{
			 fmt.Fprintf(w,"Error: Error in authentication")
			 return
		}
		user.Tok3nKey = cert.UserKey
		datastore.Put(c,key,user)
	}

	fmt.Fprintf(w,"%s",resp)
}

func wsIsLoggedIn(w http.ResponseWriter, r *http.Request){
	session, _ := sessionsStore.Get(r, "logindata")
	id := session.Values["id"]
	if id == nil{
		fmt.Fprintf(w,"NO")
	}else{
		fmt.Fprintf(w,"YES")
	}
}

func wsLogin(w http.ResponseWriter, r *http.Request){
	user := r.FormValue("user")

	if user == "" {
		fmt.Fprintf(w,"ERROR: Error in parameters")
		return
	}
	c := appengine.NewContext(r)
	q := datastore.NewQuery("User").Filter("Username = ",user)// Yes we are not ussing passwords (very bad idea but for demo purposes is OK)
	found := false
	if count,_ := q.Count(c);count>0 { //exists at least one row 
		t := q.Run(c)
		var x User
        _, err := t.Next(&x)
        if err != nil {
            fmt.Fprintf(w,"Error: something happend with the database")
           	return
        }

        session, _ := sessionsStore.Get(r, "logindata")
        session.Options = &sessions.Options{ //just valid for half an hour
		    Path:   "/",
		    MaxAge: 60*30,
		}
		session.Values["id"] = x.Username
		session.Save(r, w)
		found = true
	}
	if found{
		fmt.Fprintf(w,"OK")
	}else{
		fmt.Fprintf(w,"NO")
	}
}

func wsCreateNewUser(w http.ResponseWriter, r *http.Request){
	user := r.FormValue("user")

	if user == "" {
		//empty values or variables not sended. Very basic prevention
		fmt.Fprintf(w, "ERROR: Error in arguments")
		return
	}

	c := appengine.NewContext(r)
	q := datastore.NewQuery("User").
		Filter("Username = ",user)
	if count,_ := q.Count(c);count>0 { //Verify the previous existence of the user
		fmt.Fprintf(w,"ERROR: User Exists")
		return
	}

	var theuser User
	theuser.Username = user
	theuser.Creation = time.Now()
	theuser.Tok3nKey = ""

	key := datastore.NewIncompleteKey(c, "User", nil)
	key ,err := datastore.Put(c,key,&theuser)
	if err != nil{
		fmt.Fprintf(w,"Error: %v",err)
		return
	}

	fmt.Fprintf(w, "%s",user)
}

func wsUserHasTok3n(w http.ResponseWriter, r *http.Request){
	session, _ := sessionsStore.Get(r, "logindata")
	id := session.Values["id"]
	if id == nil{
		fmt.Fprintf(w,"ERROR: user not loggedin")
		return 
	}
	c := appengine.NewContext(r)
	q := datastore.NewQuery("User").Filter("Username = ",id.(string))
	if count,_ := q.Count(c);count==0 { //Verify the previous existence of the user
		fmt.Fprintf(w,"ERROR: User Doesnt exist")
		return
	}
	var theusers []User
	q.GetAll(c,&theusers)
	if theusers[0].Tok3nKey == ""{
		fmt.Fprintf(w,"NO")
	}else{
		fmt.Fprintf(w,"YES")
	}

}

func wsGetUser(w http.ResponseWriter, r *http.Request){
	session, _ := sessionsStore.Get(r, "logindata")
	id := session.Values["id"]
	if id == nil{
		fmt.Fprintf(w,"ERROR: user not loggedin")
		return 
	}
	c := appengine.NewContext(r)
	q := datastore.NewQuery("User").Filter("Username = ",id.(string))
	if count,_ := q.Count(c);count==0 { //Verify the previous existence of the user
		fmt.Fprintf(w,"ERROR: User Doesnt exist")
		return
	}
	var theusers []User
	q.GetAll(c,&theusers)
	jsonstring,_ := json.Marshal(theusers[0])
	fmt.Fprintf(w,string(jsonstring))

}

func getUserAndKey(r *http.Request)(*User, *datastore.Key, error){
	session, _ := sessionsStore.Get(r, "logindata")
	id := session.Values["id"]
	if id == nil{
		return nil, nil, errors.New("no sesion")
	}
	c := appengine.NewContext(r)
	q := datastore.NewQuery("User").Filter("Username = ",id.(string))
	if count,_ := q.Count(c);count==0 { //Verify the previous existence of the user
		return nil, nil, errors.New("no user")
	}
	var theusers []User
	keys, err := q.GetAll(c,&theusers)
	return &theusers[0], keys[0], err
}