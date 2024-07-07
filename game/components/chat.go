package components

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/teeworlds-go/protocol/messages7"
	"github.com/teeworlds-go/protocol/network7"
	"github.com/teeworlds-go/protocol/teeworlds7"
)

var (
	MplusFaceSource *text.GoTextFaceSource
	MplusNormalFace *text.GoTextFace
)

type ChatMsg struct {
	Name string
	Msg  *messages7.SvChat
}

type Chat struct {
	Messages []ChatMsg
}

func (c *Chat) OnInit() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	MplusFaceSource = s

	MplusNormalFace = &text.GoTextFace{
		Source: MplusFaceSource,
		Size:   14,
	}

}

func (c *Chat) RenderMsg(screen *ebiten.Image, offsetY int, msg *ChatMsg) {
	gray := color.RGBA{0x80, 0x80, 0x80, 0xff}
	const lineHeight = 20
	const x = 10
	y := float32(lineHeight * offsetY)
	line := fmt.Sprintf("%s: %s", msg.Name, msg.Msg.Message)
	w, h := text.Measure(line, MplusNormalFace, MplusNormalFace.Size*1.5)
	vector.DrawFilledRect(screen, x, y, float32(w), float32(h), gray, false)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = MplusNormalFace.Size * 1.5
	text.Draw(screen, line, MplusNormalFace, op)
}

func (c *Chat) OnRender(screen *ebiten.Image, client *teeworlds7.Client) {
	offset := 0
	const numDisplayLines = 6
	if len(c.Messages) > numDisplayLines+1 {
		offset = len(c.Messages) - numDisplayLines
	}
	for y, m := range c.Messages[offset:] {
		c.RenderMsg(screen, y, &m)
	}
}

func (c *Chat) OnChatMsg(msg *messages7.SvChat, client *teeworlds7.Client) {
	if msg.ClientId < 0 || msg.ClientId > network7.MaxClients {
		fmt.Printf("[chat] *** %s\n", msg.Message)
		return
	}
	name := client.Game.Players[msg.ClientId].Info.Name
	fmt.Printf("[chat] <%s> %s\n", name, msg.Message)

	c.Messages = append(c.Messages, ChatMsg{
		Name: name,
		Msg:  msg,
	})
}
