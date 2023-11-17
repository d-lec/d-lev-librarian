package main

/*
 * d-lib support functions
*/

import (
	"math"
	"fmt"
	"strings"
	"os"
	"path/filepath"
	"math/rand"
)

type param_t struct {
	ptype int
	plabel string
	pname string
}

// these are in preset / profile / slot order
var pro_params = []param_t {  
	{0x01, "50Hz",    "s_p0_ds"},	// 0
	{0x07, "Dith",    "s_p1_ds"},	// 1
	{0x03, "P<>V",    "s_p2_ds"},	// 2
	{0x01, "Erev",    "s_p3_ds"},	// 3
	{0x07, "Dith",    "s_p4_ds"},	// 4
	{0x24, "LCD ",    "s_p5_ds"},	// 5
	{0xcb, "Pcal",    "p_p0_ds"},	// 6
	{0xc0, "Lin ",    "p_p1_ds"},	// 7
	{0x27, "Ofs-",    "p_p2_ds"},	// 8
	{0xcb, "Sens",    "p_p3_ds"},	// 9	
	{0x27, "Ofs+",    "p_p4_ds"},	// 10
	{0xfd, "Cent",    "p_p5_ds"},	// 11
	{0xcb, "Vcal",    "v_p0_ds"},	// 12
	{0xc0, "Lin ",    "v_p1_ds"},	// 13
	{0x27, "Ofs-",    "v_p2_ds"},	// 14
	{0xcb, "Sens",    "v_p3_ds"},	// 15
	{0x24, "Drop",    "v_p4_ds"},	// 16
	{0x27, "Ofs+",    "v_p5_ds"},	// 17
	{0x31, "Mon ",    "v_p6_ds"},	// 18
	{0x01, "Out ",    "v_p7_ds"},	// 19
	{0x20, "LED ",    "t_p0_ds"},	// 20
	{0x04, "Qant",    "t_p1_ds"},	// 21
	{0x03, "Post",    "t_p2_ds"},	// 22
	{0xab, "Note",    "t_p3_ds"},	// 23
	{0xaf, "Oct ",    "t_p4_ds"},	// 24
	{0xc5, "Bass",   "eq_p0_ds"},	// 25
	{0xc5, "Treb",   "eq_p1_ds"},	// 26
	{0x31, "Line",    "v_p8_ds"},	// 27
	{0x7d, "Wait",    "s_p6_ds"},	// 28
	{0x24, "Lift",    "p_p6_ds"},	// 29
	{0x7e, "Auto",    "p_p7_ds"},	// 30
}

// these are in sequence
var not_params = []param_t {  
	{0x7f, "    ", "menu_pg_ds"},	// 31 - NOT stored in *.dlp !  MENU_PG_IDX!
	{0x7e, "load",   "ps_p0_ds"},	// 32 - NOT stored in *.dlp !
	{0x7e, "stor",   "ps_p1_ds"},	// 33 - NOT stored in *.dlp !
	{0x05, "Load",   "ps_p2_ds"},	// 34 - NOT stored in *.dlp !
	{0x05, "Stor",   "ps_p3_ds"},	// 35 - NOT stored in *.dlp !
}

