// d-lib_preset.cpp
// Preset and support functions
#pragma once

struct param_tlp {
	uint32_t type;  // type
	string lbl;  // label
	string ptr;  // pointer name
};

const param_tlp SYS_TLP[] = {
	{0x01, "50Hz",    "s_p0_ds"},	// 0
	{0x07, "Dith",    "s_p1_ds"},	// 1
	{0x03, "P<>V",    "s_p2_ds"},	// 2
	{0x01, "Erev",    "s_p3_ds"},	// 3
	{0x07, "Dith",    "s_p4_ds"},	// 4
	{0x24, "LCD ",    "s_p5_ds"},	// 5
	{0xcb, "Pcal",    "p_p0_ds"},	// 6
	{0xc0, "Lin ",    "p_p1_ds"},	// 7
	{0x26, "Ofs-",    "p_p2_ds"},	// 8
	{0xcb, "Sens",    "p_p3_ds"},	// 9
	{0x26, "Ofs+",    "p_p4_ds"},	// 10
	{0xfd, "Cent",    "p_p5_ds"},	// 11
	{0xcb, "Vcal",    "v_p0_ds"},	// 12
	{0xc0, "Lin ",    "v_p1_ds"},	// 13
	{0x26, "Ofs-",    "v_p2_ds"},	// 14
	{0xcb, "Sens",    "v_p3_ds"},	// 15
	{0x24, "Drop",    "v_p4_ds"},	// 16
	{0x26, "Ofs+",    "v_p5_ds"},	// 17
	{0x35, "Mon ",    "v_p6_ds"},	// 18
	{0x01, "Out ",    "v_p7_ds"},	// 19
	{0x20, "LED ",    "t_p0_ds"},	// 20
	{0x04, "Qant",    "t_p1_ds"},	// 21
	{0x03, "Post",    "t_p2_ds"},	// 22
	{0xab, "Note",    "t_p3_ds"},	// 23
	{0xaf, "Oct ",    "t_p4_ds"},	// 24
	{0xc5, "Bass",   "eq_p0_ds"},	// 25
	{0xc5, "Treb",   "eq_p1_ds"},	// 26
	{0x35, "Line",    "v_p8_ds"},	// 27
	{0x7e, "Wait",    "s_p6_ds"},	// 28
	{0x24, "Lift",    "p_p6_ds"},	// 29
	{0x7e, "Auto",    "p_p7_ds"},	// 30
	{0x7f, "    ", "menu_pg_ds"},	// 31 - NOT stored in *.dlp !  MENU_PG_IDX!
	{0xfe, "load",   "ps_p0_ds"},	// 32 - NOT stored in *.dlp !
	{0x7e, "stor",   "ps_p1_ds"},	// 33 - NOT stored in *.dlp !
	{0xa7, "Load",   "ps_p2_ds"},	// 34 - NOT stored in *.dlp !
	{0x07, "Stor",   "ps_p3_ds"}	// 35 - NOT stored in *.dlp !
};

const param_tlp USR_TLP[] = {
	{0x35, "osc ",  "o_p0_ds"},	// 0
	{0x24, "odd ",  "o_p1_ds"},	// 1
	{0x24, "harm",  "o_p2_ds"},	// 2
	{0xca, "pmod",  "o_p3_ds"},	// 3
	{0xca, "vmod",  "o_p4_ds"},	// 4
	{0xa7, "oct ",  "o_p5_ds"},	// 5
	{0xf0, "offs",  "o_p6_ds"},	// 6
	{0x24, "xmix",  "o_p7_ds"},	// 7
	{0x24, "fm  ",  "o_p8_ds"},	// 8
	{0x70, "freq",  "o_p9_ds"},	// 9
	{0x76, "reso", "o_p10_ds"},	// 10
	{0xa4, "mode", "o_p11_ds"},	// 11
	{0xca, "pmod", "o_p12_ds"},	// 12
	{0xca, "vmod", "o_p13_ds"},	// 13
	{0xc5, "bass", "o_p14_ds"},	// 14
	{0xc5, "treb", "o_p15_ds"},	// 15
	{0xf0, "hmul", "o_p16_ds"},	// 16
	{0xf0, "hmul", "o_p17_ds"},	// 17
	{0xf0, "offs", "o_p18_ds"},	// 18
	{0x35, "sprd", "o_p19_ds"},	// 19
	{0x24, "xmix", "o_p20_ds"},	// 20
	{0x74, "nois",  "n_p0_ds"},	// 21
	{0x70, "freq",  "n_p3_ds"},	// 22
	{0x76, "reso",  "n_p4_ds"},	// 23
	{0xa4, "mode",  "n_p5_ds"},	// 24
	{0xca, "pmod",  "n_p6_ds"},	// 25
	{0xca, "vmod",  "n_p7_ds"},	// 26
	{0xca, "pmod",  "n_p8_ds"},	// 27
	{0xca, "vmod",  "n_p9_ds"},	// 28
	{0x34, "puls", "n_p10_ds"},	// 29
	{0xc5, "bass", "n_p11_ds"},	// 30
	{0x24, "xmix", "n_p12_ds"},	// 31
	{0xc5, "treb", "n_p13_ds"},	// 32
	{0x24, "duty", "n_p14_ds"},	// 33
	{0xf1, "reso",  "r_p0_ds"},	// 34
	{0xc6, "harm",  "r_p1_ds"},	// 35
	{0x72, "freq",  "r_p2_ds"},	// 36
	{0xc6, "tap ",  "r_p3_ds"},	// 37
	{0x71, "hpf ",  "r_p4_ds"},	// 38
	{0xc5, "xmix",  "r_p5_ds"},	// 39
	{0xa2, "mode",  "r_p6_ds"},	// 40
	{0x70, "freq",  "f_p0_ds"},	// 41
	{0x35, "levl",  "f_p1_ds"},	// 42
	{0x70, "freq",  "f_p2_ds"},	// 43
	{0x35, "levl",  "f_p3_ds"},	// 44
	{0x70, "freq",  "f_p4_ds"},	// 45
	{0x35, "levl",  "f_p5_ds"},	// 46
	{0xca, "pmod",  "f_p6_ds"},	// 47
	{0xca, "vmod",  "f_p7_ds"},	// 48
	{0xca, "pmod",  "f_p8_ds"},	// 49
	{0xca, "vmod",  "f_p9_ds"},	// 50
	{0xca, "pmod", "f_p10_ds"},	// 51
	{0xca, "vmod", "f_p11_ds"},	// 52
	{0x76, "reso", "f_p12_ds"},	// 53
	{0x76, "reso", "f_p13_ds"},	// 54
	{0x70, "freq", "f_p14_ds"},	// 55
	{0x35, "levl", "f_p15_ds"},	// 56
	{0x70, "freq", "f_p16_ds"},	// 57
	{0x35, "levl", "f_p17_ds"},	// 58
	{0x70, "freq", "f_p18_ds"},	// 59
	{0x35, "levl", "f_p19_ds"},	// 60
	{0x76, "reso", "f_p20_ds"},	// 61
	{0x70, "freq", "f_p22_ds"},	// 62
	{0x35, "levl", "f_p23_ds"},	// 63
	{0x70, "freq", "f_p24_ds"},	// 64
	{0x35, "levl", "f_p25_ds"},	// 65
	{0xca, "pmod", "f_p26_ds"},	// 66
	{0xca, "vmod", "f_p27_ds"},	// 67
	{0x76, "reso", "f_p28_ds"},	// 68
	{0x24, "cntr", "pc_p0_ds"},	// 69
	{0x24, "rate", "pc_p1_ds"},	// 70
	{0x44, "span", "pc_p2_ds"},	// 71
	{0x24, "levl", "pc_p3_ds"},	// 72
	{0xc6, "vloc", "pc_p4_ds"},	// 73
	{0x25, "kloc",  "e_p0_ds"},	// 74
	{0x44, "knee",  "e_p1_ds"},	// 75
	{0x76, "fall",  "e_p2_ds"},	// 76
	{0x75, "rise",  "e_p3_ds"},	// 77
	{0x34, "velo",  "e_p4_ds"},	// 78
	{0x73, "damp",  "e_p5_ds"},	// 79
	{0x25, "dloc",  "e_p6_ds"},	// 80
	{0x74, "prev", "pp_p0_ds"},	// 81
	{0xc5, "harm", "pp_p1_ds"},	// 82
	{0xa7, "oct ", "pp_p2_ds"},	// 83
	{0xca, "pmod", "pp_p3_ds"},	// 84
	{0xca, "vmod", "pp_p4_ds"},	// 85
	{0xc5, "treb", "pp_p5_ds"},	// 86
	{0xc5, "bass", "pp_p6_ds"},	// 87
	{0xb0, "chan",  "m_p0_ds"},	// 88
	{0x25, "vloc",  "m_p1_ds"},	// 89
	{0x42, "bend",  "m_p2_ds"},	// 90
	{0xa7, "oct ",  "m_p3_ds"},	// 91
	{0x34, "velo",  "m_p4_ds"},	// 92
	{0xf2, "cc  ",  "m_p5_ds"},	// 93
	{0x45, "cloc",  "m_p6_ds"},	// 94
	{0x8b, "prvw", "pp_p7_ds"},	// 95
	{0xa3, "bank",  "b_p0_ds"}	// 96
};
	
