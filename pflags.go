package pflags

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/chahal-p/pflags/errors"
	"github.com/chahal-p/pflags/flagdef"
	"github.com/chahal-p/pflags/parse"
)

const (
	FlagHelpIdentifier              = "{FLAGS}"
	EscapedFlagHelpIdentifier       = "\\{\\FLAGS\\}"
	tmpEscapedFlagHelpIdentifier    = "\\{\\--FLAGS--\\}"
	DoubleEscapedFlagHelpIdentifier = "\\{\\{\\FLAGS\\}\\}"
)

type Option func(*Pflags)

func AllowUnrecognizedFlags() Option {
	return func(o *Pflags) {
		o.allowUnrecognizedFlags = true
	}
}

type Pflags struct {
	desc                   string
	allowUnrecognizedFlags bool
	flags                  []*flagdef.FlagDef
	result                 *parse.Result
	parsedBytes            []byte
}

func New(desc string, opts ...Option) *Pflags {
	obj := &Pflags{
		desc: desc,
	}
	for _, opt := range opts {
		opt(obj)
	}
	return obj
}

func (o *Pflags) UsageHelp() string {
	if !strings.Contains(o.desc, "{FLAGS}") {
		return o.desc
	}

	flagContent := ""

	maxShortName := float64(0)
	maxLongName := float64(0)
	for _, f := range o.flags {
		maxShortName = math.Max(maxShortName, float64(len(f.ShortName())+1))
		maxLongName = math.Max(maxLongName, float64(len(f.LongName())+2))
	}
	flagContentFormat := "%" + strconv.Itoa(int(maxShortName)) + "s%s  %" + strconv.Itoa(int(maxLongName)) + "s    "
	for _, f := range o.flags {
		sn := ""
		ln := ""
		if f.ShortName() != "" {
			sn = "-" + f.ShortName()
		}
		if f.LongName() != "" {
			ln = "--" + f.LongName()
		}
		sep := " "
		if sn != "" && ln != "" {
			sep = ","
		}
		line := fmt.Sprintf(flagContentFormat, sn, sep, ln)
		flagContent += line + indent(f.Description(), len(line)) + "\n"
	}
	indentSize := 0
	for l := range strings.SplitSeq(o.desc, "\n") {
		indx := strings.Index(l, "{FLAGS}")
		if indx != -1 {
			indentSize = indx
			break
		}
	}
	flagContent = indent(flagContent, indentSize)
	usageHelp := strings.Replace(o.desc, "{FLAGS}", flagContent, 1)
	usageHelp = strings.Join(strings.Split(usageHelp, DoubleEscapedFlagHelpIdentifier), tmpEscapedFlagHelpIdentifier)
	usageHelp = strings.Join(strings.Split(usageHelp, EscapedFlagHelpIdentifier), FlagHelpIdentifier)
	usageHelp = strings.Join(strings.Split(usageHelp, tmpEscapedFlagHelpIdentifier), EscapedFlagHelpIdentifier)

	usageHelp = strings.Join(strings.Split(usageHelp, "\\n"), "\n")

	return usageHelp
}

func indent(content string, size int) string {
	if size <= 0 {
		return content
	}
	ind := ""
	for range size {
		ind += " "
	}
	var result []string
	lines := strings.Split(content, "\n")
	for _, l := range lines {
		if strings.Trim(l, " ") != "" {
			result = append(result, ind+l)
		}
	}
	return strings.Trim(strings.Join(result, "\n"), " ")
}

func (o *Pflags) Add(shortName, longName string, flagType flagdef.FlagType, opts ...flagdef.Option) *errors.Error {
	if shortName == "" && longName == "" {
		return errors.NewError(errors.INVALID_USAGE, "Both short and long name can not be empty for a flag.")
	}
	flag, err := flagdef.New(shortName, longName, flagType, opts...)
	if err != nil {
		return err
	}
	o.flags = append(o.flags, flag)
	return nil
}

func (o *Pflags) Parse(cmdArgs []string) *errors.Error {
	o.result = nil
	o.parsedBytes = nil
	res, err := parse.Parse(o.flags, cmdArgs, o.allowUnrecognizedFlags)
	if err != nil {
		return err
	}
	resBytes, mErr := json.Marshal(res)
	if mErr != nil {
		return errors.NewError(errors.INVALID_USAGE, mErr.Error())
	}
	o.result = res
	o.parsedBytes = resBytes
	// var buf bytes.Buffer
	// json.Indent(&buf, resBytes, "", "    ")
	// println(buf.String())
	return nil
}

func (o *Pflags) Get(name string) ([]string, *errors.Error) {
	return parse.Get(name, o.result)
}

func (o *Pflags) Parsed() []byte {
	return o.parsedBytes
}

func (o *Pflags) NonFlagArgs() []string {
	return parse.NonFlagArgs(o.result)
}

func GetFromParsedBytes(name string, bytes []byte) ([]string, *errors.Error) {
	res := &parse.Result{}
	err := json.Unmarshal(bytes, res)
	if err != nil {
		return nil, errors.NewError(errors.INVALID_USAGE, err.Error())
	}
	return parse.Get(name, res)
}
