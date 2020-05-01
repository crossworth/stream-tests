## RTSP -> Web: Example of converting an RTSP stream to a web friendly stream format.

**NOTE:** None of this code is production ready. Its meant to be used as a starting point or
to comparison of the different formats/support.


If you would like to test this code you should update the stream endpoints on the `html` files inside each folder.

To run it you must have GO installed, than you can type `go run dash/main.go` (replace the folder for the format).
Its important that you run from the directory that contains this README file.

By default its use port 80, if you would like to use another port  edit each `main.go` file.

It uses a test RTSP stream that I found on the internet, use `VLC` to test if it still works before if you find any problem.
