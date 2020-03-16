#!/bin/sh

cat qwertyMatrixLasalleSp29b+numpad.json \
  | ./txtkbd2kla shiftAltGr  qwertyMatrixLasalle29b.txt $1-shiftAltGr.txt \
  | ./txtkbd2kla altGr       qwertyMatrixLasalle29b.txt $1-altGr.txt \
  | ./txtkbd2kla shift       qwertyMatrixLasalle29b.txt $1-shift.txt \
  | ./txtkbd2kla primary     qwertyMatrixLasalle29b.txt $1-main.txt \
  | ./txtkbd2kla 'author:' "phil quesnel" \
  | ./txtkbd2kla 'label:' "$1" 