const string PAGE_PTRS[PAGES][PAGE_KNOBS] = {
	{"v_p6_ds",    "v_p7_ds",   "v_p0_ds",   "p_p0_ds", "ps_p1_ds",   "b_p0_ds",  "ps_p0_ds", "menu_pg_ds"},
	{"v_p6_ds",    "v_p8_ds",  "pp_p0_ds",  "eq_p1_ds",  "o_p0_ds",  "eq_p0_ds",   "n_p0_ds", "menu_pg_ds"},
	{"pp_p4_ds",  "pp_p3_ds",  "pp_p0_ds",  "pp_p5_ds", "pp_p1_ds",  "pp_p6_ds",  "pp_p2_ds", "menu_pg_ds"},
	{"m_p1_ds",    "m_p4_ds",   "m_p6_ds",   "m_p5_ds",  "m_p2_ds",   "m_p0_ds",   "m_p3_ds", "menu_pg_ds"},
	{"e_p0_ds",    "e_p3_ds",   "e_p1_ds",   "e_p2_ds",  "e_p4_ds",   "e_p5_ds",   "e_p6_ds", "menu_pg_ds"},
	{"pc_p4_ds",  "pc_p3_ds",  "pc_p1_ds",  "pc_p2_ds", "pc_p0_ds",  "pp_p7_ds",   "t_p2_ds", "menu_pg_ds"},
	{"n_p9_ds",    "n_p8_ds",   "n_p0_ds",  "n_p13_ds", "n_p10_ds",  "n_p11_ds",  "n_p14_ds", "menu_pg_ds"},
	{"n_p7_ds",    "n_p6_ds",   "n_p3_ds",   "n_p0_ds",  "n_p5_ds",  "n_p12_ds",   "n_p4_ds", "menu_pg_ds"},
	{"o_p4_ds",    "o_p3_ds",   "o_p2_ds",  "o_p15_ds",  "o_p1_ds",  "o_p14_ds",   "o_p5_ds", "menu_pg_ds"},
	{"o_p6_ds",   "o_p18_ds",  "o_p16_ds",  "o_p17_ds",  "o_p8_ds",  "o_p19_ds",   "o_p7_ds", "menu_pg_ds"},
	{"o_p13_ds",  "o_p12_ds",   "o_p9_ds",   "o_p0_ds", "o_p11_ds",  "o_p20_ds",  "o_p10_ds", "menu_pg_ds"},
	{"r_p2_ds",    "r_p3_ds",   "r_p4_ds",   "r_p1_ds",  "r_p6_ds",   "r_p5_ds",   "r_p0_ds", "menu_pg_ds"},
	{"f_p7_ds",    "f_p6_ds",   "f_p0_ds",   "f_p1_ds", "f_p14_ds",  "f_p15_ds",  "f_p12_ds", "menu_pg_ds"},
	{"f_p9_ds",    "f_p8_ds",   "f_p2_ds",   "f_p3_ds", "f_p16_ds",  "f_p17_ds",  "f_p13_ds", "menu_pg_ds"},
	{"f_p11_ds",  "f_p10_ds",   "f_p4_ds",   "f_p5_ds", "f_p22_ds",  "f_p23_ds",  "f_p20_ds", "menu_pg_ds"},
	{"f_p27_ds",  "f_p26_ds",  "f_p18_ds",  "f_p19_ds", "f_p24_ds",  "f_p25_ds",  "f_p28_ds", "menu_pg_ds"},
	{"v_p0_ds",    "v_p4_ds",   "v_p1_ds",   "s_p4_ds",  "v_p2_ds",   "v_p5_ds",   "v_p3_ds", "menu_pg_ds"},
	{"p_p0_ds",    "p_p6_ds",   "p_p1_ds",   "s_p1_ds",  "p_p2_ds",   "p_p4_ds",   "p_p3_ds", "menu_pg_ds"},
	{"t_p0_ds",    "p_p5_ds",   "s_p5_ds",   "t_p3_ds",  "t_p1_ds",   "t_p4_ds",   "t_p2_ds", "menu_pg_ds"},
	{"s_p6_ds",    "p_p7_ds",   "s_p2_ds",   "s_p0_ds", "ps_p3_ds",   "s_p3_ds",  "ps_p2_ds", "menu_pg_ds"}
};


