package ui

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/derailed/k9s/internal/config"
	"github.com/derailed/k9s/internal/resource"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
)

func init() {
	initKeys()
}

const (
	menuIndexFmt = " [key:bg:b]<%d> [fg:bg:d]%s "
	maxRows      = 7
)

var menuRX = regexp.MustCompile(`\d`)

// MenuView represents menu options.
type MenuView struct {
	*tview.Table

	styles *config.Styles
}

// NewMenuView returns a new menu.
func NewMenuView(styles *config.Styles) *MenuView {
	v := MenuView{Table: tview.NewTable(), styles: styles}
	v.SetBackgroundColor(styles.BgColor())

	return &v
}

// HydrateMenu populate menu ui from hints.
func (v *MenuView) HydrateMenu(hh Hints) {
	v.Clear()
	sort.Sort(hh)
	t := v.buildMenuTable(hh)
	for row := 0; row < len(t); row++ {
		for col := 0; col < len(t[row]); col++ {
			if len(t[row][col]) == 0 {
				continue
			}
			c := tview.NewTableCell(t[row][col])
			c.SetBackgroundColor(v.styles.BgColor())
			v.SetCell(row, col, c)
		}
	}
}

func (v *MenuView) buildMenuTable(hh Hints) [][]string {
	table := make([][]Hint, maxRows+1)

	colCount := (len(hh) / maxRows) + 1
	for row := 0; row < maxRows; row++ {
		table[row] = make([]Hint, colCount+1)
	}

	var row, col int
	firstCmd := true
	maxKeys := make([]int, colCount+1)
	for _, h := range hh {
		isDigit := menuRX.MatchString(h.mnemonic)
		if !isDigit && firstCmd {
			row, col, firstCmd = 0, col+1, false
		}
		if maxKeys[col] < len(h.mnemonic) {
			maxKeys[col] = len(h.mnemonic)
		}
		table[row][col] = h
		row++
		if row >= maxRows {
			col++
			row = 0
		}
	}

	strTable := make([][]string, maxRows+1)
	for r := 0; r < len(table); r++ {
		strTable[r] = make([]string, len(table[r]))
	}
	for row := range strTable {
		for col := range strTable[row] {
			strTable[row][col] = v.formatMenu(table[row][col], maxKeys[col])
		}
	}

	return strTable
}

// ----------------------------------------------------------------------------
// Helpers...

func toMnemonic(s string) string {
	if len(s) == 0 {
		return s
	}

	return "<" + strings.ToLower(s) + ">"
}

func (v *MenuView) formatMenu(h Hint, size int) string {
	i, err := strconv.Atoi(h.mnemonic)
	if err == nil {
		return formatNSMenu(i, h.description, v.styles.Frame())
	}

	return formatPlainMenu(h, size, v.styles.Frame())
}

func formatNSMenu(i int, name string, styles config.Frame) string {
	fmat := strings.Replace(menuIndexFmt, "[key", "["+styles.Menu.NumKeyColor, 1)
	fmat = strings.Replace(fmat, ":bg:", ":"+styles.Title.BgColor+":", -1)
	fmat = strings.Replace(fmat, "[fg", "["+styles.Menu.FgColor, 1)
	return fmt.Sprintf(fmat, i, resource.Truncate(name, 14))
}

func formatPlainMenu(h Hint, size int, styles config.Frame) string {
	menuFmt := " [key:bg:b]%-" + strconv.Itoa(size+2) + "s [fg:bg:d]%s "
	fmat := strings.Replace(menuFmt, "[key", "["+styles.Menu.KeyColor, 1)
	fmat = strings.Replace(fmat, "[fg", "["+styles.Menu.FgColor, 1)
	fmat = strings.Replace(fmat, ":bg:", ":"+styles.Title.BgColor+":", -1)
	return fmt.Sprintf(fmat, toMnemonic(h.mnemonic), h.description)
}

// -----------------------------------------------------------------------------
// Key mapping Constants

// Defines numeric keys for container actions
const (
	Key0 int32 = iota + 48
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
)

// Defines AltKeys
const (
	KeyAltA tcell.Key = 4 * (iota + 97)
	KeyAltB
	KeyAltC
	KeyAltD
	KeyAltE
	KeyAltF
	KeyAltG
	KeyAltH
	KeyAltI
	KeyAltJ
	KeyAltK
	KeyAltL
	KeyAltM
	KeyAltN
	KeyAltO
	KeyAltP
	KeyAltQ
	KeyAltR
	KeyAltS
	KeyAltT
	KeyAltU
	KeyAltV
	KeyAltW
	KeyAltX
	KeyAltY
	KeyAltZ
)