// these are in preset / profile / slot order
var pre_params = []param_t {  
	// oscillators:
	{0x31, "osc ",  "o_p0_ds"},	// 0
	{0x24, "odd ",  "o_p1_ds"},	// 1
	{0xcd, "harm",  "o_p2_ds"},	// 2
	{0xca, "pmod",  "o_p3_ds"},	// 3
	{0xca, "vmod",  "o_p4_ds"},	// 4
	{0xa7, "oct ",  "o_p5_ds"},	// 5
	{0xf0, "offs",  "o_p6_ds"},	// 6
	{0xcd, "xmix",  "o_p7_ds"},	// 7
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
	{0x31, "sprd", "o_p19_ds"},	// 19
	{0xcd, "xmix", "o_p20_ds"},	// 20
	// noise gen:
	{0x31, "nois",  "n_p0_ds"},	// 21
	{0x70, "freq",  "n_p3_ds"},	// 22
	{0x76, "reso",  "n_p4_ds"},	// 23
	{0xa4, "mode",  "n_p5_ds"},	// 24
	{0xca, "pmod",  "n_p6_ds"},	// 25
	{0xca, "vmod",  "n_p7_ds"},	// 26
	{0xca, "pmod",  "n_p8_ds"},	// 27
	{0xca, "vmod",  "n_p9_ds"},	// 28
	{0x30, "puls", "n_p10_ds"},	// 29
	{0xc5, "bass", "n_p11_ds"},	// 30
	{0xcd, "xmix", "n_p12_ds"},	// 31
	{0xc5, "treb", "n_p13_ds"},	// 32
	{0x24, "duty", "n_p14_ds"},	// 33
	// resonator:
	{0xf1, "reso",  "r_p0_ds"},	// 34
	{0xc6, "harm",  "r_p1_ds"},	// 35
	{0x72, "freq",  "r_p2_ds"},	// 36
	{0xc6, "tap ",  "r_p3_ds"},	// 37
	{0x71, "hpf ",  "r_p4_ds"},	// 38
	{0xcd, "xmix",  "r_p5_ds"},	// 39
	{0xa2, "mode",  "r_p6_ds"},	// 40
	// formants:
	{0x70, "freq",  "f_p0_ds"},	// 41
	{0xd2, "levl",  "f_p1_ds"},	// 42
	{0x70, "freq",  "f_p2_ds"},	// 43
	{0xd2, "levl",  "f_p3_ds"},	// 44
	{0x70, "freq",  "f_p4_ds"},	// 45
	{0xd2, "levl",  "f_p5_ds"},	// 46
	{0xca, "pmod",  "f_p6_ds"},	// 47
	{0xca, "vmod",  "f_p7_ds"},	// 48
	{0xca, "pmod",  "f_p8_ds"},	// 49
	{0xca, "vmod",  "f_p9_ds"},	// 50
	{0xca, "pmod", "f_p10_ds"},	// 51
	{0xca, "vmod", "f_p11_ds"},	// 52
	{0x76, "reso", "f_p12_ds"},	// 53
	{0x76, "reso", "f_p13_ds"},	// 54
	{0x70, "freq", "f_p14_ds"},	// 55
	{0xd2, "levl", "f_p15_ds"},	// 56
	{0x70, "freq", "f_p16_ds"},	// 57
	{0xd2, "levl", "f_p17_ds"},	// 58
	{0x70, "freq", "f_p18_ds"},	// 59
	{0xd2, "levl", "f_p19_ds"},	// 60
	{0x76, "reso", "f_p20_ds"},	// 61
	{0x70, "freq", "f_p22_ds"},	// 62
	{0xd2, "levl", "f_p23_ds"},	// 63
	{0x70, "freq", "f_p24_ds"},	// 64
	{0xd2, "levl", "f_p25_ds"},	// 65
	{0xca, "pmod", "f_p26_ds"},	// 66
	{0xca, "vmod", "f_p27_ds"},	// 67
	{0x76, "reso", "f_p28_ds"},	// 68
	// pitch correction:
	{0x24, "cmod", "pc_p0_ds"},	// 69
	{0x24, "rate", "pc_p1_ds"},	// 70
	{0x44, "span", "pc_p2_ds"},	// 71
	{0x24, "corr", "pc_p3_ds"},	// 72
	{0xc9, "vmod", "pc_p4_ds"},	// 73
	// envelope gen:
	{0x25, "kloc",  "e_p0_ds"},	// 74
	{0x44, "knee",  "e_p1_ds"},	// 75
	{0x76, "fall",  "e_p2_ds"},	// 76
	{0x75, "rise",  "e_p3_ds"},	// 77
	{0xd1, "velo",  "e_p4_ds"},	// 78
	{0x73, "damp",  "e_p5_ds"},	// 79
	{0x25, "dloc",  "e_p6_ds"},	// 80
	// pitch preview:
	{0x31, "prev", "pp_p0_ds"},	// 81
	{0xc5, "harm", "pp_p1_ds"},	// 82
	{0xa7, "oct ", "pp_p2_ds"},	// 83
	{0xca, "pmod", "pp_p3_ds"},	// 84
	{0x0b, "mode", "pp_p4_ds"},	// 85
	{0xc5, "tone", "pp_p5_ds"},	// 86
	{0xca, "vmod", "pp_p6_ds"},	// 87
	// midi:
	{0xb0, "chan",  "m_p0_ds"},	// 88
	{0x25, "vloc",  "m_p1_ds"},	// 89
	{0x42, "bend",  "m_p2_ds"},	// 90
	{0xa7, "oct ",  "m_p3_ds"},	// 91
	{0x30, "velo",  "m_p4_ds"},	// 92
	{0xfc, "cc  ",  "m_p5_ds"},	// 93
	{0x45, "cloc",  "m_p6_ds"},	// 94
	// misc:
	{0x30, "cvol",  "e_p7_ds"},	// 95
	{0xa3, "bank",  "b_p0_ds"},	// 96
}

// these are in UI page order (hcl rk & wk knob order)
var knob_pnames = []string {  
	"v_p6_ds",  "v_p7_ds",  "v_p0_ds",  "p_p0_ds",  "ps_p1_ds", "b_p0_ds",  "ps_p0_ds", "menu_pg_ds",  // [0:7] D-LEV
	"v_p6_ds",  "v_p8_ds",  "pp_p0_ds", "eq_p1_ds", "o_p0_ds",  "eq_p0_ds", "n_p0_ds",  "menu_pg_ds",  // [8:15] LEVELS
	"pp_p6_ds", "pp_p3_ds", "pp_p0_ds", "pp_p4_ds", "pp_p1_ds", "pp_p5_ds", "pp_p2_ds", "menu_pg_ds",  // [16:23] PREVIEW : vmod, pmod, prev, mode, harm, tone, oct
	"m_p1_ds",  "m_p4_ds",  "m_p6_ds",  "m_p5_ds",  "m_p2_ds",  "m_p0_ds",  "m_p3_ds",  "menu_pg_ds",  // [24:31] MIDI
	"e_p0_ds",  "e_p3_ds",  "e_p1_ds",  "e_p2_ds",  "e_p4_ds",  "e_p5_ds",  "e_p6_ds",  "menu_pg_ds",  // [32:39] VOLUME
	"pc_p4_ds", "pc_p0_ds", "pc_p1_ds", "e_p7_ds", "pc_p3_ds", "pc_p2_ds",  "t_p2_ds",  "menu_pg_ds",  // [40:47] PITCH
	"n_p9_ds",  "n_p8_ds",  "n_p0_ds",  "n_p13_ds", "n_p10_ds", "n_p11_ds", "n_p14_ds", "menu_pg_ds",  // [48:55] NOISE
	"n_p7_ds",  "n_p6_ds",  "n_p3_ds",  "n_p0_ds",  "n_p5_ds",  "n_p12_ds", "n_p4_ds",  "menu_pg_ds",  // [56:63] FLT_NOISE
	"o_p4_ds",  "o_p3_ds",  "o_p2_ds",  "o_p15_ds", "o_p1_ds",  "o_p14_ds", "o_p5_ds",  "menu_pg_ds",  // [64:71] 0_OSC
	"o_p6_ds",  "o_p18_ds", "o_p16_ds", "o_p17_ds", "o_p8_ds",  "o_p19_ds", "o_p7_ds",  "menu_pg_ds",  // [72:79] 1_OSC
	"o_p13_ds", "o_p12_ds", "o_p9_ds",  "o_p0_ds",  "o_p11_ds", "o_p20_ds", "o_p10_ds", "menu_pg_ds",  // [80:87] FLT_OSC
	"r_p2_ds",  "r_p3_ds",  "r_p4_ds",  "r_p1_ds",  "r_p6_ds",  "r_p5_ds",  "r_p0_ds",  "menu_pg_ds",  // [95:95] RESON
	"f_p7_ds",  "f_p6_ds",  "f_p0_ds",  "f_p1_ds",  "f_p14_ds", "f_p15_ds", "f_p12_ds", "menu_pg_ds",  // [96:103] 0_FORM
	"f_p9_ds",  "f_p8_ds",  "f_p2_ds",  "f_p3_ds",  "f_p16_ds", "f_p17_ds", "f_p13_ds", "menu_pg_ds",  // [104:111] 1_FORM
	"f_p11_ds", "f_p10_ds", "f_p4_ds",  "f_p5_ds",  "f_p22_ds", "f_p23_ds", "f_p20_ds", "menu_pg_ds",  // [112:119] 2_FORM
	"f_p27_ds", "f_p26_ds", "f_p18_ds", "f_p19_ds", "f_p24_ds", "f_p25_ds", "f_p28_ds", "menu_pg_ds",  // [120:127] 3_FORM
	"v_p0_ds",  "v_p4_ds",  "v_p1_ds",  "s_p4_ds",  "v_p2_ds",  "v_p5_ds",  "v_p3_ds",  "menu_pg_ds",  // [128:135] V_FIELD
	"p_p0_ds",  "p_p6_ds",  "p_p1_ds",  "s_p1_ds",  "p_p2_ds",  "p_p4_ds",  "p_p3_ds",  "menu_pg_ds",  // [136:143] P_FIELD
	"t_p0_ds",  "p_p5_ds",  "s_p5_ds",  "t_p3_ds",  "t_p1_ds",  "t_p4_ds",  "t_p2_ds",  "menu_pg_ds",  // [144:151] DISPLAY
	"s_p6_ds",  "p_p7_ds",  "s_p2_ds",  "s_p0_ds",  "ps_p3_ds",  "s_p3_ds", "ps_p2_ds", "menu_pg_ds",  // [152:159] SYSTEM
}

