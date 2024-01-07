package main

/*
 * d-lib support functions
*/

import (
	"math/bits"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"math/rand"
)

// librarian commands
var lib_cmds = []string {  
	"acal",
	"bank",
	"batch",
	"copy",
	"diff",
	"dump",
	"help",
	"hcl",
	"join",
	"knob",
	"loop",
	"match",
	"menu",
	"morph",
	"port",
	"pump",
	"reset",
	"split",
	"stats",
	"ver",
	"view",
}

// bad commmand, give hint if match, exit
func cmd_hint(cmd string) {
	str := fmt.Sprint("Command '", cmd, "' not found")
	list := fuzzy_list(cmd, lib_cmds, 1, true)
	if len(list) > 0 { str += " (did you mean '" + list[0] + "' ?)" }
	error_exit(str)
}

// check for error, exit program if true
func err_chk(err error) {
	if err != nil { error_exit(err.Error()) }
}

// check slot limits
func slot_lim_chk(slot int, pro bool) {
	slots := PRE_SLOTS
	if pro { slots = PRO_SLOTS }
	if slot < 0 || slot >= slots { error_exit(fmt.Sprint("Slot out of range: ", slot)) }
}

// convert slot string, check lim
func slot_int_chk(slot_str string, pro bool) (int) {
	slot_int, err := strconv.Atoi(slot_str)
	if err != nil { error_exit("Bad slot number: " + slot_str) }
	slot_lim_chk(slot_int, pro)
	return slot_int
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

// given page:knob:value return relevant info
func pkv_parse(pkv string) (params int, page_name, knob_name string, knob_idx, val int) {
	pkv_split := strings.Split(pkv, ":")
	params = len(pkv_split)
	if params < 1 { error_exit("Too few parameters") }
	if params > 3 { error_exit("Too many parameters: " + pkv) }
	pkv_split = strs_trim(pkv_split)
	// parse page
	page_idx := 0
	page_idx, page_name = fuzzy_page_lookup(pkv_split[0])
	if page_idx < 0 { error_exit("Bad page name: " + pkv_split[0]) }
	knob_idx = (page_idx * KNOBS) + PAGE_SEL_KNOB  // default is page name
	// parse knob
	if params > 1 { 
		if strings.ToLower(pkv_split[1]) == "all" { 
			knob_name = "all"
		} else {
			var err error
			knob_idx, err = strconv.Atoi(pkv_split[1]); 
			if err != nil { error_exit("Bad knob index: " + pkv_split[1]) }
			if knob_idx < 0 || knob_idx >= KNOBS-1 { error_exit("Knob index out of range: " + pkv_split[1]) }
			knob_idx += page_idx * KNOBS  // full idx
			_, knob_name, _, _, _ = pname_lookup(knob_pnames[knob_idx])
			knob_name = strings.TrimSpace(knob_name)
		}
	}
	// parse value
	if params > 2 { 
		var err error
		val, err = strconv.Atoi(pkv_split[2]); 
		if err != nil { error_exit("Bad value: " + pkv_split[2]) }
	}
	return
}

// print librarian version & help info
func help_cmd(verbose_f bool) {
	fmt.Println("= D-LEV LIBRARIAN VERSION", ver_tbl[0].lib, ":", ver_tbl[0].date, "=") 
	fmt.Print(help_commands_str)
	if verbose_f { 
		fmt.Print(help_notes_str)
		fmt.Print(strings.Replace(help_examples_str, "LIB_EXE", os.Args[0], -1)) 
	}
}

// do processor reset
func reset_cmd() {
	sp := sp_open()
	sp_tx_rx(sp, "0 0xff000000 wr ", false)
	sp.Close()
	fmt.Println("> reset processor")
}

// do ACAL
func acal_cmd() {
	sp := sp_open()
	sp_tx_rx(sp, "acal ", false)
	sp.Close()
	fmt.Println("> acal issued")
}

// do HCL command
func hcl_cmd() {
	if len(os.Args) < 3 { error_exit("Command line is blank") }
	wr_str := strings.Join(os.Args[2:], " ") + " "
	sp := sp_open()
	rx_str := sp_tx_rx(sp, wr_str, false)
	sp.Close()
	fmt.Print(strings.TrimSpace(rx_str))
	fmt.Println(" hcl", wr_str)
}

// do loop command
func loop_cmd() {
	if len(os.Args) < 3 { error_exit("Loop text is blank") }
	wr_str := strings.Join(os.Args[2:], " ")
	sp := sp_open()
	rx_str := sp_tx_rx(sp, wr_str + ">", false)
	sp.Close()
	fmt.Println("> TX:", wr_str)
	fmt.Println("> RX:", strings.TrimSuffix(strings.TrimSpace(rx_str), ">"))
}

// check D-Lev software crc
func installed_crc_chk() (bool) {
	sp := sp_open()
	rx_str := sp_tx_rx(sp, "crc ", false)
	sp.Close()
	rx_str = decruft_hcl(rx_str)
	if !str_is_hex(rx_str) { error_exit(fmt.Sprint("Trouble getting CRC: ", rx_str)) }
	if rx_str == CRC { return true }
	return false
}

// return installed software version
func installed_ver() (string) {
	sp := sp_open()
	rx_str := sp_tx_rx(sp, "ver ", false)
	sp.Close()
	rx_str = decruft_hcl(rx_str)
	if !str_is_hex(rx_str) { error_exit(fmt.Sprint("Trouble getting version: ", rx_str)) }
	return rx_str
}

// check hive error reg
func hive_errors_chk() {
	err_str := func(err uint64) string {
		str := ""
		thd := 0
		for err != 0 {
			if err & 0x1 != 0 { str = " " + strconv.Itoa(thd) + str }
			thd++
			err >>= 1
		}
		return str
	}
	//
	sp := sp_open()
	rx_str := sp_tx_rx(sp, REG_ERROR + " rr ", false)
	sp.Close()
	rx_str = decruft_hcl(rx_str)
	rx_uint, err := strconv.ParseUint(rx_str, 16, 32)
	if err != nil { error_exit(fmt.Sprint("Trouble getting processor errors: ", rx_str)) }
	if rx_uint == 0 { fmt.Println("> no processor errors")
	} else { 
		irq_err := (rx_uint >> 24) & 0xff
		opc_err := (rx_uint >> 16) & 0xff
		psh_err := (rx_uint >> 8) & 0xff
		pop_err := rx_uint & 0xff
		fmt.Println("> -ERRORS!-")
		if irq_err != 0 { fmt.Print("> IRQ:", err_str(irq_err), "\n") }
		if opc_err != 0 { fmt.Print("> OPC:", err_str(opc_err), "\n") }
		if psh_err != 0 { fmt.Print("> PSH:", err_str(psh_err), "\n") }
		if pop_err != 0 { fmt.Print("> POP:", err_str(pop_err), "\n") }
	}
}

// check axis frequency
func freq_chk(vol_f bool) (uint32) {
	reg := REG_PITCH
	if vol_f { reg = REG_VOLUME }
	sp := sp_open()
	rx_str := sp_tx_rx(sp, reg + " rr ", false)
	sp.Close()
	rx_str = decruft_hcl(rx_str)
	rx_uint, err := strconv.ParseUint(rx_str, 16, 32)
	if err != nil { error_exit(fmt.Sprint("Trouble getting frequency: ", rx_str)) }
	f_clk := 196666666.6666667
	nco_w := 32
	f_shl := 7
	return uint32(float64(rx_uint) * f_clk / float64(uint64(1) << (nco_w + f_shl)))
}

// return operational stats
func stats_cmd(p_hz, v_hz, h_er bool) {
	all_f := !p_hz && !v_hz && !h_er
	if p_hz || all_f { fmt.Println("> P_FIELD Hz:", freq_chk(false)) }
	if v_hz || all_f { fmt.Println("> V_FIELD Hz:", freq_chk(true)) }
	if h_er || all_f { hive_errors_chk() }
}

// return file crc
func file_crc_chk(file string) (bool){
	sw_strs := file_sw_strs(file)
	sw_len, err := strconv.ParseUint(sw_strs[0], 16, 32); err_chk(err)
	ver_idx := int(sw_len)/EE_RW_BYTES
	if len(sw_strs) < ver_idx { error_exit(fmt.Sprint("Trouble getting ", file, " CRC")) }
	crc := uint32(0xffffffff)  // init to -1
	for i:=0; i<ver_idx; i++ {
		data, err := strconv.ParseUint(sw_strs[i], 16, 32); err_chk(err)
		crc = crc_32(uint32(data) ^ crc)
	}
	if strconv.FormatUint(uint64(crc), 16) == CRC { return true }
	return false
}

// return file software version
func file_ver(file string) (string){
	sw_strs := file_sw_strs(file)
	sw_len, err := strconv.ParseUint(sw_strs[0], 16, 32); err_chk(err)
	ver_idx := int(sw_len)/EE_RW_BYTES
	if len(sw_strs) < ver_idx { error_exit(fmt.Sprint("Trouble getting ", file, " version")) }
	return sw_strs[ver_idx-1]
}

// get versions & check stuff
func ver_cmd(file string, crc_chk, pre_chk, pro_chk bool) (bool) {
	sw_ver := ""
	sw_crc := false
	sw_upd := false
	prn_str := "> ver"
	if file != "" { 
		prn_str += " file[" + file + "] software: "
		sw_ver = file_ver(file)
		if crc_chk { sw_crc = file_crc_chk(file) }
	} else {
		prn_str += " installed software: "
		sw_ver = installed_ver()
		if crc_chk { sw_crc = installed_crc_chk() }
	}
	prn_str += sw_ver
	sw_date := sw_date_lookup(sw_ver)
	prn_str += "\n> ver software date: " + sw_date
	sw_lib := sw_lib_lookup(sw_ver)
	prn_str += "\n> ver associated librarian: " + sw_lib
	switch sw_ver {
	case ver_tbl[0].sw :
		prn_str += "\n> software is CURRENT"
	default :
		sw_upd = true
		if file == "" { prn_str += "\n> Software may be OLD, you may want to UPDATE it!" }
	}
	if crc_chk {
		if sw_crc { 
			prn_str += "\n> software PASSED the CRC check"
		} else {  
			sw_upd = true
			prn_str += "\n> software FAILED the CRC check!"
			if file == "" { prn_str += "\n> you may need to RE-UPLOAD or UPDATE your software" }
		}
	}
	if pre_chk {
		switch sw_ver {
		case ver_tbl[0].sw, ver_tbl[1].sw :
			prn_str += "\n> presets should be OK"
		default :
			prn_str += "\n> presets cannot be UPDATED using this version of the librarian,\n"
			prn_str += "\n> you can REPLACE them, or contact Eric for further options"
		}
	}
	if pro_chk {
		switch sw_ver {
		case ver_tbl[0].sw, ver_tbl[1].sw, ver_tbl[2].sw :
			prn_str += "\n> profiles should be OK" 
		default :
			prn_str += "\n> profiles cannot be UPDATED using this version of the librarian,\n"
			prn_str += "\n> you can REPLACE them, or contact Eric for further options"
		}
	}
	fmt.Println(prn_str)
	return sw_upd
}

// list free serial ports / set port
func port_cmd(set_str string, list bool) {
	// flags & stuff
	_l := list
	_s := (set_str != "")
	if !_l && !_s { error_exit("Nothing to do") }
	prn_str := "> port"
	port_list := sp_list()
	ports := len(port_list)
	if ports == 0 { error_exit("No serial ports found") }
	port_str := cfg_get("port")
	port_idx := str_idx(port_list, port_str)
	if _s {
		set_idx, err := strconv.Atoi(set_str)
		if err != nil { error_exit("Bad port number: " + set_str) }
		if set_idx < 0 || set_idx >= ports { error_exit("Port number out of range: " + set_str) }
		port_idx = set_idx
		port_str = port_list[set_idx]
		cfg_set("port", port_str)
		prn_str += fmt.Sprint(" set[", set_idx, "] (", port_str, ")")
	}
	if _l {
		prn_str += " list"
		for p_idx, p_str := range port_list { 
			prn_str += fmt.Sprint("\n [", p_idx, "] ", p_str)
			if p_idx == port_idx { prn_str += " <= ACTIVE" }
		}
		if port_str == "" { prn_str += "\n> Active port is not assigned!" }
		if port_idx < 0 { prn_str += "\n> Active port " + port_str + " doesn't exist!" }
	}
	fmt.Println(prn_str)
}

// view knobs, slot, DLP file, slot in PRE|PRO|EEPROM file
func view_cmd(file string, pro, knobs bool, slot string, mark int) {
	mode := "pre"; if pro { mode = "pro" }
	// flags & stuff
	_k := knobs
	_s := (slot != "")
	_f := (file != "")
	var ui_strs []string
	prn_str := "> view "
	// useful cases:
	switch {
	case  _k && !_s && !_f:  // knobs
		knob_str := sp_rx_knobs_str()
		ui_strs = knob_ui_strs(knob_str)
		prn_str += "knobs"
	case !_k &&  _s && !_f:  // slot
		slot_int := slot_int_chk(slot, pro)
		slot_str := spi_rd_slot_str(slot_int, pro)
		ui_strs = pre_ui_strs(slot_str, pro)
		prn_str += mode + " slot[" + slot + "]"
	case !_k && !_s &&  _f:  // file
		file = file_ext_chk(file, ".dlp")
		file_str := file_read_str(file)
		ui_strs = pre_ui_strs(file_str, pro)
		prn_str += mode + " file[" + file + "]"
	case !_k &&  _s &&  _f:  // slot in file
		slot_str := ""
		slot_str, mode = file_slot_str(file, slot, mode)
		if mode == "pro" { pro = true }
		ui_strs = pre_ui_strs(slot_str, pro)
		prn_str += mode + " slot[" + slot + "] in file[" + file + "]"
	default:
		error_exit("Nothing to do")
	}
	fmt.Println(ui_prn_str(ui_strs, mark))
	fmt.Println(prn_str)
}

// optionally write, then read knob, I/O is display type int
func knob_wr_rd(knob, wr_int int, wr bool) (int) {
	if (knob < 0) || (knob >= KNOBS_TOTAL-1) { error_exit("Knob index out of range") }
	ptype, _, _, _, _ := pname_lookup(knob_pnames[knob])
	tx_str := strconv.Itoa(knob)
	if wr { tx_str += " " + strconv.Itoa(freq_pint(wr_int, ptype)) + " wk "
	} else { tx_str += " rk " }
	sp := sp_open()
	rx_str := sp_tx_rx(sp, tx_str, false)
	sp.Close()
	rx_str = decruft_hcl(rx_str)
	pint, err := strconv.ParseUint(rx_str, 16, 32); err_chk(err)
	return pint_freq(int(pint), ptype)
}

// read / write, one knob or all page knobs
func knob_cmd(pkv string, view bool) {
	// process pkv string
	params, page_name, knob_name, knob_idx, val := pkv_parse(pkv) 
	prn_str := "> knob " + page_name
	if params > 1 {
		prn_str += ":" + knob_name
		if knob_name != "all" {
			rd_int := knob_wr_rd(knob_idx, val, false)
			prn_str += fmt.Sprint("[", rd_int, "]")
		}
	}
	if params == 3 {
		if knob_name == "all" {  // write all knobs
			for i:=1; i<KNOBS; i++ {
				knob_wr_rd(knob_idx-i, val, true)
			}
			prn_str += fmt.Sprint("=>[", val, "]")
		} else {  // write one knob
			wr_int := knob_wr_rd(knob_idx, val, true)
			prn_str += fmt.Sprint("=>[", wr_int, "]")
		}
	}	
	if view {
		view_cmd("", false, true, "", knob_idx) 
	}
	fmt.Println(prn_str)
}

// diff DLP file(s) / slot(s) / knobs
func diff_cmd(file, file2 string, pro, knobs bool, slot, slot2 string) {
	mode := "pre"; if pro { mode = "pro" }  // mode
	// flags
	_f := (file != "")
	_s := (slot != "")
	_f2 := (file2 != "")
	_s2 := (slot2 != "")
	_k := knobs
	// useful cases:
	base := ""
	comp := ""
	switch {
	case  _f && !_s && !_f2 && !_s2 &&  _k: base = "f";  comp = "k"     // DLP vs. knobs
	case  _f &&  _s && !_f2 && !_s2 && !_k: base = "f";  comp = "s"     // DLP vs. slot
	case  _f && !_s &&  _f2 && !_s2 && !_k: base = "f";  comp = "f2"    // DLP vs. DLP
	case  _f && !_s &&  _f2 &&  _s2 && !_k: base = "f";  comp = "f2s2"  // DLP vs. PRE|PRO|EEPROM slot
	case  _f &&  _s && !_f2 && !_s2 &&  _k: base = "fs"; comp = "k"     // PRE|PRO|EEPROM slot vs. knobs
	case  _f &&  _s && !_f2 &&  _s2 && !_k: base = "fs"; comp = "s2"    // PRE|PRO|EEPROM slot vs. slot
	case  _f &&  _s &&  _f2 && !_s2 && !_k: base = "fs"; comp = "f2"    // PRE|PRO|EEPROM slot vs. DLP
	case  _f &&  _s &&  _f2 &&  _s2 && !_k: base = "fs"; comp = "f2s2"  // PRE|PRO|EEPROM slot vs. PRE|PRO|EEPROM slot
	case !_f &&  _s && !_f2 && !_s2 &&  _k: base = "s";  comp = "k"     // slot vs. knobs
	case !_f &&  _s && !_f2 &&  _s2 && !_k: base = "s";  comp = "s2"    // slot vs. slot
	case !_f &&  _s &&  _f2 &&  _s2 && !_k: base = "s";  comp = "f2s2"  // slot vs. PRE|PRO|EEPROM slot
	case !_f && !_s &&  _f2 && !_s2 &&  _k: base = "k";  comp = "f2"    // knobs vs. DLP
	case !_f && !_s && !_f2 &&  _s2 &&  _k: base = "k";  comp = "s2"    // knobs vs. slot
	case !_f && !_s &&  _f2 &&  _s2 &&  _k: base = "k";  comp = "f2s2"  // knobs vs. PRE|PRO|EEPROM slot
	default: error_exit("Nothing to do")
	}
	// giant helper function	
	get_str := func(key string) (string, string) {
		// cases:
		switch {
		case key == "k":
			return knobs_pre_str(sp_rx_knobs_str(), pro), "knobs"
		case key == "f":
			file = file_ext_chk(file, ".dlp")
			return file_read_str(file), "file[" + file + "]"
		case key == "f2":
			file2 = file_ext_chk(file2, ".dlp")
			return file_read_str(file2), "file[" + file2 + "]"
		case key == "s":
			slot_int := slot_int_chk(slot, pro)
			return spi_rd_slot_str(slot_int, pro), "slot[" + slot + "]"
		case key == "s2":
			slot2_int := slot_int_chk(slot2, pro)
			return spi_rd_slot_str(slot2_int, pro), "slot[" + slot2 + "]"
		case key == "fs":
			var file_strs []string
			file_strs, mode = file_dlp_strs(file, mode)
			if mode == "pro" { pro = true }
			slot_int := slot_int_chk(slot, pro)
			return file_strs[slot_int], "slot[" + slot + "] in file[" + file + "]"
		case key == "f2s2":
			var file_strs []string
			file_strs, mode = file_dlp_strs(file2, mode)
			if mode == "pro" { pro = true }
			slot_int := slot_int_chk(slot2, pro)
			return file_strs[slot_int], "slot[" + slot + "] in file[" + file2 + "]"
		default:
			error_exit("Internal error")  // can't happen
		}
		return "Internal error", ""  // bogus
	}
	base_str, base_prn_str := get_str(base)
	comp_str, comp_prn_str := get_str(comp)
	fmt.Println(diff_prn_str(diff_pres(base_str, comp_str, pro)))
	fmt.Println("> diff", mode, base_prn_str, "-VS-", comp_prn_str )
}

// match slots|DLP files in dir2|"slots" in PRE|PRO|EEPROM file w/ DLP files in dir, list
func match_cmd(dir, dir2, file string, pro, hdr, guess, slots bool) {
	dir_chk(dir)
	name_strs, data_strs := dir_read_strs(dir, ".dlp", true)
	mode := "pre"; if pro { mode = "pro" }
	if len(data_strs) == 0 { error_exit(fmt.Sprint("No ", mode, " files in dir:", dir)) }
	// flags & stuff
	_s := slots
	_f := (file != "")
	_d2 := (dir2 != "")
	prn_str := "> match "
	// useful cases:
	switch {
	case  _s && !_f && !_d2:  // slots
		slots_strs := spi_rd_slots_strs(pro)
		fmt.Print(slots_prn_str(comp_file_data(slots_strs, name_strs, data_strs, pro, guess), pro, hdr))
		prn_str += mode + " slots"
	case !_s &&  _f && !_d2: // file
		file_chk(file)
		var file_strs []string
		file_strs, mode = file_dlp_strs(file, mode)
		if mode == "pro" { pro = true }
		fmt.Print(slots_prn_str(comp_file_data(file_strs, name_strs, data_strs, pro, guess), pro, hdr))
		prn_str += mode + " slots in file[" + file + "]"
	case !_s && !_f &&  _d2: // dir2
		dir_chk(dir2)
		name2_strs, data2_strs := dir_read_strs(dir2, ".dlp", true)
		if len(data2_strs) == 0 { error_exit(fmt.Sprint("No ", mode, " files in dir", dir2)) }
		fmt.Print(files_prn_str(name2_strs, comp_file_data(data2_strs, name_strs, data_strs, pro, guess)))
		prn_str += mode + " files in dir[" + dir2 + "]"
	default: error_exit("Nothing to do")
	}
	prn_str += " to " + mode + " files in dir[" + dir + "]"
	if guess { prn_str += " -guess" }
	if hdr { prn_str += " -hdr" }
	fmt.Println(prn_str)
}

// download to file
func dump_cmd(file, slot string, knobs, pro, yes bool) {
	// flags
	_k := knobs
	_s := (slot != "")
	_f := (file != "")
	prn_str := "> dump "
	//
	mode := "pre"; if pro { mode = "pro" }
	// useful cases:
	switch {
	case  _k && !_s && _f:  // knobs to DLP file
		file = file_ext_chk(file, ".dlp")
		pints := sp_rx_knobs_pints(mode)
		if file_write_str(file, ints_to_hexs(pints, 4), yes) {
			prn_str += mode + " knobs to file[" + file + "]"
		}
	case !_k &&  _s && _f:  // slot to DLP file
		file = file_ext_chk(file, ".dlp")
		slot_int := slot_int_chk(slot, pro)
		if file_write_str(file, spi_rd_slot_str(slot_int, pro), yes) {
			prn_str += mode + " slot[" + slot + "] to file[" + file + "]"
		}
	case !_k && !_s && _f:  // bulk dump
		file_blank_chk(file)
		ext := filepath.Ext(file)
		switch ext {
			case ".pre", ".pro", ".spi", ".eeprom" : // these are OK
			case "" : error_exit(fmt.Sprint("Missing file extension"))
			default : error_exit(fmt.Sprint("Wrong file extension: ", ext))
		}
		addr, end := spi_bulk_addrs(ext)
		rx_str := spi_rd(addr, end - 1, true)
		if file_write_str(file, rx_str, yes) {
			prn_str += "to file[" + file + "]"
		}
	default: error_exit("Nothing to do")
	}
	fmt.Println(prn_str)
}	

// upload from file
func pump_cmd(file, slot string, knobs, pro bool) {
	// flags & stuff
	_k := knobs
	_s := (slot != "")
	_f := (file != "")
	prn_str := "> pump "
	//
	mode := "pre"; if pro { mode = "pro" }
	// useful cases:
	switch {
	case  _k && !_s && _f:  // DLP file to knobs
		file = file_ext_chk(file, ".dlp")
		file_str := file_read_str(file)
		pints := hexs_to_ints(file_str, 4)
		if len(pints) < SLOT_BYTES { error_exit("Bad file info") }
		sp_tx_knobs_pints(pints, mode)
		prn_str += mode + " file[" + file + "] to knobs"
	case !_k &&  _s && _f:  // DLP file to slot
		file = file_ext_chk(file, ".dlp")
		file_str := file_read_str(file)
		slot_int := slot_int_chk(slot, pro)
		addr := spi_slot_addr(slot_int, pro)
		spi_wr(addr, file_str, false)
		prn_str += mode + " file[" + file + "] to slot[" + slot + "]"
	case !_k && !_s && _f:  // bulk pump
		ext := filepath.Ext(file)
		switch ext {
			case ".pre", ".pro", ".spi", ".eeprom" : // these are OK
			case "" : error_exit(fmt.Sprint("Missing file extension"))
			default : error_exit(fmt.Sprint("Wrong file extension: ", ext))
		}
		file_str := file_read_str(file)
		addr, _ := spi_bulk_addrs(ext)
		spi_wr(addr, file_str, true)
		prn_str += "from file[" + file + "]"
		if ext == ".spi" || ext == ".eeprom" || ext == ".pro" { reset_cmd() }
	default: error_exit("Nothing to do")
	}
	fmt.Println(prn_str)
}

// *.bnk => *.dlps => slots | *.pre | *.pro
func bank_cmd(slot, bank_file, target_file string, pro bool) {
	slot_int := slot_int_chk(slot, pro)
	bank_strs := file_bank_strs(bank_file)
	slots := PRE_SLOTS; if pro { slots = PRO_SLOTS }
	mode := "pre"; if pro { mode = "pro" }
	// pump to slots
	if target_file == "" {
		for idx, file := range bank_strs {
			pump_cmd(file, strconv.Itoa(slot_int+idx), false, pro)
		}
		fmt.Print("> bank ", mode, " file[" + bank_file + "] to slot[", slot_int, "]=>[", slot_int+len(bank_strs)-1, "]", "\n")
	// copy to slots in PRE|PRO file
	} else {
		ext := filepath.Ext(target_file)
		switch ext {
			case ".pre" : slots = PRE_SLOTS; pro = false; mode = "pre"
			case ".pro" : slots = PRO_SLOTS; pro = true; mode = "pro"
			case "" : error_exit(fmt.Sprint("Missing file extension"))
			default : error_exit(fmt.Sprint("Wrong file extension: ", ext))
		}
		// get dlp strs
		var dlp_strs []string
		for _, file := range bank_strs {
			dlp_strs = append(dlp_strs, file_read_str(file))
		}
		rd_str := ""
		wr_str := ""
		if path_exists(target_file) { 
			rd_str = file_read_str(target_file)  // PRE|PRO exists, read str
		} else { 
			rd_str = strings.Repeat("0\n", slots * DLP_LINES)  // else make blank str
		}
		rd_strs := split_pre_pro_str(rd_str)  // split PRE|PRO file str
		// replace rd strs with dlp strs
		for idx, str := range rd_strs {
			dlp_idx := idx - slot_int
			if (dlp_idx >= 0) && (dlp_idx < len(dlp_strs)) {
				wr_str += dlp_strs[dlp_idx] + "\n"
			} else {
				wr_str += str + "\n"
			}
		}
		if file_write_str(target_file, wr_str, false) {
			fmt.Print("> bank ", mode, " file[" + bank_file + "] to slot[", slot_int, "]=>[", slot_int+len(dlp_strs)-1, "] in file[" + target_file + "]\n")
		}
	}
}

// split file containers into sub containers
func split_cmd(file string, yes bool) {
	ext := filepath.Ext(file)
	dir, base := filepath.Split(file)
	base = strings.TrimSuffix(base, ext)
	file_str := file_read_str(file)
	switch ext {
	case ".eeprom" :
		pre_str, pro_str, spi_str := split_eeprom_str(file_str)
		pre_file := base + ".pre"
		pro_file := base + ".pro"
		spi_file := base + ".spi"
		pre_path := filepath.Join(dir, pre_file)
		pro_path := filepath.Join(dir, pro_file)
		spi_path := filepath.Join(dir, spi_file)
		wr_f := false
		prn_str := ""
		if file_write_str(pre_path, pre_str, yes) { prn_str += " " + pre_file; wr_f = true }
		if file_write_str(pro_path, pro_str, yes) { prn_str += " " + pro_file; wr_f = true }
		if file_write_str(spi_path, spi_str, yes) { prn_str += " " + spi_file; wr_f = true }
		if wr_f { fmt.Print("> split file[", file, "] => ", prn_str, "\n") }
	case ".pre", ".pro" :
		files := 0
		dlp_strs := split_pre_pro_str(file_str)
		for file_num, dlp_str := range dlp_strs {
			dlp_name := fmt.Sprintf("%03d", file_num) + ".dlp"
			if ext == ".pro" { dlp_name = "pro_" + dlp_name }
			dlp_file := filepath.Join(dir, dlp_name)
			if file_write_str(dlp_file, dlp_str, yes) { files++ }
		}
		fmt.Print("> split file[", file, "] => ", files, " numbered DLP files\n" )
	case "" : error_exit(fmt.Sprint("Missing file extension"))
	default : error_exit(fmt.Sprint("Wrong file extension: ", ext))
	}
}

// join sub containers to container
func join_cmd(file string, yes bool) {
	file_blank_chk(file)
	ext := filepath.Ext(file)
	dir, base := filepath.Split(file)
	base = strings.TrimSuffix(base, ext)
	switch ext {
	case ".eeprom" :
		base_path := filepath.Join(dir, base)
		pre_path := base_path + ".pre"
		pre_str := file_read_str(pre_path)
		pro_path := base_path + ".pro"
		pro_str := file_read_str(pro_path)
		spi_path := base_path + ".spi"
		spi_str := file_read_str(spi_path)
		file_str := pre_str + "\n"
		file_str += pro_str + "\n"
		file_str += spi_str
		if file_write_str(file, file_str, yes) {
			fmt.Print("> join ", pre_path, " ", pro_path, " ", spi_path, " => file[", file, "]\n" )
		}
	case ".pre", ".pro" :
		file_str := ""
		files := PRE_SLOTS
		if ext == ".pro" { files = PRO_SLOTS }
		for file_num := 0; file_num < files; file_num++ {
			dlp_name := fmt.Sprintf("%03d", file_num) + ".dlp"
			if ext == ".pro" { dlp_name = "pro_" + dlp_name }
			dlp_path := filepath.Join(dir, dlp_name)
			dlp_str := file_read_str(dlp_path)
			file_str += dlp_str + "\n"
		}
		if file_write_str(file, file_str, yes) {
			fmt.Print("> join ", files, " numbered DLP files => file[", file, "]\n")
		}
	case "" : error_exit(fmt.Sprint("Missing file extension"))
	default : error_exit(fmt.Sprint("Wrong file extension: ", ext))
	}
}

// morph page knobs|slot|DLP|PRE|EEPROM file to knobs
func morph_cmd(file string, knobs bool, slot string, seed int64, pkv_strs []string) {
	rand.Seed(seed)
	prn_str := ""
	var pints []int
	var knob_idxs []int
	var vals []int
	// flags
	_k := knobs
	_s := (slot != "")
	_f := (file != "")
	_v := false
	// parse pkv flags
	for _, pkv_str := range pkv_strs {
		params, page_name, knob_name, knob_idx, val := pkv_parse(pkv_str) 
		if params < 3 { error_exit("Too few params: " + pkv_str) }
		if knob_name == "all" {
			for kidx:=knob_idx; kidx<knob_idx+KNOBS-1; kidx++ { 
				knob_idxs = append(knob_idxs, kidx)
				vals = append(vals, val)
			}
		} else {
			knob_idxs = append(knob_idxs, knob_idx)
			vals = append(vals, val)
		}
		prn_str += fmt.Sprint(" ", page_name, ":", knob_name, "[", val, "]")
		_v = true
	}
	// useful cases:
	switch {
	case  _k && !_s && !_f && _v:  // knobs
		pints = sp_rx_knobs_pints("pre")
		prn_str = " knobs" + prn_str
	case !_k &&  _s && !_f && _v:  // slot
		slot_int := slot_int_chk(slot, false)
		slot_str := spi_rd_slot_str(slot_int, false)
		pints = hexs_to_ints(slot_str, 4)
		prn_str = " slot[" + slot + "]" + prn_str
	case !_k && !_s &&  _f && _v:  // DLP file
		file = file_ext_chk(file, ".dlp")
		file_str := file_read_str(file)
		pints = hexs_to_ints(file_str, 4)
		prn_str = " file[" + file + "]" + prn_str
	case !_k &&  _s &&  _f && _v:  // slot in EEPROM|PRE file
		slot_str, _ := file_slot_str(file, slot, "pre")
		pints = hexs_to_ints(slot_str, 4)
		prn_str = " slot[" + slot + "] in file[" + file + "]" + prn_str
	default: 
		error_exit("Nothing to do")
	}
	prn_str = "> morph" + prn_str
	pints = morph_pints(pints_signed(pints, false), knob_idxs, vals)
	prn_str += fmt.Sprint(" init[", seed, "]")
	sp_tx_knobs_pints(pints, "pre")
	fmt.Println(prn_str)
}

/*
// process page:morph strings
func dev_cmd(pm_strs []string) {
	for _, pm_str := range pm_strs {
		p_str, m_str, _ok := (strings.Cut(pm_str, ":"))
		page_idx, page_name := fuzzy_page_lookup(p_str)
		if page_idx < 0 { error_exit("Bad page name") }
		morph, err := strconv.ParseFloat(m_str, 64); 
		if err != nil { error_exit("Bad morph value") }
		fmt.Println(page_name, morph, _ok)
	}
}
*/
