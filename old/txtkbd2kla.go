// txtkbd2kla project  
// Copyright 2018 Philippe Quesnel  
// Licensed under the Academic Free License version 3.0  

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/phques/txt2autokey/kbdRdr"
)

type Key struct {
	Primary    int `json:"primary"`
	Shift      int `json:"shift"`
	AltGr      int `json:"altGr"`
	ShiftAltGr int `json:"shiftAltGr"`
	Finger     int `json:"finger"`
	Id         int `json:"id"`
}

type KlaKbd struct {
	Label        string         `json:"label"`
	Author       string         `json:"author"`
	MoreInfoUrl  string         `json:"moreInfoUrl"`
	MoreInfoText string         `json:"moreInfoText"`
	FingerStart  map[string]int `json:"fingerStart"`
	KeyboardType string         `json:"keyboardType"`
	Keys         []Key          `json:"keys"`
}

const NoMapChar = 0xFE // this char indicates not to map / empty pos


func (kbd KlaKbd) changeChar(kbdDest *KlaKbd, fromb, tob byte) bool {
	from := int(fromb)
	to := int(tob)
	// go through keys, change found assigned char when found
	for i, key := range kbd.Keys {
		if key.Primary == from {
			kbdDest.Keys[i].Primary = to
			return true
		}
		if key.Shift == from {
			kbdDest.Keys[i].Shift = to
			return true
		}
	}
	// !found !
	fmt.Fprintf(os.Stderr,"char !found %c/%d\n", fromb, from)
	return false
}

func (kbd KlaKbd) changeCharPrimary(kbdDest *KlaKbd, fromb, tob byte) bool {
	from := int(fromb)
	to := int(tob)
    // special char signalling to no tmap this one
    if to == NoMapChar {
        return true
    }
	// go through keys, change found assigned char when found
	for i, key := range kbd.Keys {
		if key.Primary == from {
			kbdDest.Keys[i].Primary = to
			return true
		}
	}
	// !found !
	fmt.Fprintf(os.Stderr,"primary char !found %c/%d\n", fromb, from)
	return false
}

func (kbd KlaKbd) changeCharShift(kbdDest *KlaKbd, fromb, tob byte) bool {
	from := int(fromb)
	to := int(tob)
    // special char signalling to no tmap this one
    if to == NoMapChar {
        return true
    }
	// go through keys, change found assigned char when found
	for i, key := range kbd.Keys {
		if key.Shift == from {
			kbdDest.Keys[i].Shift = to
			return true
		}
	}
	// !found !
	fmt.Fprintf(os.Stderr,"shift char !found %c/%d\n", fromb, from)
	return false
}

func (kbd KlaKbd) changeCharAltGr(kbdDest *KlaKbd, fromb, tob byte) bool {
	from := int(fromb)
	to := int(tob)
    // special char signalling to no tmap this one
    if to == NoMapChar {
        return true
    }
	// go through keys, change found assigned char when found
	for i, key := range kbd.Keys {
		if key.Primary == from {
			kbdDest.Keys[i].AltGr = to
			return true
		}
	}
	// !found !
	fmt.Fprintf(os.Stderr,"altGr char !found %c/%d\n", fromb, from)
	return false
}

func (kbd KlaKbd) changeCharShiftAltGr(kbdDest *KlaKbd, fromb, tob byte) bool {
	from := int(fromb)
	to := int(tob)
    // special char signalling to no tmap this one
    if to == NoMapChar {
        return true
    }
	// go through keys, change found assigned char when found
	for i, key := range kbd.Keys {
		if key.Shift == from {
			kbdDest.Keys[i].ShiftAltGr = to
			return true
		}
	}
	// !found !
	fmt.Fprintf(os.Stderr,"shiftAltGr char !found %c/%d\n", fromb, from)
	return false
}

func (kbd *KlaKbd) mapKbd(kbdDest *KlaKbd, fromKbd, toKbd *kbdRdr.Keyboard) {
	for i := 0; i < len(fromKbd.LowerCase); i++ {
		lowerRowFrom := fromKbd.LowerCase[i]
		upperRowFrom := fromKbd.UpperCase[i]
		lowerRowTo := toKbd.LowerCase[i]
		upperRowTo := toKbd.UpperCase[i]
		for j := 0; j < len(lowerRowFrom); j++ {
			kbd.changeCharPrimary(kbdDest, lowerRowFrom[j], lowerRowTo[j])
			// only map Shift if different than primary
			if upperRowTo[j] != lowerRowTo[j] {
                kbd.changeCharShift(kbdDest, upperRowFrom[j], upperRowTo[j])
            }
		}
	}
}

