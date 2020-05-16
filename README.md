# qrc

	$ ./qrc -h
	Usage: ./qrc [options] [text]
	If text is not given as an argument, input is taken from stdin until EOF.
	Options:
	  -e string
			error correction level: L, M, Q, or H (default "L")
	  -o string
			output file (default "-")

Example usage (assuming you have ImageMagick's `display` installed):
	
	$ some-command | qrc | display

qrc uses [rsc/qr](https://github.com/rsc/qr) interally.
