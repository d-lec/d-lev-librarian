package main

/*
 * d-lib support functions
*/

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"math/rand"
)


// Calculate Levenshtein distance between two strings
// snagged from https://www.golangprograms.com/golang-program-for-implementation-of-levenshtein-distance.html
func levenshtein(str1, str2 string) int {
	// helper func, return min of 3
	min := func(a, b, c int) int {
		if a < b {
			if a < c { return a }
		} else {
			if b < c { return b }
		}
		return c
	}
	s1_len := len(str1)
	s2_len := len(str2)
	col := make([]int, len(str1)+1)
	for y := 1; y <= s1_len; y++ {
		col[y] = y
	}
	for x := 1; x <= s2_len; x++ {
		col[0] = x
		lastkey := x - 1
		for y := 1; y <= s1_len; y++ {
			oldkey := col[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			col[y] = min(col[y]+1, col[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	return col[s1_len]
}

// librarian commands
var lib_cmds = []string {  
	"menu",
	"help",
	"ports",
	"view",
	"match",
	"diff",
	"bank",
	"dump",
	"pump",
	"split",
	"join",
	"morph",
	"batch",
	"knob",
	"ver",
	"hcl",
	"loop",
	"acal",
	"reset",
	"stats",
}

// return commmand hint
func cmd_hint (cmd_in string) string {
	cmd_out := lib_cmds[0]
	ld_lowest := levenshtein(cmd_in, lib_cmds[0])
	for _, lib_cmd := range lib_cmds {
		if strings.HasPrefix(lib_cmd, cmd_in) { return lib_cmd }
		ld := levenshtein(cmd_in, lib_cmd)
		if ld < ld_lowest { 
			ld_lowest = ld
			cmd_out = lib_cmd 
		}
	}
	return cmd_out
}

// return first word from user input line
func user_word() (string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fields := strings.Fields(scanner.Text())
	if len(fields) != 0 { return fields[0] }
	return ""
}

// pause and ask user yes | no | quit question
func user_prompt(prompt string, yes, space bool) (bool) {
	if yes { return true }
	if space { fmt.Println() }
	fmt.Print("> ", prompt, " <y|N|q>: ")
	input := user_word()
	switch input {
		case "q": quit_exit()
		case "y": return true 
		default: fmt.Println("> -CANCELED-")

	}
	return false	
}

// pause and ask user text | quit question
func user_input(prompt string, space bool) (string) {
	if space { fmt.Println() }
	fmt.Print("> ", prompt, " (q=quit): ")
	input := user_word()
	if input == "q" { quit_exit() }
	return input
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
	fmt.Println("> Issued processor reset")
}

// do ACAL
func acal_cmd() {
	sp := sp_open()
	sp_tx_rx(sp, "acal ", false)
	sp.Close()
	fmt.Println(" Issued ACAL")
}

// do HCL command
func hcl_cmd() {
	if len(os.Args) < 3 { error_exit("Command line is blank") }
	wr_str := ""
	for _, cmd := range os.Args[2:] {
		wr_str += cmd + " "
	}
	sp := sp_open()
	rx_str := sp_tx_rx(sp, wr_str, false)
	sp.Close()
	fmt.Print(strings.TrimSpace(rx_str))
	fmt.Println(" Issued hcl command:", wr_str)
}

// do loop command
func loop_cmd() {
	if len(os.Args) < 3 { error_exit("Loop text is blank") }
	wr_str := ""
	for _, arg := range os.Args[2:] {
		wr_str += arg + " "
	}
	wr_str = strings.TrimSpace(wr_str)
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
	if err != nil { error_exit(fmt.Sprint("Trouble getting errors: ", rx_str)) }
	if rx_uint == 0 { fmt.Println("> No errors.")
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
func ver_cmd(file string, pre_chk, pro_chk bool) (bool) {
	sw_ver := ""
	sw_crc := false
	sw_upd := false
	if file != "" { 
		fmt.Println("> File:", file)
		sw_ver = file_ver(file)
		sw_crc = file_crc_chk(file)
	} else {
		sw_ver = installed_ver()
		sw_crc = installed_crc_chk()
	}
	sw_date := sw_date_lookup(sw_ver)
	sw_lib := sw_lib_lookup(sw_ver)
	fmt.Println("> Software version:", sw_ver)
	fmt.Println("> Software date & associated librarian:", sw_date, "version", sw_lib)
	switch sw_ver {
	case ver_tbl[0].sw :
		fmt.Println("> Software is CURRENT.") 
	default :
		sw_upd = true
		if file == "" { fmt.Println("> Software may be OLD, you may want to UPDATE it.") }
	}
	if sw_crc { 
		fmt.Println("> Software PASSED the CRC check.") 
	} else {  
		sw_upd = true
		fmt.Println("> Software FAILED the CRC check!") 
		if file == "" { fmt.Println("> You may need to RE-UPLOAD or UPDATE your software.") }
	}
	if pre_chk {
		switch sw_ver {
		case ver_tbl[0].sw, ver_tbl[1].sw :
			fmt.Println("> Presets should be OK.") 
		default :
			fmt.Println("> Presets cannot be UPDATED using this version of the librarian,") 
			fmt.Println("> You can REPLACE them, or contact Eric for further options.")
		}
	}
	if pro_chk {
		switch sw_ver {
		case ver_tbl[0].sw :
			fmt.Println("> Profiles should be OK.") 
		case ver_tbl[1].sw :
			fmt.Println("> Profiles may need manual adjustment.") 
		default :
			fmt.Println("> Profiles cannot be UPDATED using this version of the librarian,") 
			fmt.Println("> You can REPLACE them, or contact Eric for further options.")
		}
	}
	return sw_upd
}

// list free serial ports / set port
func ports_cmd(port_new string) {
	port := cfg_get("port")
	port_list := sp_list()
	port_idx := str_exists(port_list, port)
	if len(port_list) == 0 {
		fmt.Println("> No serial ports found!")
		return
	} else {
		fmt.Println("> Available serial ports:")
		for p_num, p_str := range port_list { fmt.Printf(" [%v] %v\n", p_num, p_str) }
	}
	switch {
	case port_new != "":
		port_num, err := strconv.Atoi(port_new)
		if err != nil || port_num < 0 || port_num >= len(port_list) { error_exit(fmt.Sprint("Bad port number: ", port_new)) }
		port = port_list[port_num]
		cfg_set("port", port)
		fmt.Print("> Set port to: [", port_num, "] ", port, "\n")
	case len(os.Args) > 2:
		error_exit("Use the -p flag to set the port")
	case port == "":
		fmt.Println("> Current port is not assigned!")
	case port_idx < 0:
		fmt.Println("> Current port:", port, "doesn't exist!")
	default: 
		fmt.Print("> Current port: [", port_idx, "] ", port, "\n")
	}
}

// view knobs, slot, DLP file, slot in PRE|PRO|EEPROM file
func view_cmd(file string, pro, knobs bool, slot string, mark int) {
	mode := "pre"; if pro { mode = "pro" }
	// flags
	_k := knobs
	_s := (slot != "")
	_f := (file != "")
	// useful cases:
	switch {
	case  _k && !_s && !_f:  // knobs
		knob_str := sp_rx_knobs_str()
		fmt.Println(ui_prn_str(knob_ui_strs(knob_str), mark))
		if mark < 0 { fmt.Println("> knobs") }
	case !_k &&  _s && !_f:  // slot
		slot_int := slot_int_chk(slot, pro)
		slot_str := spi_rd_slot_str(slot_int, pro)
		fmt.Println(ui_prn_str(pre_ui_strs(slot_str, pro), mark))
		fmt.Println(">", mode, "slot", slot)
	case !_k && !_s &&  _f:  // file
		file = file_ext_chk(file, ".dlp")
		file_str := file_read_str(file)
		fmt.Println(ui_prn_str(pre_ui_strs(file_str, pro), mark))
		fmt.Println(">", mode, "file", file)
	case !_k &&  _s &&  _f:  // slot in file
		var file_strs []string
		file_strs, mode = file_dlp_strs(file, mode)
		if mode == "pro" { pro = true }
		slot_int := slot_int_chk(slot, pro)
		fmt.Println(ui_prn_str(pre_ui_strs(file_strs[slot_int], pro), mark))
		fmt.Println(">", mode, "slot", slot, "in file", file )	
	default:
		error_exit("Nothing to do")
	}
}

// read, set, offset, min one knob or all knobs
func knob_cmd(knob, ofs, set string, min, view bool) {
	// flags & stuff
	_kp := false  // knob page=ok 
	_ki := false  // knob idx=int
	_ka := false  // knob idx=all
	page_name := ""
	page_idx := 0
	knob_int := 0
	prn_str := ""
	// process knob string
	if knob != "" {  
		kp_str, kia_str, _kia := (strings.Cut(strings.TrimSpace(knob), ":"))
		page_name, page_idx = page_lookup(kp_str)
		if page_idx < 0 { error_exit("Bad page name") }
		_kp = true
		if _kia {  // knob idx string not empty
			kia_str = strings.ToLower(strings.TrimSpace(kia_str))
			if kia_str == "all" {
				_ka = true 
			} else {
				var err error
				knob_int, err = strconv.Atoi(kia_str); 
				if err != nil { error_exit("Bad knob index") }
				if knob_int < 0 || knob_int >= KNOBS-1 { error_exit("Knob index out of range") }
				_ki = true
			}
		}
	}
	// base index at this point
	knob_idx := page_idx * KNOBS
	// set UI LCD to page
	set_ui_page := func() {
		sp := sp_open()
		sp_tx_rx(sp, strconv.Itoa(PAGE_SEL_KNOB) + " " + strconv.Itoa(page_idx) + " wk ", false)
		sp.Close()
	}
	// useful cases:
	switch {
	case _kp && !_ki && !_ka:  // page:(w/ no idx)
		prn_str = fmt.Sprint("> ", page_name)
		knob_idx += PAGE_SEL_KNOB  // bracket page name
	case _kp && _ka:  // page:all
		prn_str = fmt.Sprint("> ", page_name, ":all")
		if min {
			sp := sp_open()
			for i:=0; i<KNOBS-1; i++ {
				sp_tx_rx(sp, strconv.Itoa(knob_idx+i) + " 0 wk ", false)
			}
			sp.Close()
			prn_str += fmt.Sprint("=>[min]")
		}
		knob_idx += PAGE_SEL_KNOB  // bracket page name
	case _kp && _ki:  // page:idx
		knob_idx += knob_int
		ptype, plabel, _, _ := pname_lookup(knob_pnames[knob_idx])
		sp := sp_open()
		rx_str := sp_tx_rx(sp, strconv.Itoa(knob_idx) + " rk ", false)
		sp.Close()
		rd_uint32, _ := strconv.ParseUint(decruft_hcl(rx_str), 16, 32)
		prn_str = fmt.Sprint("> ", page_name, ":", knob_int, " ", strings.TrimSpace(plabel), "[", strings.TrimSpace(pint_disp(int(rd_uint32), ptype)), "]")
		if min || ofs != "" || set != "" {
			rw_pint := pint_freq(int(rd_uint32), ptype)
			if set != "" { 
				set_int, err := strconv.Atoi(set)
				if err != nil { error_exit("Bad set value") }
				rw_pint = set_int
			}
			if ofs != "" { 
				ofs_int, err := strconv.Atoi(ofs)
				if err != nil { error_exit("Bad offset value") }
				rw_pint += ofs_int
			}
			if min { rw_pint = 0 }
			rw_pint = freq_pint(rw_pint, ptype)
			sp := sp_open()
			sp_tx_rx(sp, strconv.Itoa(knob_idx) + " " + strconv.Itoa(rw_pint) + " wk ", false)
			sp.Close()
			prn_str += fmt.Sprint("=>[", strings.TrimSpace(pint_disp(rw_pint, ptype)), "]")
		}
	case view:
		sp := sp_open()  // get current page
		rx_str := sp_tx_rx(sp, strconv.Itoa(PAGE_SEL_KNOB) + " rk ", false)
		sp.Close()
		rd_uint32, _ := strconv.ParseUint(decruft_hcl(rx_str), 16, 32)
		knob_idx = int(rd_uint32) * KNOBS + PAGE_SEL_KNOB
		prn_str = "> " + strings.TrimSpace(page_names[knob_idx / 8])
	default:
		error_exit("Nothing to do")
	}
	if view {
		view_cmd("", false, true, "", knob_idx) 
		set_ui_page()
	}
	fmt.Println(prn_str)
}


// diff DLP file(s) / slot(s) / knobs
func diff_cmd(file, file2 string, pro, knobs bool, slot, slot2 string) {
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
		mode := "pre"; if pro { mode = "pro" }  // mode
		// cases:
		switch {
		case key == "k":
			return knob_pre_str(sp_rx_knobs_str(), pro), "knobs"
		case key == "f":
			file = file_ext_chk(file, ".dlp")
			return file_read_str(file), mode + " file " + file
		case key == "f2":
			file2 = file_ext_chk(file2, ".dlp")
			return file_read_str(file2), mode + " file " + file2
		case key == "s":
			slot_int := slot_int_chk(slot, pro)
			return spi_rd_slot_str(slot_int, pro), mode + " slot " + slot
		case key == "s2":
			slot2_int := slot_int_chk(slot2, pro)
			return spi_rd_slot_str(slot2_int, pro), mode + " slot " + slot2
		case key == "fs":
			var file_strs []string
			file_strs, mode = file_dlp_strs(file, mode)
			if mode == "pro" { pro = true }
			slot_int := slot_int_chk(slot, pro)
			return file_strs[slot_int], mode + " slot " + slot + " in file " + file
		case key == "f2s2":
			var file_strs []string
			file_strs, mode = file_dlp_strs(file2, mode)
			if mode == "pro" { pro = true }
			slot_int := slot_int_chk(slot2, pro)
			return file_strs[slot_int], mode + " slot " + slot2 + " in file " + file2
		default:
			error_exit("Internal error")  // can't happen
		}
		return "Internal error", ""  // bogus
	}
	base_str, base_prn_str := get_str(base)
	comp_str, comp_prn_str := get_str(comp)
	fmt.Println(diff_prn_str(diff_pres(base_str, comp_str, pro)))
	fmt.Println(">", base_prn_str, "-VS-", comp_prn_str )
}

// match slots | DLP files in dir2 | "slots" in PRE|PRO|EEPROM file w/ DLP files in dir, list
func match_cmd(dir, dir2, file string, pro, hdr, guess, slots bool) {
	dir_chk(dir)
	name_strs, data_strs := dir_read_strs(dir, ".dlp")
	mode := "pre"; if pro { mode = "pro" }
	if len(data_strs) == 0 { error_exit(fmt.Sprint("No ", mode, " files in directory ", dir)) }
	// flags
	_s := slots
	_f := (file != "")
	_d2 := (dir2 != "")
	// useful cases:
	switch {
	case  _s && !_f && !_d2:  // slots
		slots_strs := spi_rd_slots_strs(pro)
		fmt.Print(slots_prn_str(comp_file_data(slots_strs, name_strs, data_strs, pro, guess), pro, hdr))
		fmt.Print("> ", mode, " slots")
	case !_s &&  _f && !_d2: // file
		file_chk(file)
		var file_strs []string
		file_strs, mode = file_dlp_strs(file, mode)
		if mode == "pro" { pro = true }
		fmt.Print(slots_prn_str(comp_file_data(file_strs, name_strs, data_strs, pro, guess), pro, hdr))
		fmt.Print("> ", mode, " slots in file ", file)
	case !_s && !_f &&  _d2: // dir2
		dir_chk(dir2)
		name2_strs, data2_strs := dir_read_strs(dir2, ".dlp")
		if len(data2_strs) == 0 { error_exit(fmt.Sprint("No ", mode, " files in ", dir2)) }
		fmt.Print(files_prn_str(name2_strs, comp_file_data(data2_strs, name_strs, data_strs, pro, guess)))
		fmt.Print("> ", mode, " files in ", dir2)
	default: error_exit("Nothing to do")
	}
	fmt.Println(" -MATCHED- to", mode, "files in", dir)
}

// download to file
func dump_cmd(file, slot string, knobs, pro, yes bool) {
	// flags
	_k := knobs
	_s := (slot != "")
	_f := (file != "")
	//
	mode := "pre"; if pro { mode = "pro" }
	// useful cases:
	switch {
	case  _k && !_s && _f:  // knobs to DLP file
		file = file_ext_chk(file, ".dlp")
		pints := sp_rx_knobs_pints(mode)
		if file_write_str(file, ints_to_hexs(pints, 4), yes) {
			fmt.Println("> downloaded", mode, "knobs =>", mode, "file", file) 
		}
	case !_k &&  _s && _f:  // slot to DLP file
		file = file_ext_chk(file, ".dlp")
		slot_int := slot_int_chk(slot, pro)
		if file_write_str(file, spi_rd_slot_str(slot_int, pro), yes) {
			fmt.Println("> downloaded", mode, "slot", slot_int, "=>", mode, "file", file) 
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
			fmt.Println("> dumped to", file) 
		}
	default: error_exit("Nothing to do")
	}
}	

// upload from file
func pump_cmd(file, slot string, knobs, pro bool) {
	// flags
	_k := knobs
	_s := (slot != "")
	_f := (file != "")
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
		fmt.Println("> uploaded", mode, "file", file, "=>", mode, "knobs") 
	case !_k &&  _s && _f:  // DLP file to slot
		file = file_ext_chk(file, ".dlp")
		file_str := file_read_str(file)
		slot_int := slot_int_chk(slot, pro)
		addr := spi_slot_addr(slot_int, pro)
		spi_wr(addr, file_str, false)
		fmt.Println("> uploaded", mode, "file", file, "=>", mode, "slot", slot_int) 
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
		fmt.Println("> pumped from", file)
		if ext == ".spi" || ext == ".eeprom" || ext == ".pro" { reset_cmd() }
	default: error_exit("Nothing to do")
	}
}

// *.bnk => *.dlps => slots
func bank_cmd(slot, file string, pro bool) {
	slot_int := slot_int_chk(slot, pro)
	file = file_ext_chk(file, ".bnk")
	bnk_str := file_read_str(file)

	bnk_str = strip_cmnt(bnk_str)  // strip C & CPP style comments

	bnk_split := strings.Split(bnk_str, "\n")
	dir, _ := filepath.Split(file)
	for _, line := range bnk_split {
		line := strings.TrimSpace(line);

		// if !strings.HasPrefix(line, "//") {  // skip commented lines
		if line != "" {  // skip blank lines

			file := filepath.Join(dir, line)
			pump_cmd(file, strconv.Itoa(slot_int), false, pro)
			slot_int++
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
		if wr_f { fmt.Print("> split ", file, " =>", prn_str, "\n") }
	case ".pre", ".pro" :
		files := 0
		dlp_strs := split_pre_pro_str(file_str)
		for file_num, dlp_str := range dlp_strs {
			dlp_name := fmt.Sprintf("%03d", file_num) + ".dlp"
			if ext == ".pro" { dlp_name = "pro_" + dlp_name }
			dlp_file := filepath.Join(dir, dlp_name)
			if file_write_str(dlp_file, dlp_str, yes) { files++ }
		}
		fmt.Println("> split", file, "=>", files, "numbered DLP files" )
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
			fmt.Println("> joined", pre_path, pro_path, spi_path, "=>", file )
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
			fmt.Println("> joined", files, "numbered DLP files", "=>", file)
		}
	case "" : error_exit(fmt.Sprint("Missing file extension"))
	default : error_exit(fmt.Sprint("Wrong file extension: ", ext))
	}
}

func morph_cmd(file string, knobs bool, slot string, seed int64, mo, mn, me, mf, mr, ms float64) {
	rand.Seed(seed)
	prn_str := ""
	var pints []int
	// flags
	_k := knobs
	_s := (slot != "")
	_f := (file != "")
	_m := (mo != 0) || (mn != 0) || (me != 0) || (mf != 0) || (mr != 0)
	// useful cases:
	switch {
	case  _k && !_s && !_f && _m:  // knobs
		pints = pints_signed(sp_rx_knobs_pints("pre"), false)
		prn_str = fmt.Sprint("> morphed knobs")
	case !_k &&  _s && !_f && _m:  // slot
		slot_int := slot_int_chk(slot, false)
		slot_str := spi_rd_slot_str(slot_int, false)
		pints = pints_signed(hexs_to_ints(slot_str, 4), false)
		prn_str = fmt.Sprint("> morphed slot ", slot_int)
	case !_k && !_s &&  _f && _m:  // file
		file = file_ext_chk(file, ".dlp")
		file_str := file_read_str(file)
		pints = pints_signed(hexs_to_ints(file_str, 4), false)
		prn_str = fmt.Sprint("> morphed file ", file)
	default: 
		error_exit("Nothing to do")
	}
	prn_str += fmt.Sprint(" (seed:", seed, ")")
	pints = morph_pints(pints, mo, mn, me, mf, mr, ms)
	sp_tx_knobs_pints(pints, "pre")
	fmt.Println(prn_str)
}

// do a bunch of update stuff via interactive menu
func menu_cmd(dir_work string) {
	path_exe, err := os.Executable(); err_chk(err)
	dir_exe := filepath.Dir(path_exe)
	dir_all := filepath.Join(dir_exe, PRESETS_DIR)
//	dir_work_pro := filepath.Join(dir_work, PRO_DIR)
	//
	file_spi := filepath.Join(PRESETS_DIR, ver_tbl[0].date + ".spi")
	file_factory := filepath.Join(PRESETS_DIR, ver_tbl[0].date + ".eeprom")
	file_bank := filepath.Join(PRESETS_DIR, ver_tbl[0].date + ".bnk")
	file_bank_new := filepath.Join(PRESETS_DIR, ver_tbl[0].date + "_new.bnk")
	//
	path_spi := filepath.Join(dir_exe, file_spi)
	path_factory := filepath.Join(dir_exe, file_factory)
	path_bank := filepath.Join(dir_exe, file_bank)
	path_bank_new := filepath.Join(dir_exe, file_bank_new)	
	//
	path_pre_dl := filepath.Join(dir_work, "download.pre")
	path_pre_ul := filepath.Join(dir_work, "upload.pre")
/*
	path_pro_dl := filepath.Join(dir_work_pro, "download.pro")
	path_pro_ul := filepath.Join(dir_work_pro, "upload.pro")
*/
	lib_ver := ver_tbl[0].lib
	first_f := true
	for {
		if !first_f {
			user_input("Press <ENTER> to return to the MENU", true)
		}
		first_f = false
		fmt.Println("\n")  // 2 blank lines
		fmt.Println(" ---------------------------------")
		fmt.Println(" | D-LEV LIBRARIAN - VERSION", lib_ver, " |")
		fmt.Println(" | SOFTWARE & PRESET UPDATE MENU |")
		fmt.Println(" ---------------------------------")
		fmt.Println("  0. README!")
		fmt.Println("  1. Serial port setup & check.")
		fmt.Println("  2. Backup your system to an EEPROM file.")
		fmt.Println("  3. Check & update the D-Lev software.")
		fmt.Println("  4. Download all D-Lev presets to the", dir_work, "directory.")
		fmt.Println("  5. Update all presets in the", dir_work, "directory.")
		fmt.Println("  6. Upload all D-Lev presets from the", dir_work, "directory.")
		fmt.Println("  7. Upload the latest new presets.")
		fmt.Println("  8. Convert all presets in the", dir_work, "directory to MONO.")
		fmt.Println("  9. Overwrite all D-Lev preset slots with presets from the", PRESETS_DIR, "directory.")
		fmt.Println(" 10. Factory Reset: Overwrite EVERYTHING with the latest factory EEPROM file.")
		fmt.Println(" 11. Pump some updated presets to their default slots.")
		fmt.Println("  h. List of commands.")
		fmt.Println(" hv. List of commands with notes and typical examples of their use.")
		menu_sel := user_input("Please select a MENU option", true)
		switch menu_sel {
		case "0" :
			fmt.Println()
			fmt.Print(menu_readme_str)
		case "1" :
			fmt.Println()
			ports_cmd("")
			if user_prompt("Do you want to CHANGE current port?", false, true) {
				port_new := user_input("Please input PORT number", false)
				if port_new == "" {
					fmt.Println("> -CANCEL-")
				} else {
					fmt.Println()
					ports_cmd(port_new)
				}
			}
			if user_prompt("Do you want to TEST the port (do a CTRL-C if it hangs)?", false, true) {
				installed_ver()
				fmt.Println("> Port seems to be OK!")
			}
		case "2" :
			file_backup := date_hm() + "_back.eeprom"
			if user_prompt("Do you want to BACKUP your ENTIRE D-Lev to the FILE: " + file_backup + "?", false, true) {
				dump_cmd(file_backup, "", false, false, false)
			}
		case "3" :
			fmt.Println()
			sw_upd := ver_cmd("", true, true)
			if (sw_upd)  {
				if user_prompt("Do you want to UPDATE your D-Lev SOFTWARE with the FILE: "+ file_spi + "?", false, true) {
					pump_cmd(path_spi, "", false, false)
					fmt.Println()
					ver_cmd("", false, false)
				}
			}
		case "4" :
			if user_prompt("Do you want to DOWNLOAD your D-Lev presets to " + dir_work +"?", false, true) {
				dump_cmd(path_pre_dl, "", false, false, true)
				split_cmd(path_pre_dl, true)
			}
/*
			if user_prompt("Do you want to DOWNLOAD your D-Lev profiles to " + dir_work_pro +"?", false, true) {
				dump_cmd(path_pro_dl, "", false, true, true)
				split_cmd(path_pro_dl, true)
			}
*/
		case "5" :
			if user_prompt("Do you want to UPDATE the presets in " + dir_work + "?", false, true) {
				process_dlps(dir_work, dir_work, false, false, true, false, true, false)
			}
/*
			if user_prompt("Do you want to UPDATE the profiles in " + dir_work_pro + "?", false, true) {
				process_dlps(dir_work_pro, dir_work_pro, true, false, true, false, true, false)
			}
*/
		case "6" :
			if user_prompt("Do you want to UPLOAD the presets in "+ dir_work + "?", false, true) {
				join_cmd(path_pre_ul, true)
				pump_cmd(path_pre_ul, "", false, false)
			}
/*
			if user_prompt("Do you want to UPLOAD the profiles in "+ dir_work_pro + "?", false, true) {
				join_cmd(path_pro_ul, true)
				pump_cmd(path_pro_ul, "", false, true)
			}
*/
		case "7" :
			fmt.Println("> Here is a LIST of the latest NEW presets:")
			fmt.Println()
			file_ext_chk(path_bank_new, ".bnk")
			file_str := file_read_str(path_bank_new)
			fmt.Println(file_str)
			if user_prompt("Do you want to EXAMINE your current D-Lev presets?", false, true) {
				match_cmd(dir_all, "", "", false, false, true, true) 
			}
			if user_prompt("Do you want to UPLOAD the latest NEW presets?", false, true) {
				slot := user_input("What SLOT do you want to START the upload?", true)
				if slot != "" {
					bank_cmd(slot, path_bank_new, false)
				}
				if user_prompt("Do you want to EXAMINE your current D-Lev presets?", false, true) {
					match_cmd(dir_all, "", "", false, false, true, true) 
				}
			}
		case "8" :
			if user_prompt("Do you want to CONVERT all of the presets in " + dir_work + " to MONO?", false, true) {
				process_dlps(dir_work, dir_work, false, true, false, false, true, false)
			}
		case "9" :
			if user_prompt("Do you want to OVERWRITE all D-Lev preset slots with presets in " + PRESETS_DIR + "?", false, true) {
				bank_cmd("0", path_bank, false)
				if user_prompt("Do you want to EXAMINE your current D-Lev presets?", false, true) {
					match_cmd(dir_all, "", "", false, false, true, true) 
				}
			}
		case "10" :
			if user_prompt("Do you want to OVERWRITE ABSOLUTELY EVERYTHING in your D-Lev with the latest EEPROM file?", false, true) {
				pump_cmd(path_factory, "", false, false)
				fmt.Println()
				ver_cmd("", false, false)
			}

		case "11" :
			if user_prompt("Do you want to EXAMINE your current D-Lev presets?", false, true) {
				match_cmd(dir_all, "", "", false, false, true, true) 
			}
			if user_prompt("Do you want to PUMP some modified presets to their default slots?", false, true) {
				dlps_pump(dir_all, false)
			}
		case "h" :
			fmt.Println()
			help_cmd(false)
		case "hv" :
			fmt.Println()
			help_cmd(true)
		case "robs" :
			if user_prompt("Do you want to apply Rob Schwimmer Pitch Preview to the presets in " + dir_work + "?", false, true) {
				process_dlps(dir_work, dir_work, false, false, false, true, true, false)
			}
		case "" :
			// do nothing
		default:
			fmt.Println("> Invalid menu selection!")
		}
	}
}

