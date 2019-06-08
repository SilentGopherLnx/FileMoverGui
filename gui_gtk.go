package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"

	"github.com/gotk3/gotk3/gtk"
	//"github.com/gotk3/gotk3/gdk"
)

var win *gtk.Window

var header *gtk.HeaderBar

var lbl_src_size *gtk.Label
var lbl_dst_free *gtk.Label
var lbl_src_files *gtk.Label
var lbl_speed *gtk.Label
var lbl_timepassed *gtk.Label
var lbl_timeleft *gtk.Label
var lbl_done *gtk.Label
var lbl_current *gtk.Label

var progress *gtk.ProgressBar
var spinner *gtk.Spinner

var title_prev_perc = "?"

func GUI_Init() {
	gtk.Init(nil)
}

func GUI_Create() {
	var err error
	win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		Prln("Unable to create window:") // + err)
	}
	win.SetDefaultSize(400, 200)
	win.SetPosition(gtk.WIN_POS_CENTER)

	win.Connect("destroy", func() {
		AppExit(0)
	})

	appdir := FolderLocation_App()
	win.SetIconFromFile(appdir + "gui/icon_operation.png")

	icon_oper, _ := gtk.ImageNewFromFile(appdir + "gui/icon_operation.png")
	btn_icon, _ := gtk.ButtonNewWithLabel("")
	btn_icon.SetImage(icon_oper)
	btn_icon.SetProperty("always-show-image", true)
	btn_icon.Connect("button-press-event", func() {
		Dialog_About(win, AppVersion(), AppAuthor(), AppMail(), AppRepository(), GetFlag_Russian())
	})

	//img := GTK_Image_From_File(appdir+"gui/button_abort.png", "png")
	img1 := GTK_Image_From_Name("process-stop", gtk.ICON_SIZE_BUTTON)

	btn_close, _ := gtk.ButtonNewWithLabel("")
	btn_close.SetImage(img1)
	btn_close.SetProperty("always-show-image", true)
	btn_close.Connect("button-press-event", func() {
		AppExit(0)
	})

	header, _ = gtk.HeaderBarNew()
	header.PackStart(btn_icon)
	header.PackEnd(btn_close)
	win.SetTitlebar(header)

	title_func(0)

	// ========================

	box_src_disk, _ := GTK_LabelPair("Disk SRC: ", src_disk)
	box_dst_disk, _ := GTK_LabelPair("Disk DST: ", dst_disk)

	// ========================

	lbl_src_folder_title, _ := gtk.LabelNew("SRC:")
	lbl_src_folder_title.SetMarkup("<b>Source folder:</b>")
	lbl_src_folder_title.SetHExpand(true)
	lbl_src_folder_title.SetVExpand(false)
	lbl_src_folder_title.SetHAlign(gtk.ALIGN_START)

	lbl_src_folder, _ := gtk.LabelNew(src_folder)
	lbl_src_folder.SetHExpand(true)
	lbl_src_folder.SetHAlign(gtk.ALIGN_START)
	//lbl_src_folder.SetJustify(gtk.JUSTIFY_LEFT)
	lbl_src_folder.SetVAlign(gtk.ALIGN_START)
	GTK_LabelWrapMode(lbl_src_folder, 1)

	// ========================

	lbl_src_title, _ := gtk.LabelNew("Selected files:")
	lbl_src_title.SetMarkup("<b>Selected files:</b>")
	lbl_src_title.SetHExpand(true)
	lbl_src_title.SetHAlign(gtk.ALIGN_START)

	lbl_src, _ := gtk.LabelNew(src_names)
	lbl_src.SetHExpand(true)
	lbl_src.SetVAlign(gtk.ALIGN_START)
	lbl_src.SetHAlign(gtk.ALIGN_START)
	//lbl_src.SetJustify(gtk.JUSTIFY_LEFT)
	GTK_LabelWrapMode(lbl_src, MAXI(1, len(path_src)))

	scroll_scr, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll_scr.SetVExpand(true)
	scroll_scr.SetHExpand(true)
	scroll_scr.Add(lbl_src)
	//scroll_scr.SetOverlayScrolling(true)
	//scroll_scr.SetSizeRequest()

	frame, _ := gtk.FrameNew("Selected files:")
	frame.SetLabelWidget(lbl_src_title)
	frame.Add(scroll_scr)

	// ========================

	lbl_dst_folder_title, _ := gtk.LabelNew("DST:")
	lbl_dst_folder_title.SetMarkup("<b>Destination folder:</b>")
	lbl_dst_folder_title.SetHExpand(true)
	lbl_dst_folder_title.SetVExpand(false)
	lbl_dst_folder_title.SetHAlign(gtk.ALIGN_START)

	lbl_dst_folder, _ := gtk.LabelNew(path_dst.GetVisual())
	lbl_dst_folder.SetHExpand(true)
	lbl_dst_folder.SetHAlign(gtk.ALIGN_START)
	//lbl_dst.SetJustify(gtk.JUSTIFY_LEFT)
	lbl_dst_folder.SetVAlign(gtk.ALIGN_START)
	GTK_LabelWrapMode(lbl_dst_folder, 1)

	// ========================

	lbl_separator1, _ := gtk.LabelNew(" ")

	calc := "calculating..."

	var box_src_size, box_src_files, box_dst_free *gtk.Box
	box_src_size, lbl_src_size = GTK_LabelPair("Selected total size: ", calc)
	box_src_files, lbl_src_files = GTK_LabelPair("Selected objects: ", calc)
	box_dst_free, lbl_dst_free = GTK_LabelPair("Destination free space: ", calc)

	lbl_separator2, _ := gtk.LabelNew(" ")

	// ========================

	var box_speed, box_timepassed, box_timeleft, box_done *gtk.Box
	box_speed, lbl_speed = GTK_LabelPair("Speed: ", calc)
	box_timepassed, lbl_timepassed = GTK_LabelPair("Time passed: ", calc)
	box_timeleft, lbl_timeleft = GTK_LabelPair("Time left: ", calc)
	box_done, lbl_done = GTK_LabelPair("Done: ", calc)

	// ========================

	progress, _ = gtk.ProgressBarNew()

	lbl_current, _ = gtk.LabelNew(" ")
	lbl_current.SetHExpand(true)
	lbl_current.SetHAlign(gtk.ALIGN_START)
	GTK_LabelWrapMode(lbl_current, 1)

	spinner, _ = gtk.SpinnerNew()
	spinner.Start()

	box_current, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	box_current.Add(spinner)
	box_current.Add(lbl_current)

	// ========================

	grid, _ := gtk.GridNew()
	grid.SetColumnSpacing(10)
	grid.SetColumnHomogeneous(true)

	gui_w := 1
	if oper_single {
		gui_w = 2
	}

	grid.Attach(box_src_disk, 0, 0, gui_w, 1)
	//grid.Attach(lbl_src_name, 0, 1, gui_w, 1)
	//grid.Attach(frame, 0, 1, gui_w, 2)
	grid.Attach(lbl_src_folder_title, 0, 1, 1, 1)
	grid.Attach(lbl_src_folder, 0, 2, 1, 1)
	grid.Attach(frame, 0, 3, 2, 1)

	if !oper_single {
		grid.Attach(box_dst_disk, 1, 0, 1, 1)
		grid.Attach(lbl_dst_folder_title, 1, 1, 1, 1)
		grid.Attach(lbl_dst_folder, 1, 2, 1, 1)
		//grid.Attach(box_vdst, 1, 1, 1, 2)
	}

	grid.Attach(lbl_separator1, 0, 3, gui_w, 1)
	grid.Attach(box_src_size, 0, 4, gui_w, 1)
	grid.Attach(box_src_files, 0, 5, gui_w, 1)
	if !oper_single {
		grid.Attach(box_dst_free, 1, 4, 1, 1)
	}

	grid.Attach(lbl_separator2, 0, 10, 2, 1)

	grid.Attach(box_timepassed, 0, 11, 1, 1)
	grid.Attach(box_timeleft, 1, 11, 1, 1)
	grid.Attach(box_done, 0, 12, 1, 1)
	grid.Attach(box_speed, 1, 12, 1, 1)

	grid.Attach(progress, 0, 15, 2, 1)
	//grid.Attach(btn_abort, 0, 16, 2, 1)
	grid.Attach(box_current, 0, 16, 2, 1)

	win.Add(grid)

	win.ShowAll()

}