// Defines char keystrokes
const (
	KeyA tcell.Key = iota + 97
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyHelp  = 63
	KeySlash = 47
	KeyColon = 58
)

// Define Shift Keys
const (
	KeyShiftA tcell.Key = iota + 65
	KeyShiftB
	KeyShiftC
	KeyShiftD
	KeyShiftE
	KeyShiftF
	KeyShiftG
	KeyShiftH
	KeyShiftI
	KeyShiftJ
	KeyShiftK
	KeyShiftL
	KeyShiftM
	KeyShiftN
	KeyShiftO
	KeyShiftP
	KeyShiftQ
	KeyShiftR
	KeyShiftS
	KeyShiftT
	KeyShiftU
	KeyShiftV
	KeyShiftW
	KeyShiftX
	KeyShiftY
	KeyShiftZ
)

// NumKeys tracks number keys.
var NumKeys = map[int]int32{
	0: Key0,
	1: Key1,
	2: Key2,
	3: Key3,
	4: Key4,
	5: Key5,
	6: Key6,
	7: Key7,
	8: Key8,
	9: Key9,
}

func initKeys() {
	tcell.KeyNames[tcell.Key(KeyHelp)] = "?"
	tcell.KeyNames[tcell.Key(KeySlash)] = "/"

	initNumbKeys()
	initStdKeys()
	initShiftKeys()
	initAltKeys()
}

func initNumbKeys() {
	tcell.KeyNames[tcell.Key(Key0)] = "0"
	tcell.KeyNames[tcell.Key(Key1)] = "1"
	tcell.KeyNames[tcell.Key(Key2)] = "2"
	tcell.KeyNames[tcell.Key(Key3)] = "3"
	tcell.KeyNames[tcell.Key(Key4)] = "4"
	tcell.KeyNames[tcell.Key(Key5)] = "5"
	tcell.KeyNames[tcell.Key(Key6)] = "6"
	tcell.KeyNames[tcell.Key(Key7)] = "7"
	tcell.KeyNames[tcell.Key(Key8)] = "8"
	tcell.KeyNames[tcell.Key(Key9)] = "9"
}

func initStdKeys() {
	tcell.KeyNames[tcell.Key(KeyA)] = "a"
	tcell.KeyNames[tcell.Key(KeyB)] = "b"
	tcell.KeyNames[tcell.Key(KeyC)] = "c"
	tcell.KeyNames[tcell.Key(KeyD)] = "d"
	tcell.KeyNames[tcell.Key(KeyE)] = "e"
	tcell.KeyNames[tcell.Key(KeyF)] = "f"
	tcell.KeyNames[tcell.Key(KeyG)] = "g"
	tcell.KeyNames[tcell.Key(KeyH)] = "h"
	tcell.KeyNames[tcell.Key(KeyI)] = "i"
	tcell.KeyNames[tcell.Key(KeyJ)] = "j"
	tcell.KeyNames[tcell.Key(KeyK)] = "k"
	tcell.KeyNames[tcell.Key(KeyL)] = "l"
	tcell.KeyNames[tcell.Key(KeyM)] = "m"
	tcell.KeyNames[tcell.Key(KeyN)] = "n"
	tcell.KeyNames[tcell.Key(KeyO)] = "o"
	tcell.KeyNames[tcell.Key(KeyP)] = "p"
	tcell.KeyNames[tcell.Key(KeyQ)] = "q"
	tcell.KeyNames[tcell.Key(KeyR)] = "r"
	tcell.KeyNames[tcell.Key(KeyS)] = "s"
	tcell.KeyNames[tcell.Key(KeyT)] = "t"
	tcell.KeyNames[tcell.Key(KeyU)] = "u"
	tcell.KeyNames[tcell.Key(KeyV)] = "v"
	tcell.KeyNames[tcell.Key(KeyW)] = "w"
	tcell.KeyNames[tcell.Key(KeyX)] = "x"
	tcell.KeyNames[tcell.Key(KeyY)] = "y"
	tcell.KeyNames[tcell.Key(KeyZ)] = "z"
}

