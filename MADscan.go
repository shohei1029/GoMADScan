package main

import (
	"fmt"
	// "github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	delimName = []string{"Select delimiter", "tab (\\t)", "space ( )", "comma (,)", "period (.)", "colon (:)"}
	delimList = []string{"\t", "\t", " ", ",", ".", ":"}
)

// minimum size for the array size of keywords
const minSize = 50

func authors() []string {
	return []string{"carushi<l.cawaguchi(at)gmail.com>\n\nReference:"}
}

type arguments struct {
	column     int
	inputPath  string
	filterPath string
	outputPath string
	ignoreCase bool
	perfectMatch bool
	delim      string
}

func isAlphaNum(c byte) bool {
	if !(('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9')) {
		return false
	}
	return true
}

func (a *arguments) setDelimiter() {
	if len(a.delim) > 1 {
		for i, delim := range delimName {
			if a.delim == delim {
				a.delim = delimList[i]
				return
			}
		}
		a.delim = delimList[0]
	}
}

func returnNewlineChar(lines string) string {
	newline := "\n"
	if !strings.Contains(lines, "\n") {
		newline = "\r"
	}
	return newline
}

func getKeywords(key string, ignoreCase bool) ([]string, error) {
	lines, err := ioutil.ReadFile(key)
	keyList := make([]string, 0, minSize)
	if err != nil {
		return keyList, err
	}
	newline := returnNewlineChar(string(lines))
	for _, gene := range strings.Split(string(lines), newline) {
		gene = strings.TrimRight(gene, "\r")
		if len(gene) == 0 {
			continue
		}
		if ignoreCase {
			gene = strings.ToUpper(gene)
		}
		keyList = append(keyList, gene)
		symbol := false
		for c := 0; c < len(gene); {
			if !isAlphaNum(gene[c]) {
				symbol = true
				gene = strings.Replace(gene, string(gene[c]), "", 1)
				continue
			}
			c++
		}
		if symbol {
			keyList = append(keyList, gene)
		}
	}
	return keyList, err
}

func searchKeywords(column int, inputPath string, delim string, keyList []string, ignoreCase bool, perfectMatch bool) ([]string, error) {
	matchedLines := make([]string, 0, minSize)
	lines, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return matchedLines, err
	}
	newline := returnNewlineChar(string(lines))
	for _, line := range strings.Split(string(lines), newline) {
		s := strings.Split(strings.TrimRight(line, "\r"), delim)
		notMatched := true
		for i := int(math.Max(float64(column-1), float64(0))); i < len(s) && notMatched; i++ {
			if ignoreCase {
				s[i] = strings.ToUpper(s[i])
			}
			for _, key := range keyList {
				if (perfectMatch && s[i] == key) ||	(!perfectMatch && strings.Index(s[i], key) > -1) {
					matchedLines = append(matchedLines, line)
					notMatched = false
					break
				}
			}
			if column > 0 {
				break
			}
		}
	}
	return matchedLines, err
}

