#!/bin/bash
echo "GLOX:"

for file in ./*.lox
do
  echo $file
  for i in 1 2 3 4 5
  do
    start=$(date +%s.%N)
    ../glox $file
    end=$(date +%s.%N)
    echo $end - $start | bc
    echo ""
  done
done

echo "CLOX:"

for file in ./*.lox
do
  echo $file
  for i in 1 2 3 4 5
  do
    start=$(date +%s.%N)
    ../../craftinginterpreters/clox $file
    end=$(date +%s.%N)
    echo $end - $start | bc
    echo ""
  done
done

echo "PLOX:"

for file in ./*.lox
do
  echo $file
  for i in 1 2 3 4 5
  do
    start=$(date +%s.%N)
    plox $file
    end=$(date +%s.%N)
    echo $end - $start | bc
    echo ""
  done
done

echo "GLOX2:"

for file in ./*.lox
do
  echo $file
  for i in 1 2 3 4 5
  do
    start=$(date +%s.%N)
    ../../glox2/glox2 -filePath $file
    end=$(date +%s.%N)
    echo $end - $start | bc
    echo ""
  done
done
