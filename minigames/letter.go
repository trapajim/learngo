package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	Refresh                = 60 * time.Millisecond
	Color                  = termbox.ColorCyan
	BorderColor            = termbox.ColorWhite
	DefaultColor           = termbox.ColorDefault
	GameTime               = 15
	Menu         GameState = iota
	Playing
	Quit
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

type GameState int

type Vector2d struct {
	x, y int
}

type World struct {
	width  int
	height int
}

type GameElement struct {
	coords Vector2d
	letter rune
}

type Context struct {
	gameEl    []GameElement
	world     World
	counter   int
	time      int
	ticker    *time.Ticker
	gameState GameState
}

func (v *GameElement) Draw(left, bottom int) {
	termbox.SetCell(left+v.coords.x, bottom-v.coords.y, v.letter, Color, DefaultColor)
}
func (a World) Draw(top, bottom, left int) {

	for i := top; i < bottom; i++ {
		termbox.SetCell(left-1, i, '│', BorderColor, DefaultColor)
		termbox.SetCell(left+a.width, i, '│', BorderColor, DefaultColor)
	}
	termbox.SetCell(left-1, top, '┌', BorderColor, DefaultColor)
	termbox.SetCell(left-1, bottom, '└', BorderColor, DefaultColor)
	termbox.SetCell(left+a.width, top, '┐', BorderColor, DefaultColor)
	termbox.SetCell(left+a.width, bottom, '┘', BorderColor, DefaultColor)

	fill(left, top, a.width, 1, termbox.Cell{Ch: '─'})
	fill(left, bottom, a.width, 1, termbox.Cell{Ch: '─'})
}

func (ctx *Context) Draw() {
	termbox.Clear(DefaultColor, DefaultColor)
	world := ctx.world
	var (
		w, h   = termbox.Size()
		midY   = h / 2
		left   = (w - world.width) / 2
		right  = (w + world.width) / 2
		top    = midY - (world.height / 2)
		bottom = midY + (world.height / 2) + 1
	)
	ctx.world.Draw(top, bottom, left)
	ctx.gameEl[len(ctx.gameEl)-1].Draw(left, bottom)

	counter := fmt.Sprintf("%d%s%d", ctx.counter, "/", len(ctx.gameEl))
	tbprint(left, bottom, DefaultColor, DefaultColor, counter)
	tbprint(right, bottom, DefaultColor, DefaultColor, strconv.Itoa(ctx.time))
	termbox.Flush()
}
func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, BorderColor, cell.Bg)
		}
	}
}

func NewGameElement(x, y int, letter rune) GameElement {
	return GameElement{
		coords: Vector2d{x: x, y: y},
		letter: letter,
	}
}

func NewContext() *Context {
	world := World{60, 30}
	return &Context{
		world:     world,
		counter:   0,
		gameState: Menu,
		time:      GameTime,
	}
}

func update(ctx *Context, events chan termbox.Event) {
	select {
	case ev := <-events:
		if ev.Type == termbox.EventKey {
			switch {
			case ev.Ch == ctx.gameEl[len(ctx.gameEl)-1].letter:
				ctx.addgameEl()
				ctx.counter++
			case ev.Key == termbox.KeyEsc:
				ctx.gameState = Quit
			case ev.Ch != ctx.gameEl[len(ctx.gameEl)-1].letter:
				ctx.addgameEl()
				ctx.counter--
			}
		}
	default:
		ctx.Draw()
		time.Sleep(Refresh)
	}
}

func (ctx *Context) startTimer() {
	ctx.time = GameTime
	ctx.ticker = time.NewTicker(1 * time.Second)
	go func() {
		for range ctx.ticker.C {
			ctx.time--
			if ctx.time <= 0 {
				ctx.ticker.Stop()
				ctx.gameState = Menu
			}
		}
	}()
}

func drawMenu(ctx *Context, events chan termbox.Event) {

	select {
	case ev := <-events:
		if ev.Type == termbox.EventKey {
			switch {
			case ev.Key == termbox.KeyEsc:
				ctx.gameState = Quit
			default:
				ctx.gameState = Playing
				ctx.counter = 0
				ctx.gameEl = nil
				ctx.startTimer()
				ctx.addgameEl()
			}
		}
	default:
		world := ctx.world
		var (
			w, h   = termbox.Size()
			midY   = h / 2
			left   = (w - world.width) / 2
			bottom = midY + (world.height / 2) + 1
		)
		termbox.Clear(DefaultColor, DefaultColor)
		if len(ctx.gameEl) > 0 {
			finalScoreText := fmt.Sprintf("%s%d%s%d%s", "You've got ", ctx.counter, " of ", len(ctx.gameEl), " correct")
			tbprint(left+7, midY-3, DefaultColor, DefaultColor, "Game over!")
			tbprint(left, midY, DefaultColor, DefaultColor, finalScoreText)
		}
		tbprint(left, bottom, DefaultColor, DefaultColor, "Press any button to start")
		termbox.Flush()
		time.Sleep(Refresh)
	}
}

func (ctx *Context) addgameEl() {
	ctx.gameEl = append(ctx.gameEl, NewGameElement(rand.Intn(ctx.world.width-1), rand.Intn(ctx.world.height-1), randStringRunes()))
}

func randStringRunes() rune {
	return letters[rand.Intn(len(letters))]
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()

	ctx := NewContext()
	rand.Seed(time.Now().Unix())
	for ctx.gameState != Quit {
		if ctx.gameState == Playing {
			update(ctx, events)
		} else {
			drawMenu(ctx, events)
		}
	}

}