// these are in UI screens order
var page_names = []string {
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
	"   SYSTEM",
}

// return ptype max value
func ptype_max(ptype int) int {
	max := 0
	switch {
		case ptype < 0x20 : // 0 thru 31
			max = ptype
		case ptype < 0x70 :  // 31, 63, 127, 255
			max = (1 << ((ptype & 0x3) + 5)) - 1
		case ptype < 0x72 : 
			max = 192
		case ptype == 0x72 :
			max = 127
		case ptype < 0x78 :
			max = 63
		case ptype == 0x7d :
			max = 99
		case ptype == 0x7e :
			max = 249
		case ptype == 0x7f :
			max = 19
		case ptype < 0xc0 :
			max = ptype & 0x1f
		case ptype < 0xf0 :  // 15, 31, 63, 127
			max = (1 << ((ptype & 0x3) + 4)) - 1
		case ptype == 0xf0 :
			max = 127
		case ptype == 0xf1 :
			max = 63
		case ptype == 0xfc :
			max = 31
		case ptype == 0xfd :
			max = 99
		default:
			max = 0
	}
	return max
	}

// return type min value
func ptype_min(ptype int) int {
	min := 0
	switch {
		case ptype < 0x80 :
			min = 0
		case ptype < 0xa0 :  // -1, -2, ..., -16
			min = -(ptype_max(ptype) + 1)
		case ptype == 0xfc :
			min = -127
		default:
			min = -ptype_max(ptype)
	}
	return min
}

// return limited value of type
func pint_lim(pint, ptype int) int {
	min := ptype_min(ptype)
	max := ptype_max(ptype)
	if pint < min { return min }
	if pint > max { return max }
	return pint
}

// return signed value of type
func pint_signed(pint, ptype int) int {
	if ptype_min(ptype) < 0 { return int(int8(uint8(pint))) }
	return pint
}

// make preset data signed
func pints_signed(pints []int, pro bool) ([]int) {
	if pro {
		for pidx, param := range pro_params {
			pints[pidx] = pint_signed(pints[pidx], param.ptype)
		}
	} else {
		for pidx, param := range pre_params {
			pints[pidx] = pint_signed(pints[pidx], param.ptype)
		}
	}
	return pints
}

// return filter freq value (type 0x70, 0x71)
// 7041 * EXP2((ENC * (2^27) / 24) + 3/4)
// [0:192] => [27:7040] (Hz)
func filt_freq(pint int) (int) {
	enc_mo := float64(int64(pint) * ((1 << 27) / 24) + 0xc0000000)
	return int(7041 * (math.Pow(2, enc_mo / math.Pow(2, 27)) / math.Pow(2, 32)))
}

// return reson freq value (type 0x72)
// 48001 / ((((((~(ENC<<25))^4)*0.871)+((~(ENC<<25))>>3))>>22)+4)
// [0:127] => [46:9600] (Hz)
func reson_freq(pint int) (int) {
	fs_rev := uint64(^(uint32(pint) << 25))
	sq := (fs_rev * fs_rev) >> 32
	qd := (sq * sq) >> 32
	return int(48001 / ((((uint64(float64(qd) * 0.871) + (fs_rev >> 3)) >> 22) + 4)))
}

// given pint and type, return frequency
func pint_freq(pint, ptype int) (int) {
	switch ptype {
		case 0x70, 0x71 : return filt_freq(pint)
		case 0x72 : return reson_freq(pint)
		default : return pint_signed(pint, ptype)
	}
}