func initShiftKeys() {
	tcell.KeyNames[tcell.Key(KeyShiftA)] = "Shift-A"
	tcell.KeyNames[tcell.Key(KeyShiftB)] = "Shift-B"
	tcell.KeyNames[tcell.Key(KeyShiftC)] = "Shift-C"
	tcell.KeyNames[tcell.Key(KeyShiftD)] = "Shift-D"
	tcell.KeyNames[tcell.Key(KeyShiftE)] = "Shift-E"
	tcell.KeyNames[tcell.Key(KeyShiftF)] = "Shift-F"
	tcell.KeyNames[tcell.Key(KeyShiftG)] = "Shift-G"
	tcell.KeyNames[tcell.Key(KeyShiftH)] = "Shift-H"
	tcell.KeyNames[tcell.Key(KeyShiftI)] = "Shift-I"
	tcell.KeyNames[tcell.Key(KeyShiftJ)] = "Shift-J"
	tcell.KeyNames[tcell.Key(KeyShiftK)] = "Shift-K"
	tcell.KeyNames[tcell.Key(KeyShiftL)] = "Shift-L"
	tcell.KeyNames[tcell.Key(KeyShiftM)] = "Shift-M"
	tcell.KeyNames[tcell.Key(KeyShiftN)] = "Shift-N"
	tcell.KeyNames[tcell.Key(KeyShiftO)] = "Shift-O"
	tcell.KeyNames[tcell.Key(KeyShiftP)] = "Shift-P"
	tcell.KeyNames[tcell.Key(KeyShiftQ)] = "Shift-Q"
	tcell.KeyNames[tcell.Key(KeyShiftR)] = "Shift-R"
	tcell.KeyNames[tcell.Key(KeyShiftS)] = "Shift-S"
	tcell.KeyNames[tcell.Key(KeyShiftT)] = "Shift-T"
	tcell.KeyNames[tcell.Key(KeyShiftU)] = "Shift-U"
	tcell.KeyNames[tcell.Key(KeyShiftV)] = "Shift-V"
	tcell.KeyNames[tcell.Key(KeyShiftW)] = "Shift-W"
	tcell.KeyNames[tcell.Key(KeyShiftX)] = "Shift-X"
	tcell.KeyNames[tcell.Key(KeyShiftY)] = "Shift-Y"
	tcell.KeyNames[tcell.Key(KeyShiftZ)] = "Shift-Z"
}

func initAltKeys() {
	tcell.KeyNames[tcell.Key(KeyAltA)] = "Alt-A"
	tcell.KeyNames[tcell.Key(KeyAltB)] = "Alt-B"
	tcell.KeyNames[tcell.Key(KeyAltC)] = "Alt-C"
	tcell.KeyNames[tcell.Key(KeyAltD)] = "Alt-D"
	tcell.KeyNames[tcell.Key(KeyAltE)] = "Alt-E"
	tcell.KeyNames[tcell.Key(KeyAltF)] = "Alt-F"
	tcell.KeyNames[tcell.Key(KeyAltG)] = "Alt-G"
	tcell.KeyNames[tcell.Key(KeyAltH)] = "Alt-H"
	tcell.KeyNames[tcell.Key(KeyAltI)] = "Alt-I"
	tcell.KeyNames[tcell.Key(KeyAltJ)] = "Alt-J"
	tcell.KeyNames[tcell.Key(KeyAltK)] = "Alt-K"
	tcell.KeyNames[tcell.Key(KeyAltL)] = "Alt-L"
	tcell.KeyNames[tcell.Key(KeyAltM)] = "Alt-M"
	tcell.KeyNames[tcell.Key(KeyAltN)] = "Alt-N"
	tcell.KeyNames[tcell.Key(KeyAltO)] = "Alt-O"
	tcell.KeyNames[tcell.Key(KeyAltP)] = "Alt-P"
	tcell.KeyNames[tcell.Key(KeyAltQ)] = "Alt-Q"
	tcell.KeyNames[tcell.Key(KeyAltR)] = "Alt-R"
	tcell.KeyNames[tcell.Key(KeyAltS)] = "Alt-S"
	tcell.KeyNames[tcell.Key(KeyAltT)] = "Alt-T"
	tcell.KeyNames[tcell.Key(KeyAltU)] = "Alt-U"
	tcell.KeyNames[tcell.Key(KeyAltV)] = "Alt-V"
	tcell.KeyNames[tcell.Key(KeyAltW)] = "Alt-W"
	tcell.KeyNames[tcell.Key(KeyAltX)] = "Alt-X"
	tcell.KeyNames[tcell.Key(KeyAltY)] = "Alt-Y"
	tcell.KeyNames[tcell.Key(KeyAltZ)] = "Alt-Z"
}
