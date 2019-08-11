FROM golang:1.12

WORKDIR /go/src/GoMADScan
COPY . /go/src/GoMADScan

RUN apt-get update && apt-get install -y \
    libgtk2.0-dev \
    libglib2.0-dev \
    libgtksourceview2.0-dev

RUN go get github.com/mattn/go-gtk/gtk
#RUN go get github.com/mattn/go-gtk/glib

CMD ["go", "run", "./GoMADScan.go"]

#(GoMADScan:14747): Gtk-WARNING **: 07:25:33.137: cannot open display:
