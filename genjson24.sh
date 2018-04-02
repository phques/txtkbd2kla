#!/bin/sh

# call this with the base name of the new layout files
# it will also be used as the label in the json file
# For example, with these files
#  mcb24v1.0-altGr.txt
#  mcb24v1.0-main.txt
#  mcb24v1.0-shift.txt
#  mcb24v1.0-shiftAltGr.txt
# each one describing a layer,
# call the script: genjson24 mcb24v1.0 > mcb24v1.0.json
cat qwertyMatrix24+digits.json \
  | ./txtkbd2kla shiftAltGr  qwertyMatrix24.txt       $1-shiftAltGr.txt \
  | ./txtkbd2kla altGr       qwertyMatrix24.txt       $1-altGr.txt \
  | ./txtkbd2kla shift       qwertyMatrix24.txt       $1-shift.txt \
  | ./txtkbd2kla primary     qwertyMatrix24+digits.txt $1-main.txt \
  | ./txtkbd2kla 'author:' "phil quesnel" \
  | ./txtkbd2kla 'label:' "$1" 