// given pint and type, return display string[5]
func pint_disp(pint, ptype int) (string) {
	return fmt.Sprintf("%5v", pint_freq(pint, ptype))
}

// return closest filter pint given freq (type 0x70, 0x71)
// [27:7040] (Hz) => [0:192]
func filt_pint(freq int) (int) {
	p_best := 0
	err_best := 0
	first := true
	for p := ptype_min(0x70); p <= ptype_max(0x70); p++ {
		err := (freq - filt_freq(p)) * (freq - filt_freq(p))
		if err < err_best || first { 
			err_best = err
			p_best = p
		}
		first = false
	}
	return p_best
}

// return closest reson pint given freq (type 0x72)
// [46:9600] (Hz) => [0:127]
func reson_pint(freq int) (int) {
	p_best := 0
	err_best := 0
	first := true
	for p := ptype_min(0x72); p <= ptype_max(0x72); p++ {
		err := (freq - reson_freq(p)) * (freq - reson_freq(p))
		if err < err_best || first { 
			err_best = err
			p_best = p
		}
		first = false
	}
	return p_best
}

// given frequency and type, return limited pint
func freq_pint(freq, ptype int) (int) {
	switch ptype {
		case 0x70, 0x71 : return filt_pint(freq)
		case 0x72 : return reson_pint(freq)
		default : return pint_lim(freq, ptype)
	}
}

// given pname, return ptype, plabel, pidx, pgroup
func pname_lookup(pname string) (int, string, int, string) {
	for pidx, param := range pre_params {
		if pname == param.pname { return param.ptype, param.plabel, pidx, "pre" }
	}
	for pidx, param := range pro_params {
		if pname == param.pname { return param.ptype, param.plabel, pidx, "pro" }
	}
	for pidx, param := range not_params {
		if pname == param.pname { return param.ptype, param.plabel, pidx, "not" }
	}
	return 0, "", 0, ""
}

// given pidx & pro, return kidx, kflg
func knob_lookup(pidx int, pro bool) (int, bool) {
	if pro {
		if pidx >= len(pro_params) { return 0, false }
		for kidx, kname := range knob_pnames {
			if kname == pro_params[pidx].pname { return kidx, true }
		}
	} else {
		if pidx >= len(pre_params) { return 0, false }
		for kidx, kname := range knob_pnames {
			if kname == pre_params[pidx].pname { return kidx, true }
		}
	}
	return 0, false
}

// given partial page str, return name & idx
// gotcha: for ""  HasPrefix always returns true!
func page_lookup(page string) (string, int) {
	if page != "" {
		page = strings.ToUpper(page)
		for idx, name := range page_names {
			name = strings.TrimSpace(name)
			if strings.HasPrefix(name, page) { return name, idx }
		}
	}
	return "", -1  // default
}

// put knob ints in preset / slot order
func knob_pre_order(kints []int, mode string) ([]int) {
	pints := make([]int, SLOT_BYTES)
	for kidx, kname := range knob_pnames {
		_, _, pidx, pmode := pname_lookup(kname)
		if mode == pmode {
			pints[pidx] = kints[kidx]
		}
	}
	return pints
}

// put knob hex str in preset / slot order, return hex string
func knob_pre_str(knob_str string, pro bool) (string) {
	str_split := (strings.Split(strings.TrimSpace(knob_str), "\n"))
	if len(str_split) < KNOBS_TOTAL { error_exit("Bad knob info") }
	hex_str := ""
	line_str := ""
	for pidx:=0; pidx<SLOT_BYTES; pidx++ {
		kidx, kflg := knob_lookup(pidx, pro)
		if kflg { line_str = fmt.Sprintf("%02s", str_split[kidx]) + line_str } else { line_str = "00" + line_str }
		if pidx % 4 == 3 { 
			hex_str += line_str + "\n" 
			line_str = ""
		}
	}
	return hex_str
}

// split eeprom str into pre / pro / spi strings
func split_eeprom_str(eeprom_str string) (string, string, string) {
	str_split := strings.Split(strings.TrimSpace(eeprom_str), "\n")
	pre_str := ""
	pro_str := ""
	spi_str := ""
	const lines = SLOT_BYTES/EE_RW_BYTES
	for line, str := range str_split {
		if line < PRE_SLOTS*lines { 
			pre_str += str + "\n"
		} else if line < SLOTS*lines { 
			pro_str += str + "\n"
		} else { 
			spi_str += str + "\n"
		}
	}
	return pre_str, pro_str, spi_str
}

// split pre/pro str into dlp strings
func split_pre_pro_str(pre_pro_str string) ([]string) {
	str_split := strings.Split(strings.TrimSpace(pre_pro_str), "\n")
	var dlp_strs []string
	dlp_str := ""
	const lines = SLOT_BYTES/EE_RW_BYTES
	for line, str := range str_split {
		dlp_str += str + "\n"
		if line % lines == lines - 1 { 
			dlp_strs = append(dlp_strs, dlp_str)
			dlp_str = ""
		}
	}
	return dlp_strs
}

// read PRE|PRO|EEPROM file and extract dlp strings
func file_dlp_strs(file, mode string) ([]string, string) {
	file_str := ""
	ext := filepath.Ext(file)
	switch ext {
		case ".pre" :
			file_str = file_read_str(file)
		case ".pro" :
			file_str = file_read_str(file)
			mode = "pro"
		case ".eeprom" :
			file_str = file_read_str(file)
			switch mode {
				case "pro" : _, file_str, _ = split_eeprom_str(file_str)
				default : file_str, _, _ = split_eeprom_str(file_str)
			}
		case "" : error_exit(fmt.Sprint("Missing file extension"))
		default : error_exit(fmt.Sprint("Wrong file extension: ", ext))
	}
	return split_pre_pro_str(file_str), mode
}

