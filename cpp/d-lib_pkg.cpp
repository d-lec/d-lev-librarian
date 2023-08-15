// d-lib_pkg.cpp
#pragma once

/////////////
// DEFINES //
/////////////

// !!! VERSION !!!
#define LIBV 0x129			// librarian version

// address stuff
#define BOOT_ADDR 0x3F00			// boot loader mem start addr (256 bytes)
#define EEPROM_BYTES 0x20000		// eeprom bytes total (1Mb/8 = 131,072B)
#define EEPROM_PG_BYTES 0x100		// eeprom bytes per page
#define EEPROM_SW_ADDR 0x10000		// eeprom code start addr
#define EEPROM_SW_BYTES 0x4000		// eeprom 16kB code space
#define EEPROM_PRE_ADDR 0x0			// eeprom preset addr
#define EEPROM_PRE_BYTES 0x10000	// eeprom preset[-120:120] & profile[-7:7] space
#define EEPROM_PRO_ADDR 0x7900		// eeprom profile addr
#define EEPROM_PRO_BYTES 0xF00		// eeprom profile[-7:7] space

// serial port stuff
#define PRE_UPD_MDLY 50		// preset update delay
#define SPI_MDLY 6			// SPI delay
#define RX_CRC_MDLY 50		// CRC delay
#define RX_POLL_MDLY 2		// RX polling delay
#define RX_BURST_MDLY 10	// RX burst delay (for 230.4k baud, 10ms gives 230 chars to buffer)
#define RX_BUF_SZ 1024		// rx buffer max chars
#define RX_MAX 500			// rx max reads before bailing

// preset stuff
#define PRE_PARAMS 256		// preset params
#define KNOB_COLS 2			// knob columns
#define KNOB_ROWS 4			// knob rows
#define PAGE_COLS 4			// editor page columns
#define PAGE_ROWS 5			// editor page rows
#define PAGES PAGE_COLS*PAGE_ROWS			// page count
#define PAGE_KNOBS KNOB_COLS*KNOB_ROWS		// knobs per page
#define KNOBS PAGES*PAGE_KNOBS	// total knobs in editor
#define MENU_PG_IDX 31		// menu page index in sys params (exclude top 5 knobs!)
#define SLOT_MIN_USR -120	// min slot #
#define SLOT_MAX_USR 120	// max slot #
#define SLOT_MIN_SYS 121	// min sys slot #
#define SLOT_MAX_SYS 135	// max sys slot #
#define PRE_NO_MATCH " -?- "	// tomb of the unknown preset


// file extensions
#define DLP_EXT "dlp"
#define BNK_EXT "bnk"
#define SPI_EXT "spi"
#define PRE_EXT "pre"
#define PRO_EXT "pro"
#define EEPROM_EXT "eeprom"



///////////////
// CONSTANTS //
///////////////

// console size
const short CON_W				= 85;
const short CON_H				= 30;


//////////////////////
// HELPER FUNCTIONS //
//////////////////////

// sleep for milliseconds
void ms_sleep(uint32_t ms_i) {
	this_thread::sleep_for(chrono::milliseconds(ms_i));
	return;
}

// returns ceiling(log2(n)) of input
// (why isn't this in std?)
uint32_t clog2(uint32_t a_i) {
	uint32_t clog2a = 0;
	uint32_t i = a_i;
	while(i > 1) {
		i = (i + 1) >> 1;
		clog2a++;
	}
	return(clog2a);
};

// do "correct" modulo of a%b
// which is: (b+(a%b))%b
// so: -4%4=0, -3%4=1, -2%4=2, -1%4=3, 0%4=0, 1%4=1, 2%4=2, 3%4=3, 4%4=0, etc.
int32_t mod_int(int32_t a_i, int32_t b_i) {
	return ((b_i + (a_i % b_i)) % b_i);
}

// do "correct" div of a/b
// which is: floor(a/double(b))
// so: -5/4=-2, -4/4=-1, -3/4=-1, -2/4=-1, -1/4=-1, 0/4=0, 1/4=0, 2/4=0, 3/4=0, 4/4=1, etc.
int32_t div_int(int32_t a_i, int32_t b_i) {
	return (floor(a_i/double(b_i)));
}

// test for valid general whitespace char - tab, newline, vertical tab, form feed, carriage return, space
bool ch_is_white(int32_t ch_i) {
	return ((ch_i == '\t') || (ch_i == '\n') || (ch_i == '\r') || (ch_i == '\v') || (ch_i == '\f') || (ch_i == ' '));
}

