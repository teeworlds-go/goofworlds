package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/teeworlds-go/protocol/messages7"
	"github.com/teeworlds-go/protocol/teeworlds7"
)

type Component interface {
	OnInit()
	OnRender(screen *ebiten.Image, client *teeworlds7.Client)
	OnChatMsg(msg *messages7.SvChat, client *teeworlds7.Client)
}

