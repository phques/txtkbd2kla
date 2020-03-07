#!/bin/sh

cat qwertyMatrixLasalle20+digits.json \
  | ./txtkbd2kla shiftAltGr  qwertyMatrixLasalle20.txt $1-shiftAltGr.txt \
  | ./txtkbd2kla altGr       qwertyMatrixLasalle20.txt $1-altGr.txt \
  | ./txtkbd2kla shift       qwertyMatrixLasalle20.txt $1-shift.txt \
  | ./txtkbd2kla primary     qwertyMatrixLasalle20.txt $1-main.txt \
  | ./txtkbd2kla 'author:' "phil quesnel" \
  | ./txtkbd2kla 'label:' "$1" 