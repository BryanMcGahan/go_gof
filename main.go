package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	GameActiveState bool      `json:"GameActiveState"`
	Generation      int       `json:"Generation"`
	BoardWidth      int32     `json:"BoardWith"`
	BoardHeight     int32     `json:"BoardHeight"`
	Cells           [][]*Cell `json:"Cells"`
}

type Cell struct {
	Position  rl.Vector2 `json:"Position"`
	Size      rl.Vector2 `json:"Size"`
	LifeState bool       `json:"LifeState"`
	Next      bool       `json:"Next"`
}

const (
	padding      = 20
	gridLineSize = 1
)

var screenWidth int32 = 700
var screenHeight int32 = 700
var previousScreenHeight int32
var previousScreenWidth int32
var cellSize int32 = 10
var numOfRows int32 = 10
var numOfCols int32 = 10

var backgroundColor rl.Color = rl.Color{R: 37, G: 40, B: 61, A: 255}
var gridColor rl.Color = rl.Color{R: 218, G: 210, B: 216, A: 75}
var liveColor rl.Color = rl.Color{R: 143, G: 57, B: 133, A: 255}

func main() {

	var game Game
	args := os.Args[1:]

	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&game); err != nil {
			log.Fatal(err)
		}
	} else {
		game.Init(false)
	}

	rl.SetConfigFlags(
		rl.FlagWindowHighdpi,
	)

	rl.InitWindow(screenWidth, screenHeight, "Game of Life")
	defer rl.CloseWindow()

	rl.SetTargetFPS(30)

	numOfRows = (screenHeight - (padding * 2)) / cellSize
	numOfCols = (screenWidth - (padding * 2)) / cellSize

	for !rl.WindowShouldClose() {

		previousScreenWidth = screenWidth
		previousScreenHeight = screenHeight

		screenWidth = int32(rl.GetScreenWidth())
		screenHeight = int32(rl.GetScreenHeight())

		if screenHeight != previousScreenHeight || screenWidth != previousScreenWidth {
			numOfRows = (screenHeight - (padding * 2)) / cellSize
			numOfCols = (screenWidth - (padding * 2)) / cellSize
			game.Cells = getCells(false)
		}

		if game.GameActiveState {
			game.Update()
		}

		if game.Generation > 10000 {
			numOfRows = (screenHeight - (padding * 2)) / cellSize
			numOfCols = (screenWidth - (padding * 2)) / cellSize
			game.Cells = getCells(false)
			game.Generation = 0
		}

		game.Input()
		game.Draw()
	}
}

func (game *Game) Init(empty bool) {
	game.Cells = getCells(empty)
	game.GameActiveState = false
	game.Generation = 1
}

func (game *Game) CheckClick(x, y int32) {
	for i := 0; i < len(game.Cells); i++ {
		for j := 0; j < len(game.Cells[i]); j++ {
			cell := game.Cells[i][j]
			if int32(cell.Position.X) < x && int32(cell.Position.X)+cellSize > x && int32(cell.Position.Y) < y && int32(cell.Position.Y)+cellSize > y {
				game.Cells[i][j].LifeState = !game.Cells[i][j].LifeState
			}
		}
	}
}

func (game *Game) Update() {
	// len(game.Cells) = number of rows -> based on height
	// len(game.Cells[row]) = number of cols -> based on width
	for row := 0; row < len(game.Cells); row++ {
		for col := 0; col < len(game.Cells[row]); col++ {

			var aliveNeighbors int = 0

			// Checking neighbors
			for x := -1; x < 2; x++ {
				if row+x < 0 || row+x > len(game.Cells)-1 {
					continue
				}

				for y := -1; y < 2; y++ {
					if x == 0 && y == 0 {
						continue
					}

					if col+y < 0 || col+y > len(game.Cells[row])-1 {
						continue
					}

					if game.Cells[row+x][col+y].LifeState {
						aliveNeighbors++
					}
				}
			}

			if game.Cells[row][col].LifeState {
				if aliveNeighbors < 2 {
					game.Cells[row][col].Next = false
				} else if aliveNeighbors > 3 {
					game.Cells[row][col].Next = false
				} else {
					game.Cells[row][col].Next = true
				}
			} else if !game.Cells[row][col].LifeState {
				if aliveNeighbors == 3 {
					game.Cells[row][col].Next = true
				}
			}
		}
	}

	for x := 0; x < len(game.Cells); x++ {
		for y := 0; y < len(game.Cells[x]); y++ {
			game.Cells[x][y].LifeState = game.Cells[x][y].Next
		}
	}

	game.Generation++
}

