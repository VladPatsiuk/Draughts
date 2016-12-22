package main

import (
	//"bufio"
	"fmt"
	"os"
	"os/exec"
	//"strconv"
	"github.com/eiannone/keyboard"
	"github.com/shiena/ansicolor"
)

const (
	HEIGHT           = 8
	WIDTH            = 8
	WHITE_BACKGROUND = "\x1b[47m"
	//i'll add other colors in future
)

type Field struct {
	xPos         int
	yPos         int
	color        string // ' ' - white, x - black, a - active
	active       bool   // this fiels is active now
	hasChecker   bool
	checkerColor int // 0 - white, 1 - black
}

type Board struct {
	bord         [HEIGHT][WIDTH]Field
	activeFieldX int
	activeFieldY int

	destX int
	destY int

	whiteCheckersCount int
	blackCheckersCount int
}

type GameProcess struct {
	turn     int  // 0 - white, 1 - black
	gameOver bool // true - end of the game
}

func (g *GameProcess) InitGame() {
	g.turn = 0
	g.gameOver = false
}

func (g *GameProcess) Game() {
	b := Board{}
	b.CreateField()
	b.Fill()
	for g.gameOver == false {
		ClearScreen()
		if b.GameState() == 1 {
			g.gameOver = true
		}
		b.Draw()
		if b.ChooseChecker(g.turn) == true {
			g.NextTurn()
		}
	}
	fmt.Println("GAME OVER")

}

func (g *GameProcess) NextTurn() {
	g.turn = 1 - g.turn
}

func (f *Field) SetChecker(color int) {
	f.hasChecker = true
	f.checkerColor = color
}

func (b *Board) CreateField() {
	var i, j int
	for i = 0; i < WIDTH; i++ {
		for j = 0; j < HEIGHT; j++ {
			if (i+j)%2 == 0 {
				b.bord[i][j].color = " "
			} else {
				b.bord[i][j].color = "x"
			}
		}
	}
	b.activeFieldX = 5
	b.activeFieldY = 5
	b.bord[b.activeFieldX][b.activeFieldY].active = true

	b.destX = 0
	b.destY = 0

	b.whiteCheckersCount = 12
	b.blackCheckersCount = 12
}

func (b *Board) Fill() {
	var i, j int
	for i = 0; i < WIDTH; i++ {
		for j = 0; j < HEIGHT; j++ {
			if i < 3 && (i+j)%2 == 1 {
				b.bord[i][j].hasChecker = true
				b.bord[i][j].checkerColor = 0
			} else if i > 4 && (i+j)%2 == 1 {
				b.bord[i][j].hasChecker = true
				b.bord[i][j].checkerColor = 1
			} else {
				b.bord[i][j].hasChecker = false
			}
		}
	}
}

func (b *Board) Draw() {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	var i, j int

	for i = 0; i < WIDTH; i++ {
		for j = 0; j < HEIGHT; j++ {

			if b.bord[i][j].hasChecker == true {
				if b.bord[i][j].active == true {
					if b.bord[i][j].checkerColor == 0 {
						fmt.Fprintf(w, "%s%so%s", "\x1b[37m", "\x1b[45m", "\x1b[0m")
					} else {
						fmt.Fprintf(w, "%s%so%s", "\x1b[30m", "\x1b[45m", "\x1b[0m")
					}
				} else {
					if b.bord[i][j].checkerColor == 1 {
						fmt.Fprintf(w, "%s%so%s", "\x1b[30m", "\x1b[44m", "\x1b[0m")
					} else {
						fmt.Fprintf(w, "%s%so%s", "\x1b[37m", "\x1b[44m", "\x1b[0m")
					}
				}
			} else if b.bord[i][j].active == true {
				fmt.Fprintf(w, "%s %s", "\x1b[45m", "\x1b[0m")
			} else if b.bord[i][j].color == " " {
				fmt.Fprintf(w, "%s %s", WHITE_BACKGROUND, "\x1b[0m")
			} else if b.bord[i][j].color == "x" {
				fmt.Fprintf(w, "%s %s", "\x1b[44m", "\x1b[0m")
			}
		}
		fmt.Println()
	}
}

func (b *Board) ChooseChecker(color int) bool {
	s := ""
	if color == 0 {
		s += "white"
	} else {
		s += "black"
	}
	fmt.Println("" + s + " player turn")
	char, _, err := keyboard.GetSingleKey()
	if err != nil {
		panic(err)
	}
	command := fmt.Sprintf("%c", char)

	b.bord[b.activeFieldX][b.activeFieldY].active = false
	if command == "w" && b.activeFieldX > 0 {
		b.activeFieldX--
	}
	if command == "a" && b.activeFieldY > 0 {
		b.activeFieldY--
	}
	if command == "s" && b.activeFieldX < 7 {
		b.activeFieldX++
	}
	if command == "d" && b.activeFieldY < 7 {
		b.activeFieldY++
	}
	if command == "e" && b.bord[b.activeFieldX][b.activeFieldY].hasChecker == true {
		if b.bord[b.activeFieldX][b.activeFieldY].checkerColor == color {
			if b.MoveChecker(b.activeFieldX, b.activeFieldY) == true {
				return true
			}

		} else {
			fmt.Println("Its not your checker. Choose checker with other color")
			fmt.Scanln()
		}
	}
	if command == "e" && b.bord[b.activeFieldX][b.activeFieldY].hasChecker == false {
		fmt.Println("This is empty field. Choose other")
		fmt.Scanln()
	}
	b.bord[b.activeFieldX][b.activeFieldY].active = true
	return false
}

