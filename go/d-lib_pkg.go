package main

/*
 * d-lev constants & helper functions
*/

type ver_tbl_t struct {
	lib string
	sw string
	date string
}

// librarian & software versions, dates
// current @ [0]
var ver_tbl = []ver_tbl_t {  
	{"12",		"94e67227",	"2024-01-08"}, // 0
	{"11",		"cabfa8fe",	"2023-11-02"}, // 1
	{"10",		"f1c279cc",	"2023-10-02"}, // 2
	{"9",		"6be9394f",	"2023-07-26"}, // 3
	{"8",		"7bbb846b",	"2023-06-20"}, // 4
	{"7",		"73c6c3d7",	"2023-05-24"}, // 5
	{"6",		"27c263bf",	"2023-01-31"}, // 6
	{"5",		"2d58f653",	"2023-01-01"}, // 7
	{"2",		"add46826",	"2022-10-06"}, // 8
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
	SPI_BYTES = 0x4000									// eeprom code size : 16kB code space
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
	EE_SPI_END = EE_SPI_ADDR + SPI_BYTES				// eeprom code end addr
	//
	EE_START = EE_PRE_ADDR								// eeprom start addr
	EE_END = EE_SPI_END									// eeprom end addr
	EE_WR_MS = 6										// eeprom write wait time (ms)
	//
	DLP_LINES = SLOT_BYTES / EE_RW_BYTES				// lines in DLP file
	SPI_LINES = SPI_BYTES / EE_RW_BYTES					// lines in SPI file
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
	return "???"
}

// given sw_ver, return librarian version
func sw_lib_lookup(sw_ver string) (string) {
	for _, entry := range ver_tbl {
		if sw_ver == entry.sw { return entry.lib }
	}
	return "???"
}