const string PAGE_TITLES[PAGES] = {
	"    D-LEV",
	"   LEVELS",
	"  PREVIEW",
	"     MIDI",
	"   VOLUME",
	"    PITCH",
	"    NOISE",
	"FLT_NOISE",
	"    0_OSC",
	"    1_OSC",
	"  FLT_OSC",
	"    RESON",
	"   0_FORM",
	"   1_FORM",
	"   2_FORM",
	"   3_FORM",
	"  V_FIELD",
	"  P_FIELD",
	"  DISPLAY",
	"   SYSTEM"
};


// this is bs!
const int32_t SYS_PARAMS = sizeof(SYS_TLP)/sizeof(*SYS_TLP);
const int32_t USR_PARAMS = sizeof(USR_TLP)/sizeof(*USR_TLP);

// editor preset type
struct ed_pre_t {
	int32_t sys_enc[PRE_PARAMS] = { 0 };
	int32_t usr_enc[PRE_PARAMS] = { 0 };
};

// knob table type
struct knob_tbl_t {
	uint32_t type[KNOBS] = { 0 };  // type
	uint32_t idx[KNOBS] = { 0 };  // position in preset
	bool sys_f[KNOBS] = { 0 };  // profile flag
	bool sgn_f[KNOBS] = { 0 };  // signed flag
};

// editor preset
ed_pre_t ed_pre;

// knob table
knob_tbl_t knob_tbl;

// slot range s/b [-120:120]
bool slot_usr(int32_t slot_i) {
	if ((slot_i < SLOT_MIN_USR) || (slot_i > SLOT_MAX_USR)) { return(false); }
	else { return(true); }
}

// slot range s/b [121:135]
bool slot_sys(int32_t slot_i) {
	if ((slot_i < SLOT_MIN_SYS) || (slot_i > SLOT_MAX_SYS)) { return(false); }
	else { return(true); }
}

// slot range s/b [-120:135]
bool slot_error(int32_t slot_i) {
	if (slot_usr(slot_i) || slot_sys(slot_i)) { return(false); }
	else { return(true); }
}

// knob range s/b [0:159] (20 pages * 8 page_knobs = 160)
bool knob_error(int32_t knob_i) {
	if ((knob_i < 0) || (knob_i >= KNOBS)) { return(true); }
	else { return(false); }
}

// return index and group given ptr string
// return true if no match
bool ptr_to_idx(const string& ptr_i, uint32_t& idx_o, bool& sys_o) {
	bool match_f = false;  // flag
	// search sys params
	sys_o = true;  // default (sys)
	idx_o = 0;  // init
	while (!match_f && (idx_o < SYS_PARAMS)) {
		if (ptr_i == SYS_TLP[idx_o].ptr) { match_f = true; }
		else { idx_o++; }
	}
	if (match_f) { return(!match_f); }
	// else, search user params
	sys_o = false;
	idx_o = 0;
	while (!match_f && (idx_o < USR_PARAMS)) {
		if (ptr_i == USR_TLP[idx_o].ptr) { match_f = true; }
		else { idx_o++; }
	}
	return(!match_f);
}

// return type maximum value
int32_t type_max(uint32_t type_i) {
	int32_t max = 0;  // default
	if (type_i < 0x20) { max = type_i & 0x1f; }  // 0 thru 31
	else if (type_i < 0x70) { max = (1 << ((type_i & 0x3) + 5)) - 1; }  // 31, 63, 127, 255
	else if (type_i < 0x72) { max = 192; }  // 192
	else if (type_i == 0x72) { max = 127; }  // 127
	else if (type_i < 0x78) { max = 63; }  // 63
	else if (type_i == 0x7e) { max = 120; }  // 120
	else if (type_i == 0x7f) { max = 19; }  // 19
	else if (type_i < 0xc0) { max = type_i & 0x1f; }  // 0 thru 31
	else if (type_i < 0xf0) { max = (1 << ((type_i & 0x3) + 4)) - 1; }  // 15, 31, 63, 127
	else if (type_i == 0xf0) { max = 127; }  // 127
	else if (type_i == 0xf1) { max = 63; }  // 63
	else if (type_i == 0xf2) { max = 31; }  // 31
	else if (type_i == 0xfd) { max = 99; }  // 99
	else if (type_i == 0xfe) { max = 120; }  // 120
	else { max = 0; }
	return(max);
}

// return type minimum value
int32_t type_min(uint32_t type_i) {
	int32_t min = 0;  // default
	if (type_i < 0x80) { min = 0; }  // 0
	else if (type_i < 0xa0) { min = -(type_max(type_i) + 1); }  // -1, -2, ..., -16
	else if (type_i == 0xf2) { min = -127; }  // -127
	else { min = -type_max(type_i); }
	return(min);
}

// return type signed
bool type_signed(uint32_t type_i) {
	if (type_min(type_i) < 0) { return(true); }
	else { return(false); }
}

// check index
// return true if idx_i out of range
bool idx_error(uint32_t idx_i, bool sys_i) {
	uint32_t idx_max = (sys_i) ? SYS_PARAMS : USR_PARAMS;
	if (idx_i > idx_max) { return(true); }
	else { return(false); }
}

// return index label
// return empty string if idx_i out of range
string idx_lbl(uint32_t idx_i, bool sys_i) {
	if (idx_error(idx_i, sys_i)) { return(""); }
	if (sys_i) { return(SYS_TLP[idx_i].lbl); }
	else { return(USR_TLP[idx_i].lbl); }
}

// return type value
// return 0 if idx_i out of range
uint32_t idx_type(uint32_t idx_i, bool sys_i) {
	if (idx_error(idx_i, sys_i)) { return(0); }
	if (sys_i) { return(SYS_TLP[idx_i].type); }
	else { return(USR_TLP[idx_i].type); }
}

