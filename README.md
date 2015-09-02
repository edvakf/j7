# j7

Filter stdin to stdout with JavaScript

[![Build Status](https://travis-ci.org/edvakf/j7.svg)](https://travis-ci.org/edvakf/j7)

## Usage

```
$ echo '[1,2,3]' | j7 'function main(input) {return JSON.parse(input).join("\n")}'
1
2
3
```

## Command line options

```
  -j  JSON mode (input is JSON.pars-ed and output is JSON.stringify-ed)
  -l  line-by-line mode (input is filtered line by line)
  -m string
      JS entry point function (default "main")
  -n int
      output size in number of bytes, only effective in conjugation with -j
```

You can have multiple JS expressions and files to be evaluated. An argument starting with @ is treated as a file path.

```
j7 'window={}' @/some/library.js 'function main(){return window.foo}'
```

### Examples

```
$ echo '[1,2,3]' | j7 -j 'function main(input) {return input.filter(function(a) {return a%2===1})}'
[1,3]
```

```
$ echo "[1,2,3]\n[4,5,6]" | j7 -l -j 'function main(input){return input[0]}'
1
4
```

## About V7

j7 makes use of the [V7 JavaScript engine](https://github.com/cesanta/v7),
a very small JavaScript engine for embedded devices, and it's go binding [gov7](https://github.com/edvakf/gov7).

The JavaScript may therefore behave differently from what runs in web browsers.
For example, to my knowledge, `JSON.stringify` only outputs the first 100 bytes.
The `-j` option increases the limit heuristically, but you may sometimes want to specify it with the `-n` option by yourself.

## LICENSE

Because V7 is distributed under the GPLv2, j7 has also adopted GPLv2.
