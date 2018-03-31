// txtkbd2kla project  
// Copyright 2018 Philippe Quesnel  
// Licensed under the Academic Free License version 3.0  

Very simple Go program to create KLA (keyboard layout Analyzer) JSON keyboard layout files from a simple text file description of the layout.  

KLA : http://shenafu.com/code/keyboard/klatest/#/main  

The program reads a text qwerty layout as reference, then a similar text layout that gives the new character.  
It then replaces the corresponding character in the template KLA file.  
Once all characters are mapped, the resulting JSON is output  to stdout.

parameters: klaRefQwertyLayoutJson qwertyLayoutFile newLayoutFile [newLayoutFileAltGr]  

klaRefQwertyLayoutJson is the template QWERTY KLA layout  
qwertyLayoutFile is the text qwerty reference  
newLayoutFile is the new layout that we want  
newLayoutFileAltGr is an optional text file for the AltGr layer of the new layout 

Example of mapping the left hand part of a keyboard to a MTGAP layout  

qwertyLayoutFile: reference qwerty (uppercase rows 1st)  

    Q W E R T
    A S D F G
    Z X C V B 
            
    q w e r t
    a s d f g
    z x c v b 

newLayoutFile   

    B L O U ?
    H R E A /
    K X < > Z

    b l o u :
    h r e a ;
    k x , . z

All the charaters from qwertyLayoutFile are mapped to the characters of newLayoutFile, into the template KLA JSON file which is then output to the screen. The resulting JSON text can be loaded into KLA and should now have a MTGAP layout in the left hand.  
