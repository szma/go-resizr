go-resizr
=========

Resizr is a small utility that resizes for all jpg images in a directory (and its subfolders)
and saves them in a different destination directory.

This is useful for me to have a small version of the pictures on a USB stick, or show them on my phone from a NAS over the internet or even on TV in the local LAN. 

This could easily be done with a bash script as well, but I wanted to learn a little bit of Go and this was my toy example. See Disclaimer.

Disclaimer
----------

This is my first Go program. It might be full of errors and the programming style might not be good at all.
I used Go 1.9, but I'm pretty sure it will run with lower version's as well.

Feel free to 

Installation
------------

```bash
$ go get github.com/szma/go-resizr
```

If the $GOPATH/bin folder is in your $PATH run


```bash
$ resizr help
$ resizr --dest /home/user/previewpictures --size 1024 /home/user/pictures
```

This stores all your jpgs from /home/user/pictures in /home/user/previewpictures resized with maximum width/height of 1024.

