package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/inancgumus/screen"
)

var branch = "░"
var light = "█"
var trunk = "▓"
var ground = "▔"
var flake = "*"

var red = color.New(color.FgRed, color.BgBlack).PrintfFunc()
var blue = color.New(color.FgBlue, color.BgBlack).PrintfFunc()
var yellow = color.New(color.FgYellow, color.BgBlack).PrintfFunc()
var green = color.New(color.FgGreen, color.BgBlack).PrintfFunc()
var magenta = color.New(color.FgMagenta, color.BgBlack).PrintfFunc()
var black = color.New(color.FgHiBlack, color.BgBlack).PrintfFunc()
var white = color.New(color.FgWhite, color.BgBlack).PrintfFunc()
var bgBlack = color.New(color.BgBlack).PrintfFunc()

func loadTree(txtFile string) (string, error) {
	bytes, err := os.ReadFile(txtFile)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func printLightChar(char rune, i int) {
	c := string(char)

	if i == 1 {
		red(c)
	} else if i == 2 {
		blue(c)
	} else if i == 3 {
		yellow(c)
	} else if i == 4 {
		magenta(c)
	} else {
		green(c)
	}
}

func printConstColoredChar(char rune) {
	c := string(char) // cast rune to string

	switch c {
	case branch:
		green(c)
	case trunk:
		black(c)
	case ground:
		white(c)
	case flake:
		white(c)
	default:
		bgBlack(c)
	}
}

func countdown() {
	if time.Now().Month().String() != "December" {
		return
	}

	xmasDay := 25
	currentDay := time.Now().Local().Day()
	calendar := " "

	if currentDay < xmasDay {
		for i := 0; i < currentDay; i++ {
			calendar += "●"
		}
		for i := 0; i < xmasDay-currentDay; i++ {
			calendar += "○"
		}
		calendar += " "
	} else {
		calendar += "●●●●●●●●●●●●●●●●●●●●●●●●● "
	}

	white(calendar)
}

func newSnowRow() string {
	blankRow := "                           "
	snowRow := ""
	for _, c := range blankRow {
		if rand.Intn(20) == 1 {
			snowRow += flake
			continue
		}
		snowRow += string(c)
	}
	return snowRow
}

func letItSnow(cleanTreeString string, lastTreeString string) string {
	cleanTree := strings.Split(cleanTreeString, "\n")
	lastTree := strings.Split(lastTreeString, "\n")
	var currentTree []string

	// create a placeholder first row
	currentTree = append(currentTree, newSnowRow())

	// iterate over tree rows, skipping the first one
	for i := 1; i < len(cleanTree); i++ {
		// if character in the clean row is blank, check if corresponding one in the previous last row was a flake so it can get moved down
		// get the current clean row
		currCleanRow := cleanTree[i]
		// get the above row from the last tree
		aboveLastRow := lastTree[i-1]
		// init current row in current tree
		var currCurrRow string
		// iterate through and if above has a flake copy it down into current
		for i := 0; i < utf8.RuneCountInString(currCleanRow); i++ {
			cleanStr := strings.Split(currCleanRow, "")[i]
			if cleanStr != " " {
				// character is not blank so it can't hold a snowflake, return
				currCurrRow += cleanStr
				continue
			}
			// character is blank, check if above last row had a flake
			aboveLastStr := strings.Split(aboveLastRow, "")[i]
			if aboveLastStr == "*" {
				// copy down snowflake
				currCurrRow += aboveLastStr
			} else {
				// current and last are both blank
				currCurrRow += cleanStr
			}
		}

		currentTree = append(currentTree, currCurrRow)
	}

	return strings.Join(currentTree, "\n")
}

func run(t string) {
	// clear terminal and move cursor to origin
	screen.Clear()
	screen.MoveTopLeft()

	// find positions of the lights in the tree
	var lightPos []int
	for pos, char := range t {
		if string(char) == light {
			lightPos = append(lightPos, pos)
		}
	}

	// randomly generate light colors
	var lightColor []int
	for i := 0; i < len(lightPos); i++ {
		lightColor = append(lightColor, rand.Intn(5))
	}

	// generate snow frame
	bwFrame := letItSnow(t, t)

	// print colorized tree with snow
	for pos, char := range bwFrame {
		if slices.Contains(lightPos, pos) {
			i := slices.Index(lightPos, pos)
			printLightChar(char, lightColor[i])
		} else {
			printConstColoredChar(char)
		}
	}
	fmt.Printf("\n")

	counter := 0

	for {
		// generate next snow frame
		bwFrame = letItSnow(t, bwFrame)

		// increment counter
		counter += 1

		// update color
		if counter == 4 {
			// empty slice
			lightColor = []int{}
			// fill with new colors
			for i := 0; i < len(lightPos); i++ {
				lightColor = append(lightColor, rand.Intn(5))
			}
			// reset counter
			counter = 0
		}

		// move cursor back to origin
		screen.MoveTopLeft()

		// print colorized tree with snow
		for pos, char := range bwFrame {
			if slices.Contains(lightPos, pos) {
				i := slices.Index(lightPos, pos)
				printLightChar(char, lightColor[i])
			} else {
				printConstColoredChar(char)
			}
		}
		fmt.Printf("\n")

		// print countdown
		countdown()
		fmt.Printf("\n")

		// wait before next iteration
		time.Sleep(300 * time.Millisecond)
	}
}

func main() {
	t, err := loadTree("tree.txt")

	if err != nil {
		log.Fatal(err)
	}

	run(t)
}
