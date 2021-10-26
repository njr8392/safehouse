package main

import (
	"encoding/gob"
	"flag"
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

func SetOrigin(f *os.File) error {
	_, err := f.Seek(0, 0)
	if err != nil {
		return err
	}
	return nil
}

var (
	flagGet    = flag.String("g", "", "get desired info")
	flagList   = flag.Bool("l", false, "list all decrypted information")
	addCmd     = flag.NewFlagSet("add", flag.ExitOnError)
	alias      = addCmd.String("alias", "", "Name of site")
	url        = addCmd.String("url", "", "Url of site")
	pword      = addCmd.String("password", "", "Password for the site")
	flagChange = flag.String("c", "", "change password of site")
	encryptCmd = flag.NewFlagSet("encrypt", flag.ExitOnError)
	enkey      = encryptCmd.String("k", "", "key")
	decryptCmd = flag.NewFlagSet("decrypt", flag.ExitOnError)
	deckey     = decryptCmd.String("k", "", "key")
)

func run() {
	//when I get to encryption I would like to have a the option to generate a password.
	//start small add character/length contraints later
	var list List

	file, err := DB()
	defer file.Close()
//	err = list.DecodeAll(file)
	if err != nil {
		fmt.Printf("error opening file ---- %s\n", err)
		file.Close()
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//bug after added to the list and writing to the file.  if the command requires a write to the store
	// I should clear the store to a blank file then rewrite. should still be fast enough.. most people have <500 password
	if os.Args[1] == "add" {
		addCmd.Parse(os.Args[2:])
		np := NewPassword(*alias, *url, *pword)
		fmt.Println(np)
		list.Append(np)

		err := SetOrigin(file)
		if err != nil {
			fmt.Printf("error setting file origin to 0 ----- %s\n", err)
		}
		err = list.EncodeAll(file)
		if err != nil {
			fmt.Printf("error encoding data ---- %s", err)
		}
	}
	if os.Args[1] == "encrypt" {
		encryptCmd.Parse(os.Args[2:])
		if *enkey == "" {
			fmt.Printf("Invalid key. Must pass key with -k flag")
			return
		}
		fmt.Println(*enkey, len(*enkey))
		cipher, err := Encrypt([]byte(*enkey), file)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(cipher)

		//could maybe use file.WriteAt instead of calling offset (2 syscalls)
		_, err = file.WriteAt(cipher,0)
		if err != nil {
			fmt.Println(err)
		}
		return

	}
	if os.Args[1] == "decrypt" {
		decryptCmd.Parse(os.Args[2:])
		//add in check. key must be 32 bytes!
		if *deckey == "" {
			fmt.Printf("Invalid key. Must pass key with -k flag")
			return
		}
		txt, err := Decrypt([]byte(*deckey), file)
		if err != nil {
			fmt.Println(err)
			return
		}

		//could maybe use file.WriteAt instead of calling offset (2 syscalls)
		_, err = file.WriteAt(txt,0)
		if err != nil {
			fmt.Println(err)
		}
		return

	}
	// add list here so you don't have to declare one inside every if

	if *flagGet != "" {
		var list List
		err := list.DecodeAll(file)
		if err != nil{
			fmt.Println(err)
			return
		}
		pass := list.Get(*flagGet)
		fmt.Printf("%v\n", pass)
	}

	if *flagList {
		list.ListAll()
	}

	if *flagChange != "" && *flagGet != "" {
		pass := list.Get(*flagGet)
		if pass == nil {
			fmt.Printf("No password for %s found\n", *flagGet)
			return
		}
		pass.Password, pass.Modified = *flagChange, time.Now()
		SetOrigin(file)
		err := list.EncodeAll(file)
		if err != nil {
			fmt.Printf("error encoding data ---- %s", err)
		}
		fmt.Printf("Password for %s changed to %s\n", pass.Name, pass.Password)
	}
}

//can not decode two different list into the same list. will have to read and write the whole thing
func main() {
	flag.Parse()
	run()
}
