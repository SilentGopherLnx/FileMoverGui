# About this:
File copy and move command-line operations with GTK-based GUI status information

This app is used by this project https://github.com/SilentGopherLnx/GopherFileManager

I separated file's copy, move operations into this app for stable work

This app do file operations and shows progress by GTK-gui.

It also make check is you can do this. App can ask you some classic questions (also by GTK-gui):
- This file already exist in destination folder. Replace it?
- Can't read (or write) file. Try again or ignore?
- Destination foler already exist. Merge this folders?
- There are not enought space on disk. Try again?

and so on...

![test picture](https://github.com/SilentGopherLnx/screenshots_and_binaries/blob/master/SCREENS_GopherFileManagerFileMoverGui/mover_01.png)
![test picture](https://github.com/SilentGopherLnx/screenshots_and_binaries/blob/master/SCREENS_GopherFileManagerFileMoverGui/mover_02.png)

# Dependencies for GOPATH:
1) Golang GTK3 lib
https://github.com/gotk3/gotk3
also for gtk:
> sudo apt-get install libgtk-3-dev
>
> sudo apt-get install libcairo2-dev
>
> sudo apt-get install libglib2.0-dev
2) my Framework
https://github.com/SilentGopherLnx/easygolang
(this package also have some unnessesary sub-packages)

# Status:
App is under development and looks freaky (you can see version if run with "-v" argument)

**Not all functions are implemented and realised as planned!** such as:
- ask if folders gonna be merged, else use "os.Rename()" for full folder tree (not file by file!)
- detail information for compare existing files when ask for replace
- copy "modified time" files
- writing files to temporary name (half-copied files will not be looked as "ok")
- free space check
- ask for options for symlink...

**News**
- 0.2.0 - multi-language support (english default, russian on config file too), renaming now here too
- 0.2.1 - updated for **gotk3 0.6.1, golang 1.17**
- 0.2.2 - percent before name on title

# Platform & License:
**Only Linux!** Tested only on amd64 on Cinnamon desktop of Linux Mint.

Windows support is NOT planned

**License type is GPL3**

# Usage example for copy of two files:
**golang:**
> exec.Command("FileMoverGui", "-cmd", "copy", "-src", "file:///path_src/file1" + "\n" + "file:///path_src/file2", "-dst", "file:///path_dst/")

**terminal:**
> $ ./FileMoverGui -cmd copy -src "file:///path_src/file1
>
> \> file:///path_src/file2" -dst "file:///path_dst/"
(yes. This is new line charater. You can input this if you have unclosed brakets)

# Args for running:
**-cmd**
> copy - copy files list from **-src** to **-dst** folder
>
> move - rename files list from **-src** to -dst folder or move files if -src and -dst are on different disks or network folders
>
> delete - delete files list from **-src**
>
> clear - delete files in **-src** folders (files in list will be ignored, folders will be cleared) (not implemented!)
>
> rename - rename dialog window for  **-src** file

**-src** is files and/or folders list, separated by new line symbol (not "\n" string)

**-dst** is always folder

**-buf** buffer size of bytes for file copy operations. Value will be multiplied by 1024 

**-lang** language code (in config file, "en" (english) is default)

**-v** you can see version of this app if run with only this one argument

# -src & -dst path format:
I tried to use url file scheme from Copy/Paste and Drag&Drop operations.

Disk files have "file://" prefix

Network samba share have "smb://" prefix

Also are supportted: 
- "mtp://" (smartphones), "gphoto2://" (photos)
- "ftp://", "dav://" & "ftps://", "davs://" (ftp protocol, webdaw) 

Space characters are replaed to "%20"

Non-english language charactes are escaped like url with % and code like in http urls

You can test format by copying file in your file manager and printing result from "xclip"
> xclip -o -selection clipboard -t "x-special/gnome-copied-files"