func title_func(perc float64) {
	new_txt := ""
	if perc > 0.0 {
		new_txt = " - " + F2S(perc, 2) + "%"
	}
	if new_txt != title_prev_perc {
		txt := StringTitle(operation) + new_txt
		win.SetTitle(txt)
		header.SetTitle(txt)
		title_prev_perc = new_txt
	}
}

func GUI_Iteration() {
	gtk.MainIteration()

	sizebytes := src_size.Get()
	lbl_src_size.SetText(FileSizeNiceString(sizebytes) + " (" + I2Ss(sizebytes) + " bytes)")
	if !oper_single {
		freebytes := dst_free.Get()
		lbl_dst_free.SetText(FileSizeNiceString(freebytes)) //+ " (" + I2Ss(freebytes) + " bytes)")
	}

	sel_files := src_files.Get()
	sel_folders := src_folders.Get()
	lbl_src_files.SetText(I2S64(sel_files+sel_folders) + " (" + I2S64(sel_files) + " files & " + I2S64(sel_folders) + " folders)")

	passed := 0.0
	w := work.Get()
	if w && time_start == nil {
		time_ := TimeNow()
		time_start = &time_
	}
	if time_start != nil && !w && time_end == nil {
		time_ := TimeNow()
		time_end = &time_
	}
	if time_start != nil {
		if time_end == nil {
			passed = TimeSeconds(*time_start)
		} else {
			passed = TimeSecondsSub(*time_start, *time_end)
		}
	}

	sizedone := done_bytes.Get()
	perc := float64(sizedone) / float64(sizebytes)

	tleft := "0"
	if w && sizedone > 0.0 && passed > 1.0 {
		tleft = F2S(float64(sizebytes-sizedone)*passed/float64(sizedone), 0)
		speed := F2S(float64(sizedone)/passed/float64(BytesInMb), 2)
		lbl_speed.SetText(speed + " MB/s")
	}

	lbl_timepassed.SetText(I2S(int(passed)) + " seconds")
	lbl_timeleft.SetText(tleft + " seconds")

	lbl_done.SetText(StringEnd("   "+F2S(perc*100.0, 2), 6) + "% (" + I2Ss(sizedone) + " bytes) " + I2S64(done_files.Get()) + " objects")

	current_text := current_file.Get()
	lbl_current.SetText(current_text)

	if perc >= 0.0 {
		progress.SetFraction(perc)
		title_func(perc * 100.0)
	}

	if !work.Get() && spinner != nil && GTK_SpinnerActive(spinner, true) {
		Prln("spinner off!")
		spinner.Stop()
	}

}

