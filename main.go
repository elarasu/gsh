package main

import (
	"bytes"
	"fmt"
	"gsh/banner"
	"io"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/chzyer/readline"
	"github.com/go-resty/resty"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
)

// setup commandline params
var (
	debug   = kingpin.Flag("debug", "Enable debug mode.").Bool()
	timeout = kingpin.Flag("timeout", "Timeout for http call.").Default("5s").OverrideDefaultFromEnvar("HTTP_TIMEOUT").Short('t').Duration()
	url     = kingpin.Arg("url", "URL address.").Required().URL()
)

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("mode",
		readline.PcItem("vi"),
		readline.PcItem("emacs"),
	),
	readline.PcItem("bye"),
	readline.PcItem("get"),
	readline.PcItem("help"),
)

// setup log levels
func init() {
	// set this if environment is production
	//	log.SetFormatter(&log.JSONFormatter{})

	// set this for non prod
	log.SetFormatter(&log.TextFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Debug("debug statement")
	log.Info("info text")
	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")

	// command line params
	kingpin.Version("0.0.1")
	kingpin.Parse()
	//	fmt.Printf("Would ping: %s with timeout %s and count %d", *ip, *timeout, *count)

	// configuration parsing
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var yamlExample = []byte(`
Hacker: true
name: steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: denim
age: 35
eyes : brown
beard: true
`)

	viper.ReadConfig(bytes.NewBuffer(yamlExample))
	fmt.Println(viper.AllKeys())
	fmt.Println(viper.Get("name")) // this would be "steve"
	fmt.Println(viper.Get("hobbies"))
	banner.Print("gsh")
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "mode "):
			switch line[5:] {
			case "vi":
				l.SetVimMode(true)
			case "emacs":
				l.SetVimMode(false)
			default:
				println("invalid mode:", line[5:])
			}
		case line == "mode":
			if l.IsVimMode() {
				println("current mode: vim")
			} else {
				println("current mode: emacs")
			}
		case strings.HasPrefix(line, "get"):
			line := strings.TrimSpace(line[3:])
			geturl := (*url).String() + "/" + line
			println("get ", geturl)
			resp, _ := resty.R().Get(geturl)
			// explore response object
			// fmt.Printf("\nError: %v", err)
			// fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
			// fmt.Printf("\nResponse Status: %v", resp.Status())
			// fmt.Printf("\nResponse Time: %v", resp.Time())
			// fmt.Printf("\nResponse Recevied At: %v", resp.ReceivedAt())
			fmt.Println(string(resp.Body())) // or resp.String() or string(resp.Body())

		case line == "help":
			usage(l.Stderr())
		case line == "bye":
			goto exit
		case line == "":
		default:
			log.Println("you said:", strconv.Quote(line))
		}
	}
exit:
}
