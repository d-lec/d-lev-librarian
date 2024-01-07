package main

/*
 * d-lib support functions
*/

import (
	"strings"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// print error message and exit program
func error_exit(error_str string) {
	fmt.Println("> -ERROR-", error_str, "!")
	os.Exit(0) 
}

// print quit message and exit program
func quit_exit() {
	fmt.Println("> -QUIT- exiting program...")
	os.Exit(0) 
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
func user_input(prompt string, newln bool) (string) {
	if newln { fmt.Println() }
	fmt.Print("> ", prompt, " (q=quit): ")
	input := user_word()
	if input == "q" { quit_exit() }
	return input
}


var menu_readme_str = `
//////////////////
// README START //
//////////////////
D-LEV SOFTWARE CHANGES:
- NOISE:nois[0] now routes the oscillators thru the noise EQ & SVF.
- VOLUME:damp[1:63] sets MIDI note_off @ VOLUME:dloc (h/t Vincent).
- VOLUME:damp[1:63] won't allow MIDI note_on below dloc threshold (h/t Vincent).
- VOLUME:velo removed velo<0 velocity rectify option.
- MIDI:velo scaled up 2x.

D-LIB LIBRARIAN CHANGES:
- New menu option to list, select, & pump EEPROM file (system restore).
- New menu option to display pitch & volume field frequencies.
- Command bank can now generate / operate on PRE|PRO files too.
- Command morph can now use PRE|EEPROM files as source.
- Command morph is now multiple -pkv page:knob:value based (& back to integer values).
- Command morph -sm flag removed: negative values gives 50/50 chance of sign flip.
- Command knob flag change: -k => -pk
- Fuzzy matching for page lookup & command hinting.

UPDATE PROCESSING:
- There are no preset or profile updates needed with this release.
- Please note: The preset / profile updating process is always entirely algorithmic,
  and so doesn't rely at all on which presets / profiles are in which slots.

GENERAL:
- To UPDATE the software and get new presets: Do 1, 2, 3, 5.
- To UPDATE the software and OVERWRITE ALL OF THE PRESET SLOTS: Do 1, 2, 3, 6.
- TO UPDATE & OVERWRITE ABSOLUTELY EVERYTHING INCLUDING PROFILE SLOTS: Do 1, 2, 8.
- TO RESTORE & OVERWRITE ABSOLUTELY EVERYTHING INCLUDING PROFILE SLOTS: Do 1, 2, 9.
- If you run into trouble, RESTORE your backup EEPROM file created in step 2 via step 9.
- Valid prompt responses: y=yes, ENTER=no, q=quit the program.
- If unresponsive, do a CTRL-C (hold down the CONTROL key and press the C key).
- DO NOT turn or press any D-Lev knobs during the upload / download process!
////////////////
// README END //
////////////////
`

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
//	path_pre_dl := filepath.Join(dir_work, "download.pre")
//	path_pre_ul := filepath.Join(dir_work, "upload.pre")
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
		fmt.Println("  0. !PLEASE README!")
		fmt.Println("  1. Serial Port: Setup & check.")
		fmt.Println("  2. System Backup: Download EVERYTHING to an EEPROM file.")
		fmt.Println("  3. Software Update: Check version & upload SPI file.")
		
//		fmt.Println(" 4A. Download Your Presets: To the", dir_work, "directory.")
//		fmt.Println(" 4B. Update Your Presets: In the", dir_work, "directory.")
//		fmt.Println(" 4C. Upload Your Presets: From the", dir_work, "directory.")

		fmt.Println("  5. Upload New Presets: Install the latest new presets.")
		fmt.Println("  6. Upload All Presets: Overwrite ALL presets from the", PRESETS_DIR, "directory.")
		fmt.Println("  7. Factory Reset: Overwrite EVERYTHING with the latest FACTORY EEPROM file.")
		fmt.Println("  8. System Restore: Overwrite EVERYTHING from a BACKUP EEPROM file.")

//		fmt.Println("  9. Convert All presets in the", dir_work, "directory to MONO.")
//		fmt.Println(" 10. Pump some updated presets to their default slots.")

		fmt.Println(" hz. Stats: Report pitch & volume field frequencies.")
		fmt.Println("  h. Help: List of commands.")
		fmt.Println(" hv. Help: List of commands with notes and typical examples of their use.")
		menu_sel := user_input("Please select a MENU option", true)
		switch menu_sel {
		case "0" :
			fmt.Println()
			fmt.Print(menu_readme_str)
		case "1" :
			fmt.Println()
			port_cmd("", true)
			if user_prompt("Do you want to CHANGE the ACTIVE active port?", false, true) {
				port_new := user_input("Please input PORT number", false)
				if port_new == "" {
					fmt.Println("> -CANCEL-")
				} else {
					fmt.Println()
					port_cmd(port_new, true)
				}
			}
			if user_prompt("Do you want to TEST the ACTIVE port (do a CTRL-C if it hangs)?", false, true) {
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
			sw_upd := ver_cmd("", true, true, true)
			if (sw_upd)  {
				if user_prompt("Do you want to UPDATE your D-Lev SOFTWARE with the FILE: "+ file_spi + "?", false, true) {
					pump_cmd(path_spi, "", false, false)
					fmt.Println()
					ver_cmd("", true, false, false)
				}
			}
/*
		case "4a" :
			if user_prompt("Do you want to DOWNLOAD your D-Lev presets to " + dir_work +"?", false, true) {
				dump_cmd(path_pre_dl, "", false, false, true)
				split_cmd(path_pre_dl, true)
			}
*/
/*
			if user_prompt("Do you want to DOWNLOAD your D-Lev profiles to " + dir_work_pro +"?", false, true) {
				dump_cmd(path_pro_dl, "", false, true, true)
				split_cmd(path_pro_dl, true)
			}
*/
/*
		case "4b" :
			if user_prompt("Do you want to UPDATE the presets in " + dir_work + "?", false, true) {
				process_dlps(dir_work, dir_work, false, false, true, false, true, false)
			}
/*
			if user_prompt("Do you want to UPDATE the profiles in " + dir_work_pro + "?", false, true) {
				process_dlps(dir_work_pro, dir_work_pro, true, false, true, false, true, false)
			}
*/
/*
		case "4c" :
			if user_prompt("Do you want to UPLOAD the presets in "+ dir_work + "?", false, true) {
				join_cmd(path_pre_ul, true)
				pump_cmd(path_pre_ul, "", false, false)
			}
*/
/*
			if user_prompt("Do you want to UPLOAD the profiles in "+ dir_work_pro + "?", false, true) {
				join_cmd(path_pro_ul, true)
				pump_cmd(path_pro_ul, "", false, true)
			}
*/
		case "5" :
			fmt.Println()
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
					bank_cmd(slot, path_bank_new, "", false)
				}
				if user_prompt("Do you want to EXAMINE your current D-Lev presets?", false, true) {
					match_cmd(dir_all, "", "", false, false, true, true) 
				}
			}
		case "6" :
			if user_prompt("Do you want to OVERWRITE all D-Lev preset slots with presets in " + PRESETS_DIR + "?", false, true) {
				bank_cmd("0", path_bank, "", false)
				if user_prompt("Do you want to EXAMINE your current D-Lev presets?", false, true) {
					match_cmd(dir_all, "", "", false, false, true, true) 
				}
			}
		case "7" :
			if user_prompt("Do you want to OVERWRITE ABSOLUTELY EVERYTHING in your D-Lev with the FACTORY EEPROM file?", false, true) {
				pump_cmd(path_factory, "", false, false)
				fmt.Println()
				ver_cmd("", true, false, false)
			}
		case "8" :
			fmt.Println()
			fmt.Println("> Here is a LIST of EEPROM files in the", dir_exe, "directory:")
			fmt.Println()
			eeprom_files, _ := dir_read_strs(dir_exe, ".eeprom", false)
			if len(eeprom_files) < 1 {
				fmt.Println("No EEPROM files!")
			} else { 
				for idx, file := range eeprom_files {
					fmt.Print(" [", idx, "] ", file, "\n")
				}
				if user_prompt("Do you want to OVERWRITE ABSOLUTELY EVERYTHING in your D-Lev with one of these BACKUP EEPROM files?", false, true) {
					file_sel := user_input("Please input FILE number", false)
					file_idx, err := strconv.Atoi(file_sel)
					if file_sel == "" {
						fmt.Println("> -CANCEL-")
					} else if (err != nil) || (file_idx < 0) || file_idx >= len(eeprom_files) { error_exit("Bad file selection: " + file_sel)
					} else {
						pump_cmd(filepath.Join(dir_exe, eeprom_files[file_idx]), "", false, false)
					}
				}
			}
/*
		case "8" :
			if user_prompt("Do you want to CONVERT all of the presets in " + dir_work + " to MONO?", false, true) {
				process_dlps(dir_work, dir_work, false, true, false, false, true, false)
			}
*/
/*
		case "10" :
			if user_prompt("Do you want to EXAMINE your current D-Lev presets?", false, true) {
				match_cmd(dir_all, "", "", false, false, true, true) 
			}
			if user_prompt("Do you want to PUMP some modified presets to their default slots?", false, true) {
				dlps_pump(dir_all, false)
			}
*/
		case "hz" :
			fmt.Println()
			stats_cmd(true, true, false)
		case "h" :
			fmt.Println()
			help_cmd(false)
		case "hv" :
			fmt.Println()
			help_cmd(true)
/*
		case "robs" :
			if user_prompt("Do you want to apply Rob Schwimmer Pitch Preview to the presets in " + dir_work + "?", false, true) {
				process_dlps(dir_work, dir_work, false, false, false, true, true, false)
			}
*/
		case "" :
			// do nothing
		default:
			fmt.Println("> Invalid menu selection!")
		}
	}
}


