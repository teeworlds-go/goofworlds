package game

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/teeworlds-go/goofworlds/game/components"
	"github.com/teeworlds-go/protocol/messages7"
	"github.com/teeworlds-go/protocol/teeworlds7"
)

const (
	ScreenWidth  = 960
	ScreenHeight = 540
)

var (
	TeeSprite *ebiten.Image

	MplusFaceSource *text.GoTextFaceSource
	MplusNormalFace *text.GoTextFace
)

type Camera struct {
	X int
	Y int
}

type CameraOffset struct {
	X int
	Y int
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func init() {
	// Preload images
	img, err := getImageFromFilePath("img/tee.png")
	if err != nil {
		panic(err)
	}
	TeeSprite = ebiten.NewImageFromImage(img)

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	MplusFaceSource = s

	MplusNormalFace = &text.GoTextFace{
		Source: MplusFaceSource,
		Size:   24,
	}
}

func getCameraOffset(camera Camera) CameraOffset {
	wc := ScreenWidth / 2
	hc := ScreenHeight / 2
	x := -camera.X + wc
	y := -camera.Y + hc
	return CameraOffset{
		X: x,
		Y: y,
	}
}

type TextInput struct {
	runes          []rune
	text           string
	counter        int
	submitCallback func(line string)
	active         bool
}

func (inp *TextInput) IsActive() bool {
	return inp.active
}

func (inp *TextInput) Activate() {
	inp.active = true
}

func (inp *TextInput) Deactivate() {
	inp.active = false
}

func NewTextInput(submitCallback func(line string)) *TextInput {
	return &TextInput{
		submitCallback: submitCallback,
	}
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func (inp *TextInput) Clear() {
	inp.text = ""
}

func (inp *TextInput) Draw(screen *ebiten.Image) {
	if inp.IsActive() == false {
		return
	}
	t := inp.text
	if inp.counter%60 < 30 {
		t += "_"
	}
	ebitenutil.DebugPrintAt(screen, t, 10, ScreenHeight-20)
}

func (inp *TextInput) Update() error {
	if inp.IsActive() == false {
		return nil
	}
	// Add runes that are input by the user by AppendInputChars.
	// Note that AppendInputChars result changes every frame, so you need to call this
	// every frame.
	inp.runes = ebiten.AppendInputChars(inp.runes[:0])
	inp.text += string(inp.runes)

	// Adjust the string to be at most 10 lines.
	ss := strings.Split(inp.text, "\n")
	if len(ss) > 10 {
		inp.text = strings.Join(ss[len(ss)-10:], "\n")
	}

	// If the enter key is pressed, add a line break.
	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		inp.submitCallback(inp.text)
		inp.text = ""
	}

	// If the backspace key is pressed, remove one character.
	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(inp.text) >= 1 {
			inp.text = inp.text[:len(inp.text)-1]
		}
	}

	inp.counter++
	return nil
}

type Game struct {
	Client     teeworlds7.Client
	Camera     Camera
	Fullscreen bool
	Ip         string
	Port       int
	Components []components.Component
	ChatInp    *TextInput
}

func (g *Game) Update() error {
	g.Client.Game.Input.Direction = 0
	g.Client.Game.Input.Jump = 0
	g.Client.Game.Input.Hook = 0
	g.Client.Game.Input.Fire = 0
	g.Client.Game.Input.PrevWeapon = 0
	g.Client.Game.Input.NextWeapon = 0
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Client.Game.Input.Direction = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Client.Game.Input.Direction = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Client.Game.Input.Jump = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) && g.ChatInp.IsActive() == false {
		g.Client.SendMessage(&messages7.CtrlClose{})
		os.Exit(0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyT) {
		g.ChatInp.Activate()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.Fullscreen = !g.Fullscreen
		ebiten.SetFullscreen(g.Fullscreen)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.Client.Game.Input.Fire++
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.Client.Game.Input.Hook = 1
	}
	_, dy := ebiten.Wheel()
	// this is super cursed fast weapon switch
	if dy > 0 {
		g.Client.Game.Input.PrevWeapon = 1
	}
	if dy < 0 {
		g.Client.Game.Input.NextWeapon = 1
	}

	g.ChatInp.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, character := range g.Client.Game.Snap.Characters {
		offset := getCameraOffset(g.Camera)

		screenX := float32(character.X) + float32(offset.X)
		screenY := float32(character.Y) + float32(offset.Y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.069, 0.069)
		op.GeoM.Translate(float64(screenX)-32, float64(screenY)-32)
		screen.DrawImage(TeeSprite, op)

		name := g.Client.Game.Players[character.Id()].Info.Name
		gray := color.RGBA{0x80, 0x80, 0x80, 0xff}

		{
			x := screenX - 64
			y := screenY - 64
			w, h := text.Measure(name, MplusNormalFace, MplusNormalFace.Size*1.5)
			vector.DrawFilledRect(screen, x, y, float32(w), float32(h), gray, false)
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			op.LineSpacing = MplusNormalFace.Size * 1.5
			text.Draw(screen, name, MplusNormalFace, op)
		}
	}

	for _, c := range g.Components {
		c.OnRender(screen, &g.Client)
	}

	g.ChatInp.Draw(screen)

	x, y := ebiten.CursorPosition()
	aimX := x
	aimY := y
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d  AimX: %d, AimY: %d", x, y, aimX, aimY))

	g.Client.Game.Input.TargetX = aimX - ScreenWidth/2
	g.Client.Game.Input.TargetY = aimY - ScreenHeight/2

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
