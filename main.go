package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/edvakf/gov7"
)

var mainFunc string
var jsonMode bool
var outLen int
var lineMode bool

func init() {
	flag.StringVar(&mainFunc, "m", "main", "JS entry point function")
	flag.BoolVar(&jsonMode, "j", false, "JSON mode (input is JSON.pars-ed and output is JSON.stringify-ed)")
	flag.IntVar(&outLen, "n", 0, "output size in number of bytes, only effective in conjugation with -j")
	flag.BoolVar(&lineMode, "l", false, "line-by-line mode (input is filtered line by line)")
}

func main() {

	flag.Parse()

	v7 := gov7.New()
	defer v7.Destroy()

	err := loadJS(v7)
	if err != nil {
		exit(err)
	}

	if lineMode {

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			err = execute(v7, scanner.Text())
			if err != nil {
				exit(err)
			}
		}
		if err := scanner.Err(); err != nil {
			exit(err)
		}

	} else {

		input, err := readStdin(v7)
		if err != nil {
			exit(err)
		}
		err = execute(v7, input)
		if err != nil {
			exit(err)
		}

	}
}

func loadJS(v7 *gov7.V7) error {
	for _, js := range flag.Args() {
		if js[:1] == "@" {
			file := js[1:]
			content, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			_, err = v7.Exec(string(content))
			if err != nil {
				return err
			}
		} else {
			_, err := v7.Exec(js)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func readStdin(v7 *gov7.V7) (string, error) {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	input := string(bytes)
	return input, nil
}

func execute(v7 *gov7.V7, input string) error {
	global := v7.GetGlobalObject()

	if !regexp.MustCompile(`^[a-zA-Z_$][a-zA-Z0-9_$]*$`).MatchString(mainFunc) {
		return errors.New("invalid function name specified")
	}

	fn := v7.Get(global, mainFunc)
	if v7.IsUndefined(fn) {
		return fmt.Errorf("global function `%s` is not defined", mainFunc)
	}

	if jsonMode {
		// if jsonMode is set, use the input string itself as the argument
		// that has the same effect as JSON.parse-ing it

		if !isJSON(input) {
			return errors.New("input stream is not JSON")
		}
		result, err := v7.Exec(mainFunc + "(" + input + ")")
		if err != nil {
			return err
		}

		// huristicly decided output length
		l := max(len(input)*4, 4096)
		// override it if -n option is provided
		if outLen > 0 {
			l = outLen
		}

		output := v7.ToJSON(result, l)
		fmt.Println(output)

	} else {
		// if jsonMode is not set, we set a value $$input to the global object
		// and call the main function passing it
		// note: there is v7.Apply(), but there is no way to check if it's thrown
		// https://github.com/cesanta/v7/issues/501

		attrs := gov7.PROPERTY_READ_ONLY | gov7.PROPERTY_DONT_ENUM | gov7.PROPERTY_DONT_DELETE
		err := v7.Set(global, "$$input", attrs, v7.CreateString(input))
		if err != nil {
			return err
		}

		result, err := v7.Exec(mainFunc + "($$input)+''")
		if err != nil {
			return err
		}

		output, err := v7.ToString(result)
		if err != nil {
			return err
		}
		fmt.Println(output)

	}
	return nil
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

var rx1 = regexp.MustCompile(`^[\],:{}\s]*$`)
var rx2 = regexp.MustCompile(`\\(?:["\\\/bfnrt]|u[0-9a-fA-F]{4})`)
var rx3 = regexp.MustCompile(`"[^"\\\n\r]*"|true|false|null|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?`)
var rx4 = regexp.MustCompile(`(?:^|:|,)(?:\s*\[)+`)

func isJSON(str string) bool {
	// https://github.com/douglascrockford/JSON-js/blob/master/json2.js#L490-L497
	return rx1.MatchString(rx4.ReplaceAllString(rx3.ReplaceAllString(rx2.ReplaceAllString(str, "@"), "]"), ""))
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}