func output(outputPath string, matchedLines []string) {
	err := ioutil.WriteFile(outputPath, []byte(strings.Join(matchedLines[:], "\n")), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getKeysearchWords(arg arguments) (int, error) {
	arg.setDelimiter()
	fmt.Println("Read: " + arg.filterPath)
	keyList, err := getKeywords(arg.filterPath, arg.ignoreCase)
	if err != nil {
		return -1, err
	}
	fmt.Println("Read: " + arg.inputPath)
	matchedLines, err := searchKeywords(arg.column, arg.inputPath, arg.delim, keyList, arg.ignoreCase, arg.perfectMatch)
	fmt.Printf("Write: %s, n=", arg.outputPath)
	fmt.Println(len(matchedLines))
	if len(matchedLines) == 0 || err != nil {
		return 0, err
	}
	output(arg.outputPath, matchedLines)
	return len(matchedLines), err
}

func main() {
	gtk.Init(nil)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("MADscan")
	window.SetIconName("MADscan-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		// fmt.Println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "")

	//--------------------------------------------------------
	// GtkVBox
	//--------------------------------------------------------
	vbox := gtk.NewVBox(false, 1)

	//--------------------------------------------------------
	// GtkMenuBar
	//--------------------------------------------------------
	menubar := gtk.NewMenuBar()
	vbox.PackStart(menubar, false, false, 0)

	//--------------------------------------------------------
	// GtkVPaned
	//--------------------------------------------------------
	vpaned := gtk.NewVPaned()
	vbox.Add(vpaned)

	//--------------------------------------------------------
	// GtkFrame
	//--------------------------------------------------------
	frame1 := gtk.NewFrame("")
	framebox1 := gtk.NewVBox(false, 1)
	frame1.Add(framebox1)

	frame2 := gtk.NewFrame("Column position")
	framebox2 := gtk.NewVBox(false, 1)
	frame2.Add(framebox2)

	vpaned.Pack1(frame1, false, false)
	vpaned.Pack2(frame2, false, false)

	//--------------------------------------------------------
	// GtkImage
	//--------------------------------------------------------
	dir := os.Getenv("GOPATH")
	imagefile := filepath.Join(dir, "/src/github.com/carushi/MADscan/image/logo.png")
	label := gtk.NewLabel("Modification associated database scanner")
	label.ModifyFontEasy("DejaVu Serif 15")
	framebox1.PackStart(label, false, true, 0)
	image := gtk.NewImageFromFile(imagefile)
	framebox1.Add(image)

	//--------------------------------------------------------
	// Data input and output filename
	//--------------------------------------------------------
	arg := arguments{
		0,
		filepath.Join(dir, "/src/github.com/carushi/MADscan/data/Acetylation_site_dataset"),
		filepath.Join(dir, "/src/github.com/carushi/MADscan/data/Ras_gene_list.txt"),
		filepath.Join(dir, "/src/github.com/carushi/MADscan/data/output.txt"),
		false,
		true,
		"\t"}

	//--------------------------------------------------------
	// GtkScale
	//--------------------------------------------------------
	scale := gtk.NewHScaleWithRange(0, 20, 1)
	scale.Connect("value-changed", func() {
		arg.column = int(scale.GetValue())
		// fmt.Println("scale:", int(scale.GetValue()))
	})
	framebox2.Add(scale)

	//--------------------------------------------------------
	// InputArea
	//--------------------------------------------------------
	ientry := gtk.NewEntry()
	ientry.SetText(arg.inputPath)
	inputs := gtk.NewHBox(false, 1)
	button := gtk.NewButtonWithLabel("Choose input file")
	button.Clicked(func() {
		//--------------------------------------------------------
		// GtkFileChooserDialog
		//--------------------------------------------------------
		filechooserdialog := gtk.NewFileChooserDialog(
			"Choose File...",
			button.GetTopLevelAsWindow(),
			gtk.FILE_CHOOSER_ACTION_OPEN,
			gtk.STOCK_OK,
			gtk.RESPONSE_ACCEPT)
		filter := gtk.NewFileFilter()
		filter.AddPattern("*")
		filechooserdialog.AddFilter(filter)
		filechooserdialog.Response(func() {
			arg.inputPath = filechooserdialog.GetFilename()
			ientry.SetText(arg.inputPath)
			filechooserdialog.Destroy()
		})
		filechooserdialog.Run()
	})
	inputs.Add(button)
	inputs.Add(ientry)
	framebox2.PackStart(inputs, false, false, 0)

	//--------------------------------------------------------
	// FilterArea
	//--------------------------------------------------------

	oentry := gtk.NewEntry()
	oentry.SetText(arg.outputPath)
	inputs = gtk.NewHBox(false, 1)
	button = gtk.NewButtonWithLabel("Choose output file")
	button.Clicked(func() {
		//--------------------------------------------------------
		// GtkFileChooserDialog
		//--------------------------------------------------------
		filechooserdialog := gtk.NewFileChooserDialog(
			"Choose File...",
			button.GetTopLevelAsWindow(),
			gtk.FILE_CHOOSER_ACTION_OPEN,
			gtk.STOCK_OK,
			gtk.RESPONSE_ACCEPT)
		filter := gtk.NewFileFilter()
		filter.AddPattern("*")
		filechooserdialog.AddFilter(filter)
		filechooserdialog.Response(func() {
			arg.outputPath = filechooserdialog.GetFilename()
			oentry.SetText(arg.outputPath)
			filechooserdialog.Destroy()
		})
		filechooserdialog.Run()
	})
	inputs.Add(button)
	inputs.Add(oentry)
	framebox2.PackStart(inputs, false, false, 0)

	//--------------------------------------------------------
	// FilterArea
	//--------------------------------------------------------

	fentry := gtk.NewEntry()
	fentry.SetText(arg.filterPath)
	inputs = gtk.NewHBox(false, 1)
	button = gtk.NewButtonWithLabel("Choose keyword file")
	button.Clicked(func() {
		//--------------------------------------------------------
		// GtkFileChooserDialog
		//--------------------------------------------------------
		filechooserdialog := gtk.NewFileChooserDialog(
			"Choose File...",
			button.GetTopLevelAsWindow(),
			gtk.FILE_CHOOSER_ACTION_OPEN,
			gtk.STOCK_OK,
			gtk.RESPONSE_ACCEPT)
		filter := gtk.NewFileFilter()
		filter.AddPattern("*")
		filechooserdialog.AddFilter(filter)
		filechooserdialog.Response(func() {
			arg.filterPath = filechooserdialog.GetFilename()
			fentry.SetText(arg.filterPath)
			filechooserdialog.Destroy()
		})
		filechooserdialog.Run()
	})
	inputs.Add(button)
	inputs.Add(fentry)
	framebox2.PackStart(inputs, false, false, 0)

	buttons := gtk.NewHBox(false, 1)

	//--------------------------------------------------------
	// GtkCheckButton
	//--------------------------------------------------------
	checkbutton := gtk.NewCheckButtonWithLabel("Ignore lower/upper case")
	checkbutton.Connect("toggled", func() {
		if checkbutton.GetActive() {
			arg.ignoreCase = true
		} else {
			arg.ignoreCase = false
		}
	})
	buttons.Add(checkbutton)

	checkMatchButton := gtk.NewCheckButtonWithLabel("Partial matching / Perfect matching")
	checkMatchButton.Connect("toggled", func() {
		if checkMatchButton.GetActive() {
			arg.perfectMatch = false
		} else {
			arg.perfectMatch = true
		}
	})
	buttons.Add(checkMatchButton)

	combobox := gtk.NewComboBoxText()
	for _, delim := range delimName {
		combobox.AppendText(delim)
	}
	combobox.SetActive(0)
	combobox.Connect("changed", func() {
		fmt.Println("value:", combobox.GetActiveText())
		arg.delim = combobox.GetActiveText()
	})
	buttons.Add(combobox)

	//--------------------------------------------------------
	// GtkTextView
	//--------------------------------------------------------
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)
	textview := gtk.NewTextView()
	// var start, end gtk.TextIter
	var end gtk.TextIter
	buffer := textview.GetBuffer()
	swin.Add(textview)

	//--------------------------------------------------------
	// Run button
	//--------------------------------------------------------
	runbutton := gtk.NewButtonWithLabel("Run")
	runbutton.Clicked(func() {
		num, err := getKeysearchWords(arg)
		buffer.GetStartIter(&end)
		if err != nil {
			log.Println(err)
			buffer.Insert(&end, err.Error()+"\n")
		} else {
			buffer.Insert(&end, "Results n="+strconv.Itoa(num)+"\n")
		}
	})
	buttons.Add(runbutton)
	framebox2.PackStart(buttons, false, false, 0)

	//--------------------------------------------------------
	// GtkVSeparator
	//--------------------------------------------------------
	vsep := gtk.NewVSeparator()
	framebox2.PackStart(vsep, false, false, 0)

	//--------------------------------------------------------
	// GtkTextView
	//--------------------------------------------------------
	framebox2.Add(swin)

	// buffer.Connect("changed", func() {
	// 	// fmt.Println("changed")
	// })

	//--------------------------------------------------------
	// GtkMenuItem
	//--------------------------------------------------------
	cascademenu := gtk.NewMenuItemWithMnemonic("_File")
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	var menuitem *gtk.MenuItem
	menuitem = gtk.NewMenuItemWithMnemonic("_Exit")
	menuitem.Connect("activate", func() {
		gtk.MainQuit()
	})
	submenu.Append(menuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_View")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	checkmenuitem := gtk.NewCheckMenuItemWithMnemonic("_Disable")
	checkmenuitem.Connect("activate", func() {
		vpaned.SetSensitive(!checkmenuitem.GetActive())
	})
	submenu.Append(checkmenuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_Help")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	menuitem = gtk.NewMenuItemWithMnemonic("_About")
	menuitem.Connect("activate", func() {
		dialog := gtk.NewAboutDialog()
		dialog.SetName("MADscan")
		dialog.SetProgramName("MADscan")
		dialog.SetAuthors(authors())
		dialog.SetLicense("GPL v3")
		dialog.SetWrapLicense(true)
		dialog.Run()
		dialog.Destroy()
	})
	submenu.Append(menuitem)

	//--------------------------------------------------------
	// GtkStatusbar
	//--------------------------------------------------------
	statusbar := gtk.NewStatusbar()
	context_id := statusbar.GetContextId("MADscan v0")
	statusbar.Push(context_id, "Simple search GUI")

	framebox2.PackStart(statusbar, false, false, 0)

	//--------------------------------------------------------
	// Event
	//--------------------------------------------------------
	window.Add(vbox)
	window.SetSizeRequest(600, 600)
	window.ShowAll()
	gtk.Main()
}