// zero out editor
void ed_clear(bool sys_i) {
	if (sys_i) { fill(begin(ed_pre.sys_enc), end(ed_pre.sys_enc), 0); }
	else { fill(begin(ed_pre.usr_enc), end(ed_pre.usr_enc), 0); }
	return;
}

// initialize knob idx array
// return true if error
bool knob_tbl_init() {
	bool error_f = false;
	// loop over PAGE_PTRS to get idx & sys_f
	uint32_t page = 0;
	while (!error_f && (page < PAGES)) {
		uint32_t page_knob = 0;
		while (!error_f && (page_knob < PAGE_KNOBS)) {
			uint32_t knob = (page * PAGE_KNOBS) + page_knob;
			uint32_t idx = 0;
			bool sys_f = false;
			error_f = ptr_to_idx(PAGE_PTRS[page][page_knob], idx, sys_f);
			if (!error_f) {
				uint32_t type = idx_type(idx, sys_f);
				bool sgn_f = type_signed(type);
				// write to table
				knob_tbl.type[knob] = type;
				knob_tbl.idx[knob] = idx;
				knob_tbl.sys_f[knob] = sys_f;
				knob_tbl.sgn_f[knob] = sgn_f;
			}
			page_knob++;
		}
		page++;
	}
	return(error_f);
}

// return editor encoder value
// return true if error
bool ed_idx_rd(int32_t& enc_o, uint32_t idx_i, bool sys_i) {
	bool error_f = idx_error(idx_i, sys_i);
	enc_o = 0;  // init
	if (!error_f) { 
		enc_o = (sys_i) ? ed_pre.sys_enc[idx_i] : ed_pre.usr_enc[idx_i]; 
	}
	return(error_f);
}

// constrain encoder to type limits
// return true if error
bool enc_lim(int32_t& enc_io, uint32_t idx_i, bool sys_i) {
	bool error_f = idx_error(idx_i, sys_i);
	if (!error_f) {
		uint32_t type = idx_type(idx_i, sys_i);
		int32_t min = type_min(type);
		int32_t max = type_max(type);
		if (enc_io > max) { enc_io = max; }
		if (enc_io < min) { enc_io = min; }
	}
	return(error_f);
}

// write sys/usr encoder value to editor preset
// constrain to type min/max limits
// return true if error
bool ed_idx_wr(int32_t enc_i, uint32_t idx_i, bool sys_i) {
	bool error_f = idx_error(idx_i, sys_i);
	int32_t enc = enc_i;
	if (!error_f) { error_f = enc_lim(enc, idx_i, sys_i); }
	if (!error_f) { 
		if (sys_i) { ed_pre.sys_enc[idx_i] = enc; }
		else { ed_pre.usr_enc[idx_i] = enc; }
	}
	return(error_f);
}

// return filter freq value (type 0xf0, 0xf1)
// 7041 * EXP2((ENC * (2^27) / 24) + 3/4)
// input: [0:192]
// output: 27 to 7040 (Hz)
int32_t enc_to_filt_freq(int32_t enc_i) {
	uint32_t mult = pow(2, 27) / 24;
	uint32_t offs = 0xc0000000;
	uint32_t enc_mo = (enc_i * mult) + offs;
	return(7041 * (pow(2, (enc_mo / pow(2, 27))) / pow(2, 32)));
}

// return nearest encoder value for a given filter frequency
uint8_t filt_freq_to_enc(uint64_t freq_i) {
	uint32_t error_best = 1<<31;  // some big value
	uint32_t enc_o = 0;
	for(int i=0; i<=192; i++) {
		uint32_t freq = enc_to_filt_freq(i);  // test freq
		uint32_t error = labs(int64_t(freq_i) - int64_t(freq));
		if (error < error_best) {
			error_best = error;
			enc_o = i;
		}
	}
	return(enc_o);
}

// return reson freq value (type 0xf2)
// 48001 / ((((((~(ENC<<25))^4)*0.871)+((~(ENC<<25))>>3))>>22)+4)
// input: [0:127]
// output: 46 to 9600 (Hz)
int32_t enc_to_reson_freq(int32_t enc_i) {
	uint32_t enc_fs_rev = ~(enc_i << 25);
	uint64_t efr = enc_fs_rev;
	uint64_t sq = (efr * efr) >> 32;
	uint64_t qd = (sq * sq) >> 32;
	return(48001 / ((((uint32_t(qd * 0.871) + (efr / 8)) >> 22) + 4)));
}

// return nearest encoder value for a given reson frequency
uint8_t reson_freq_to_enc(uint64_t freq_i) {
	uint32_t error_best = 1<<31;  // some big value
	uint32_t enc_o = 0;
	for(int i=0; i<=127; i++) {
		uint32_t freq = enc_to_reson_freq(i);  // test freq
		uint32_t error = labs(int64_t(freq_i) - int64_t(freq));
		if (error < error_best) {
			error_best = error;
			enc_o = i;
		}
	}
	return(enc_o);
}

// return idx display string
string idx_disp(uint32_t idx_i, bool sys_i) {
	stringstream str_o;
	int32_t enc = 0;
	bool error_f = ed_idx_rd(enc, idx_i, sys_i);
	if (!error_f) {
		int32_t rtn = 0;
		int32_t type = idx_type(idx_i, sys_i);
		if ((type == 0x70) || (type == 0x71)) { rtn = enc_to_filt_freq(enc); }
		else if (type == 0x72) { rtn = enc_to_reson_freq(enc); }
		else { rtn = enc; }
		if (type == 0) { str_o << setw(5) << "";}
		else { str_o << setw(5) << rtn; }
	}
	return(str_o.str());
}

// return true if x_i out of range
bool x_error(int32_t x_i) { return((x_i < 0) || (x_i >= (PAGE_COLS * KNOB_COLS))); }

// return true if y_i out of range
bool y_error(int32_t y_i) { return((y_i < 0) || (y_i >= (PAGE_ROWS * KNOB_ROWS))); }

// return true if x_i or y_i out of range
bool xy_error(int32_t x_i, int32_t y_i) { 
	bool error_f = x_error(x_i);
	error_f |= x_error(x_i);
	return(error_f); 
}

// return page given UI(x_i,y_i)
int32_t xy_to_page(int32_t x_i, int32_t y_i) {
	return((x_i / KNOB_COLS) + ((y_i / KNOB_ROWS) * KNOB_ROWS));
}

// return page_knob given UI(x_i,y_i)
int32_t xy_to_page_knob(int32_t x_i, int32_t y_i) {
	return((x_i % KNOB_COLS) + ((y_i % KNOB_ROWS) * KNOB_COLS));
}

