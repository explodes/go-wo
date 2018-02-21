go-wo
=====

Go-lang graphical scene library for games.

Built around [pixel](http://github.com/faiface/pixel), this library provides an api for creating an running scenes in a 
game loop.

### Installation and setup


```sh

# Debian
sudo apt install libasound2-dev libgl1-mesa-dev libxcursor-dev libx11-dev libxinerama-dev libxi-dev libxrandr-dev
go get -u github.com/jteeuwen/go-bindata/...
go get github.com/explodes/go-wo
```


### Running examples

```sh
# Set up vendor directory
cd $GOPATH/src/github.com/explodes/go-wo
dep ensure

# Run an example (pick one)
cd $GOPATH/src/github.com/explodes/go-wo/examples/{soccer,aliens,guns,flappy}
make run
```