func GUI_Warn_SrcUnread(pre_read_errs string) {
	err_txt := "Not all files can be read:"
	dial := gtk.MessageDialogNew(nil, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, err_txt+"\n"+pre_read_errs)
	dial.SetMarkup("<b>" + HtmlEscape(err_txt) + "</b>\n" + HtmlEscape(pre_read_errs))
	dial.SetTitle("Some problems?")
	dial.Connect("destroy", func() {
		AppExit(0)
	})
	resp := dial.Run()
	if resp == gtk.RESPONSE_OK {
		dial.Close()
		AppExit(0)
	} else {
		AppExit(0)
	}
}

func GUI_Warn_SrcDstEqual(path_folder string) {
	err_txt := "Can't move folder inside itself:"
	dial := gtk.MessageDialogNew(nil, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, err_txt+"\n"+path_folder)
	dial.SetMarkup("<b>" + HtmlEscape(err_txt) + "</b>\n" + HtmlEscape(path_folder))
	dial.SetTitle("Invalid operation!")
	dial.Connect("destroy", func() {
		AppExit(0)
	})
	resp := dial.Run()
	if resp == gtk.RESPONSE_OK {
		dial.Close()
		AppExit(0)
	} else {
		AppExit(0)
	}
}

func GUI_Warn_SrcDelete(path_src_folder string, path_src_names string, clear_mode bool) {
	msg := "Delete"
	if clear_mode {
		msg = "Clear"
	}
	reported := NewAtomicBool(false, [2]string{"", ""})
	dial := gtk.MessageDialogNew(nil, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK_CANCEL, msg+"? :")
	dial.SetTitle(msg + " files?")
	//dial.SetMarkup("<b>" + msg + "</b> :\n" + HtmlEscape(path_src_names) + "<b>?</b>\nin folder:\n" + HtmlEscape(path_src_folder))
	dial.SetMarkup("<b>" + msg + " at path:</b>\n" + HtmlEscape(FolderPathEndSlash(path_src_folder)))
	dial.SetDefaultResponse(gtk.RESPONSE_OK)
	dial.Connect("destroy", func() {
		if !reported.Get() {
			AppExit(0)
		}
	})

	this_txt := "this folder(s)/file(s)? :"
	lbl_name, _ := gtk.LabelNew(this_txt)
	lbl_name.SetMarkup("<b>" + HtmlEscape(this_txt) + "</b>")
	frame, _ := gtk.FrameNew(this_txt)
	frame.SetLabelWidget(lbl_name)
	//frame.SetHExpand(true)
	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	//scroll.SetVExpand(true)
	//scroll.SetHExpand(true)
	text, _ := gtk.LabelNew(path_src_names)
	//text.SetHExpand(true)
	scroll.Add(text)
	frame.Add(scroll)
	area, _ := dial.GetContentArea()
	area.SetSpacing(0)
	area.Add(frame)
	area.ShowAll()

	resp := dial.Run()
	if resp == gtk.RESPONSE_OK {
		reported.Set(true)
		dial.Close()
	} else {
		AppExit(0)
	}
}