// read SPI|EEPROM file and extract software strings
func file_sw_strs(file string) ([]string) {
	ext := filepath.Ext(file)
	file_str := ""
	switch ext {
		case ".spi" :
			file_str = file_read_str(file)
		case ".eeprom" :
			file_str = file_read_str(file)
			_, _, file_str = split_eeprom_str(file_str)
		case "" : error_exit(fmt.Sprint("Missing file extension"))
		default : error_exit(fmt.Sprint("Wrong file extension: ", ext))
	}
	return strings.Split(strings.TrimSpace(file_str), "\n")
}

// generate knob ui display strings
func knob_ui_strs(hex_str string) ([]string) {
	kints := hexs_to_ints(hex_str, 1)
	if len(kints) < KNOBS_TOTAL { error_exit("Bad knob info") }
	var strs []string
	for kidx, kname := range knob_pnames {
		ptype, plabel, _, _ := pname_lookup(kname)
		if kidx % KNOBS == PAGE_SEL_KNOB { 
			strs = append(strs, page_names[kidx / KNOBS])
		} else { 
			strs = append(strs, plabel + pint_disp(kints[kidx], ptype)) 
		}
	}
	return strs
}

// generate pre / pro / slot ui display strings
func pre_ui_strs(hex_str string, pro bool) ([]string) {
	pints := hexs_to_ints(hex_str, 4)
	if len(pints) < SLOT_BYTES { error_exit("Bad file / slot info") }
	var strs []string
	for idx, pname := range knob_pnames {
		ptype, plabel, pidx, pgroup := pname_lookup(pname)
		if idx % KNOBS == PAGE_SEL_KNOB { 
			strs = append(strs, page_names[idx / KNOBS])
		} else { 
			if pro == (pgroup == "pro") && pgroup != "not" {
				strs = append(strs, plabel + pint_disp(pints[pidx], ptype)) 
			} else {
				strs = append(strs, plabel + "     ") 
			}
		}
	}
	return strs
}

// render ui display strings to printable string
func ui_prn_str(strs []string, mark int) (string) {
	if len(strs) < len(knob_pnames) { error_exit("Bad input info") }
	h_line_sub := "+" + strings.Repeat("-", 22);
	h_line := strings.Repeat(h_line_sub, PAGES_COLS) + "+\n";
	prn_str := h_line
	for prow:=0; prow<PAGES_ROWS; prow++ {
		for uirow:=0; uirow<KNOBS_ROWS; uirow++ {
			for pcol:=0; pcol<PAGES_COLS; pcol++ {
				idx := (prow * KNOBS * PAGES_COLS) + (uirow * KNOBS_COLS) + (pcol * KNOBS)
				prn_str += "|"
				for i:=0; i<KNOBS_COLS; i++ {
					if idx+i == mark { 
						prn_str += "[" + strs[idx+i] + "]"
					} else {
						prn_str += " " + strs[idx+i] + " "
					}
				}
			}
			prn_str += "|\n"
		}
		prn_str += h_line
	}
	return strings.TrimSpace(prn_str)
}

// generate preset diff display strings
func diff_pres(pre_str0, pre_str1 string, pro bool) ([]string, []string, []bool) {
	pints0 := hexs_to_ints(pre_str0, 4)
	pints1 := hexs_to_ints(pre_str1, 4)
	if (len(pints0) < SLOT_BYTES) || (len(pints1) < SLOT_BYTES) { error_exit("Bad preset info") }
	var strs0 []string
	var strs1 []string
	var diffs []bool
	for kidx, pname := range knob_pnames {
		ptype, plabel, pidx, pgroup := pname_lookup(pname)
		if kidx % KNOBS == PAGE_SEL_KNOB { 
			strs0 = append(strs0, page_names[kidx / KNOBS])
			strs1 = append(strs1, page_names[kidx / KNOBS])
			diffs = append(diffs, false)
		} else { 
			if pro == (pgroup == "pro") && pgroup != "not" {
				strs0 = append(strs0, plabel + pint_disp(pints0[pidx], ptype)) 
				if pints0[pidx] != pints1[pidx] {
					strs1 = append(strs1, plabel + pint_disp(pints1[pidx], ptype)) 
					diffs = append(diffs, true)
				} else {
					strs1 = append(strs1, plabel + "     ") 
					diffs = append(diffs, false)
				}
			} else {
				strs0 = append(strs0, plabel + "     ") 
				strs1 = append(strs1, plabel + "     ") 
				diffs = append(diffs, false)
			}
		}
	}
	return strs0, strs1, diffs
}

// render ui display strings to printable string
func diff_prn_str(strs0, strs1 []string, diffs []bool) (string) {
	if (len(strs0) < len(knob_pnames)) || (len(strs1) < len(knob_pnames)) || (len(diffs) < len(knob_pnames)) { error_exit("Bad input info") }
	h_line_sub := "+" + strings.Repeat("-", 22);
	h_line := strings.Repeat(h_line_sub, 2) + "+\n";
	prn_str := ""
	chgs := 0
	for scrn:=0; scrn<PAGES; scrn++ {
		page_str := ""
		chg_f := false
		for row:=0; row<KNOBS_ROWS; row++ {
			idx := scrn*KNOBS + row*KNOBS_COLS
			page_str += "| " + strs0[idx] + "  " + strs0[idx+1] + " "
			page_str += "| " + strs1[idx] + "  " + strs1[idx+1] + " "
			page_str += "|\n"
			if diffs[idx] { chg_f = true; chgs++ }
			if diffs[idx+1] { chg_f = true; chgs++ }
		}
		page_str += h_line
		if chg_f { prn_str += page_str }
	}
	if chgs != 0 { prn_str = h_line + prn_str }  // top line
	prn_str += fmt.Sprintln("> differences:", chgs)
	return strings.TrimSpace(prn_str)
}

