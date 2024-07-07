package main

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
	"github.com/teeworlds-go/protocol/snapshot7"
	"github.com/teeworlds-go/protocol/teeworlds7"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	teeSprite *ebiten.Image

	mplusFaceSource *text.GoTextFaceSource
	mplusNormalFace *text.GoTextFace
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
	teeSprite = ebiten.NewImageFromImage(img)

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	mplusNormalFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}
}

func getCameraOffset(camera Camera) CameraOffset {
	wc := screenWidth / 2
	hc := screenHeight / 2
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
}

func (g *Game) Update() error {
	g.Client.Game.Input.Direction = 0
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Client.Game.Input.Direction = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Client.Game.Input.Direction = 1
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
		screen.DrawImage(teeSprite, op)

		name := g.Client.Game.Players[character.Id()].Info.Name
		gray := color.RGBA{0x80, 0x80, 0x80, 0xff}

		{
			x := screenX - 64
			y := screenY - 64
			w, h := text.Measure(name, mplusNormalFace, mplusNormalFace.Size*1.5)
			vector.DrawFilledRect(screen, x, y, float32(w), float32(h), gray, false)
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			op.LineSpacing = mplusNormalFace.Size * 1.5
			text.Draw(screen, name, mplusNormalFace, op)
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("goofworlds")

	game := &Game{}
	game.Client = *teeworlds7.NewClient()
	game.Client.Name = "goofy"

	game.Client.OnSnapshot(func(snap *snapshot7.Snapshot, defaultAction teeworlds7.DefaultAction) {
		char, err := game.Client.SnapFindCharacter(game.Client.LocalClientId)
		if err == nil {
			game.Camera.X = char.X
			game.Camera.Y = char.Y
		}
	})

	go func() {
		game.Client.Connect("127.0.0.1", 8303)
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
