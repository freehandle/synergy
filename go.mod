module github.com/lienkolabs/synergy

go 1.19

replace github.com/lienkolabs/breeze => ../breeze

replace github.com/lienkolabs/axeprotocol => ../axeprotocol

require (
	github.com/gomarkdown/markdown v0.0.0-20230922112808-5421fefb8386
	github.com/lienkolabs/axeprotocol v0.0.0-00010101000000-000000000000
	github.com/lienkolabs/breeze v0.0.0-00010101000000-000000000000
)

require (
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/term v0.9.0 // indirect
)