// return knob given UI(x_i,y_i)
uint32_t xy_to_knob(int32_t x_i, int32_t y_i) {
	int32_t page = xy_to_page(x_i, y_i);
	int32_t page_knob = xy_to_page_knob(x_i, y_i);
	uint32_t knob = (page * PAGE_KNOBS) + page_knob;
	return(knob);
}

// return true if UI(x_i,y_i) at LCD screen title
bool xy_on_title(int32_t x_i, int32_t y_i) {
	if ((mod_int(x_i, KNOB_COLS) == KNOB_COLS - 1) && (mod_int(y_i, KNOB_ROWS) == KNOB_ROWS - 1)) { return(true); }
	else { return(false); }
}

// move x/y position
// stop at individual x & y limits
void xy_move(int32_t& x_io, int32_t& y_io, int32_t dx_i, int32_t dy_i) {
	int32_t x_max = (PAGE_COLS * KNOB_COLS) - 1;
	int32_t y_max = (PAGE_ROWS * KNOB_ROWS) - 1;
	// offset
	x_io += dx_i;
	y_io += dy_i;
	// limit
	if (x_io < 0) { x_io = 0; }
	if (x_io > x_max) { x_io = x_max; }
	if (y_io < 0) { y_io = 0; }
	if (y_io > y_max) { y_io = y_max; }
	return;
}

// move x/y position to home
void xy_home(int32_t& x_io, int32_t& y_io) {
	x_io = KNOB_COLS - 1;
	y_io = KNOB_ROWS - 1;
	return;
}

// move x/y position to end
void xy_end(int32_t& x_io, int32_t& y_io) {
	x_io = (PAGE_COLS * KNOB_COLS) - 1;
	y_io = (PAGE_ROWS * KNOB_ROWS) - 1;
	return;
}

// draw editor UI screen
void ed_draw(int32_t x_i, int32_t y_i, bool ed_i) {
	for (int32_t y=0; y<PAGE_ROWS*KNOB_ROWS; y++) {
		if (!(y % KNOB_ROWS)) { con_font('h'); cout << setw(CON_W) << ""; con_font('n'); cout << endl; }  // row divider
		for (int32_t x=0; x<PAGE_COLS*KNOB_COLS; x++) {
			int32_t page = xy_to_page(x, y);
			int32_t page_knob = xy_to_page_knob(x, y);
			if (!(x % KNOB_COLS)) { con_font('h'); cout << " "; con_font('n'); }  // col divider
			// draw screen title
			if (page_knob == PAGE_KNOBS - 1) {  
				if ((x == x_i) && (y == y_i) && ed_i) { con_font('s'); }  // hilite
				else { con_font('o'); }
				cout << PAGE_TITLES[page]; 
				con_font('n'); 
			}
			// draw screen parameter
			else {  
				uint32_t knob = xy_to_knob(x, y);
				uint32_t idx = knob_tbl.idx[knob];
				bool sys_f = knob_tbl.sys_f[knob];
				if ((x == x_i) && (y == y_i) && ed_i) { con_font('s'); }  // hilite
				cout << idx_lbl(idx, sys_f);
				cout << idx_disp(idx, sys_f);
				con_font('n');
				if (!(x % KNOB_COLS)) { cout << "  "; }
			}
		}
		con_font('h'); cout << " " << endl; con_font('n');  // ending col divider
	}
	con_font('o'); cout << setw(50) << right << "F1 - EDITOR" << setw(35) << ""; 
	con_font('n'); 
	cout << endl;  // ending row divider
	return;
}

// print editor UI screen => txt string
string ed_print_to_str() {
	string h_line_sub = '+' + string(22, '-');
	string h_line = " " + h_line_sub + h_line_sub + h_line_sub + h_line_sub + "+ \n";
	stringstream ss;
	for (int32_t y=0; y<PAGE_ROWS*KNOB_ROWS; y++) {
		if (!(y % KNOB_ROWS)) { ss << h_line; }  // row divider
		for (int32_t x=0; x<PAGE_COLS*KNOB_COLS; x++) {
			int32_t page = xy_to_page(x, y);
			int32_t page_knob = xy_to_page_knob(x, y);
			if (!(x % KNOB_COLS)) { ss << " | "; }  // col divider
			if (page_knob == PAGE_KNOBS - 1) { ss << PAGE_TITLES[page]; }
			else { 
				uint32_t knob = xy_to_knob(x, y);
				uint32_t idx = knob_tbl.idx[knob];
				bool sys_f = knob_tbl.sys_f[knob];
				ss << idx_lbl(idx, sys_f);
				ss << idx_disp(idx, sys_f);
				con_font('n');
				if (!(x % KNOB_COLS)) { ss << "  "; }
			}
		}
		ss << " | " << endl;  // ending col divider
	}
	ss << h_line;  // ending row divider
	return(ss.str());
}

// print editor UI screen => txt file
// return true if error
bool ed_print_to_file(const string& tfname_i, const string& psname_i, bool append_i) {
	string str = ed_print_to_str();
	bool error_f = str_to_file(tfname_i, str, append_i);
	if (!error_f) {
		stringstream ss;
		ss << "preset: " << psname_i << endl << endl;
		error_f = str_to_file(tfname_i, ss.str(), true);
	}
	return(error_f);
}

// preset unsigned byte array => editor
// don't copy page menu, load/stor, Load/Stor encoders
void uba_to_ed(uint8_t pre_uba_i[], bool sys_i) {
	uint32_t idx = 0;
	while (idx < PRE_PARAMS) {
		int32_t enc = 0;
		if (type_signed(idx_type(idx, sys_i))) { enc = int8_t(pre_uba_i[idx]); }
		else { enc = uint8_t(pre_uba_i[idx]); }
		if (sys_i) { if (idx < MENU_PG_IDX) { ed_pre.sys_enc[idx] = enc; } }
		else { ed_pre.usr_enc[idx] = enc; }
		idx++;
	}
	return;
}

// preset file => editor
// return true if error
bool ed_dlp_rd(const string& dlp_i, bool sys_i) {
	string pre_str = "";
	uint8_t pre_uba[PRE_PARAMS] = { 0 };
	bool error_f = file_to_str(dlp_i, pre_str);
	if (!error_f) { error_f = hex_str_to_uba(pre_str, pre_uba, 4, PRE_PARAMS); }
	if (!error_f) { uba_to_ed(pre_uba, sys_i); }
	return(error_f);
}

