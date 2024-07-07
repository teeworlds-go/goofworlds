package components

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/teeworlds-go/protocol/messages7"
	"github.com/teeworlds-go/protocol/network7"
	"github.com/teeworlds-go/protocol/teeworlds7"
)

type Chat struct {
}

func (c *Chat) OnInit() {
}

func (c *Chat) OnRender(screen *ebiten.Image, client *teeworlds7.Client) {
}

func (c *Chat) OnChatMsg(msg *messages7.SvChat, client *teeworlds7.Client) {
	if msg.ClientId < 0 || msg.ClientId > network7.MaxClients {
		fmt.Printf("[chat] *** %s\n", msg.Message)
		return
	}
	name := client.Game.Players[msg.ClientId].Info.Name
	fmt.Printf("[chat] <%s> %s\n", name, msg.Message)
}

