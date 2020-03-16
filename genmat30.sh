#!/bin/sh

cat qwertyMatrix30+digits.json \
  | ./txtkbd2kla shiftAltGr  qwertyMain30.txt $1-shiftAltGr.txt \
  | ./txtkbd2kla altGr       qwertyMain30.txt $1-altGr.txt \
  | ./txtkbd2kla shift       qwertyMain30.txt $1-shift.txt \
  | ./txtkbd2kla primary     qwertyMain30.txt $1-main.txt \
  | ./txtkbd2kla 'author:' "phil quesnel" \
  | ./txtkbd2kla 'label:' "$1" 