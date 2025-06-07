# Moving objects in the browser

Moving objects in the browser, users can see the movement remotely

see the scheme

install GCC:

sudo apt -y install build-essential

sudo apt install libx11-dev xorg-dev libxtst-dev

sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev

run the program:

go run main.go

press q to start

open 2 browsers on a screen, move red square in one and see the movement in other:

http://127.0.0.1:1234

http://127.0.0.1:1234/client