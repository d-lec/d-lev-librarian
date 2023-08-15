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
	fmt.Print(help_str) 
	if verbose_f { fmt.Print(help_verbose_str) }  // print verbose help
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
func ver_cmd(file string, pre_chk bool) (bool) {
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
	fmt.Println("> Software version:", sw_ver)
	fmt.Println("> Software date:", sw_date_lookup(sw_ver))
	switch sw_ver {
		case ver_tbl[0].sw, ver_tbl[1].sw :
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
			case ver_tbl[0].sw, ver_tbl[1].sw, ver_tbl[2].sw :
				fmt.Println("> Presets should be OK.") 
			case ver_tbl[3].sw :
				fmt.Println("> Presets can be UPDATED with this version of the librarian.") 
			default :
				fmt.Println("> Presets cannot be UPDATED using this version of the librarian,") 
				fmt.Println("> You can REPLACE your presets, or contact Eric for further options.")
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
	} else {
		fmt.Println("> Available serial ports:")
		for p_num, p_str := range port_list { fmt.Printf(" [%v] %v\n", p_num, p_str) }
	}
	if port_new != "" { 
		port_num, err := strconv.Atoi(port_new)
		if err != nil || port_num < 0 || port_num >= len(port_list) { error_exit(fmt.Sprint("Bad port number: ", port_new)) }
		port = port_list[port_num]
		cfg_set("port", port)
		fmt.Print("> Set port to: [", port_num, "] ", port, "\n")
	} else if len(os.Args) > 2 { 
		error_exit("Use the -p flag to set the port")
	} else if port == "" {
		fmt.Println("> Current port is not assigned!")
	} else if port_idx < 0 {
		fmt.Println("> Current port:", port, "doesn't exist!")
	} else {
		fmt.Print("> Current port: [", port_idx, "] ", port, "\n")
	}
}

// view knobs, slot, DLP file, slot in PRE|PRO|EEPROM file
func view_cmd(file string, pro, knobs bool, slot string, mark int) {
	mode := "pre"; if pro { mode = "pro" }
	if knobs {  // view current knobs
		knob_str := sp_rx_knobs_str()
		fmt.Println(ui_prn_str(knob_ui_strs(knob_str), mark))
		if mark < 0 { fmt.Println("> knobs") }
	} else if file != "" || slot != "" {
		if slot == "" {  // view a *.dlp file
			file = file_ext_chk(file, ".dlp")
			file_str := file_read_str(file)
			fmt.Println(ui_prn_str(pre_ui_strs(file_str, pro), mark))
			fmt.Println(">", mode, "file", file)
		} else if file == "" {  // view a slot
			slot_int := slot_int_chk(slot, pro)
			slot_str := spi_rd_slot_str(slot_int, pro)
			fmt.Println(ui_prn_str(pre_ui_strs(slot_str, pro), mark))
			fmt.Println(">", mode, "slot", slot_int)
		} else {  // view a slot in a *.pre, *.pro, or *.eeprom file
			var file_strs []string
			file_strs, mode = file_dlp_strs(file, mode)
			if mode == "pro" { pro = true }
			slot_int := slot_int_chk(slot, pro)
			fmt.Println(ui_prn_str(pre_ui_strs(file_strs[slot_int], pro), mark))
			fmt.Println("> file", file, mode, "slot", slot_int)	
		}
	} else {
		error_exit("Nothing to do")
	}
}

