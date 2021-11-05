package term

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/pterm/pterm"
)

type PrintableList interface {
	Get(i int) (string, string)
	Size() int
}

func List(items PrintableList) (int, error) {
	area, err := pterm.DefaultArea.WithRemoveWhenDone(true).Start()
	if err != nil {
		return -1, err
	}

	//linesToRemove := 0
	var data pterm.TableData
	for i := 0; i < items.Size(); i++ {
		name, desc := items.Get(i)
		data = append(data, []string{strconv.Itoa(i), name, desc})
		//linesToRemove++
	}

	tbl, err := pterm.DefaultTable.WithData(data).Srender()
	if err != nil {
		return -1, err
	}
	//pterm.Println(tbl)
	area.Update(tbl)

	linesToRemove := 0
	pos := -1
	for pos == -1 {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		in := scanner.Text()
		linesToRemove++

		pos, err = strconv.Atoi(in)
		if err != nil {
			fmt.Println("Please input a valid index")
			linesToRemove++
			pos = -1
		}
	}

	area.Stop()

	// Remove any lines that were not added through pterm
	for i := 0; i < linesToRemove; i++ {
		fmt.Print("\033[A")
		fmt.Print("\033[2K")
	}
	//fmt.Println(linesToRemove)

	name, _ := items.Get(pos)
	pterm.Println(name)

	return pos, nil
}

//var (
//shortlist = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
//)

//func paintList(items PrintableList) error {
//app := tview.NewApplication()

//list := tview.NewList()

//sel := func(name string) func() {
//return func() {
//app.Stop()
//fmt.Println(name)
//}
//}

//for i := 0; i < items.Size(); i++ {
//name, desc := items.Get(i)
//list.AddItem(name, desc, shortlist[i], sel(name))
//}

//if err := app.SetRoot(list, true).Run(); err != nil {
//return err
//}

//return nil
//}
