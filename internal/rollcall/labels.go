package rollcall

import (
	"fmt"
	"strings"
)

func title() []byte {
	const (
		msg     = "Sistema de Chamada"
		msgLen  = len(msg)
		msgHalf = msgLen / 2
		msgSub  = msgLen + (msgHalf / 2)
		msgPart = msgLen + msgHalf
	)

	var title string
	title += strings.Repeat("#", msgPart) + "\n"
	title += fmt.Sprintf("%*s\n", msgSub, msg)
	title += strings.Repeat("#", msgPart) + "\n"

	return []byte(title)
}