func knob_cmd(knob, ofs, set string, min, view bool) {
	str_split := (strings.Split(strings.TrimSpace(knob), ":"))
	pg_name, pg_idx := page_lookup(str_split[0])
	if pg_idx < 0 { error_exit("Bad page name") }
	knob_idx := pg_idx * KNOBS  // base index at this point
	prn_str := ""
	if len(str_split) < 2 { 
		prn_str = fmt.Sprint("> ", pg_name)
		knob_idx += PAGE_SEL_KNOB  // bracket page name
	} else if strings.ToLower(strings.TrimSpace(str_split[1])) == "all" {
		prn_str = fmt.Sprint("> ", pg_name, ":all")
		if min || set == "0" {
			sp := sp_open()
			for i:=0; i<KNOBS-1; i++ {
				sp_tx_rx(sp, strconv.Itoa(knob_idx+i) + " 0 wk ", false)
			}
			sp.Close()
			prn_str += fmt.Sprint("=>[min]")
		}
		knob_idx += PAGE_SEL_KNOB  // bracket page name
	} else {
		knob_int, err := strconv.Atoi(str_split[1]); if err != nil { error_exit("Bad knob index") }
		if knob_int < 0 || knob_int >= KNOBS-1 { error_exit("Knob index out of range") }
		knob_idx += knob_int
		ptype, plabel, _, _ := pname_lookup(knob_pnames[knob_idx])
		sp := sp_open()
		rx_str := sp_tx_rx(sp, strconv.Itoa(knob_idx) + " rk ", false)
		sp.Close()
		rd_uint32, _ := strconv.ParseUint(decruft_hcl(rx_str), 16, 32)
		prn_str = fmt.Sprint("> ", pg_name, ":", knob_int, " ", strings.TrimSpace(plabel), "[", strings.TrimSpace(pint_disp(int(rd_uint32), ptype)), "]")
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
	}
	if view { 
		// set LCD to edited screen
		sp := sp_open()
		sp_tx_rx(sp, strconv.Itoa(PAGE_SEL_KNOB) + " " + strconv.Itoa(pg_idx) + " wk ", false)
		sp.Close()
		view_cmd("", false, true, "", knob_idx) 
	}
	fmt.Println(prn_str)
}

// diff DLP file(s) / slot(s) / knobs
func diff_cmd(file, file2 string, pro, knobs bool, slot, slot2 string) {
	mode := "pre"; if pro { mode = "pro" }
	if file != "" {  // compare to a *.dlp file
		file = file_ext_chk(file, ".dlp")
		file_str := file_read_str(file)
		if knobs {  // file vs. knobs
			knob_str := knob_pre_str(sp_rx_knobs_str(), pro)
			fmt.Println(diff_prn_str(diff_pres(file_str, knob_str, pro)))
			fmt.Println(">", mode, "file", file, "vs. knobs" )
		} else if file2 != "" {  // file vs. file2
			file2 = file_ext_chk(file2, ".dlp")
			file2_str := file_read_str(file2)
			fmt.Println(diff_prn_str(diff_pres(file_str, file2_str, pro)))
			fmt.Println(">", mode, "file", file, "vs.", file2 )
		} else if slot != "" {  // file vs. slot
			slot_int := slot_int_chk(slot, pro)
			slot_str := spi_rd_slot_str(slot_int, pro)
			fmt.Println(diff_prn_str(diff_pres(file_str, slot_str, pro)))
			fmt.Println(">", mode, "file", file, "vs. slot", slot )
		} else {
			error_exit("Nothing to do")
		}
	} else if slot != "" {  // compare to a slot
		slot_int := slot_int_chk(slot, pro)
		slot_str := spi_rd_slot_str(slot_int, pro)
		if knobs {  // slot vs. knobs
			knob_str := knob_pre_str(sp_rx_knobs_str(), pro)
			fmt.Println(diff_prn_str(diff_pres(slot_str, knob_str, pro)))
			fmt.Println(">", mode, "slot", slot, "vs. knobs" )
		} else if slot2 != "" {  // slot vs. slot2
			slot2_int := slot_int_chk(slot2, pro)
			slot2_str := spi_rd_slot_str(slot2_int, pro)
			fmt.Println(diff_prn_str(diff_pres(slot_str, slot2_str, pro)))
			fmt.Println(">", mode, "slot", slot, "vs. slot", slot2 )
		} else {
			error_exit("Nothing to do")
		}
	} else {
		error_exit("Nothing to do")
	}
}

