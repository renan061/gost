# gost

O pacote `gost` (*go server template*) possui uma série de interfaces que parametrizam blocos básicos utilizados na criação de REST APIs em Go. As interfaces podem ser implementadas por você conforme suas necessidades específicas ou, caso um comportamento padrão seja suficiente, você pode utilizar as implementações contidas no pacote `gosti` (*gost implementations*).

A pasta `exs` apresenta um conjunto de chamadas exemplificando a utilização do `gost` com as implementações default do `gosti`.

### Authenticator
Um `gost.Authenticator` obtém as informações relevantes para a autenticação da requisição e valida as mesmas mediante um conjunto de informações de autenticação (`AuthInfo`). O Authenticator retorna um conjunto de informações definidas como `Claims` obtidas a partir da autenticação e um booleano para indicar a ocorrência de erros (os erros devem ser tratados internamente usando o `http.ResponseWriter` passado).
```
type AuthInfo interface{}

type Claims interface{}

type Authenticator interface {
	Authenticate(http.ResponseWriter, *http.Request, AuthInfo) (Claims, bool)
}
```

### Decoder
Um `gost.Decoder` lê as informações do corpo da requisição e preenche um objeto `RequestBody` com as mesmas. Um objeto que satisfaz a interface *RequestBody* deve implementar uma função que identifica o objeto como válido. Essa função é chamada durante o processo de decodificação. O *Decoder* retorna um booleano para indicar a ocorrência de erros (os erros devem ser tratados internamente usando o `http.ResponseWriter` passado).
```
type RequestBody interface {
	Valid() error
}

type Decoder interface {
	Decode(http.ResponseWriter, *http.Request, RequestBody) bool
}
```

### Responder
Um `gost.Responder` recebe algum tipo de resposta parametrizada e escreve a mesma no `http.ResponseWriter` fornecido. O *Responder* retorna um booleano para indicar a ocorrência de erros (os erros devem ser tratados internamente).
```
type Response interface{}

type Responder interface {
	Respond(http.ResponseWriter, Response) bool
}
```
