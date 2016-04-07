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
	* [Please visit a public website for go language](https://golang.org)

* go-gtk (including GTK-Development-Packages)
	* [Please refer go-gtk repository](https://github.com/mattn/go-gtk)

### Example

1. Choose input file such as tab delimited file (specified in ``Select delimiter'' button)
2. Choose output file (default: MADscan/data/output.txt)
3. Choose keyword file (containing each keyword in a single line such as gene id)
4. Set column position (0: search all columns, 1: search 1st column, ...)
5. Check ``Ignore lower/upper case'' for more flexible search
6. Push Run button!
	
### ToDo

* filename editing for output file
* header skipping or including



