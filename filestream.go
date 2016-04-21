/*
FileStream.go by Allen J. Mills
    4.21.16

    An asycronous file writer.

    This work is under Creative Commons Attribution 4.0 International (CC BY 4.0)
    Information about this license can be found here: https://creativecommons.org/licenses/by/4.0/
*/

// Package filestream provides a simple to use asyncronous file writer.
package filestream

import (
	"errors"
	"os"
	"time"
)

var (
	// ErrStreamClosed is returned after an active stream is explicitly
	// indicated to close through FileStream.Quit
	ErrStreamClosed = errors.New("FS Error: Stream is Closed.")
)

// FileStream,
// This is the basic structure for functionality with this package.
// Once given to an active stream OK will have either nil or an error to indicate if the stream is live.
// Once live, any information sent into Write will be written to the stream's file.
// To explicitly close the stream and file, send any integer information into Quit.
type FileStream struct {
	Write chan string
	Quit  chan int
	OK    chan error
}

// NewFileStream,
// This will return a blank FileStream ready to be given to a stream.
func NewFileStream() FileStream {
	return FileStream{make(chan string), make(chan int), make(chan error)}
}

// StartStream,
// filename: the name of a new file for output.
// fs: A FileStream for information passing.
//
// This function will attempt to create a file of filename. If this fails,
// FileStream.OK will have the error passed from os.OpenFile. If this succeeds,
// FileStream.OK will be given nil to indicate that the stream is ready to use.
// Once this happens, you may add information to FileStream.Write and this function
// write the incoming string information to a file.
//
// Warning: This function will block execution if not placed inside a go routine.
func StartStream(filename string, fs *FileStream) {
	f, fileErr := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if fileErr != nil {
		fs.OK <- fileErr
		return
	}
	defer f.Close()

	fs.OK <- nil // Free to use, unblock now.
	for {
		select {
		case out := <-fs.Write: // incoming information
			f.WriteString(out)
		case <-fs.Quit: // I've been explicitly asked to close.
			fs.OK <- ErrStreamClosed
			return
		default:
			time.Sleep(1 * time.Second) // we'll wait.
		}
	}
	return
}
