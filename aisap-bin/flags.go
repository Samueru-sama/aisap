package main

import (
    "fmt"
    "os"

    flag "github.com/spf13/pflag"
)

type arrayFlags []string

var (
    // Normal flags
    help      = flag.BoolP("help",       "h", false, "")
    verbose   = flag.BoolP("verbose",    "v", false, "")
    listPerms = flag.BoolP("list-perms", "l", false, "")

    // Long-only flags
    permFile = flag.String("profile", "", "")
    level    = flag.Int("level",        -1, "")

	// Flags that can be called multiple times
	addFile  arrayFlags
	addDev   arrayFlags
	addSoc   arrayFlags
	addShare arrayFlags
)

// Initialization of global variables and help menu
func init() {
    var present bool

    flag.Var(&addFile,  "file",   "")
    flag.Var(&addDev,   "device", "")
    flag.Var(&addSoc,   "socket", "")
    flag.Var(&addShare, "share",  "")

	// Prefer AppImage-provided variable `ARGV0` if present
    if argv0, present = os.LookupEnv("ARGV0"); !present {
        argv0 = os.Args[0]
    }

    flag.Usage = func() {
		g := "\033[32m" // Green
		y := "\033[33m" // Yellow
		r := "\033[31m" // Red
		z := "\033[0m"  // Reset

        fmt.Printf("Usage: %s%s%s [OPTIONS] [APPIMAGE]\n\n", g, argv0, z)
        fmt.Printf("Easily sandbox AppImages in BubbleWrap\n")
        fmt.Printf("With no PERMFILE, read permissions directly from AppImage\n")
        fmt.Printf("Sandbox level of 0 only changes data directory, not actually sandboxed!\n\n")
        fmt.Printf("%sNormal options:\n", y)
        fmt.Printf("%s  -h, --help    %sDisplay this help menu\n", g, z)
        fmt.Printf("%s  -v, --verbose %sBe more verbose (NEI)\n\n", g, z)
        fmt.Printf("%sLong-only options:\n", y)
        fmt.Printf("%s  --file    %sAdd file to sandbox\n", g, z)
        fmt.Printf("%s  --socket  %sAllow access to additional sockets\n", g, z)
	fmt.Printf("%s  --share   %sAdd share to sandbox (eg: network)\n", g, z)
        fmt.Printf("%s  --device  %sAllow access to additional /dev files\n", g ,z)
        fmt.Printf("%s  --level   %sChange the base security level of the sandbox (min: 0, max: 3)\n", g, z)
        fmt.Printf("%s  --profile %sLook for permissions in this entry instead of the AppImage\n\n", g, z)
        fmt.Printf("%sExamples:%s\n", y, z)
        fmt.Printf("  %s%s --profile%s=./f.desktop -- ./f.app\n", g, argv0, z)
        fmt.Printf("    Sandbox `f.app` using permissions from `f.desktop`\n\n")
        fmt.Printf("  %s%s ./f.app --level%s=2\n", g, argv0, z)
        fmt.Printf("    Tighten `f.app` sandbox to level 2 (default: 1)\n\n")
        fmt.Printf("  %s%s --file%s=./f.txt %s--file%s ./other.bin ./f.app\n", g, argv0, z, g, z)
        fmt.Printf("    Allow sandbox to access files `f.txt` and `other.bin`\n\n")
        fmt.Printf("%sWARNING:%s No sandbox is impossible to escape! This is to *aid* security, not\n", r, z)
        fmt.Printf("guarentee safety when downloading sketchy stuff online. Don't be stupid!\n\n")
        fmt.Printf("Plus, this is ALPHA software! Very little testing has been done;\n")
        fmt.Printf("%sUSE AT YOUR OWN RISK!%s\n", r, z)
        os.Exit(0)
    }

    flag.Parse()

    if *help || len(os.Args) < 2 {
        flag.Usage()
    }
}

func (i *arrayFlags) Set(value string) error {
    *i = append(*i, value)
    return nil
}

func (i *arrayFlags) String() string {
    return ""
}

func (i *arrayFlags) Type() string {
    return ""
}