// test for valid general non-whitespace char - surprisingly contiguous!
bool ch_is_nonwhite(int32_t ch_i) {
	return ((ch_i >= '!') && (ch_i <= '~'));
}

// test for upper case alpha char
bool ch_is_upper(int32_t ch_i) {
	return ((ch_i >= 'A') && (ch_i <= 'Z'));
}

// test for lower case alpha char
bool ch_is_lower(int32_t ch_i) {
	return ((ch_i >= 'a') && (ch_i <= 'z'));
}

// change alpha char to upper case
int32_t ch_to_upper(int32_t ch_i) {
	return (ch_is_lower(ch_i) ? ch_i - 32 : ch_i);
}

// change alpha char to lower case
int32_t ch_to_lower(int32_t ch_i) {
	return (ch_is_upper(ch_i) ? ch_i + 32 : ch_i);
}

// test for 0-9
bool ch_is_digit(int32_t ch_i) {
	return ((ch_i >= '0') && (ch_i <= '9'));
}

// test for 0-9, A-F, a-f hex char
bool ch_is_hex(int32_t ch_i) {
	return (ch_is_digit(ch_i) || ((ch_to_upper(ch_i) >= 'A') && (ch_to_upper(ch_i) <= 'F')));
}

// convert ch to hex int value
int32_t ch_to_hex(int32_t ch_i) {
	if (ch_is_digit(ch_i)) { return (ch_i - '0'); }
	else if (ch_is_lower(ch_i)) { return (ch_i + 0xA - 'a'); }
	else if (ch_is_upper(ch_i)) { return (ch_i + 0xA - 'A'); }
	else { return (0); }
}

// convert hex int value to ch
char hex_to_ch(int32_t int_i) {
	if (int_i >= 0 && int_i <= 9) { return (int_i + '0'); }
	else if (int_i >= 0xa && int_i <= 0xf) { return (int_i - 0xA + 'A'); }
	else { return (' '); }
}

// trim whitespace at ends of string
void str_trim(string& str_io) {
    while (ch_is_white(str_io.front())) { str_io.erase(str_io.begin()); }
    while (ch_is_white(str_io.back())) { str_io.pop_back(); }
   	return;
}

// kill all '\r' chars in a string
void str_dos_to_unix(string& str_io) {
	str_io.erase(remove(str_io.begin(), str_io.end(), '\r'), str_io.end());
	return;
}

// string to tokens
vector<string> str_to_tokens(const string& str_i) {
	stringstream str_ss(str_i);
	vector<string> vtoken;
	string token;
	while (str_ss >> token) { vtoken.push_back(token); }
	return (vtoken);
}

// tokens to string
string tokens_to_str(const vector<string>& vtoken_i) {
	string str_buf;
	for (uint32_t i = 0; i < vtoken_i.size(); i++) { 
		str_buf += vtoken_i[i] + " ";
	}
	return (str_buf);
}

// test string for signed int
// white space at start will return false
// true start conditions (remaining not examined):
// #, -#, +#
bool str_is_int(const string& str_i) {
	if (str_i.empty()) { return(false); }  // empty = false
	if (isdigit(str_i[0])) { return(true); }  // # start = true
	if (str_i.length() < 2)  { return(false); }  // too short = false
	if (((str_i[0] == '-') || (str_i[0] == '+')) && isdigit(str_i[1])) { return(true); }  // +/-# = true
	else { return(false); }  // else false
}

// test string for hex int
// white space at start will return false
// only first ch is tested (remaining not examined)
bool str_is_hex(const string& str_i) {
	if (str_i.empty()) { return(false); }  // empty = false
	if (ch_is_hex(str_i[0])) { return(true); }  // hex # start = true
	else { return(false); }  // else false
}

// return true if file exists
bool file_exists(const string& fname_i) {
	ifstream f(fname_i.c_str());
	return(f.good());
}

// delete file, return true if error
bool file_delete(const string& fname_i) {
	return(remove(fname_i.c_str()));
}

// rename file, return true if error
bool file_rename(const string& old_fname_i, const string& new_fname_i) {
	return(rename(old_fname_i.c_str(), new_fname_i.c_str()));
}

// copy file, return true if error
bool file_copy(const string& src_fname_i, const string& dst_fname_i) {
	bool error_f = false;
	ifstream src(src_fname_i.c_str(), ios::binary);
	if (!error_f) { error_f = !src.is_open(); }
	ofstream dst(dst_fname_i.c_str(), ios::binary);
	if (!error_f) { error_f = !dst.is_open(); }
	if (!error_f) { dst << src.rdbuf(); }
	if (src.is_open()) { src.close(); }
	if (dst.is_open()) { dst.close(); }
	return(error_f);
}

