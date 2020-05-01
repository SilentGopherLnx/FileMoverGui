package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easylinux"

	"flag"
	"os"
)

const OPERATION_DEMO = "demo"
const OPERATION_COPY = "copy"
const OPERATION_MOVE = "move"
const OPERATION_DELETE = "delete"
const OPERATION_CLEAR = "clear"
const OPERATION_RENAME = "rename"

//const OPERATION_NEWFILE = "newfile"
//const OPERATION_NEWFOLDER = "newfolder"

var operation string
var oper_single = false
var oper_demo = false

var path_src []*LinuxPath = []*LinuxPath{}
var files_src []os.FileInfo = []os.FileInfo{}
var path_src_visual string = ""
var path_dst *LinuxPath = NewLinuxPath(false)

var src_size *AInt64 = NewAtomicInt64(0)
var src_files *AInt64 = NewAtomicInt64(0)
var src_folders *AInt64 = NewAtomicInt64(0)
var src_unread *AInt64 = NewAtomicInt64(0)
var src_irregular *AInt64 = NewAtomicInt64(0)
var src_mount *AInt64 = NewAtomicInt64(0)
var src_symlinks *AInt64 = NewAtomicInt64(0)
var dst_free *AInt64 = NewAtomicInt64(0)
var done_bytes *AInt64 = NewAtomicInt64(0)
var done_fobjects *AInt64 = NewAtomicInt64(0)
var current_file *AString = NewAtomicString("")

var work = NewAtomicBool(false)

var time_start *Time = nil
var time_end *Time = nil

var BUFFER_SIZE = 1024 * 64

var mount_list [][2]string

var gui_chan_cmd chan FileInteractiveResponse = make(chan FileInteractiveResponse)
var gui_chan_ask chan FileInteractiveRequest = make(chan FileInteractiveRequest)

var src_disk, dst_disk string

var src_folder = ""
var src_names = ""

var langs *LangArr

func init() {
	Prln("file mover  - inited")
	AboutVersion(AppVersion())
	langs = InitLang(FolderPathEndSlash(FolderLocation_App()+"localization") + "translation_mover.cfg")
}