func (game *Game) Input() {
	if rl.IsKeyPressed(rl.KeySpace) {
		game.GameActiveState = !game.GameActiveState
	}

	if rl.IsKeyPressed(rl.KeyR) {
		numOfRows = (screenHeight - (padding * 2)) / cellSize
		numOfCols = (screenWidth - (padding * 2)) / cellSize
		game.Cells = getCells(false)
		game.GameActiveState = false
	}

	if rl.IsKeyPressed(rl.KeyC) {
		numOfRows = (screenHeight - (padding * 2)) / cellSize
		numOfCols = (screenWidth - (padding * 2)) / cellSize
		game.Cells = getCells(true)
		game.GameActiveState = false
	}

	if rl.IsKeyPressed(rl.KeyS) {
		game.GameActiveState = false
		game.SaveGame()
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		game.CheckClick(rl.GetMouseX(), rl.GetMouseY())
	}
}

func (game *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(backgroundColor)

	for i := 0; i < len(game.Cells); i++ {
		for j := 0; j < len(game.Cells[i]); j++ {
			currentCell := game.Cells[i][j]
			if currentCell.LifeState {
				rl.DrawRectangleV(currentCell.Position, currentCell.Size, liveColor)
			} else {
				rl.DrawRectangleV(currentCell.Position, currentCell.Size, backgroundColor)
			}
		}
	}

	for y := padding; y < int(screenHeight-padding)+1; y += int(cellSize) {
		rl.DrawLineV(
			rl.Vector2{
				X: padding,
				Y: float32(y),
			},
			rl.Vector2{
				X: float32(screenWidth - padding),
				Y: float32(y),
			},
			gridColor,
		)
	}

	for x := padding; x < int(screenWidth-padding)+1; x += int(cellSize) {
		rl.DrawLineV(
			rl.Vector2{
				X: float32(x),
				Y: padding,
			},
			rl.Vector2{
				X: float32(x),
				Y: float32(screenHeight - padding),
			},
			gridColor,
		)
	}

	rl.EndDrawing()
}

func (game *Game) SaveGame() {
	json, err := json.MarshalIndent(game, "", "	")
	if err != nil {
		fmt.Println(err)
	}

	now := time.Now().Format(time.RFC3339Nano)
	filePath := fmt.Sprintf("./data/%s.json", now)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	if bytesWrote, err := file.Write(json); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Bytes Wrote: %d\n", bytesWrote)
	}

}

func getCells(empty bool) [][]*Cell {
	var board [][]*Cell = make([][]*Cell, numOfRows)

	for i := range board {
		board[i] = make([]*Cell, numOfCols)
	}

	var currentStartX int = padding
	var currentStartY int = padding

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			var randomNum int = rand.Intn(10)

			board[i][j] = &Cell{
				Position: rl.Vector2{
					X: float32(currentStartX),
					Y: float32(currentStartY),
				},
				Size: rl.Vector2{
					X: float32(cellSize),
					Y: float32(cellSize),
				},
			}

			if !empty {
				if randomNum >= 7 {
					board[i][j].LifeState = true
				} else {
					board[i][j].LifeState = false
				}

			} else {
				board[i][j].LifeState = false
			}

			currentStartX += int(cellSize)
		}

		currentStartX = padding
		currentStartY += int(cellSize)
	}

	return board
}