func (b *Board) ChooseDestination() {
	char, _, err := keyboard.GetSingleKey()
	if err != nil {
		panic(err)
	}
	command := fmt.Sprintf("%c", char)
	b.bord[b.activeFieldX][b.activeFieldY].active = false
	if command == "w" && b.activeFieldX > 0 {
		b.activeFieldX--
	}
	if command == "a" && b.activeFieldY > 0 {
		b.activeFieldY--
	}
	if command == "s" && b.activeFieldX < 7 {
		b.activeFieldX++
	}
	if command == "d" && b.activeFieldY < 7 {
		b.activeFieldY++
	}
	if command == "e" {
		b.destX = b.activeFieldX
		b.destY = b.activeFieldY
	}
	b.bord[b.activeFieldX][b.activeFieldY].active = true
}

func (b *Board) MoveChecker(x, y int) bool {
	var xCoord, yCoord int
	b.destX = 0
	b.destY = 0

	for b.destX == 0 && b.destY == 0 {
		ClearScreen()
		b.Draw()
		b.ChooseDestination()
	}
	xCoord = b.destX
	yCoord = b.destY

	//flag := CheckMove(x, y, yCoord, xCoord-1)
	flag := b.CheckMove(x, y, xCoord, yCoord)

	beat := b.CheckBeat(x, y)

	if beat == true {
		flag = false
		//fmt.Print("You should beat enemy checker")
		fmt.Scanln()
		b.Beat(x, y, xCoord, yCoord)
		if b.CheckBeat(xCoord, yCoord) == false {
			return true
		}
	} else if flag == true {
		b.bord[x][y].hasChecker = false
		color := b.bord[x][y].checkerColor

		b.bord[xCoord][yCoord].hasChecker = true
		b.bord[xCoord][yCoord].checkerColor = color
		return true
	} else {
		fmt.Print("Wrong position. Try again")
		fmt.Scanln()
	}
	return false
}

func (b *Board) CheckMove(x, y, xx, yy int) bool {
	if b.bord[xx][yy].hasChecker == false {
		if b.bord[x][y].checkerColor == 0 {
			if xx == x+1 && (yy == y+1 || yy == y-1) {
				return true
			}
		} else if b.bord[x][y].checkerColor == 1 {
			if xx == x-1 && (yy == y+1 || yy == y-1) {
				if 0 <= xx && xx < 8 && 0 <= yy && yy < 8 {
					return true
				}
			}
		}
	}
	if b.bord[(x+xx)/2][(y+yy)/2].hasChecker == true && b.bord[(x+xx)/2][(y+yy)/2].checkerColor == (1-b.bord[x][y].checkerColor) {
		return true
	}
	return false
}

func (b *Board) Beat(x, y, xx, yy int) {
	if b.bord[(x+xx)/2][(y+yy)/2].hasChecker == true && b.bord[(x+xx)/2][(y+yy)/2].checkerColor == (1-b.bord[x][y].checkerColor) {
		b.bord[x][y].hasChecker = false
		color := b.bord[x][y].checkerColor
		b.bord[(x+xx)/2][(y+yy)/2].hasChecker = false
		b.bord[xx][yy].hasChecker = true
		b.bord[xx][yy].checkerColor = color

		if b.bord[(x+xx)/2][(y+yy)/2].checkerColor == 0 {
			b.whiteCheckersCount--
		} else {
			b.blackCheckersCount--
		}
	}

}

func (b *Board) CheckBeat(x, y int) bool {

	if (x+1) < 7 && (y+1) < 7 && b.bord[x+1][y+1].hasChecker == true && b.bord[x+1][y+1].checkerColor == (1-b.bord[x][y].checkerColor) && b.bord[x+2][y+2].hasChecker == false {
		return true
	} else if (x+1) < 7 && (y-1) >= 1 && b.bord[x+1][y-1].hasChecker == true && b.bord[x+1][y-1].checkerColor == (1-b.bord[x][y].checkerColor) && b.bord[x+2][y-2].hasChecker == false {
		return true
	} else if (x-1) >= 1 && (y+1) < 7 && b.bord[x-1][y+1].hasChecker == true && b.bord[x-1][y+1].checkerColor == (1-b.bord[x][y].checkerColor) && b.bord[x-2][y+2].hasChecker == false {
		return true
	} else if (x-1) >= 1 && (y-1) >= 1 && b.bord[x-1][y-1].hasChecker == true && b.bord[x-1][y-1].checkerColor == (1-b.bord[x][y].checkerColor) && b.bord[x-2][y-2].hasChecker == false {
		return true
	}
	return false
}

func (b *Board) GameState() int {
	if b.whiteCheckersCount == 0 || b.blackCheckersCount == 0 {
		return 1
	}
	return 0
}

// ------------------

func ClearScreen() {
	c := exec.Command("cmd", "/c", "cls")
	c.Stdout = os.Stdout
	c.Run()
}

func main() {
	gp := GameProcess{}
	gp.InitGame()
	gp.Game()

}