// match slots | DLP files in dir2 | "slots" in PRE|PRO|EEPROM file w/ DLP files in dir, list
func match_cmd(dir, dir2, file string, pro, hdr, guess, slots bool) {
	dir_chk(dir)
	name_strs, data_strs := dir_read_strs(dir, ".dlp")
	mode := "pre"; if pro { mode = "pro" }
	if len(data_strs) == 0 { error_exit(fmt.Sprint("No ", mode, " files in directory ", dir)) }
	if slots {
		slots_strs := spi_rd_slots_strs(pro)
		fmt.Print(slots_prn_str(comp_file_data(slots_strs, name_strs, data_strs, pro, guess), pro, hdr))
		fmt.Println("> matched", mode, "slots to", mode, "files in", dir)
	} else if file != "" {
		file_chk(file)
		var file_strs []string
		file_strs, mode = file_dlp_strs(file, mode)
		if mode == "pro" { pro = true }
		fmt.Print(slots_prn_str(comp_file_data(file_strs, name_strs, data_strs, pro, guess), pro, hdr))
		fmt.Println("> matched file", file, "to", mode, "files in", dir)
	} else if dir2 != "" {
		dir_chk(dir2)
		name2_strs, data2_strs := dir_read_strs(dir2, ".dlp")
		if len(data2_strs) == 0 { error_exit(fmt.Sprint("No ", mode, " files in ", dir2)) }
		fmt.Print(files_prn_str(name2_strs, comp_file_data(data2_strs, name_strs, data_strs, pro, guess)))
		fmt.Println("> matched", mode, "files in", dir2, "to", mode, "files in", dir)
	} else {
		error_exit("Nothing to do")
	}
}

// download to file
func dump_cmd(file, slot string, knobs, pro, yes bool) {
	if knobs || slot != "" {  // knobs or slot to DLP file
		file = file_ext_chk(file, ".dlp")
		mode := "pre"; if pro { mode = "pro" }
		if knobs && slot == "" {  // knobs to DLP file
			pints := sp_rx_knobs_pints(mode)
			if file_write_str(file, ints_to_hexs(pints, 4), yes) {
				fmt.Println("> downloaded", mode, "knobs =>", mode, "file", file) 
			}
		} else if !knobs && slot != "" {  // slot to DLP file
			slot_int := slot_int_chk(slot, pro)
			if file_write_str(file, spi_rd_slot_str(slot_int, pro), yes) {
				fmt.Println("> downloaded", mode, "slot", slot_int, "=>", mode, "file", file) 
			}
		} else {
			error_exit("Nothing to do")
		}
	} else {  // bulk dump
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
	}
}

// upload from file
func pump_cmd(file, slot string, knobs, pro bool) {
	if knobs || slot != "" {  // DLP file to knobs or slot
		file = file_ext_chk(file, ".dlp")
		file_str := file_read_str(file)
		mode := "pre"; if pro { mode = "pro" }
		if knobs && slot == "" {  // DLP file to knobs
			pints := hexs_to_ints(file_str, 4)
			if len(pints) < SLOT_BYTES { error_exit("Bad file info") }
			sp_tx_knobs_pints(pints, mode)
			fmt.Println("> uploaded", mode, "file", file, "=>", mode, "knobs") 
		} else if !knobs && slot != "" {  // DLP file to slot
			slot_int := slot_int_chk(slot, pro)
			addr := spi_slot_addr(slot_int, pro)
			spi_wr(addr, file_str, false)
			fmt.Println("> uploaded", mode, "file", file, "=>", mode, "slot", slot_int) 
		} else {
			error_exit("Nothing to do")
		}
	} else {  // bulk pump
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
		if ext == ".spi" || ext == ".eeprom" { reset_cmd() }
	}
}