func (kbd *KlaKbd) mapKbdAltGr(kbdDest *KlaKbd, fromKbd, toKbd *kbdRdr.Keyboard) {

	for i := 0; i < len(fromKbd.LowerCase); i++ {
		lowerRowFrom := fromKbd.LowerCase[i]
		upperRowFrom := fromKbd.UpperCase[i]
		lowerRowTo := toKbd.LowerCase[i]
		upperRowTo := toKbd.UpperCase[i]
		for j := 0; j < len(lowerRowFrom); j++ {
			kbd.changeCharAltGr(kbdDest, lowerRowFrom[j], lowerRowTo[j])
			// only map ShiftAltGr if different than AltGr
			if upperRowTo[j] != lowerRowTo[j] {
				kbd.changeCharShiftAltGr(kbdDest, upperRowFrom[j], upperRowTo[j])
			}
		}
	}
}

func main() {
	if len(os.Args) != 4 && len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr,"parameters: klaRefQwertyLayoutJson qwertyLayoutFile newLayoutFile [newLayoutFileAltGr]")
		os.Exit(-1)
	}

	// get variables from command line
	fromJsonFilename := os.Args[1]
	fromFilename := os.Args[2]
	toFilename := os.Args[3]
	toFilenameAltGr := ""
	//	toFilenameAltGr := nil
	if len(os.Args) == 5 {
		toFilenameAltGr = os.Args[4]
	}

	// read JSON KLA template/ref qwerty
	jsonstr, err := ioutil.ReadFile(fromJsonFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(-1)
	}

	// read 'from' keyboard def (eg. qwery30Main)
	fmt.Fprintf(os.Stderr, "reading %s\n", fromFilename)
	fromKbd := new(kbdRdr.Keyboard)
	err = fromKbd.ReadKeyboardFile(fromFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(-1)
	}

	// read 'to' keyboard def
	fmt.Fprintf(os.Stderr, "reading %s\n", toFilename)
	toKbd := new(kbdRdr.Keyboard)
	err = toKbd.ReadKeyboardFile(toFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(-1)
	}

	// read optional altGr kbd layer
	var toKbdAltGr *kbdRdr.Keyboard
	if toFilenameAltGr != "" {
		fmt.Fprintf(os.Stderr, "reading %s\n", toFilenameAltGr)
		toKbdAltGr = new(kbdRdr.Keyboard)
		err = toKbdAltGr.ReadKeyboardFile(toFilenameAltGr)
		if err != nil {
			fmt.Fprintln(os.Stderr,err)
			os.Exit(-1)
		}
	}

	// check that they have the same layout
	if fromKbd.LayoutString() != toKbd.LayoutString() {
		fmt.Fprintln(os.Stderr,"Error, the from and to keyboards must have the same layout")
		fmt.Fprintf(os.Stderr, "%s <-> %s\n", 
            fromKbd.LowerCase.LayoutString(), 
            toKbd.LowerCase.LayoutString())
		os.Exit(-1)
	}

	if toKbdAltGr != nil {
		if fromKbd.LayoutString() != toKbdAltGr.LayoutString() {
			fmt.Fprintln(os.Stderr,"Error, the from and toAltGr keyboards must have the same layout")
            fmt.Fprintf(os.Stderr, "%s <-> %s\n", 
                fromKbd.LowerCase.LayoutString(), 
                toKbdAltGr.LowerCase.LayoutString())
			os.Exit(-1)
		}
	}

	// create a KlaKbd from the json Qwerty template
	var klaKbd KlaKbd
	//	jsonstr := getKbdJsonString()
	err = json.Unmarshal([]byte(jsonstr), &klaKbd)
	if err != nil {
		fmt.Fprintln(os.Stderr,"failed to unmarshall json kbd: ", err)
		os.Exit(-1)
	}

	// make a 2nd one, this one will receive the changes
	var klaKbdDest KlaKbd
	err = json.Unmarshal([]byte(jsonstr), &klaKbdDest)

	// map the keyboard!
	klaKbd.mapKbd(&klaKbdDest, fromKbd, toKbd)

	// map the AltGr(s) layers of the keyboard!
	if toKbdAltGr != nil {
		// Only map ShiftAltGr if the shift layer is different than the primary
		//		doShiftAltGr := !toKbdAltGr.LowerCase.Equal(toKbdAltGr.UpperCase)
		klaKbd.mapKbdAltGr(&klaKbdDest, fromKbd, toKbdAltGr)
	}

	resJson, err := json.MarshalIndent(klaKbdDest, "", "  ")
	fmt.Printf("%s\n", resJson)
}
