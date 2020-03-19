#!/bin/sh

cat qwertyLaSalle28SpAZ+numpad.json \
  | ./txtkbd2kla shiftAltGr  qwertyLasalle28AZ.txt $1-shiftAltGr.txt \
  | ./txtkbd2kla altGr       qwertyLasalle28AZ.txt $1-altGr.txt \
  | ./txtkbd2kla shift       qwertyLasalle28AZ.txt $1-shift.txt \
  | ./txtkbd2kla primary     qwertyLasalle28AZ.txt $1-main.txt \
  | ./txtkbd2kla 'author:' "phil quesnel" \
  | ./txtkbd2kla 'label:' "$1" 