// editor enc to preset usigned byte array
// don't copy page menu, load/stor, Load/Stor encoders
void ed_to_uba(uint8_t pre_uba_o[], bool sys_i) {
	uint32_t idx = 0;
	while (idx < PRE_PARAMS) {
		if (sys_i) { if (idx < MENU_PG_IDX) { pre_uba_o[idx] = ed_pre.sys_enc[idx]; } }
		else { pre_uba_o[idx] = ed_pre.usr_enc[idx]; }
		idx++;
	}
	return;
}

// editor => preset file
// return true if error
bool ed_dlp_wr(const string& dlp_i, bool sys_i) {
	uint8_t uba[PRE_PARAMS] = { 0 };
	string str = "";
	ed_to_uba(uba, sys_i);
	uba_to_str(uba, str, PRE_PARAMS);
	return(str_to_file(dlp_i, str, false));
}

// knob editor read
// return true if error
bool ed_knob_rd(int32_t &enc_o, uint32_t knob_i) {
	bool error_f = knob_error(knob_i);
	if (!error_f) { 
		uint32_t idx = knob_tbl.idx[knob_i];
		bool sys_f = knob_tbl.sys_f[knob_i];
		error_f = ed_idx_rd(enc_o, idx, sys_f);
	}
	return(error_f);
}

// editor xy enc read
// return true if error
bool xy_rd(int32_t &enc_o, int32_t x_i, int32_t y_i) {
	bool error_f = xy_error(x_i, y_i);
	if (!error_f) { 
		uint32_t knob = xy_to_knob(x_i, y_i);
		uint32_t idx = knob_tbl.idx[knob];
		bool sys_f = knob_tbl.sys_f[knob];
		error_f = ed_idx_rd(enc_o, idx, sys_f);
	}
	return(error_f);
}

// knob value tx to serial port
// return true if error
bool knob_tx(const sp_type& sp_i, uint32_t knob_i, int32_t data_i) {
	bool error_f = sp_open_err(sp_i);  // check port
	error_f |= knob_error(knob_i);  // check knob range
	if (!error_f) {
		stringstream cmd_ss;
		cmd_ss << hex << showbase << knob_i << " " << (data_i & 0xff) << " wk " << flush;
		sp_tx(sp_i, cmd_ss.str());  // tx cmd
		error_f = sp_rx_wait(sp_i);  // rx wait
	}
	return(error_f);
}

// editor xy enc write
// also tx knob (via read value b/c it's limited)
// return true if error
bool xy_wr_tx(const sp_type& sp_i, int32_t enc_i, int32_t x_i, int32_t y_i) {
	bool error_f = xy_error(x_i, y_i);
	if (!error_f) { 
		int32_t enc;
		uint32_t knob = xy_to_knob(x_i, y_i);
		uint32_t idx = knob_tbl.idx[knob];
		bool sys_f = knob_tbl.sys_f[knob];
		if (!error_f) { error_f = ed_idx_wr(enc_i, idx, sys_f); }
		if (!error_f) { error_f = ed_idx_rd(enc, idx, sys_f); }
		if (!error_f) { error_f = knob_tx(sp_i, knob, enc); }
	}
	return(error_f);
}

// editor xy enc delta wr
// read, apply delta & write back, tx
// return true if error
bool xy_delta_wr_tx(const sp_type& sp_i, int32_t delta_i, int32_t x_i, int32_t y_i) {
	bool error_f = xy_error(x_i, y_i);
	int32_t enc;
	if (!error_f) { error_f = xy_rd(enc, x_i, y_i); }
	if (!error_f) { error_f = xy_wr_tx(sp_i, enc + delta_i, x_i, y_i); }
	return(error_f);
}

// editor xy enc input wr
// return true if error
bool xy_input_wr_tx(const sp_type& sp_i, int32_t val_i, int32_t x_i, int32_t y_i) {
	bool error_f = xy_error(x_i, y_i);
	if (!error_f) { 
		uint32_t knob = xy_to_knob(x_i, y_i);
		uint32_t type = knob_tbl.type[knob];
		int32_t enc = val_i;
		if ((type == 0x70) || (type == 0x71)) { enc = filt_freq_to_enc(val_i); } // filt freq
		else if (type == 0x72) { enc = reson_freq_to_enc(val_i); } // reson freq
		error_f = xy_wr_tx(sp_i, enc, x_i, y_i);
	}
	return(error_f);
}

// dir sv => preset sv
// return true if error
bool dir_pre_sv_rd(const vector<string>& dir_sv_i, vector<string>& pre_sv_o) {
	bool error_f = false;
	pre_sv_o.clear();
	uint32_t idx = 0;
	while (!error_f && (idx < dir_sv_i.size())) {
		string pre_str;
		error_f |= file_to_str((dir_sv_i[idx] + "." + DLP_EXT), pre_str);
		pre_sv_o.push_back(pre_str);
		idx++;
	}
	return(error_f);
}

// process string to remove extraneous stuff
// return true if error
bool str_addr_strip(string& str_i, string& str_o) {
	stringstream ss(str_i);
	string line = "";
	int32_t lines = 0;
	str_o = "";  // clear output string
	bool error_f = false;
	bool gt_f = false;
	while (!error_f && !gt_f && getline(ss, line)) {
		bool addr_f = false;
		uint32_t idx = 0;
		while (!error_f && !gt_f && !addr_f && (idx < line.size())) {
			if (line[idx] == '>') { gt_f = true; }
			if (line[idx] == ']') { addr_f = true; }
			if (addr_f) {
				if (idx+1 < line.size()) { line = line.substr(idx+1); }
				else { error_f = true; }
			}
			idx++;
		}
		if (addr_f) {  // only lines w/ addr
			lines++;
			str_o += line;
			str_o += '\n';
		}
	}
	if (!gt_f) { error_f = true; }
	return(error_f);
}


// spi => string
// return true if error
bool spi_str_rx(const sp_type& sp_i, uint32_t addr_start_i, uint32_t addr_end_i, string& str_o) {
	bool error_f = sp_open_err(sp_i);  // check port
	if (!error_f) {
		// flush input
		sp_rx_flush(sp_i);
		// read data burst to string
		string str_rx = "";
		stringstream cmd_ss;
		cmd_ss << hex << showbase << addr_start_i << " " << addr_end_i << " rs " << flush;
		sp_tx(sp_i, cmd_ss.str());  // tx cmd
		sp_rx_burst(sp_i, str_rx);  // rx str
		error_f = str_addr_strip(str_rx, str_o);  // strip cmd crap
	}
	return(error_f);
}

