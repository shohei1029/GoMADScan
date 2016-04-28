# MADscan


[![Build Status](https://drone.io/github.com/carushi/MADscan/status.png)](https://drone.io/github.com/carushi/MADscan/latest)

Modification associated database scanner based on GUI inteface using go-gtk.

![](image/window.png)

* It is intended to achieve user-friendly keyword searching for character separated files such as modification dataset available in [PhosphoSitePlus](http://www.phosphosite.org/homeAction.action).
* MADscan also supports '\n', '\r', and '\r\n' as a newline character.

### Downloads

Type a below command in your terminal.

```
go get github.com/carushi/MADscan
```

### Requires

* go language
	* Reference:[golang.org](https://golang.org)

* go-gtk (including GTK-Development-Packages)
	* Reference: [go-gtk repository](https://github.com/mattn/go-gtk)
	* I refered "demo" of go-gtk for implementation.

### GOPATH
Set \$GOPATH appropriately. I assume that \$GOPATH directory contains

* GOPATH/
	* bin/
	* pkg/
	* src/

If you set GOPATH to \$HOME/work, type the command as below.

```
mkdir $HOME/work
export GOPATH=$HOME/work
```
Reference: [golang.org](https://golang.org/doc/code.html).

### Example

1. Start MADscan
	* Type "MADscan" in your terminal
	* Or click "MADscan" in \$GOPATH/bin directory.
2. Choose input (database) file
	* default: MADscan/data/Sample\_modification\_site
	* tab or character delimited file (specified in ``Select delimiter'' button)
3. Choose output file 
	* default: MADscan/data/output.txt
4. Choose keyword file
	* default: MADscan/data/Ras\_gene\_list.txt
	* It should contain each keyword in a single line such as gene id.
5. Set column position
	* 0: search all columns, 1: search 1st column, ...
6. Check "Ignore lower/upper case" for more flexible search 
7. Push Run button!
	
### ToDo
* filename editing for output file
* header skipping or including



