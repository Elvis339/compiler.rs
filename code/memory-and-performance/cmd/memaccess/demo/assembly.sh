#!/usr/bin/env bash

go tool compile -S demo.go > go_assembly.s
go build -o demo demo.go

objdump -d demo > assembly.s

grep -A 50 "main\\..*node.*search" assembly.s > ptr.s
grep -A 50 "main\\..*contiguousBST.*search" assembly.s > array.s

diff -u ptr.s array.s > ptr_array_diff.txt

rm -rf *.s *.o demo

echo "Assembly diff saved to ptr_array_diff.txt"