// slot => preset string
// return true if error
bool slot_str_rx(const sp_type& sp_i, int32_t slot_i, string& pre_str_o) {
	bool error_f = slot_error(slot_i);
	if (!error_f) {
		uint32_t slot = slot_i & 0xff;  // to unsigned
		uint32_t addr_start = slot * EEPROM_PG_BYTES;
		uint32_t addr_end = addr_start + EEPROM_PG_BYTES - 1;
		error_f = spi_str_rx(sp_i, addr_start, addr_end, pre_str_o);
	}
	return(error_f);
}

// set spi csn high
// return true if error
bool spi_csn_hi_tx(const sp_type& sp_i) {
	sp_tx(sp_i, "6 0x100 wr ");  // SPI CS
	bool error_f = sp_rx_wait(sp_i);
	return(error_f);
}

// spi write enable
// return true if error
bool spi_wr_en_tx(const sp_type& sp_i) {
	sp_tx(sp_i, "6 6 wr ");  // SPI wr enable
	bool error_f = sp_rx_wait(sp_i);
	error_f |= spi_csn_hi_tx(sp_i);
	return(error_f);
}

// spi write wait
// return true if error
// note: sets CSN hi
bool spi_wr_wait_tx(const sp_type& sp_i) {
	bool error_f = spi_csn_hi_tx(sp_i);  // SPI CS
	ms_sleep(SPI_MDLY);  // SPI wr delay
	return(error_f);
}

// spi write protect
// return true if error
// prot=true: set protected
// prot=false: set unprotected
bool spi_wr_prot_tx(const sp_type& sp_i, bool prot_i) {
	bool error_f = spi_wr_en_tx(sp_i);  // SPI wr enable
	sp_tx(sp_i, "6 1 wr ");  // SPI WRSR reg
	error_f |= sp_rx_wait(sp_i);
	if (prot_i) { sp_tx(sp_i, "6 0xc wr "); }  // SPI prot
	else { sp_tx(sp_i, "6 0 wr "); }  // SPI unprot
	error_f |= sp_rx_wait(sp_i);
	error_f |= spi_wr_wait_tx(sp_i);
	return(error_f);
}

// str => spi
// return true if error
bool spi_str_tx(const sp_type& sp_i, uint32_t spi_addr_i, const string& str_i) {
	bool error_f = sp_open_err(sp_i);  // check port
	if (!error_f) {
		// flush input
		sp_rx_flush(sp_i);
		// process lines in file
		stringstream ss_tx(str_i);  // so we can use getline
		string line = "";
		uint32_t spi_addr = spi_addr_i;
		cout << flush;
		error_f |= spi_wr_prot_tx(sp_i, false);  // unprotect
		while (!error_f && getline(ss_tx, line)) {
			stringstream cmd_ss;
			if (!(spi_addr % EEPROM_PG_BYTES)) {  // page boundary
				error_f |= spi_wr_wait_tx(sp_i);  // SPI write wait
				error_f |= spi_wr_en_tx(sp_i);  // SPI wr enable
				cmd_ss << hex << showbase << spi_addr << " " << flush;  // addr
			}
			if (line != "0") { cmd_ss << "0x"; }  // no 0x for zero data
			cmd_ss << line << " ws " << flush;  // data
			sp_tx(sp_i, cmd_ss.str());
			error_f |= sp_rx_wait(sp_i);
			spi_addr += 4;
			if (!(spi_addr % EEPROM_PG_BYTES)) { cout << "." << flush; }  // show activity
		}
		// final SPI CS
		error_f |= spi_wr_wait_tx(sp_i);  // SPI write wait
		error_f |= spi_wr_prot_tx(sp_i, true);  // protect
		if (error_f) { sp_rx_flush(sp_i); }  // flush rx if error
	}
	return(error_f);
}

// preset string => slot
// return true if error
bool slot_str_tx(const sp_type& sp_i, int32_t slot_i, const string& pre_str_i) {
	bool error_f = sp_open_err(sp_i);  // check port
	if (!error_f) { error_f = slot_error(slot_i); }
	if (!error_f) {
		uint32_t slot = slot_i & 0xff;  // to unsigned
		if (!error_f) { error_f = spi_str_tx(sp_i, slot * EEPROM_PG_BYTES, pre_str_i); }
		if (error_f) { sp_rx_flush(sp_i); }  // flush rx if error
	}
	return(error_f);
}

// preset file => slot
// return true if error
bool dlp_slot_tx(const sp_type& sp_i, int32_t slot_i, const string& dlp_i) {
	string pre_str = "";
	bool error_f = file_to_str(dlp_i, pre_str);
	if (!error_f) { error_f = slot_str_tx(sp_i, slot_i, pre_str); }
	return(error_f);
}

// slot => preset file
// return true if error
bool dlp_slot_rx(const sp_type& sp_i, int32_t slot_i, const string& pfname_i) {
	string pre_str = "";
	bool error_f = slot_str_rx(sp_i, slot_i, pre_str);
	if (!error_f) { error_f = str_to_file(pfname_i, pre_str, false); }	
	return(error_f);
}

// editor => slot
// return true if error
bool ed_slot_tx(const sp_type& sp_i, int32_t slot_i) {
	uint8_t pre_uba[PRE_PARAMS] = { 0 };
	string pre_str = "";
	bool sys_f = slot_sys(slot_i);
	ed_to_uba(pre_uba, sys_f);
	uba_to_str(pre_uba, pre_str, PRE_PARAMS);
	bool error_f = slot_str_tx(sp_i, slot_i, pre_str);
	return(error_f);
}

// slot => editor
// return true if error
bool ed_slot_rx(const sp_type& sp_i, int32_t slot_i) {
	string pre_str = "";
	uint8_t pre_uba[PRE_PARAMS] = { 0 };
	bool sys_f = slot_sys(slot_i);
	bool error_f = slot_str_rx(sp_i, slot_i, pre_str);
	if (!error_f) { error_f = hex_str_to_uba(pre_str, pre_uba, 4, PRE_PARAMS); }
	if (!error_f) { uba_to_ed(pre_uba, sys_f); }
	return(error_f);
}

// knobs => string
// return true if error
bool knobs_str_rx(const sp_type& sp_i, string& str_o) {
	string str;
	bool error_f = sp_open_err(sp_i);  // check port
	if (!error_f) {
		// address calcs
		uint32_t end_knob = KNOBS - 1;
		// flush input
		sp_rx_flush(sp_i);
		// read data burst to string
		string str_rx = "";
		stringstream cmd_ss;
		cmd_ss << hex << showbase << 0 << " " << end_knob << " rk " << flush;
		sp_tx(sp_i, cmd_ss.str());  // tx cmd
		sp_rx_burst(sp_i, str_rx);  // rx str
		error_f = str_addr_strip(str_rx, str_o);  // strip cmd crap
	}
	return(error_f);
}

