module github.com/freehandle/synergy

go 1.19

replace github.com/freehandle/breeze => ../breeze

replace github.com/freehandle/handles => ../handles

replace github.com/freehandle/papirus => ../papirus

require (
	github.com/freehandle/breeze v0.0.0-20240119145142-68027d2c379a
	github.com/freehandle/handles v0.0.0-00010101000000-000000000000
	github.com/gomarkdown/markdown v0.0.0-20240723152757-afa4a469d4f9
)

require github.com/freehandle/papirus v0.0.0-20240109003453-7c1dc112a42b // indirect
