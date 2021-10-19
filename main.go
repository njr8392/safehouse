package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"time"
)

type PassWord struct {
	Name     string //alias for the name ie google not google.com
	Url      string
	Created  time.Time
	Modified time.Time
	Password string
}

type List []*PassWord

func NewPassword(alias, url string) *PassWord {
	return &PassWord{
		Name:    alias,
		Url:     url,
		Created: time.Now(),
	}
}

func (l *List) Append(p *Password) {
	l = append(l, p)
}

func (l *List) EncodeAll(w io.Writer) error {
	enc := gob.NewEncoder(w)
	err := enc.Encode(*l)
	if err != nil {
		return err
	}
	return nil
}

func (l *List) DecodeAll(r io.Reader) error {
	dec := gob.NewDecoder(r)
	err := dec.Decode(l)
	if err != nil {
		return err
	}
	return nil
}

//method allows us to use gob encoding for unexported fields
func (p *PassWord) gobEncode(f *os.File) ([]byte, error) {
	enc := gob.NewEncoder(f)
	err := enc.Encode(*p)
	if err != nil {
		return nil, err
	}
	return nil, err

}

func DB() (*os.File, error) {
	f, err := os.OpenFile("store", os.O_RDWR|os.O_CREATE, 0744)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func SetOrigin(f *os.file) error {
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	return nil
}

//can not decode two different list into the same list. will have to read and write the whole thing
func main() {
	// THESE TESTS WORK
	var list List

	a := NewPassword("google", "google.com")
	b := NewPassword("espn", "espn.com")
	list = append(list, a, b)
	//	b := NewPassword("twtter", "twitter.com")
	file, err := DB()
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	list.EncodeAll(file)
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
	}

	var test List
	test.DecodeAll(file)
	for _, x := range test {
		fmt.Println(*x)
	}
}