func main() {
	var path_src_url string = ""
	var path_dst_url string = ""
	var buf_1024 = 0
	var lang string = ""

	flag.StringVar(&operation, "cmd", OPERATION_DEMO, "copy,move,delete,clear")
	flag.StringVar(&path_src_url, "src", "", "source url array (new line is separator) - file:// or smb:// and ....")
	flag.StringVar(&path_dst_url, "dst", "", "destination url folder file:// or smb:// and ....")
	flag.StringVar(&lang, "lang", DEFAULT_LANG, "en or another lang")
	flag.IntVar(&buf_1024, "buf", 64, "buffer size in bytes (value will be multiplied by 1024!!)")
	flag.Parse()

	if !flag.Parsed() {
		Prln("wrong flags")
		AppExit(1)
	}
	Prln("file mover - args parsed")

	langs.SetLang(lang)

	BUFFER_SIZE = MINI(256, MAXI(1, buf_1024)) * 1024

	if operation == OPERATION_DEMO {
		oper_demo = true
		// operation = "copy"
		// path_src_url = "file:///mnt/dm-1/golang/my_code/FileMoverGui/test_dir/file1.txt"
		// path_dst_url = "file:///mnt/dm-1/golang/my_code/FileMoverGui/test_dir/to/"

		// operation = "copy"
		// path_src_url = "file:///mnt/dm-1/golang/my_code/GopherFileManager/test_dir/New Folder/"
		// path_dst_url = "file:///mnt/dm-1/golang/my_code/GopherFileManager/test_dir/щ/"
	} else {

	}

	operation = StringDown(StringTrim(operation))
	oper_arr := []string{OPERATION_COPY, OPERATION_MOVE, OPERATION_DELETE, OPERATION_CLEAR, OPERATION_RENAME}
	if StringInArray(operation, oper_arr) == -1 {
		Prln("Wrong operation command")
		AppExit(2)
	}

	path_dst_url = StringTrim(path_dst_url)
	path_dst.SetUrl(path_dst_url)

	path_src_arr := StringSplitLines(path_src_url)
	for j := 0; j < len(path_src_arr); j++ {
		new_path := NewLinuxPath(false)
		path_src_arr[j] = StringTrim(path_src_arr[j])
		if len(path_src_arr[j]) > 0 {
			new_path.SetUrl(path_src_arr[j])
			if new_path.GetParseProblems() {
				Prln("Parsing problem [src]:" + path_src_arr[j])
				AppExit(3)
			}
			path_src = append(path_src, new_path)
		}
	}
	if len(path_src) == 0 {
		Prln("File list is empty!")
		AppExit(4)
	}

	if operation == OPERATION_DELETE || operation == OPERATION_CLEAR || operation == OPERATION_RENAME {
		oper_single = true
	} else {
		if len(path_dst_url) == 0 {
			Prln("Destination is empty!")
			AppExit(5)
		}
		if path_dst.GetParseProblems() {
			Prln("Parsing problem [dst]:" + path_dst_url)
			AppExit(6)
		}
	}

	mount_list = LinuxGetMountList()
	if len(path_src) > 0 {
		src_disk, _ = LinuxFilePartition(mount_list, path_src[0].GetReal())
		for j := 1; j < len(path_src); j++ {
			src_disk_1, _ := LinuxFilePartition(mount_list, path_src[j].GetReal())
			if src_disk_1 != src_disk {
				Prln("Src file at [" + I2S(j+1) + "] position have not same disk with first file:" + src_disk)
				AppExit(7)
			}
		}
		if src_disk == "" {
			Prln("Src disk not detected")
			AppExit(8)
		}
	}
	if !oper_single {
		dp := ""
		dst_disk, dp = LinuxFilePartition(mount_list, path_dst.GetReal())
		if dst_disk == "" {
			Prln("Dst disk not detected")
			AppExit(9)
		}

		//go func() {
		dp = FilePathEndSlashRemove(dp)
		Prln(">>" + dp)
		dst_free.Set(LinuxFolderFreeSpace(dp))
		//}()
	}

	GUI_Init()

	pre_read_errs := ""
	for j := 0; j < len(path_src); j++ {
		problem := false
		real1 := path_src[j].GetReal()
		real2, err := FileEvalSymlinks(real1)
		if err == nil {
			if FilePathEndSlashRemove(real1) != FilePathEndSlashRemove(real2) {
				path_src[j].SetReal(real2)
			}
			finfo, err := FileInfo(real2, false)
			if err == nil {
				files_src = append(files_src, finfo)
			} else {
				problem = true
			}
		} else {
			problem = true
		}
		if problem {
			pre_read_errs += path_src[j].GetVisual() + "\n"
		}
	}

	if len(pre_read_errs) > 0 {
		GUI_Warn_SrcUnread(pre_read_errs)
	}

	src_names = ""
	visual_arr := []string{}
	for j := 0; j < len(path_src); j++ {
		visual_arr = append(visual_arr, path_src[j].GetVisual())
		if files_src[j].IsDir() {
			src_names += FolderPathEndSlash(path_src[j].GetLastNode()) + "\n"
		} else {
			src_names += FilePathEndSlashRemove(path_src[j].GetLastNode()) + "\n"
		}
		if FolderPathEndSlash(path_src[j].GetReal()) == FolderPathEndSlash(path_dst.GetReal()) {
			GUI_Warn_SrcDstEqual(path_src[j].GetVisual())
		}
	}
	src_folder = FolderPathEndSlash(LinuxFileGetParent(path_src[0].GetReal()))
	path_src_visual = StringJoin(visual_arr, "\n")

	if oper_single {
		if operation == OPERATION_DELETE || operation == OPERATION_CLEAR {
			GUI_Warn_SrcDelete(LinuxFileGetParent(path_src[0].GetReal()), src_names, operation == OPERATION_CLEAR)
		} else {
			if operation == OPERATION_RENAME { // || operation == OPERATION_NEWFILE ||  operation == OPERATION_NEWFOLDER{
				fpath2 := FolderPathEndSlash(src_folder)
				old_name := files_src[0].Name()
				f_rename := func(new_name string) (bool, string) {
					safe_name := new_name
					// ADD CHECK!!!!!
					// Windows (FAT32, NTFS): Any Unicode except NUL, \, /, :, *, ", <, >, |
					// Mac(HFS, HFS+): Any valid Unicode except : or /
					// Linux(ext[2-4]): Any byte except NUL or /
					if safe_name != old_name {
						ok, err_txt := FileRename(fpath2+old_name, fpath2+safe_name)
						if ok {
							return true, ""
						} else {
							return false, err_txt
						}
					}
					return false, "?"
				}
				GUI_FileRename(fpath2, old_name, f_rename)
			}
		}
	}

	if operation != OPERATION_RENAME { // || operation == OPERATION_NEWFILE ||  operation == OPERATION_NEWFOLDER{
		go oper_switch_runner()
		go speed_counter()
		GUI_Create()

		//timelast := TimeNow()
		for {
			GUI_Iteration()

			select {
			case q := <-gui_chan_ask:
				GUI_Ask_File(q, gui_chan_cmd)
			default:
				//CONTINUE!!
			}

			//q, ok := ChanGetOrSkip(gui_chan_cmd.(chan interface{}))
			// if ok {
			// 	str := q.(string)
			// 	GUI_Ask_File(str, gui_chan_cmd)
			// }

			// if TimeSeconds(timelast) > 0.5 {
			// 	//win.ShowAll()
			// 	timelast = TimeNow()
			// }
			RuntimeGosched()
		}
	}
}

