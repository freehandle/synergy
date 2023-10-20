module github.com/freehandle/synergy

go 1.19

replace github.com/freehandle/breeze => ../breeze

replace github.com/freehandle/axe => ../axe

replace github.com/freehandle/papirus => ../papirus

replace github.com/freehandle/cb => ../cb

require (
	github.com/freehandle/axe v0.0.0-00010101000000-000000000000
	github.com/freehandle/breeze v0.0.0-00010101000000-000000000000
	github.com/freehandle/cb v0.0.0-00010101000000-000000000000
	github.com/gomarkdown/markdown v0.0.0-20230922112808-5421fefb8386
)

require github.com/freehandle/papirus v0.0.0-00010101000000-000000000000 // indirect