// string => txt file
bool str_to_file(const string& fname_i, const string& str_i, bool append_i) {
	bool error_f = false;
	ofstream wr_file;
	if (append_i) { wr_file.open(fname_i, ios::app); }
	else { wr_file.open(fname_i); }
	if (!error_f) { error_f = !wr_file.is_open(); }
	if (!error_f) { wr_file << str_i; }
	if (wr_file.is_open()) { wr_file.close(); }  // close file if open
	return(error_f);
}

// string vector => txt file
bool sv_to_file(const string& fname_i, const vector<string>& sv_i, bool append_i) {
	bool error_f = false;
	ofstream wr_file;
	if (append_i) { wr_file.open(fname_i, ios::app); }
	else { wr_file.open(fname_i); }
	if (!error_f) { error_f = !wr_file.is_open(); }
	if (!error_f) { for (uint32_t idx = 0; idx < sv_i.size(); idx++) { wr_file << sv_i[idx] << endl; } }
	if (wr_file.is_open()) { wr_file.close(); }  // close file if open
	return(error_f);
}

// txt file => string
// kill dos shit
bool file_to_str(const string& fname_i, string& str_o) {
	bool error_f = false;
	// open file stream
	ifstream rd_file;
	rd_file.open(fname_i, ios::in | ios::binary);
	if (!error_f) {	error_f = !rd_file.is_open(); }
	// file => ss => str, kill crap
	if (!error_f) {
		stringstream ss;
		ss << rd_file.rdbuf();
		str_o = ss.str();
		str_dos_to_unix(str_o);  // kill all '\r' chars (@ DOS/WIN eol)
	}
	if (rd_file.is_open()) { rd_file.close(); }  // close file if open
	return(error_f);
}

// string => string vector
void str_to_sv(const string& str_i, vector<string>& sv_o) {
	stringstream ss(str_i);  // so we can use getline
	string line = "";
	while (getline(ss, line)) { sv_o.push_back(line); } 
	return;
}

// txt file => string vector
bool file_to_sv(const string& fname_i, vector<string>& sv_o, bool trim) {
	string rd_str = "";
	bool error_f = file_to_str(fname_i, rd_str);
	if (!error_f) { 
		if (trim) { str_trim(rd_str); }
		str_to_sv(rd_str, sv_o); }
	return(error_f);
}

// search sv for substring, return match count and idx of last match
// if len_i=0 do full string
// use sorted_i=1 iff sv_i is alpha sorted
uint32_t sv_str_match(const string str_i, const vector<string>& sv_i, uint32_t len_i, uint32_t& idx_o) {
	uint32_t matches = 0;
	int32_t comp = 0;
	uint32_t idx = 0;
	idx_o = 0;  // init
	while (idx < sv_i.size()) {;
		if (len_i) { comp = str_i.compare(0, len_i, sv_i[idx], 0, len_i); }
		else { comp = str_i.compare(sv_i[idx]); } // 0:match
		if (!comp) { matches++; idx_o = idx; }
		idx++;
	}
	return(matches);
}

// hex string => unsigned byte array
// width_i = bytes per line (1 to 4)
// total_i = total bytes to read
// return true if error
bool hex_str_to_uba(const string& pre_str_i, uint8_t uba_o[], uint32_t width_i, uint32_t total_i) {
	bool error_f = false;
	uint32_t idx = 0;
	stringstream ss(pre_str_i);
	string line = "";
	while (!error_f && getline(ss, line)) {
		uint64_t line_hex = 0;  // default
		error_f = !str_is_hex(line);
		if (!error_f) { 
			line_hex = strtoul(line.c_str(), nullptr, 16);  // STRTOUL!
			for (uint32_t b=0; b<width_i; b++) {
				uint64_t line_byte = line_hex & 0xff;
				uba_o[idx] = uint8_t(line_byte);
				line_hex >>= 8;
				idx++;
			}
		}
		if (idx > total_i) { error_f = true; }  // too many bytes
	}
	if (idx != total_i) { error_f = true; }  // wrong total
	return(error_f);
}

// unsigned byte array => string
// return true if error
void uba_to_str(uint8_t uba_i[], string& str_o, uint32_t size_i) {
	uint32_t idx = 0;
	stringstream ss_o;
	while (idx < size_i) {
		uint32_t line_hex = 0;
		for (uint32_t b=0; b<4; b++) {
			uint8_t enc = uba_i[idx];
			line_hex += enc << (b * 8);
			idx++;
		}
		ss_o << hex << line_hex << endl;
	}
	str_o = ss_o.str();
	return;
}
