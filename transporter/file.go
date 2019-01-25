package transporter

import (
	"fmt"
	"log"
	"os"
)

type FileTransporter struct {
	Dir      string
	FileName string
	Rotate   bool
}

func NewFileTransporter(dir string, file string, rotate bool) (*FileTransporter, error) {

	t := FileTransporter{FileName: file, Dir: dir, Rotate: rotate}
	//check if file can be opened or created
	f, err := t.getFile()
	if err != nil {
		return &t, err
	}
	defer f.Close()

	return &t, nil
}

func (t FileTransporter) Write(data []byte) (int, error) {

	f, err := t.getFile()
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := f.Write(data)

	return b, err
}

//getFile opens a file in append mode
//If the file does not exist creates a new file and returns the pointer
//If the file can not be created os.Stdout is returned
func (t FileTransporter) getFile() (*os.File, error) {

	var fileName string
	var out *os.File

	if t.Dir == "" {
		fileName = t.FileName
	} else {
		fileName = fmt.Sprintf("%s%s%s", t.Dir, string(os.PathSeparator), t.FileName)

	}
	//check if file exists
	var _, err = os.Stat(fileName)

	// create file if not exists
	if os.IsNotExist(err) {
		out, err = os.Create(fileName)
		if err != nil {
			return out, err
		}
		defer out.Close()
	}

	out, err = os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		return out, err
	}

	return out, err
}
