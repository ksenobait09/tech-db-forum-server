package error

//easyjson:json
type Error struct {
	// Текстовое описание ошибки.
	// В процессе проверки API никаких проверок на содерижимое данного описание не делается.
	//
	// Read Only: true
	Message string `json:"message,omitempty"`
}
var DefaultMessage = []byte(`{"message": "Can't find user with id #42\n"}`)
