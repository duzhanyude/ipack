package message

import (
	"ipack/com.ipack/constant"
)

var mes = *new([]Message)

type RegisterMessage struct {
}

func (r *RegisterMessage) Register(message Message) {
	mes = append(mes, message)
}
func (r *RegisterMessage) SendMessage(packDefine constant.PackDefine) {

	if len(mes) > 0 {
		for _, m := range mes {
			m.Send(packDefine)
		}
	}
}