// compare 2 pre strings, return sum of squares of differences
func comp_pres(pre_str0, pre_str1 string, pro bool) (int) {
	pints0 := pints_signed(hexs_to_ints(pre_str0, 4), pro)
	pints1 := pints_signed(hexs_to_ints(pre_str1, 4), pro)
	if (len(pints0) < SLOT_BYTES) || (len(pints1) < SLOT_BYTES) { error_exit("Bad preset info") }
	ssd := 0
	if pro {
		for pidx, _ := range pro_params {
			diff := pints0[pidx] - pints1[pidx]
			ssd += diff * diff
		}
	} else {
		for pidx, _ := range pre_params {
			diff := pints0[pidx] - pints1[pidx]
			ssd += diff * diff
		}
	}
	return ssd
}
		
// compare slot / file data to file data, return slice of names
func comp_file_data(data2_strs, name_strs, data_strs []string, pro, guess bool) ([]string) {
	var strs []string
	for _, data2_str := range data2_strs {
		first := true
		ssd_min := 0
		idx_min := 0
		for file_idx, data_str := range data_strs {
			ssd := comp_pres(data2_str, data_str, pro)
			if first || (ssd < ssd_min) {
				ssd_min = ssd
				idx_min = file_idx
				first = false
			}
		}
		if first { error_exit("No file data") }
		if ssd_min == 0 { 
			strs = append(strs, name_strs[idx_min])
		} else if guess {
			strs = append(strs, name_strs[idx_min] + " // ?(" + fmt.Sprint(math.Ceil(math.Sqrt(float64(ssd_min)))) + ")") 
		} else {
			strs = append(strs, "__????__") 
		}
	}
	return strs
}

// render slots match display strings to printable string
func slots_prn_str(strs []string, pro, hdr bool) (string) {
	//if len(strs) < SLOTS { error_exit("Bad slots info") }
	prn_str := ""  
	if pro {
		if hdr { prn_str += "/*\n pro slots [0:5]\n*/\n" }
		for row, str := range strs {
			if hdr { prn_str += strings.TrimSpace(str)
			} else { prn_str += fmt.Sprintf("[%1v] %s", row, str) }
			prn_str += "\n"
		}
	} else if hdr {
		for row, str := range strs {
			if (row % 10 == 0) {
				prn_str += fmt.Sprint("/*\n pre slots [", row, ":", row+9, "]\n*/\n")
			}
			prn_str += strings.TrimSpace(str) + "\n"
		}
	} else {
		const cols = 5
		const rows = PRE_SLOTS/cols
		// find minimum column widths
		var col_w [cols]int
		for row:=0; row<rows; row++ {
			for col:=0; col<cols; col++ {
				idx := col*rows + row
				str_len := len(strings.TrimSpace(strs[idx]))
				if str_len > col_w[col] { col_w[col] = str_len }
			}
		}
		// assemble print string
		for row:=0; row<rows; row++ {
			for col:=0; col<cols; col++ {
				idx := col*rows + row
				prn_str += fmt.Sprintf("[%2v] %-*s", idx, col_w[col]+2, strings.TrimSpace(strs[idx]))
			}
			prn_str += "\n"
		}
	}
	return prn_str
}

// render files match display strings to printable string
func files_prn_str(strs2, strs []string) (string) {
	// find minimum column width
	col_w := 0
	for _, str2 := range strs2 {
		str2_len := len(str2)
		if str2_len > col_w { col_w = str2_len }
	}
	// assemble print string
	prn_str := ""  
	for i, str := range strs {
		prn_str += fmt.Sprintf("%*s : %s", col_w+1, strs2[i], str)
		prn_str += "\n"
	}
	return prn_str
}


///////////
// morph //
///////////

// add scaled normal random value, reflect as necessary until inside limits
// flip / confine signed knobs to signed regions
func pint_rnd_refl(pint, min, max int, mult, ms float64) (int) {
	// randomly flip signed knob sign
	if (mult != 0) && (min < 0) && (rand.Float64() < ms) { pint = -pint }
	// confine to signed area
	if pint < 0 { max = 0 } else { min = 0 }
	norm := rand.NormFloat64()
	scale := mult * float64(max - min) / 100  // mult scaled 1/100 (percent)
	norm_offs := int(math.Round(scale * norm))
	pint += norm_offs
	for (pint < min) || (pint > max) {
		if pint > max { pint -= 2 * (pint - max) }
		if pint < min { pint += 2 * (min - pint) }
	}
	return pint
}

// morph most oscillator settings, all filter freqs, some resonator settings
func morph_pints(pints []int, mo, mn, me, mf, mr, ms float64) ([]int) {
	// function to interate over relevant 
	morph := func(m_pidx []int, mult float64) ([]int) {
		for _, pidx := range m_pidx {
			ptype := pre_params[pidx].ptype;
			pints[pidx] = pint_rnd_refl(pints[pidx], ptype_min(ptype), ptype_max(ptype), mult, ms) 
		}
		return pints
	}
	// OSC: odd, harm, offs, xmix, fm, hmul, hmul, offs, sprd:
	pints = morph([]int{ 1, 2, 6, 7, 8, 16, 17, 18, 19 }, mo)
	// NOISE: puls, duty:
	pints = morph([]int{ 29, 33 }, mn)
	// EQ: osc & noise bass, treb:
	pints = morph([]int{ 14, 15, 30, 32 }, me)
	// FILT: all formant freqs, osc & noise filter freqs too:
	pints = morph([]int{ 41, 55, 43, 57, 45, 62, 59, 64,  9, 22 }, mf)
	// RESON: harm, freq, tap:
	pints = morph([]int{ 35, 36, 37 }, mr)
	return pints
}


////////////
// update //
////////////

