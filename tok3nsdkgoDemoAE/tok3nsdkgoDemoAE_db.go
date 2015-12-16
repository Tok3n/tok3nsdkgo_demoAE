package tok3nsdkgoDemoAE
import (
		"time"
		)

type User struct {
    Creation time.Time
    Username string
    //With these key Tok3n identify the user in your system
    Tok3nKey string 
}

type Cert struct{
    Valid string
    CertKey string
    PublicKey string
    UserKey string
    TransactionId string
    Hash string
    Result string
}

type Monto struct{
    Monto string
    Creation time.Time
}

type Datos struct{
    Nombre string
    Direccion string
    Telefono string
    Creation time.Time
}