func oper_switch_runner() {
	Prln("what is the command: " + operation)
	for j := 0; j < len(path_src); j++ {
		FoldersRecursively_Size(mount_list, files_src[j], path_src[j].GetReal(), src_size, src_files, src_folders, src_unread, src_irregular, src_mount, src_symlinks, nil)
	}
	Prln("starting command...")
	work.Set(true)
	switch operation {
	case OPERATION_COPY:
		dst_real := path_dst.GetReal()
		cmd_saved := ""
		for j := 0; j < len(path_src); j++ {
			cmd_saved = FoldersRecursively_Copy(mount_list, files_src[j], path_src[j].GetReal(), dst_real, done_bytes, done_fobjects, BUFFER_SIZE, gui_chan_cmd, gui_chan_ask, current_file, cmd_saved)
		}
	case OPERATION_MOVE:
		dst_real := path_dst.GetReal()
		cmd_saved := ""
		for j := 0; j < len(path_src); j++ {
			//simple := false
			cmd_saved, _ = FoldersRecursively_Move(mount_list, files_src[j], path_src[j].GetReal(), dst_real, done_bytes, done_fobjects, BUFFER_SIZE, gui_chan_cmd, gui_chan_ask, current_file, cmd_saved, src_disk == dst_disk)
			// if simple {
			// 	FoldersRecursively_Size(mount_list, files_src[j], path_src[j].GetReal(), src_size, src_files, src_folders, src_unread, src_irregular, nil, nil, nil)
			// }
		}
	case OPERATION_DELETE:
		for j := 0; j < len(path_src); j++ {
			FoldersRecursively_Delete(mount_list, files_src[j], path_src[j].GetReal(), done_bytes, done_fobjects, current_file, false)
		}
	case OPERATION_CLEAR:
		for j := 0; j < len(path_src); j++ {
			FoldersRecursively_Delete(mount_list, files_src[j], path_src[j].GetReal(), done_bytes, done_fobjects, current_file, true)
		}
	}
	work.Set(false)
	if done_bytes.Get() == src_size.Get() && done_fobjects.Get() == src_files.Get()+src_folders.Get() {
		Prln("done")
		AppExit(0)
	} else {
		Prln("done, waiting..." + I2S64(done_fobjects.Get()) + "/" + I2S64(src_files.Get()))
	}
}

var done_new int64
var done_prev int64
var copy_speed *AInt64 = NewAtomicInt64(0)

func speed_counter() {
	for {
		done_new = done_bytes.Get()
		done_delta := done_new - done_prev
		copy_speed.Set(done_delta)
		done_prev = done_new
		SleepMS(1000)
	}
}