func GUI_Ask_File(q FileInteractiveRequest, cmd chan FileInteractiveResponse) {

	reported := NewAtomicBool(false, [2]string{"", ""})

	dial := gtk.MessageDialogNew(nil, gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_WARNING, gtk.BUTTONS_NONE, q.FileName)

	att := ""
	if q.Attempt > 1 {
		att = " [attempt# " + I2S(q.Attempt) + "]"
	}
	dial.SetTitle("Problem" + att)

	qs := [2]string{"", ""}
	switch q.AskType {
	case FILE_INTERACTIVE_ASK_EXIST:
		qs[0], qs[1] = "File exist at:", "Replace?"
	case FILE_INTERACTIVE_ASK_ERROR:
		qs[0], qs[1] = "File problem with:", "Try again?"
	case FILE_INTERACTIVE_ASK_PANIC:
		qs[0], qs[1] = "File problem with:", "Unsolvable =("
	}
	dial.SetMarkup("<b>" + HtmlEscape(qs[0]) + "</b>\n" + HtmlEscape(q.FileName) + "\n<b>" + HtmlEscape(qs[1]) + "</b>")

	choice, _ := gtk.CheckButtonNewWithLabel("Remember choice")

	area, _ := dial.GetContentArea()
	area.SetSpacing(0)
	if q.Attempt < 2 && q.AskType != FILE_INTERACTIVE_ASK_PANIC {
		area.Add(choice)
	}
	area.ShowAll()

	if q.AskType == FILE_INTERACTIVE_ASK_PANIC { // stop/skip?
		dial.SetDefaultResponse(gtk.RESPONSE_NO)
		dial.AddButton("OK", gtk.RESPONSE_NO)
	} else {
		dial.SetDefaultResponse(gtk.RESPONSE_YES)
		dial.AddButton("Yes", gtk.RESPONSE_YES)
		if q.AskType == FILE_INTERACTIVE_ASK_EXIST {
			dial.AddButton("Save with new name", gtk.RESPONSE_ACCEPT)
		}
		dial.AddButton("No", gtk.RESPONSE_NO)
	}
	dial.Connect("destroy", func() {
		if !reported.Get() {
			cmd <- FileInteractiveResponse{SaveChoice: false, Command: FILE_INTERACTIVE_SKIP}
		}
	})

	resp := dial.Run()

	save := choice.GetActive()

	if resp == gtk.RESPONSE_YES {
		reported.Set(true)
		cmd <- FileInteractiveResponse{SaveChoice: save, Command: FILE_INTERACTIVE_RETRY}
	}
	if resp == gtk.RESPONSE_ACCEPT {
		reported.Set(true)
		cmd <- FileInteractiveResponse{SaveChoice: save, Command: FILE_INTERACTIVE_NEWNAME}
	}
	if resp == gtk.RESPONSE_NO {
		reported.Set(true)
		cmd <- FileInteractiveResponse{SaveChoice: save, Command: FILE_INTERACTIVE_SKIP}
	}
	dial.Close()
}
