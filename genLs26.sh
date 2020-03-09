#!/bin/sh

cat qwertyMatrixLasalle26+digits.json \
  | ./txtkbd2kla shiftAltGr  qwertyMatrixLasalle26.txt $1-shiftAltGr.txt \
  | ./txtkbd2kla altGr       qwertyMatrixLasalle26.txt $1-altGr.txt \
  | ./txtkbd2kla shift       qwertyMatrixLasalle26.txt $1-shift.txt \
  | ./txtkbd2kla primary     qwertyMatrixLasalle26.txt $1-main.txt \
  | ./txtkbd2kla 'author:' "phil quesnel" \
  | ./txtkbd2kla 'label:' "$1" 