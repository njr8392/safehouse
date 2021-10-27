package internal

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"time"
)

/*
A password is how will store passwords. Everyone will have an alias and a url as well as metadata when it was created/modifed
Alias will be used for ease of querying ie instead of searching for twitter instead of twitter.com.
The URL will mainly serve as a backup in case you forget.  User will also be able to list all passwords just in case.
Just don't forget the password to decrypt the file!
*/
type PassWord struct {
	Name     string //alias for the name ie google not google.com
	Url      string
	Created  time.Time
	Modified time.Time
	Password string
}

/*
List is how our data will be store in the gob encoded file.
*/
type List []*PassWord

func NewPassword(alias, url, pword string) *PassWord {
	return &PassWord{
		Name:     alias,
		Url:      url,
		Created:  time.Now(),
		Password: pword,
	}
}

func (l *List) Append(p *PassWord) {
	*l = append(*l, p)
}

// will retrive the struct for the desired password by name or by url
func (l *List) Get(name string) *PassWord {
	for _, p := range *l {
		if p.Name == name || p.Url == name {
			return p
		}
	}
	return nil
}

//List a passwords to stdout. usage for the -l flag
func (l *List) ListAll() {
	for _, p := range *l {
		fmt.Printf("%v\n", p)
	}
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
func(p *PassWord)String()string{
	if p != nil{
	return fmt.Sprintf("Alias: %s\nUrl: %s\nPassword: %s\nCreated at: %s\nLast Modified: %s", p.Name, p.Url, p.Password, p.Created, p.Modified)
	}
	return fmt.Sprintf("Alias or Url not found")
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

func DB(name string) (*os.File, error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0744)
	if err != nil {
		return nil, err
	}
	return f, nil
}

//set origin is used to reset the posistion of the file to the begining
func SetOrigin(f *os.File) error {
	_, err := f.Seek(0, 0)
	if err != nil {
		return err
	}
	return nil
}