// batch read, process, write all *.dlp files from dir to dir2
func process_dlps(dir, dir2 string, pro, mono, update, robs, yes, dry bool) {
	dir_chk(dir)
	if dir2 == "" { dir2 = dir }  // default to same dir
	// prompt user
	if !dry {
		if dir == dir2 { 
			if !user_prompt("Overwrite DLP files in SOURCE directory " + dir + "?", yes, false) { return }
		} else if path_exists(dir2) { 
			if !user_prompt("Overwrite any DLP files in DESTINATION directory " + dir2 + "?", yes, false)  { return }
		}
	}
	files, err := os.ReadDir(dir); err_chk(err)
	upd_cnt := 0
	dlp_cnt := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".dlp" && file.IsDir() == false {
			file_name := file.Name()
			file_path := filepath.Join(dir, file_name)
			file_path2 := filepath.Join(dir2, file_name)
			dlp_cnt++
			// read in
			file_str := file_read_str(file_path)
			pints := pints_signed(hexs_to_ints(file_str, 4), pro)
			proc_str := ""
			zero_f := true
			for _, param := range pints {
				if param != 0 { zero_f = false }
			}
			//////////////
			// PROFILES //
			//////////////
			if pro && !zero_f {  // don't process blank profiles
				/////////////////////////
				// zero out high fluff //
				/////////////////////////
				nz_cnt := 0
				for idx, param := range pints {
					if idx >= len(pro_params) { 
						if param != 0 { 
							pints[idx] = 0
							nz_cnt++
						}
					}
				}
				if nz_cnt > 0 { 
					proc_str += fmt.Sprintln(" Fluffs zeroed:", nz_cnt)
				}
			/////////////
			// PRESETS //
			/////////////
			} else if !pro && !zero_f {  // don't process blank presets
				/////////////////////////
				// zero out high fluff //
				/////////////////////////
				nz_cnt := 0
				for idx, param := range pints {
					if idx >= len(pre_params) { 
						if param != 0 { 
							pints[idx] = 0
							nz_cnt++
						}
					}
				}
				if nz_cnt > 0 { 
					proc_str += fmt.Sprintln(" Fluffs zeroed:", nz_cnt)
				} 
				if update {
					///////////////////////
					// 2023-11-02 update //
					///////////////////////
					//////////////////////////
					// min dloc if damp > 0 //
					//////////////////////////
					damp := pints[79]
					dloc := pints[80]
					dloc_new := dloc
					if (damp > 0) && (dloc < 16) { dloc_new = 16 }
					if dloc != dloc_new { 
						pints[80] = dloc_new
						proc_str += fmt.Sprintln(" VOLUME:dloc", dloc, "=>", dloc_new) 
					}
				}
				if mono {
					////////////////////
					// stereo => mono //
					////////////////////
					xmix := pints[39]
					mode := pints[40]
					xmix_new := xmix
					mode_new := mode
					if xmix != 0 {  // if reson active
						if mode == 2 { 
							mode_new = 1  // parallel stereo => mono
						}
						if mode == -2 { 
							mode_new = -1  // series stereo => mono
							xmix_new = 0  // kill xmix
						}
						if mode != mode_new { 
							pints[40] = mode_new
							proc_str += fmt.Sprintln(" RESON:mode", mode, "=>", mode_new) 
						}
						if xmix != xmix_new { 
							pints[39] = xmix_new
							proc_str += fmt.Sprintln(" RESON:xmix", xmix, "=>", xmix_new) 
						}
					}
				}
				if robs {
					prev := pints[81]; 
					mode := pints[85];
					if (prev == 0) || (mode %4 != 0) {  // don't change 4th osc presets in use
						////////////////////////////
						// pp defaults for rob s. //
						////////////////////////////
						prev_new := 63
						harm := pints[82]; harm_new := 10
						oct := pints[83]; oct_new := 0
						pmod := pints[84]; pmod_new := -20
						mode_new := 7
						tone := pints[86]; tone_new := 0
						vmod := pints[87]; vmod_new := -55
						if prev != prev_new { 
							pints[81] = prev_new
							proc_str += fmt.Sprintln(" PREVIEW:prev", prev, "=>", prev_new) 
						}
						if harm != harm_new { 
							pints[82] = harm_new
							proc_str += fmt.Sprintln(" PREVIEW:harm", harm, "=>", harm_new) 
						}
						if oct != oct_new { 
							pints[83] = oct_new
							proc_str += fmt.Sprintln(" PREVIEW:oct", oct, "=>", oct_new) 
						}
						if pmod != pmod_new { 
							pints[84] = pmod_new
							proc_str += fmt.Sprintln(" PREVIEW:pmod", pmod, "=>", pmod_new) 
						}
						if mode != mode_new { 
							pints[85] = mode_new
							proc_str += fmt.Sprintln(" PREVIEW:mode", mode, "=>", mode_new) 
						}
						if tone != tone_new { 
							pints[86] = tone_new
							proc_str += fmt.Sprintln(" PREVIEW:tone", tone, "=>", tone_new) 
						}
						if vmod != vmod_new { 
							pints[87] = vmod_new
							proc_str += fmt.Sprintln(" PREVIEW:vmod", vmod, "=>", vmod_new) 
						}
					}
				}
			}
			// write file
			if !dry { file_write_str(file_path2, ints_to_hexs(pints, 4), true) }
			if proc_str != "" { 
				upd_cnt++ 
				fmt.Println("*", file_name, "changes *")
				fmt.Print(proc_str)
			} else {
				fmt.Println(file_name, "(ok)")
			}
		}
    }
	mode := "pre"; if pro { mode = "pro" }
	prn_str := fmt.Sprint("> -PROCESSED- ", upd_cnt, " of ", dlp_cnt, " ", mode, " files from ", dir, " to ", dir2)
	prn_str += " with FLAGS: "
	if mono { prn_str += "-m " }
	if update { prn_str += "-u " }
	if robs { prn_str += "-r " }
	if dry { prn_str += "-dry (no file write)" }
	fmt.Println(prn_str)
}


