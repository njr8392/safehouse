package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	sh "github.com/njr8392/safehouse/internal"
	"os"
	"time"
)

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

//can not decode two different list into the same list. will have to read and write the whole thing
func main() {
	flag.Parse()
	run()
}

func run() {
	//I would like to have a the option to generate a password.
	//start small add character/length contraints later
	var list sh.List

	file, err := sh.DB("store")
	defer file.Close()
	if err != nil {
		fmt.Printf("error opening file ---- %s\n", err)
		file.Close()
		os.Exit(1)
	}

	//should add something to check if password already exists
	//bug here adding to the will overwrite at the begging of the file.
	//tmp fix ----- read all file into the list then set origin to the beginning and overwrite the file
	if os.Args[1] == "add" {
		addCmd.Parse(os.Args[2:])
		np := sh.NewPassword(*alias, *url, *pword)
		fmt.Println(np)
		f, _ := file.Stat()
		size := f.Size()
		if size != 0 {
			err := list.DecodeAll(file)
			if err != nil {
				fmt.Printf("error decoding file data ---- %s", err)
				return
			}
		}
		list.Append(np)

		err = sh.SetOrigin(file)
		if err != nil {
			fmt.Printf("error setting file origin to 0 ----- %s\n", err)
		}
		err = list.EncodeAll(file)
		if err != nil {
			fmt.Printf("error encoding data ---- %s", err)
		}
		return
	}
	if os.Args[1] == "encrypt" {
		encryptCmd.Parse(os.Args[2:])
		if *enkey == "" {
			fmt.Printf("Invalid key. Must pass key with -k flag")
			return
		}
		key := sha256.Sum256([]byte(*enkey))
		k := sh.CopySha256(key)
		cipher, err := sh.Encrypt(k, file)
		if err != nil {
			fmt.Println(err)
			return
		}

		//could maybe use file.WriteAt instead of calling offset (2 syscalls)
		_, err = file.WriteAt(cipher, 0)
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
		key := sha256.Sum256([]byte(*deckey))
		k := sh.CopySha256(key)
		txt, err := sh.Decrypt(k, file)
		if err != nil {
			fmt.Println(err)
			return
		}

		//could maybe use file.WriteAt instead of calling offset (2 syscalls)
		_, err = file.WriteAt(txt, 0)
		if err != nil {
			fmt.Println(err)
		}
		return

	}
	err = list.DecodeAll(file)
	if err != nil {
		fmt.Printf("File must be decrypted before it can be decoded\n")
		return
	}

	if *flagGet != "" {
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
		sh.SetOrigin(file)
		err := list.EncodeAll(file)
		if err != nil {
			fmt.Printf("error encoding data ---- %s", err)
		}
		fmt.Printf("Password for %s changed to %s\n", pass.Name, pass.Password)
	}
}
