package main

/*
 * d-lev constants & helper functions
*/

import (
	"strings"
	"strconv"
	"math/bits"
	"fmt"
	"os"
	"regexp"
)

type ver_tbl_t struct {
	lib string
	sw string
	date string
}

// librarian & software versions, dates
// current @ [0]
var ver_tbl = []ver_tbl_t {  
	{"11",		"cabfa8fe",	"2023-11-02"}, // 0
	{"10",		"f1c279cc",	"2023-10-02"}, // 1
	{"9",		"6be9394f",	"2023-07-26"}, // 2
	{"8",		"7bbb846b",	"2023-06-20"}, // 3
	{"7",		"73c6c3d7",	"2023-05-24"}, // 4
	{"6",		"27c263bf",	"2023-01-31"}, // 5
	{"5",		"2d58f653",	"2023-01-01"}, // 6
	{"2",		"add46826",	"2022-10-06"}, // 7
	{"old_129",	"7bc1bd55",	"2022-07-05"},
	{"old_128",	"93152c8b",	"2022-05-10"},
	{"old_127",	"d202d35",	"2022-05-04"},
	{"old_126",	"af3f63c4",	"2022-04-30"},
	{"old_125",	"67517a97",	"2022-04-17"},
	{"old_124",	"5ba55477",	"2022-03-17"},
	{"old_121",	"7b6a0484",	"2022-01-01"},
	{"old_120",	"84f7f31c",	"2021-12-18"},
	{"old_119",	"52fe7d",	"2021-12-04"},
	{"old_115",	"240b1e68",	"2021-10-31"},
}

const (
	SLOTS = 256											// pre + pro slots
	SLOT_BYTES = 256									// bytes per slot
	PRO_SLOTS = 6										// profile[0:5]
	PRE_SLOTS = SLOTS - PRO_SLOTS						// preset[0:249]
	//
	EE_RW_BYTES = 4										// eeprom bytes per read / write cycle
	EE_PG_BYTES = 256									// eeprom bytes per page
	//
	EE_PRE_ADDR = 0x0									// eeprom pre start addr
	EE_PRE_END = EE_PRE_ADDR + (PRE_SLOTS * SLOT_BYTES)	// eeprom pre end addr
	//
	EE_PRO_ADDR = EE_PRE_END							// eeprom pro start addr
	EE_PRO_END = EE_PRO_ADDR + (PRO_SLOTS * SLOT_BYTES)	// eeprom pro end addr
	//
	EE_SPI_ADDR = EE_PRO_END							// eeprom code start addr
	EE_SPI_SZ = 0x4000									// eeprom code size : 16kB code space
	EE_SPI_END = EE_SPI_ADDR + EE_SPI_SZ				// eeprom code end addr
	//
	EE_START = EE_PRE_ADDR								// eeprom start addr
	EE_END = EE_SPI_END									// eeprom end addr
	EE_WR_MS = 6										// eeprom write wait time (ms)
	//
	PAGES_COLS = 4										// pages print columns
	PAGES_ROWS = 5										// pages print rows
	PAGES = PAGES_COLS * PAGES_ROWS						// pages
	//
	KNOBS_COLS = 2										// knob columns
	KNOBS_ROWS = 4										// knob rows
	KNOBS = KNOBS_COLS * KNOBS_ROWS						// knobs
	KNOBS_TOTAL = KNOBS * PAGES							// total knobs
	PAGE_SEL_KNOB = 7									// page selector knob
	//
	RX_BUF_BYTES = 512									// serial port rx buffer size
	CHARS_PER_DOT = 4096								// chars for each activity dot printed
	CFG_FILE = "d-lib.cfg"								// config file name
	WORK_DIR = "_WORK_"									// presets scratch dir
	PRO_DIR = "_PRO_"									// profiles scratch dir (off of work dir)
	PRESETS_DIR = "_ALL_"								// presets dir
	CRC = "debb20e3"									// good CRC
	//
	REG_ERROR = "0x2"									// hive error reg
	REG_PITCH = "0xa"									// hive pitch reg
	REG_VOLUME = "0xb"									// hive volume reg
)

// given sw_ver, return date
func sw_date_lookup(sw_ver string) (string) {
	for _, entry := range ver_tbl {
		if sw_ver == entry.sw { return entry.date }
	}
	return "unknown"
}

// given sw_ver, return librarian version
func sw_lib_lookup(sw_ver string) (string) {
	for _, entry := range ver_tbl {
		if sw_ver == entry.sw { return entry.lib }
	}
	return "unknown"
}

// convert string of multi-byte hex values to slice of ints
// hex string values on separate lines
func hexs_to_ints(hex_str string, bytes int) ([]int) {
	var ints []int
	str_split := (strings.Split(strings.TrimSpace(hex_str), "\n"))
	for _, str := range str_split {
		sh_reg, err := strconv.ParseUint(str, 16, 32); err_chk(err)
		for b:=0; b<bytes; b++ { 
			sh_byte := int(uint8(sh_reg))
			ints = append(ints, sh_byte)
			sh_reg >>= 8
		}
	}
	return ints
}

// convert slice of ints to string of multi-byte hex values
// hex string values on separate lines
func ints_to_hexs(ints []int, bytes int) (string) {
	var hex_str string
	for i:=0; i<len(ints); i+=bytes {
		var line_int int64
		for b:=0; b<bytes; b++ { 
			line_int += int64(uint8(ints[i+b])) << (b * 8)
		}
		hex_str += strconv.FormatInt(line_int, 16) + "\n"
	}
	return hex_str
}

// check for hexness
func str_is_hex(str string) bool {
	if len(str) == 0 { return false }
	for _, ch := range str {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) { return false }
	}
	return true
}

// return index of string in slice, else -1
func str_exists(strs []string, str string) (int) {
	for idx, entry := range strs { if str == entry { return idx } }
	return -1
}

// return crc of 32 bit input
func crc_32(sh_reg uint32) (uint32) {
	poly := uint32(0x6db88320)
	for i:=0; i<32; i++ { 
		sh_reg = bits.RotateLeft32(sh_reg, -1)  // >>r 1
		if sh_reg & 0x80000000 != 0 { sh_reg ^= poly }  // xor w/ poly if MSb is set
	}
	return sh_reg
}

// print quit message and exit program
func quit_exit() {
	fmt.Println("> -QUIT- exiting program...")
	os.Exit(0) 
}

// print error message and exit program
func error_exit(error_str string) {
	fmt.Println("> -ERROR-", error_str, "!")
	os.Exit(0) 
}

// check for error, exit program if true
func err_chk(err error) {
	if err != nil { error_exit(err.Error()) }
}

// remove C style block comments
func strip_c_cmnt(str string) string {
	c_cmnt, err:= regexp.Compile(`/\*[^*]*\*+(?:[^*/][^*]*\*+)*/`); err_chk(err)
	return c_cmnt.ReplaceAllString(str, string(""))
}

// remove C++ style EOL comments
func strip_cpp_cmnt(str string) string {
	cpp_cmnt, err := regexp.Compile(`//.*`); err_chk(err)
	return cpp_cmnt.ReplaceAllString(str, string(""))
}

func strip_cmnt(str string) string {
	return strip_cpp_cmnt(strip_c_cmnt(str))
}

/*
// this is a comment
/\
*  Comment? *\
/

// not a comment
"/\
* Comment? *\
/"
*/