// pump list of dlps to default slots
func dlps_pump(path string, yes bool) {
	type slot_dlp_t struct {
		slot string
		fname string
	}
	var slots_dlps = []slot_dlp_t {  
		{"51", "dobro.dlp"},
		{"93", "bowl_1.dlp"},
		{"94", "bowl_3.dlp"},
		{"95", "little_ben_0.dlp"},
		{"96", "little_ben_1.dlp"},
		{"99", "wine_1.dlp"},	
		{"115", "bowl_2.dlp"},
		{"135", "vanbelis.dlp"},
	}
	for _, slot_dlp := range slots_dlps {
		slot := slot_dlp.slot
		dlp := filepath.Join(path, slot_dlp.fname)
		if user_prompt("Do you want to PUMP " + dlp + " to slot " + slot + " ?", yes, true) {
			pump_cmd(dlp, slot, false, false)
		}
	}
}

/////////
// dev //
/////////

// find DLP files with various values
func find_dlp(dir string) {
	files, err := os.ReadDir(dir); err_chk(err)
	dlp_cnt := 0
	find_cnt := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".dlp" && !file.IsDir() {
			dlp_cnt++
			file_name := file.Name()
			file_path := filepath.Join(dir, file_name)
			// read in
			file_str := file_read_str(file_path)
			pints := pints_signed(hexs_to_ints(file_str, 4), false)


			if pints[79] > 0 {  // damp
				fmt.Println("damp", pints[79])
				fmt.Println("dloc", pints[80])
				fmt.Println("-", file.Name(), "\n")
				find_cnt++
			}
/*
			if pints[7] == 31 {  // 1_osc:xmix
				fmt.Println("xmix", pints[7])
				fmt.Println("-", file.Name(), "\n")
				find_cnt++
			}
*/
/*
			if pints[75] > 0 && pints[74] < 42 {  // knee & kloc
				fmt.Println("kloc", pints[74])
				fmt.Println("knee", pints[75])
				fmt.Println("velo", pints[78])
				fmt.Println("-", file.Name(), "\n")
				find_cnt++
			}
*/
/*
			if pints[74] == 0 && pints[78] > 0  {  // kloc, velo
				fmt.Println("kloc", pints[74])
				fmt.Println("velo", pints[78])
				fmt.Println("-", file.Name(), "\n")
				find_cnt++
			}
*/
/*
			if pints[74] == 0 {  // kloc
				fmt.Println("kloc", pints[74])
				fmt.Println("-", file.Name(), "\n")
				find_cnt++
			}
*/
/*
			if (pints[0] > 0 || pints[21] > 0) && pints[78] == 0 {  // osc, nois, velo
				fmt.Println("osc", pints[0])
				fmt.Println("nois", pints[21])
				fmt.Println("velo", pints[78])
				fmt.Println("-", file.Name(), "\n")
				find_cnt++
			}
*/
/*
			if pints[21] > 0 && pints[29] > 0 && pints[7] > 0 {  // nois, puls, 1_osc:xmix
				fmt.Println("nois", pints[21])
				fmt.Println("puls", pints[29])
				fmt.Println("xmix", pints[7])
				fmt.Println("-", file.Name(), "\n")
				find_cnt++
			}
*/
/*			
			if pints[95] != 0 {  // cvol
				fmt.Println("cvol", pints[95])
				fmt.Println("-", file.Name(), "\n")
				find_cnt++
			}
*/			
			/*
//			if pints[21] != 0 && pints[32] > 0 {  // nois
//			if pints[21] != 0 && pints[30] < 0 {  // nz nois & neg bass
			if pints[21] > 52 {  // nois > 52
//				fmt.Println("treb", pints[32])
//				fmt.Println("bass", pints[30])
				fmt.Println("nois", pints[21])
				fmt.Println("-", file.Name(), "\n")
				find_cnt++
			}
			*/
			/*
			// pitch preview:
			{0x31, "prev", "pp_p0_ds"},	// 81
			{0xc5, "harm", "pp_p1_ds"},	// 82
			{0xa7, "oct ", "pp_p2_ds"},	// 83
			{0xca, "pmod", "pp_p3_ds"},	// 84
			{0x0b, "mode", "pp_p4_ds"},	// 85
			{0xc5, "treb", "pp_p5_ds"},	// 86
			{0xc5, "bass", "pp_p6_ds"},	// 87
			*/
			/*
			if pints[81] != 0 {  // prev != 0
				fmt.Println("prev", pints[81])
				fmt.Println("mode", pints[85])
				fmt.Println("-", file.Name(), "\n")
				find_cnt++
			}
			*/

/*
			if pints[81] != 0 {  // prev != 0
				fmt.Println("prev", pints[81])
				fmt.Println("harm", pints[82])
				fmt.Println("oct",  pints[83])
				fmt.Println("pmod", pints[84])
				fmt.Println("mode", pints[85])
				fmt.Println("tone", pints[86])
				fmt.Println("vmod", pints[87])
				fmt.Println("-", file.Name(), "\n")
				find_cnt++
			}
*/
/*
			if pints[39] != 0 && (pints[40] == 2 || pints[40] == -2) {  // xmix non-zero && |mode| == 2
//				fmt.Println("freq", pints[36])
//				fmt.Println("hpf", pints[38])
				fmt.Println("mode", pints[40])
//				fmt.Println("reso", pints[34])
//				fmt.Println("tap", pints[37])
//				fmt.Println("harm", pints[35])
				fmt.Println("xmix", pints[39])
				fmt.Println("-", file.Name(), "\n")
				find_cnt++
			}
*/

			
		}
    }
	fmt.Println("> examined", dlp_cnt, "DLP files")
	fmt.Println("> found", find_cnt, "instances")
}

