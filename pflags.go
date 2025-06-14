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

type Pflags struct {
	desc        string
	flags       []*flagdef.FlagDef
	result      *parse.Result
	parsedBytes []byte
}

func New(desc string) *Pflags {
	return &Pflags{
		desc: desc,
	}
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
	return strings.Replace(o.desc, "{FLAGS}", flagContent, 1)
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
	flag, err := flagdef.New(shortName, longName, flagType, opts...)
	if err != nil {
		return err
	}
	o.flags = append(o.flags, flag)
	return nil
}

func (o *Pflags) Parse(cmdArgs []string) *errors.Error {
	res, err := parse.Parse(o.flags, cmdArgs)
	if err != nil {
		return err
	}
	resBytes, mErr := json.Marshal(res)
	if mErr != nil {
		return errors.NewError(errors.INVALID_USAGE, mErr.Error())
	}
	o.result = res
	o.parsedBytes = resBytes

	return nil
}

func (o *Pflags) Get(name string) []string {
	return parse.Get(name, o.result)
}

func (o *Pflags) Parsed() []byte {
	return o.parsedBytes
}

func GetFromParsedBytes(name string, bytes []byte) ([]string, *errors.Error) {
	res := &parse.Result{}
	err := json.Unmarshal(bytes, res)
	if err != nil {
		return nil, errors.NewError(errors.INVALID_USAGE, err.Error())
	}
	return parse.Get(name, res), nil
}
