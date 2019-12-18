package message

var mes = *new([]Message)

type RegisterMessage struct {
}

func (r *RegisterMessage) Register(message Message) {
	mes = append(mes, message)
}
func (r *RegisterMessage) SendMessage(payload interface{}) {

	if len(mes) > 0 {
		for _, m := range mes {
			m.Send(payload)
		}
	}
}