// knobs => editor
// get everything including page menu, load/stor, Load/Stor encoders
// return true if error
bool ed_knobs_rx(const sp_type& sp_i) {
	string str;
	uint8_t uba[KNOBS] = { 0 };
	bool error_f = knobs_str_rx(sp_i, str);
	if (!error_f) { 
		error_f = hex_str_to_uba(str, uba, 1, KNOBS);
		ed_clear(true); 
		ed_clear(false); 
	}
	// loop over KNOBS to fill editor preset
	uint32_t page = 0;
	while (!error_f && (page < PAGES)) {
		uint32_t page_knob = 0;
		while (!error_f && (page_knob < PAGE_KNOBS)) {
			uint32_t knob = (page * PAGE_KNOBS) + page_knob;
			uint32_t idx = knob_tbl.idx[knob];
			bool sys_f = knob_tbl.sys_f[knob];
			bool sgn_f = knob_tbl.sgn_f[knob];
			int32_t enc = (sgn_f) ? int8_t(uba[knob]) : uba[knob];
			error_f = ed_idx_wr(enc, idx, sys_f);
			page_knob++;
		}
		page++;
	}
	return(error_f);
}

// editor => knobs
// if sys don't tx page menu, load/stor, Load/Stor encoders
// return true if error
bool ed_knobs_tx(const sp_type& sp_i, bool sys_i) {
	bool error_f = false;
	// loop over KNOBS and tx
	uint32_t page = 0;
	while (!error_f && (page < PAGES)) {
		uint32_t page_knob = 0;
		while (!error_f && (page_knob < PAGE_KNOBS)) {
			uint32_t knob = (page * PAGE_KNOBS) + page_knob;
			if ((!sys_i && !knob_tbl.sys_f[knob]) || (sys_i && (knob_tbl.idx[knob] < MENU_PG_IDX))) {
				int32_t enc = 0;
				error_f = ed_knob_rd(enc, knob);
				if (!error_f) { error_f = knob_tx(sp_i, knob, enc); }
			}
			page_knob++;
		}
		page++;
	}
	return(error_f);
}

// slots => preset string vector
// return true if error
bool slots_sv_rx(const sp_type& sp_i, int32_t min_i, int32_t max_i, vector<string>& pre_sv_o) {
	bool error_f = slot_error(min_i);
	if (!error_f) { error_f = slot_error(max_i); }
	if (!error_f) { error_f = (max_i < min_i); }
	int32_t slot = min_i;
	pre_sv_o.clear();
	string pre_str;
	if (!error_f) { cout << flush; }
	while (!error_f && (slot <= max_i)) {
		error_f |= slot_str_rx(sp_i, slot, pre_str);
		pre_sv_o.push_back(pre_str);
		slot++;
		if (!(slot % 4)) { cout << "." << flush; }  // show activity
	}
	return(error_f);
}

// slots => preset string vector, match file names
// return true if error
bool slots_sv_names_rx(const sp_type& sp_i, int32_t min_i, int32_t max_i, const vector<string>& dir_sv_i, vector<string>& dl_pre_sv_o, vector<string>& dl_names_o) {
	vector<string> dir_pre_sv;
	bool error_f = dir_pre_sv_rd(dir_sv_i, dir_pre_sv);
	if (!error_f) { error_f = slots_sv_rx(sp_i, min_i, max_i, dl_pre_sv_o); }
	if (!error_f) {
		dl_names_o.clear();  // clear names
		for (uint32_t i=0; i<dl_pre_sv_o.size(); i++) {
			uint32_t idx = 0;
			if (sv_str_match(dl_pre_sv_o[i], dir_pre_sv, 0, idx)) { dl_names_o.push_back(dir_sv_i[idx]); }
			else { dl_names_o.push_back(PRE_NO_MATCH); }
		}
	}
	return(error_f);
}

// upload file to eeprom
// start at addr_i, write until end of file
// return true if error
bool pump_tx(const sp_type& sp_i, const string& fname_i, uint32_t addr_i, bool rst_i) {
	string str_tx = "";
	bool error_f = file_to_str(fname_i, str_tx);  // file => str
	if (!error_f) {
		error_f = spi_str_tx(sp_i, addr_i, str_tx);
		if (!error_f && rst_i) {
			// reset all threads
			sp_tx(sp_i, "0 0xff000000 wr ");
			error_f |= sp_rx_wait(sp_i);
			if (error_f) { sp_rx_flush(sp_i); }  // flush rx if error
		}
	}
	return(error_f);
}

// download eeprom to file
// start at addr_i, write bytes_i
// return true if error
bool dump_rx(const sp_type& sp_i, const string& fname_i, uint32_t addr_i, uint32_t bytes_i) {
	bool error_f = false;
	bool append_f = false;
	string str_rx = "";
	int32_t addr = addr_i;
	int32_t addr_inc = EEPROM_PG_BYTES;
	int32_t addr_end = addr_i + bytes_i - 1;
	while (!error_f && (addr < addr_end)) {
		str_rx = "";
		error_f = spi_str_rx(sp_i, addr, addr+addr_inc-1, str_rx);
		if (!error_f) { error_f = str_to_file(fname_i, str_rx, append_f); }
		addr += addr_inc;
		append_f = true;
		cout << "." << flush;  // show activity
	}
	return(error_f);
}

// do command, optionally get result
// return true if error
bool cmd_tx(const sp_type& sp_i, string& str_io, bool rd_f) {
	bool error_f = sp_open_err(sp_i);  // check port
	if (!error_f) {
		// flush input
		sp_rx_flush(sp_i);
		// read data burst to string
		string str_rx = "";
		sp_tx(sp_i, str_io);  // tx cmd
		ms_sleep(RX_CRC_MDLY);  // delay for CRC calc
		sp_rx_burst(sp_i, str_rx);  // rx str
		// process lines in str
		string line = "";
		stringstream ss_rx(str_rx);
		if (!getline(ss_rx, line)) { error_f = true; }
		if (str_io != line) { error_f = true; }
		str_io = "";
		if (rd_f) { if (!getline(ss_rx, str_io)) { error_f = true; } }
	}
	return(error_f);
}
