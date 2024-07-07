package game

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/teeworlds-go/protocol/messages7"
	"github.com/teeworlds-go/protocol/teeworlds7"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
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

type Game struct {
	Client     teeworlds7.Client
	Camera     Camera
	Fullscreen bool
	Ip         string
	Port       int
}

func (g *Game) Update() error {
	g.Client.Game.Input.Direction = 0
	g.Client.Game.Input.Jump = 0
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Client.Game.Input.Direction = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Client.Game.Input.Direction = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Client.Game.Input.Jump = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.Client.SendMessage(&messages7.CtrlClose{})
		os.Exit(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.Fullscreen = !g.Fullscreen
		ebiten.SetFullscreen(g.Fullscreen)
	}

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

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
