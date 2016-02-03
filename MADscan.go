package main

import (
	"fmt"
	// "github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"log"
	"io/ioutil"
	"os"
	// "os/exec"
	"path/filepath"
	// "regexp"
	// "sort"
	"strings"
	"strconv"
)

var (
	delimName = []string{"Select deliminator", "tab (\\t)", "space ( )", "comma (,)", "period (.)", "colon (:)"}
	delimList = []string{"\t", "\t", " ", ",", ".", ":"}
)

// minimum size for the array size of keywords
const minSize = 50

func authors() []string {
	return []string{"carushi<l.cawaguchi@gmail.com>\n\nReference:"}
}

type arguments struct {
	column     int
	inputPath  string
	filterPath string
	outputPath string
	ignoreCase bool
	delim      string
}

func (a *arguments) setDeliminator() {
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

func getKeywords(key string, ignoreCase bool) ([]string, error) {
	lines, err := ioutil.ReadFile(key)
	keyList := make([]string, 0, minSize)
	if err != nil {
		return keyList, err
	}
	for _, gene := range strings.Split(string(lines), "\n") {
		gene = strings.TrimRight(gene, "\r")
		if len(gene) == 0 {
			continue
		}
		if ignoreCase {
			keyList = append(keyList, strings.ToUpper(gene))
		} else {
			keyList = append(keyList, gene)
		}
	}
	return keyList, err

}

func searchKeywords(column int, inputPath string, delim string, keyList []string, ignoreCase bool) ([]string, error) {
	matchedLines := make([]string, 0, minSize)
	lines, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return matchedLines, err
	}
	for _, line := range strings.Split(string(lines), "\n") {
		s := strings.Split(strings.TrimRight(line, "\r"), delim)
		if len(s) <= column {
			continue
		}
		if ignoreCase {
			s[column] = strings.ToUpper(s[column])
		}
		for _, key := range keyList {
			if s[column] == key {
				matchedLines = append(matchedLines, line)
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
	arg.setDeliminator()
	fmt.Println("Read: "+arg.filterPath)
	keyList, err := getKeywords(arg.filterPath, arg.ignoreCase)
	if err != nil {
		return -1, err
	}
	fmt.Println("Read: "+arg.inputPath)
	matchedLines, err := searchKeywords(arg.column, arg.inputPath, arg.delim, keyList, arg.ignoreCase)
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
		fmt.Println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")

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
	dir, _ := filepath.Split(os.Args[0])
	imagefile := filepath.Join(dir, "../image/logo.png")
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
		filepath.Join(dir, "../data/Acetylation_site_dataset"),
		filepath.Join(dir, "../data/Ras_gene_list.txt"),
		filepath.Join(dir, "../data/output.txt"),
		false,
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
	menuitem = gtk.NewMenuItemWithMnemonic("E_xit")
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
		dialog.SetLicense("The library is available under the same terms and conditions as the Go, the BSD style license, and the LGPL (Lesser GNU Public License). The idea is that if you can use Go (and Gtk) in a project, you should also be able to use go-gtk.")
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
