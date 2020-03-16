// small 'filter' program to set the characters in a KLA layout JSON File
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type Key struct {
	Primary    int `json:"primary"`
	Shift      int `json:"shift"`
	AltGr      int `json:"altGr"`
	ShiftAltGr int `json:"shiftAltGr"`
	Numpad     int `json:"numpad"`
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

//------------------------

func showErrorExit() {
	fmt.Fprintln(os.Stderr, "parameters: layer qwertyTemplateFile newLayoutFile")
	fmt.Fprintln(os.Stderr, "The KLA JSON template is read from stdin, we act as a filter")
	fmt.Fprintln(os.Stderr, "layer: primary, shift, altGr, shiftAltGr")
	fmt.Fprintln(os.Stderr, "Can also be called to set label or author:")
	fmt.Fprintln(os.Stderr, "  label: label")
	fmt.Fprintln(os.Stderr, "  author: author")
	os.Exit(-1)
}

func (kbd KlaKbd) changeChar(kbdDest *KlaKbd, fromb, tob byte, layerToMap string) bool {

	from := int(fromb)
	to := int(tob)
	if to == 0 {
		to = -1
	}
	// fmt.Fprintf(os.Stderr, "%c : %c\n", from, to)

	// go through keys looking for 'from', change corresponding key in kbdDest to 'to'
	for i, key := range kbd.Keys {
		// we compare with primary..
		// so care should be taken to map primary last if calling multiple times!
		if key.Primary == from {
			switch layerToMap {
			case "primary":
				kbdDest.Keys[i].Primary = to
			case "shift":
				kbdDest.Keys[i].Shift = to
			case "altGr":
				kbdDest.Keys[i].AltGr = to
			case "shiftAltGr":
				kbdDest.Keys[i].ShiftAltGr = to
			default:
				fmt.Fprintf(os.Stderr, "not a known layer: %s\n", layerToMap)
				return false
			}
			return true
		}
	}

	// src key not found !
	fmt.Fprintf(os.Stderr, "'from' char not found %c / %d in KLA template\n", fromb, from)
	return false
}

//------------------------

// read special tokens space=? and empty=? prefixes tokens
// they are removed from newLayout
func readSpecialTokens(newLayout []byte) (retLayout []byte, emptySpotMarker, spaceMarker byte) {
	emptySpotMarker = byte(0)
	spaceMarker = byte(0)

	for {
		switch string(newLayout[0:6]) {

		// check for empty=? instruction
		// set the character used to indicate the empty spot marker
		case "empty=":
			emptySpotMarker = newLayout[6]
			newLayout = newLayout[7:]
			fmt.Fprintf(os.Stderr, "  empty = %c\n", emptySpotMarker)

		// check for space=? instruction
		// set the character used to indicate Space character
		case "space=":
			spaceMarker = newLayout[6]
			newLayout = newLayout[7:]
			fmt.Fprintf(os.Stderr, "  space = %c\n", spaceMarker)

		default:
			retLayout = newLayout
			return
		}
	}
}

// read from/to kbd definition files and map onto klaKbdDest
func mapKbd(klaKbdSrc, klaKbdDest *KlaKbd) {
	if len(os.Args) != 4 {
		showErrorExit()
	}

	layerToMap := os.Args[1]

	fmt.Fprintf(os.Stderr, "mapping layer %s\n", layerToMap)

	// read qwertyTemplateFile
	fromFilename := os.Args[2]
	fmt.Fprintf(os.Stderr, "  reading qwertyTemplateFile : %s\n", fromFilename)
	qwertyTemplate, err := ioutil.ReadFile(fromFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	// read newLayoutFile
	toFilename := os.Args[3]
	fmt.Fprintf(os.Stderr, "  reading newLayoutFile      : %s\n", toFilename)
	newLayout, err := ioutil.ReadFile(toFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	// remove all whitespace, keeps only the keys
	re := regexp.MustCompile(`\s+`)

	qwertyTemplate = re.ReplaceAllLiteral(qwertyTemplate, nil)
	newLayout = re.ReplaceAllLiteral(newLayout, nil)

	// read special tokens space=? and empty=?
	emptySpotMarker := byte(0)
	spaceMarker := byte(0)
	newLayout, emptySpotMarker, spaceMarker = readSpecialTokens(newLayout)

	// check that from / to have the same length
	if len(newLayout) != len(qwertyTemplate) {
		// error not same length
		fmt.Fprintln(os.Stderr,
			"the from / to layouts are not the same length",
			len(qwertyTemplate), len(newLayout))

		fmt.Fprintln(os.Stderr, "You can use space=? and/or empty=? at beginning of files")

		os.Exit(-1)
	}

	// go through from / to bytes
	// map the corresponding 'from' key in KLA keyboard to 'to' char (into klaKbdDest)
	for i, fromKey := range qwertyTemplate {
		toChar := newLayout[i]

		if toChar == spaceMarker {
			toChar = ' '
		}
		if toChar == emptySpotMarker {
			// mark spot empty in destination
			toChar = 0
		}
		klaKbdSrc.changeChar(klaKbdDest, fromKey, toChar, layerToMap)
	}
}

// read the KLA json template from stdin
// (this app acts as a filter, multiple calls can be chained with pipes)
func readKlaKbd() (*KlaKbd, *KlaKbd) {
	// fmt.Fprintln(os.Stderr, "reading KLA JSON fom stdin")

	// read json as []byte from stdin
	var jsonTemplateBuff bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Bytes()
		jsonTemplateBuff.Write(text)
		jsonTemplateBuff.WriteString("\n")
	}
	jsonTemplateBytes := jsonTemplateBuff.Bytes()

	// create a KlaKbd from the json template
	var klaKbdSrc KlaKbd
	err := json.Unmarshal(jsonTemplateBytes, &klaKbdSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to unmarshall json kbd: ", err)
		os.Exit(-1)
	}

	// make a 2nd one, this one will receive the changes
	var klaKbdDest KlaKbd
	json.Unmarshal(jsonTemplateBytes, &klaKbdDest)

	return &klaKbdSrc, &klaKbdDest
}

func main() {

	if len(os.Args) != 3 && len(os.Args) != 4 {
		showErrorExit()
	}

	klaKbdSrc, klaKbdDest := readKlaKbd()

	// check command to execute
	command := os.Args[1]

	switch command {
	case "label:":
		if len(os.Args) != 3 {
			showErrorExit()
		}
		klaKbdDest.Label = os.Args[2]
		fmt.Fprintln(os.Stderr, "changing label: ", klaKbdDest.Label)

	case "author:":
		if len(os.Args) != 3 {
			showErrorExit()
		}
		klaKbdDest.Author = os.Args[2]
		fmt.Fprintln(os.Stderr, "changing author: ", klaKbdDest.Author)

	case "primary":
		mapKbd(klaKbdSrc, klaKbdDest)
	case "shift":
		mapKbd(klaKbdSrc, klaKbdDest)
	case "altGr":
		mapKbd(klaKbdSrc, klaKbdDest)
	case "shiftAltGr":
		mapKbd(klaKbdSrc, klaKbdDest)

	default:
		fmt.Fprintln(os.Stderr, "unkown command :", command)
		showErrorExit()
	}

	// output th resulting KLA JSON
	resJson, err := json.MarshalIndent(klaKbdDest, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create json result: ", err)
		os.Exit(-1)
	}
	fmt.Printf("%s\n", resJson)

}
