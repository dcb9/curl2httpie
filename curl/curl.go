package curl

import (
	"fmt"
	"io/ioutil"
	"strings"
	"encoding/json"
)

var rawCmdlineOptsDir = "./cmdline-opts/"

type Tag string
type Protocol string
type LongName string
type Feature string

type Option struct {
	Short     byte
	Long      LongName
	HasArg    bool
	Arg       string
	Magic     string
	Tags      []Tag
	Protocols []Protocol
	Added     string
	Mutexed   []LongName
	Requires  []Feature
	SeeAlso   []LongName
	Help      string
	Body      string
}

func (o *Option) String() string {
	hasArg := "false"
	if o.HasArg {
		hasArg = "true"
	}

	var parts = make([]string, 3)
	parts[0] = fmt.Sprintf("%s %s %s", o.Long, string(o.Short), hasArg)

	if len(o.Requires) > 0 {
		requires := make([]byte, 0, len(o.Requires)*16)
		for _, v := range o.Requires {
			requires = append(requires, []byte(v)...)
			requires = append(requires, ' ')
		}
		parts[1] = fmt.Sprintf(" Requires: %s", string(requires))
	}

	if len(o.Requires) > 0 {
		mutexed := make([]byte, 0, len(o.Mutexed)*16)
		for _, v := range o.Mutexed {
			mutexed = append(mutexed, []byte(v)...)
		}
		parts[2] = fmt.Sprintf(" Mutexed: %s", string(mutexed))

	}
	return strings.Join(parts, "")
}

// parse a file

func DotD2Option(content string) (o Option) {
	s := strings.SplitN(content, "\n", 12)
	bIdx := 0
	for _, v := range s {
		bIdx += len(v) + 1
		if v == "---" {
			o.Body = string([]byte(content)[bIdx:])
			break
		} else {
			parseLine(v, &o)
		}
	}
	return
}

func parseLine(line string, o *Option) {
	s := strings.SplitN(line, ":", 2)
	if len(s) != 2 {
		return
	}

	key, val := strings.ToLower(strings.TrimSpace(s[0])), strings.TrimSpace(s[1])
	switch key {
	case "short":
		o.Short = []byte(val)[0]
	case "long":
		o.Long = LongName(val)
	case "arg":
		o.Arg = val
		o.HasArg = val != ""
	case "magic":
		o.Magic = val
	case "added":
		o.Added = val
	case "help":
		o.Help = val
	case "body":
		o.Body = val
	case "tags":
		s := strings.Split(val, " ")
		o.Tags = make([]Tag, len(s))
		for i, v := range s {
			o.Tags[i] = Tag(v)
		}
	case "protocols":
		s := strings.Split(val, " ")
		o.Protocols = make([]Protocol, len(s))
		for i, v := range s {
			o.Protocols[i] = Protocol(v)
		}
	case "mutexed":
		s := strings.Split(val, " ")
		o.Mutexed = make([]LongName, len(s))
		for i, v := range s {
			o.Mutexed[i] = LongName(v)
		}
	case "requires":
		s := strings.Split(val, " ")
		o.Requires = make([]Feature, len(s))
		for i, v := range s {
			o.Requires[i] = Feature(v)
		}
	case "see-also":
		s := strings.Split(val, " ")
		o.SeeAlso = make([]LongName, len(s))
		for i, v := range s {
			o.SeeAlso[i] = LongName(v)
		}
	}
}

func ContainsProtocol(sl []Protocol, v Protocol) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func ContainsLongName(sl []LongName, v LongName) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func GenerateHTTPOptions(path string) []*Option {
	files, _ := ioutil.ReadDir(path)
	httpOptions := make([]*Option, 0, len(files))
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".d") {
			file := path + f.Name()
			content, _ := ioutil.ReadFile(file)

			option := DotD2Option(string(content))
			if isHTTPOption(option) {
				httpOptions = append(httpOptions, &option)
			}
		}
	}
	return httpOptions
}

func isHTTPOption(o Option) bool {
	if len(o.Protocols) == 0 {
		return true
	}

	sl := []Protocol{"HTTP", "HTTPS", "TLS", "SSL"}

	for _, p := range sl {
		if ContainsProtocol(o.Protocols, p) {
			return true
		}
	}
	return false
}

func HTTPOptions() (options []*Option) {
	data, err := Asset("data/options.json")
	if err != nil {
		panic(err)
	}

	options = make([]*Option, 0, 256)
	err = json.Unmarshal(data, &options)
	if err != nil {
		panic(err)
	}
	return
}

func URLAndOptions(args []string) (string, []*Option) {
	var availableOptions = HTTPOptions()

	url, options := "", make([]*Option, 0, len(args))

	indices := make([]int, 0, len(args))
	for i := 0; i < len(args); {
		v := args[i]
		if strings.HasPrefix(v, "--") {
			long := strings.TrimLeft(v, "-")
			for _, o := range availableOptions {
				if string(o.Long) == long {
					// matched
					indices = append(indices, i)

					// clone struct
					option := &Option{}
					*option = *o

					if option.HasArg {
						i++
						option.Arg = args[i]
						indices = append(indices, i)
					}
					options = append(options, option)
				}
			}
		} else if strings.HasPrefix(v, "-") {
			bytesV := []byte(v)
			short := bytesV[1]
			for _, o := range availableOptions {
				if o.Short == short {
					// matched
					indices = append(indices, i)

					// clone struct
					option := &Option{}
					*option = *o

					if option.HasArg {
						i++
						option.Arg = args[i]
						indices = append(indices, i)
					}
					options = append(options, option)
				}
			}
		}

		i++
	}

	if url == "" {
		for i := range args {
			if !InIntSlice(indices, i)	 {
				url = args[i]
			}
		}
	}

	return url, options
}

func InIntSlice(s []int, v int) (found bool) {
	for _, val := range s {
		if val == v {
			found = true
			return
		}
	}
	return
}