// *.bnk => *.dlps => slots
func bank_cmd(slot, file string, pro bool) {
	slot_int := slot_int_chk(slot, pro)
	file = file_ext_chk(file, ".bnk")
	bnk_str := file_read_str(file)
	bnk_split := strings.Split(bnk_str, "\n")
	dir, _ := filepath.Split(file)
	for _, line := range bnk_split {
		line := strings.TrimSpace(line);
		if !strings.HasPrefix(line, "//") {  // skip commented lines
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

func morph_cmd(file string, knobs bool, slot string, seed int, mo, mn, me, mf, mr int) {
	rand.Seed(int64(seed))
	prn_str := ""
	var pints []int
	if mo | mn | me | mf | mr == 0 {
		error_exit("Nothing to do")
	} else if knobs {  // morph current knobs
		pints = pints_signed(sp_rx_knobs_pints("pre"), false)
		prn_str = fmt.Sprint("> morphed knobs")
	} else if file != "" {  // morph a *.dlp file
		file = file_ext_chk(file, ".dlp")
		file_str := file_read_str(file)
		pints = pints_signed(hexs_to_ints(file_str, 4), false)
		prn_str = fmt.Sprint("> morphed file ", file)
	} else if slot != "" {  // morph a slot
		slot_int := slot_int_chk(slot, false)
		slot_str := spi_rd_slot_str(slot_int, false)
		pints = pints_signed(hexs_to_ints(slot_str, 4), false)
		prn_str = fmt.Sprint("> morphed slot ", slot_int)
	} else {
		error_exit("Nothing to do")
	}
	prn_str += fmt.Sprint(" (-i=", seed, ")")
	pints = morph_pints(pints, mo, mn, me, mf, mr)
	sp_tx_knobs_pints(pints, "pre")
	fmt.Println(prn_str)
}

// do a bunch of update stuff via interactive menu
func menu_cmd(dir_work string) {
	path_exe, err := os.Executable(); err_chk(err)
	dir_exe := filepath.Dir(path_exe)
	dir_all := filepath.Join(dir_exe, PRESETS_DIR)
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
	//
	first_f := true
	for {
		if !first_f {
			user_input("Press <ENTER> to return to the MENU", true)
		}
		first_f = false
		fmt.Println()
		fmt.Println()
		fmt.Println(" ---------------------------------")
		fmt.Println(" |  D-LEV LIBRARIAN - VERSION", ver_tbl[0].lib, " |")
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
		menu_sel := user_input("Please select a MENU option", true)
		switch menu_sel {
			case "0" :
				fmt.Println()
				fmt.Println()
				fmt.Println(" ////////////")
				fmt.Println(" // README //")
				fmt.Println(" ////////////")
				fmt.Println(" - To UPDATE the software and UPDATE ALL of the preset SLOTS: Do 1 thru 7.")
				fmt.Println(" - To UPDATE the software and OVERWRITE ALL of the preset SLOTS: Do 1 thru 3, then 9.")
				fmt.Println(" - TO UPDATE & OVERWRITE ABSOLUTELY EVERYTHING INCLUDING PROFILE SLOTS: Do 1, 2, 10.")
				fmt.Println(" - To CONVERT ALL of the preset SLOTS to MONO: Do 1, 2, 4, 8, 6.")
				fmt.Println(" - If you run into trouble, quit and pump the backup EEPROM file created in step 2.")
				fmt.Println(" - Valid prompt responses: y=yes, ENTER=no, q=quit the program.")
				fmt.Println(" - If unresponsive, do a CTRL-C (hold down the CONTROL key and press the C key).")
				fmt.Println(" - DO NOT turn or press any D-Lev knobs during the upload / download process!")
				fmt.Println(" ////////////")
				fmt.Println(" // README //")
				fmt.Println(" ////////////")
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
				file_backup := date() + "_backup.eeprom"
				if user_prompt("Do you want to BACKUP your ENTIRE D-Lev to the FILE: " + file_backup + "?", false, true) {
					dump_cmd(file_backup, "", false, false, false)
				}
			case "3" :
				fmt.Println()
				sw_upd := ver_cmd("", true)
				if (sw_upd)  {
					if user_prompt("Do you want to UPDATE your D-Lev SOFTWARE with the FILE: "+ file_spi + "?", false, true) {
						pump_cmd(path_spi, "", false, false)
						fmt.Println()
						ver_cmd("", false)
					}
				}
			case "4" :
				if user_prompt("Do you want to DOWNLOAD your D-Lev presets to " + dir_work +"?", false, true) {
					dump_cmd(path_pre_dl, "", false, false, true)
					split_cmd(path_pre_dl, true)
				}
			case "5" :
				if user_prompt("Do you want to UPDATE the presets in " + dir_work + "?", false, true) {
					process_dlps(dir_work, dir_work, false, false, true, false, true)
				}
			case "6" :
				if user_prompt("Do you want to UPLOAD the presets in "+ dir_work + "?", false, true) {
					join_cmd(path_pre_ul, true)
					pump_cmd(path_pre_ul, "", false, false)
				}
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
					process_dlps(dir_work, dir_work, false, true, false, false, true)
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
					ver_cmd("", false)
				}
			case "" :
				// do nothing
			default:
				fmt.Println("> Invalid menu selection!")
		}
	}
}

