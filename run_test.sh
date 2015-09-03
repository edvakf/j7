#!/bin/bash -eux

cd $(dirname $0)

ACTUAL="$(echo "[1,2,3]" | j7 'function main(input) {return JSON.parse(input).join(":")}')"
echo $ACTUAL
EXPECTED="1:2:3"
if [ "$ACTUAL" != "$EXPECTED" ]; then
  echo "FAIL: basic usage"
  exit 1
fi

ACTUAL="$(echo "[1,2,3]" | j7 -m foo 'function foo(input) {return JSON.parse(input).join(":")}')"
echo $ACTUAL
EXPECTED="1:2:3"
if [ "$ACTUAL" != "$EXPECTED" ]; then
  echo "FAIL: -m to change entry point name"
  exit 1
fi

ACTUAL="$(echo "[1,2,3]" | j7 'window={}' @test/1.js 'function main(input) {return new window.joiner(":").join(JSON.parse(input))}')"
echo $ACTUAL
EXPECTED="1:2:3"
if [ "$ACTUAL" != "$EXPECTED" ]; then
  echo "FAIL: evaluate file"
  exit 1
fi

ACTUAL="$(cat test/2.in.txt | j7 -j -l 'function main(input){ return input[0] + input[1] + input[2]  }')"
echo $ACTUAL
EXPECTED="$(cat test/2.out.txt)"
if [ "$ACTUAL" != "$EXPECTED" ]; then
  echo "FAIL: line mode"
  exit 1
fi


