// tree package
package main

import (
	"flag"
	"fmt"
	"github.com/richardpct/pkgsrc"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
)

var destdir = flag.String("destdir", "", "directory installation")
var pkg pkgsrc.Pkg

const (
	name     = "tree"
	vers     = "1.8.0"
	ext      = "tgz"
	url      = "http://mama.indstate.edu/users/ice/tree/src"
	hashType = "sha256"
	hash     = "715d5d4b434321ce74706d0dd067505bb60c5ea83b5f0b3655dae40aa6f9b7c2"
)

func checkArgs() error {
	if *destdir == "" {
		return fmt.Errorf("Argument destdir is missing")
	}
	return nil
}

func configure() {
	fmt.Println("Waiting while configuring ...")
	f := "Makefile"
	re := regexp.MustCompile(`(?m)^(CFLAGS=.*LINUX.*)`)

	content, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}

	newContent := re.ReplaceAllString(string(content), "#"+"$1")
	err = ioutil.WriteFile(f, []byte(newContent), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func build() {
	cc := "cc"
	cflags := "-fomit-frame-pointer"
	ldflags := ""
	obj := "tree.o unix.o html.o xml.o json.o hash.o color.o file.o strverscmp.o"
	mandir := *destdir + "/share/man/man1"

	fmt.Println("Waiting while compiling and installing ...")
	cmd := exec.Command("make",
		"prefix="+*destdir,
		"CC="+cc,
		"CFLAGS="+cflags,
		"LDFLAGS="+ldflags,
		"OBJS="+obj,
		"mandir="+mandir,
		"install",
	)
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func main() {
	flag.Parse()
	if err := checkArgs(); err != nil {
		log.Fatal(err)
	}

	pkg.Init(name, vers, ext, url, hashType, hash)
	pkg.CleanWorkdir()
	if !pkg.CheckSum() {
		pkg.DownloadPkg()
	}
	if !pkg.CheckSum() {
		log.Fatal("Package is corrupted")
	}

	pkg.Unpack()
	wdPkgName := path.Join(pkgsrc.Workdir, pkg.PkgName)
	if err := os.Chdir(wdPkgName); err != nil {
		log.Fatal(err)
	}
	configure()
	build()
}
