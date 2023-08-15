/*
d-lib.cpp
- Started 2019-12-18 by Eric Wallin
- No explicit project management, just compile this file.

================
!FIXME! / TODO / NOTE
- note: DON'T USE ISSPACE()! - crashes w/ bad input (use ch_is_white)
- note: DON'T USE STOL()! - crashes w/ bad input (use strtol w/ str_is_hex() check)
- note: DON'T USE getch()! - broken, replaced with read()
- note: do a cout flush before sleep!
================
*/

// comment this out for POSIX compile
// #define _MSWIN

// c++ stuff
#include <bits/stdc++.h>
#include <stdint.h>
#ifdef _MSWIN
  #include <conio.h>
  #include <windows.h>
  #include <unistd.h>
#else
  #include <fcntl.h>
  #include <linux/serial.h>
  #include <sys/ioctl.h>
#endif

using namespace std;

// lib stuff
#include "d-lib_pkg.cpp"
#include "d-lib_console.cpp"
#include "d-lib_serial.cpp"
#include "d-lib_cmd.cpp"
#include "d-lib_preset.cpp"
#include "d-lib_preset_upd.cpp"
#include "d-lib_dir.cpp"
#include "d-lib_term.cpp"
#include "d-lib_file_view.cpp"


//////////
// MAIN //
//////////

int main(int argc, char* argv[]) {

	// if there are command line arguments
	if (argc > 1) {  
		cout << "\nd-lib uses no command line arguements";
		return(-1);
	}

#ifdef _MSWIN
	// size the MSWIN console (width +1 for WinXP console?)
	con_size(CON_W+1, CON_H);
	static const string dev_base = "\\\\.\\COM";
	// for more accurate (~1ms) timing - NOTE: add -lwinmm to build!
	timeBeginPeriod(1);
#else
	static const string dev_base = "/dev/ttyUSB";  // breaks w/o "static"?
#endif

	// constants
	const int32_t FILE_Y = 26;
	const int32_t PREV_Y = 28;
	const int32_t CMD_Y = PREV_Y + 1;
	const int32_t CMD_W = CON_W;
	// file strings
	static const string help_fname = "d-lib_help.txt";
	vector<string> help_sv; 
	string pre_fname = "";
	string prn_fname = "d-lib_print.txt";
	// command line
	const int32_t CMD_DEPTH = 32;  // cmd history depth
	bool cmd_f = false;  // cmd flag
	bool err_f = false;  // cmd error flag
	string cmd_str = "";  // cmd str
	string cmd_prev = "";  // cmd history
	string cmd_msg = "";  // cmd result message
	cmd_mem<CMD_DEPTH> cmd_hist;  // cmd buffer
	uint32_t cursor_pos = 0;
	// modes
	enum ui_modes { edit, files, slots } mode = edit; 
	bool first_f = true;  // do extra stuff @ start
	bool hl_f = false;  // hilite mode flag
	// coords & sycg
	int32_t ed_x = 0;  // edit coords
	int32_t ed_y = 0;
	int32_t dir_idx = 0;  // hilite edit idx
	int64_t key_i = 0;
	int32_t dl_min = 0;
	int32_t dl_max = 79;
	int32_t dl_idx = 0;
	// serial port
	sp_type sp;
	uint32_t dev_num = 0;
	const uint32_t dev_max = 31;
	bool dev_open[dev_max+1];  // +1 to hold all
	uint32_t devs = 0;
	uint32_t dev_last = 0;

	// output version
	cout << "-= D-Lev Preset Librarian (version 0x" << hex << LIBV << ") =-" << endl << endl; 
	cout.unsetf(ios::hex); // I hate c++
	
	// scan for open ports
	while (dev_num <= dev_max) {
		dev_open[dev_num] = sp_test_port(dev_base + to_string(dev_num));
		if (dev_open[dev_num]) { dev_last = dev_num; devs++;}
		dev_num++;
	}
	
	if (!devs) {
		cout << "NO OPEN PORTS! :-(" << endl << endl;
		dev_num = 0;
	}
	else {
		cout << "Open ports:" << endl;
		for (uint32_t i=0; i<=dev_max; i++) { if (dev_open[i]) { cout << " " << i; } }
		// query for port
		dev_num = dev_last;
		cout << endl << "Select port (hit Enter for port " << dev_num << "): ";
		string dev_str;
		getline(cin, dev_str);
		if (str_is_hex(dev_str)) { dev_num = strtol(dev_str.c_str(), NULL, 0); }
		cout << endl;
	}

	// open port & check
	sp_open_port(sp, dev_base + to_string(dev_num));
	cout << "sel  : " << dev_num << endl; 
	cout << "port : " << sp.dev << endl; 
	if (sp_open_err(sp)) {
		cout << "open : ERROR! (fd=" << sp.fd << ")" << endl;
	}
	else { 
		cout << "open : OK! (fd=" << sp.fd << ")" << endl;
		if (sp_config_port(sp)) { cout << "cfg  : ERROR!" << endl; }
		else { cout << "cfg  : OK! " << endl; }
	}

	// read help file
	if (file_to_sv(help_fname, help_sv, false)) { 
		cout << "can't read help file:" << help_fname << endl << endl;
	}

#ifndef _MSWIN
	// switch POSIX kbd mode to non-blocking
	kbd_mode(true, false);
#endif

	// pause to display addresses and port status
	cout << endl << "press any key to continue..." << endl;
	kbd_wait();

	// directory listings
	bool view_f = true;  // directory view flag
	vector<string> dir_all, dir_dlp, dir_view;
	dir_proc(dir_all, dir_dlp, dir_view, view_f);  // read, sort, filter by ext, remove ext

	// downloaded preset contents
	vector<string> dl_pre_sv;

	// downloaded slot names
	vector<string> dl_names_sv;

	// fill knob table
	knob_tbl_init();

	// loop forever
	while (true) {

		if (!first_f) {  // skip once @ start

		    //////////////
		    // key proc //
		    //////////////

			// get key code (blocking)
			key_i = kbd_getkey();
			bool key_used = false;

			// in edit hilite mode
			if ((mode == edit) && hl_f && !key_used) {
				key_used = true;  // default
				int32_t dx = 0;
				int32_t dy = 0;
				// quit hilite mode
				if (key_i == KBD_ESC) { hl_f = false; }
				// hilite navigation
				else if (key_i == KBD_UP) { dy = -1; }
				else if (key_i == KBD_DN) { dy = 1; }
				else if (key_i == KBD_LT) { dx = -1; }
				else if (key_i == KBD_RT) { dx = 1; }
				else if (key_i == KBD_CTRL_UP) { dy = -4; }
				else if (key_i == KBD_CTRL_DN) { dy = 4; }
				else if (key_i == KBD_CTRL_LT) { dx = -2; }
				else if (key_i == KBD_CTRL_RT) { dx = 2; }
				else if (key_i == KBD_HOME) { xy_home(ed_x, ed_y); }
				else if (key_i == KBD_END) { xy_end(ed_x, ed_y); }
				// inc/dec
				else if (key_i == KBD_ALT_UP) { xy_delta_wr_tx(sp, 1, ed_x, ed_y); }
				else if (key_i == KBD_ALT_DN) { xy_delta_wr_tx(sp, -1, ed_x, ed_y); }
				else if (key_i == KBD_ALT_RT) { xy_delta_wr_tx(sp, 10, ed_x, ed_y); }
				else if (key_i == KBD_ALT_LT) { xy_delta_wr_tx(sp, -10, ed_x, ed_y); }
				// zero entry
				else if (key_i == KBD_DEL) { cmd_str = "0 "; }
				// default
				else { key_used = false; }
				// recalc idx
				xy_move(ed_x, ed_y, dx, dy);
			}
			// in files hilite mode
			if ((mode == files) && hl_f && !key_used) {
				key_used = true;  // default
				int32_t dx = 0;
				int32_t dy = 0;
				// quit hilite mode
				if (key_i == KBD_ESC) { hl_f = false; }
				// hilite navigation
				else if (key_i == KBD_UP) { dy = -1; }
				else if (key_i == KBD_DN) { dy = 1; }
				else if (key_i == KBD_LT) { dx = -1; }
				else if (key_i == KBD_RT) { dx = 1; }
				else if (key_i == KBD_CTRL_UP) { dy = -4; }
				else if (key_i == KBD_CTRL_DN) { dy = 4; }
				else if (key_i == KBD_CTRL_LT) { dx = -2; }
				else if (key_i == KBD_CTRL_RT) { dx = 2; }
				else if (key_i == KBD_HOME) { dir_idx = 0; }
				else if (key_i == KBD_END) { dir_idx = dir_view.size()-1; }
				// delete file
				else if (key_i == KBD_DEL) { 
					cmd_str = dir_view[dir_idx] + " delf ";
					cursor_pos = cmd_str.size();
				}
				// copy hilited file name to cmd line
				else if (key_i == KBD_ENT) {
					cmd_str += dir_view[dir_idx];
					cursor_pos = cmd_str.size();
					hl_f = false;
				}
				// load hilited voice preset file to editor
				else if (key_i == 'e') {
					cmd_str = dir_view[dir_idx];
					cmd_str += " ptoe ";
					cursor_pos = cmd_str.size();
				}
				// load hilited system profile file to editor
				else if (key_i == 's') {
					cmd_str = dir_view[dir_idx];
					cmd_str += " ptoe_sys ";
					cursor_pos = cmd_str.size();
				}
				// view hilited file
				else if (key_i == 'v') { 
					cmd_str = dir_view[dir_idx] + " viewf ";
					cursor_pos = cmd_str.size();
				}
				// change to hilighed directory
				else if (key_i == 'c') { 
					cmd_str = dir_view[dir_idx] + " cd ";
					cursor_pos = cmd_str.size();
				}
				// update hilited voice preset file & inc to next
				else if (key_i == 'u') {
					cmd_str = dir_view[dir_idx];
					cmd_str += " up ";
					cursor_pos = cmd_str.size();
					dy = 1;
				}
				else { key_used = false; }
				// recalc idx
				dir_idx = xy_move_idx(dir_idx, dx, dy, dir_view.size()-1, DIR_ROWS);
			}
			// in slots hilite mode
			if ((mode == slots) && hl_f && !key_used) {
				key_used = true;  // default
				int32_t dx = 0;
				int32_t dy = 0;
				int32_t slot = dl_idx + dl_min;
				// quit hilite mode
				if (key_i == KBD_ESC) { hl_f = false; }
				// hilite navigation
				else if (key_i == KBD_UP) { dy = -1; }
				else if (key_i == KBD_DN) { dy = 1; }
				else if (key_i == KBD_LT) { dx = -1; }
				else if (key_i == KBD_RT) { dx = 1; }
				else if (key_i == KBD_CTRL_UP) { dy = -4; }
				else if (key_i == KBD_CTRL_DN) { dy = 4; }
				else if (key_i == KBD_CTRL_LT) { dx = -2; }
				else if (key_i == KBD_CTRL_RT) { dx = 2; }
				else if (key_i == KBD_HOME) { dl_idx = 0; }
				else if (key_i == KBD_END) { dl_idx = dl_pre_sv.size()-1; }
				// copy hilited slot name to cmd line & exit hilite mode
				else if (key_i == KBD_ENT) { 
					cmd_str += dl_names_sv[dl_idx]; 
					cursor_pos = cmd_str.size();
					hl_f = false;
				}
				// download hilited slot to editor & rename editor
				else if (key_i == 'e') {
					pre_fname = dl_names_sv[dl_idx];
					cmd_str = to_string(slot); 
					cmd_str += " stoe ";
					cursor_pos = cmd_str.size();
				}
				else { key_used = false; }
				// recalc idx
				dl_idx = xy_move_idx(dl_idx, dx, dy, max(0, int32_t(dl_pre_sv.size()-1)), SLOTS_ROWS);
			}
			// command line ops
			if (!key_used && !cmd_edit(cmd_str, cursor_pos, key_i)) {
				// edit mode
				if ((mode == edit) && ((key_i == KBD_CTRL_UP) || (key_i == KBD_CTRL_DN))) {
					hl_f = true; 
					cmd_str.clear();
					cursor_pos = cmd_str.size();
				}
				else if (key_i == KBD_F1) {
					if (mode != edit) { 
						mode = edit;
						hl_f = false; 
						con_cls();
					}
					else if (!hl_f) { 
						hl_f = true; 
						cmd_str.clear();
						cursor_pos = cmd_str.size();
					}
					else { hl_f = false; }
				}
				// files mode
				else if ((mode == files) && ((key_i == KBD_CTRL_UP) || (key_i == KBD_CTRL_DN))) {
					hl_f = true; 
				}
				else if (key_i == KBD_F2) {
					if (mode != files) { 
						mode = files;
						dir_idx = xy_move_idx(dir_idx, 0, 0, dir_view.size()-1, DIR_ROWS);
						hl_f = false; 
						con_cls();
					}
					else { hl_f = !hl_f; }
				}
				// slots mode
				else if ((mode == slots) && ((key_i == KBD_CTRL_UP) || (key_i == KBD_CTRL_DN))) {
					hl_f = true; 
				}
				else if (key_i == KBD_F3) {
					if (mode != slots) { 
						mode = slots;
						hl_f = false; 
						con_cls(); 
					}
					else { hl_f = !hl_f; }
				}
				// cmd_hist recall
				else if (key_i == KBD_UP) {
					if (cmd_hist.rd(cmd_str, true)) { cursor_pos = cmd_str.size(); }
				}
				else if (key_i == KBD_DN) {
					if (cmd_hist.rd(cmd_str, false)) { cursor_pos = cmd_str.size(); }
				}
				// file name autocomplete
				else if (key_i == KBD_TAB) {
					cmd_acomp(cmd_str, cursor_pos, dir_all); 
				}
			}
		}
		

		////////////////
		// token proc //
		////////////////

		if (!cmd_str.empty() && ch_is_white(cmd_str.back())) {
			con_xy(0+cursor_pos, CMD_Y);  // position cursor on cmd line
			vector<string> vtoken = str_to_tokens(cmd_str);
			int32_t tokens = vtoken.size();
			cmd_f = false;
			if (tokens) {
				// get last token
				string last_token = vtoken[tokens-1];
				// edit @ hilite
				if ((mode == edit) && hl_f) {
					cmd_f = true;
					err_f = !str_is_int(vtoken[0]);
					if (err_f) { cmd_msg = "in edit mode, expecting a number"; }
					if (!err_f) {
						int32_t ed_val = strtol(vtoken[0].c_str(), NULL, 0);  // STRTOL!
						if (xy_on_title(ed_x, ed_y)) {  // on title, edit all
							xy_input_wr_tx(sp, ed_val, ed_x-1, ed_y-3);
							xy_input_wr_tx(sp, ed_val, ed_x-1, ed_y-2);
							xy_input_wr_tx(sp, ed_val, ed_x-1, ed_y-1);
							xy_input_wr_tx(sp, ed_val, ed_x-1, ed_y-0);
							xy_input_wr_tx(sp, ed_val, ed_x, ed_y-3);
							xy_input_wr_tx(sp, ed_val, ed_x, ed_y-2);
							xy_input_wr_tx(sp, ed_val, ed_x, ed_y-1);
						}
						else {  // edit single
							xy_input_wr_tx(sp, ed_val, ed_x, ed_y);
						}
					}
				}
				// refresh F2 file list
				else if ((last_token == "r") && (tokens == 1)) {
					cmd_f = true;
					err_f = false;
					dir_proc(dir_all, dir_dlp, dir_view, view_f);  // update file list
					dir_idx = 0;
				}
				// toggle declutter of F2 files view
				else if ((last_token == "d") && (tokens == 1)) {
					cmd_f = true;
					err_f = false;
					if (mode == files) { 
						view_f = !view_f;
						dir_proc(dir_all, dir_dlp, dir_view, view_f);  // update file list
					}
				}
				// quit
				else if (((last_token == "quit") || (last_token == "exit")) && (tokens == 1)) {
					cmd_f = true;
					err_f = false;
					kbd_mode(false, true);  // restore orig keyboard mode
					sp_close_port(sp);  // close serial port
					#ifdef _MSWIN
						timeEndPeriod(1);
					#else
						kbd_mode(false, true);  // restore orig keyboard mode
					#endif					
					return(0);  // go bye-bye
				}
				// help mode
				else if ((last_token == "help") && (tokens == 1)) {
					cmd_f = true;
					sv_view(help_sv, CON_W, CON_H);  // call sv viewer
				}
				// view mode
				else if (last_token == "viewf") {
					cmd_f = true;
					err_f = (tokens != 2);
					if (err_f) { cmd_msg = "expecting 1 filename"; }
					if (!err_f) {
						string fnamex = file_ext_add(vtoken[0], DLP_EXT);
						err_f = file_view(fnamex, CON_W, CON_H);
						if (err_f) { cmd_msg = "can't view file:" + fnamex; }
					}
				}
				// terminal mode
				else if ((last_token == "term") && (tokens == 1)) {
					cmd_f = true;
					err_f = false;
					term(sp);  // call terminal sub
				}
				// read SW version
				else if ((last_token == "ver") && (tokens == 1)) {
					cmd_f = true;
					cmd_msg = last_token + " ";
					err_f = cmd_tx(sp, cmd_msg, true);  // call sub
					if (err_f) { cmd_msg = "can't read SW version"; }
				}
				// read CRC
				else if ((last_token == "crc") && (tokens == 1)) {
					cmd_f = true;
					cmd_msg = last_token + " ";
					err_f = cmd_tx(sp, cmd_msg, true);  // call sub
					if (err_f) { cmd_msg = "can't read CRC version"; }
				}
				// do ACAL
				else if ((last_token == "acal") && (tokens == 1)) {
					cmd_f = true;
					cmd_msg = last_token + " ";
					err_f = cmd_tx(sp, cmd_msg, false);  // call sub
					if (err_f) { cmd_msg = "can't do ACAL"; }
				}
				// delete file
				else if (last_token == "delf") {
					cmd_f = true;
					err_f = (tokens != 2);
					if (err_f) { cmd_msg = "expecting 1 filename"; }
					if (!err_f) {
						string fnamex = file_ext_add(vtoken[0], DLP_EXT);
						err_f = file_delete(fnamex);
						if (err_f) { cmd_msg = "can't delete file:" + fnamex; }
					}
					dir_proc(dir_all, dir_dlp, dir_view, view_f);  // update file list
					dir_idx = xy_move_idx(dir_idx, 0, 0, dir_view.size()-1, DIR_ROWS);
				}
				// rename | copy file
				else if ((last_token == "renf") || (last_token == "copf")) {
					cmd_f = true;
					string src_fnamex, dst_fnamex;
					err_f = (tokens != 3);
					if (err_f) { cmd_msg = "expecting 2 filenames"; }
					if (!err_f) {
						src_fnamex = file_ext_add(vtoken[0], DLP_EXT);
						err_f = !file_exists(src_fnamex);
						if (err_f) { cmd_msg = "source file doesn't exist:" + src_fnamex; }
					}
					if (!err_f) {
						dst_fnamex = file_ext_add(vtoken[1], DLP_EXT);
						err_f = file_exists(dst_fnamex);
						if (err_f) { cmd_msg = "target file exists:" + dst_fnamex; }
					}
					if (!err_f && (last_token == "renf")) {
						err_f = file_rename(src_fnamex, dst_fnamex);
						if (err_f) { cmd_msg = "can't rename file:" + src_fnamex; }
					}
					if (!err_f && (last_token == "copf")) {
						err_f = file_copy(src_fnamex, dst_fnamex);
						if (err_f) { cmd_msg = "can't copy file:" + src_fnamex; }
					}
					dir_proc(dir_all, dir_dlp, dir_view, view_f);  // update file list
					dir_idx = xy_move_idx(dir_idx, 0, 0, dir_view.size()-1, DIR_ROWS);
				}
				// change directory
				else if ((last_token == "cd") || (last_token == "chdir")) {
					cmd_f = true;
					err_f = (tokens != 2);
					if (err_f) { cmd_msg = "expecting directory name"; }
					if (!err_f) {
						err_f = chdir(vtoken[0].c_str());
						if (err_f) { cmd_msg = "can't open directory:" + vtoken[0]; }
					}
					dir_proc(dir_all, dir_dlp, dir_view, view_f);  // update file list
					dir_idx = 0;
				}
				// rename editor preset
				else if (last_token == "rene") {
					cmd_f = true;
					err_f = (tokens != 2); 
					if (err_f) { cmd_msg = "expecting a filename"; }
					if (!err_f) { pre_fname = vtoken[0].c_str(); }
				}
				// preset / profile file => editor
				else if ((last_token == "ptoe") || (last_token == "ptoe_sys")) {
					cmd_f = true;
					string rd_fname;
					err_f = (tokens > 2);
					bool sys_f = (last_token == "ptoe_sys");
					if (err_f) { cmd_msg = "expecting 0 or 1 filenames"; }
					if (!err_f) {
						rd_fname = (tokens == 2) ? vtoken[0] : pre_fname;
						string rd_fnamex = file_ext_add(rd_fname, DLP_EXT);  // add extension
						err_f = ed_dlp_rd(rd_fnamex, sys_f);
						if (err_f) { cmd_msg = "can't read file:" + rd_fnamex; }
					}
					if (!err_f) { 
						pre_fname = rd_fname; 
						mode = edit;  // edit mode
					}
				}
				// editor => preset / profile file
				else if ((last_token == "etop") || (last_token == "etop_sys")) {
					cmd_f = true;
					string wr_fname, wr_fnamex;
					err_f = (tokens > 2);
					if (err_f) { cmd_msg = "expecting 0 or 1 filenames"; }
					if (!err_f) {
						bool sys_f = (last_token == "etop_sys");
						wr_fname = (tokens == 2) ? vtoken[0] : pre_fname;
						wr_fnamex = file_ext_add(wr_fname, DLP_EXT);  // add extension
						err_f = ed_dlp_wr(wr_fnamex, sys_f);
						if (err_f) { cmd_msg = "can't write file:" + wr_fnamex; }
					}
					if (!err_f) { pre_fname = wr_fname; }
					dir_proc(dir_all, dir_dlp, dir_view, view_f);  // update file list
					dir_idx = xy_move_idx(dir_idx, 0, 0, dir_view.size()-1, DIR_ROWS);
				}
				// knobs => editor
				else if ((last_token == "ktoe") && (tokens == 1)) {
					cmd_f = true;
					err_f = ed_knobs_rx(sp);
					if (err_f) { cmd_msg = "can't download from knobs"; }
					if (!err_f) { 
						mode = edit;  // edit mode
					}
				}
				// editor => preset knobs
				else if ((last_token == "etok") && (tokens == 1)) {
					cmd_f = true;
					err_f = ed_knobs_tx(sp, false);
					if (err_f) { cmd_msg = "can't upload knobs"; }
				}
				// editor => profile knobs
				else if ((last_token == "etok_sys") && (tokens == 1)) {
					cmd_f = true;
					err_f = ed_knobs_tx(sp, true);
					if (err_f) { cmd_msg = "can't upload knobs"; }
				}
				// screen => text file
				else if (last_token == "print") {
					cmd_f = true;
					string fname;
					err_f = (tokens > 2); 
					if (err_f) { cmd_msg = "expecting 0 or 1 filenames"; }
					if (!err_f) {
						fname = prn_fname;  // safe default
						if (tokens == 2) { fname = vtoken[0]; }
						if (mode == edit) {
							err_f = ed_print_to_file(fname, pre_fname, true);  // append
						}
						else if (mode == files) {
							err_f = dir_print_to_file(fname, dir_view, true);  // append
						}
						else if (mode == slots) {
							err_f = slots_print_to_file(fname, dl_min, dl_names_sv, true);  // append
						}
						if (err_f) { cmd_msg = "can't write file:" + fname; }
					}
					if (!err_f) { prn_fname = fname; }
				}
				// preset / profile file => slot
				else if ((last_token == "ptos") || (last_token == "ptos_sys")) {
					cmd_f = true;
					int32_t slot = 0;
					string fnamex;
					err_f = (tokens != 3); 
					if (err_f) { cmd_msg = "expecting file & slot"; }
					if (!err_f) {  // check slot limits
						if (!str_is_int(vtoken[1])) { err_f = true; }
						if (err_f) { cmd_msg = "expecting a number"; }
					}
					if (!err_f) {
						slot = strtol(vtoken[1].c_str(), NULL, 0);  // STRTOL!
						if (last_token == "ptos_sys") {
							slot += 128;
							err_f = !slot_sys(slot);
						}
						else {
							err_f = slot_error(slot);
						}
						if (err_f) { cmd_msg = "slot out of range"; }
					}
					if (!err_f) {
						fnamex = file_ext_add(vtoken[0], DLP_EXT);  // add extension
						err_f = !file_exists(fnamex);
						if (err_f) { cmd_msg = "can't read file:" + fnamex;}
					}
					if (!err_f) { 
						err_f = dlp_slot_tx(sp, slot, fnamex);
						if (err_f) { cmd_msg = "can't upload to slot:" + to_string(slot); }
					}
				}
				// slot => preset / profile file
				else if ((last_token == "stop") || (last_token == "stop_sys")) {
					cmd_f = true;
					int32_t slot = 0;
					err_f = (tokens != 3); 
					if (err_f) { cmd_msg = "expecting slot & file"; }
					if (!err_f) {  // check slot limits
						err_f = !str_is_int(vtoken[0]);
						if (err_f) { cmd_msg = "expecting a number"; }
					}
					if (!err_f) {
						slot = strtol(vtoken[0].c_str(), NULL, 0);  // STRTOL!
						if (last_token == "stop_sys") {
							slot += 128;
							err_f = !slot_sys(slot);
						}
						else {
							err_f = slot_error(slot);
						}
						if (err_f) { cmd_msg = "slot out of range"; }
					}
					if (!err_f) { 
						string fnamex = file_ext_add(vtoken[1], DLP_EXT);  // add extension
						err_f = dlp_slot_rx(sp, slot, fnamex);
						if (err_f) { cmd_msg = "can't download from slot:" + to_string(slot); }
					}
					dir_proc(dir_all, dir_dlp, dir_view, view_f);  // update file list
					dir_idx = xy_move_idx(dir_idx, 0, 0, dir_view.size()-1, DIR_ROWS);
				}
				// slots => numbered preset files
				else if (last_token == "stop_num") {
					cmd_f = true;
					err_f = (tokens != 3); 
					if (err_f) { cmd_msg = "expecting slot start & end"; }
					if (!err_f) {  // check slot
						err_f = !str_is_int(vtoken[0]);
						if (err_f) { cmd_msg = "expecting a number"; }
						if (!err_f) { dl_min = strtol(vtoken[0].c_str(), NULL, 0); }  // STRTOL!
					}
					if (!err_f) {  // check slot
						err_f = !str_is_int(vtoken[1]);
						if (err_f) { cmd_msg = "expecting a number"; }
						if (!err_f) { dl_max = strtol(vtoken[1].c_str(), NULL, 0); }  // STRTOL!
					}
					if (!err_f) {  // swap slots
						int32_t tmp_min = min(dl_min, dl_max);
						dl_max = max(dl_min, dl_max);
						dl_min = tmp_min;
					}
					if (!err_f) {  // check slot limits
						err_f = slot_error(dl_min);
						err_f |= slot_error(dl_max);
						if (err_f) { cmd_msg = "slot(s) out of range"; }
					}
					cmd_f = true;
					int32_t slot = dl_min;
					while ((!err_f) && (slot <= dl_max)) { 
						char buf[5];
						sprintf (buf, "%04d", slot);  // format to 4 decimals
						string slot_num(buf);						
						string fnamex = file_ext_add(slot_num, DLP_EXT);  // add extension
						err_f = dlp_slot_rx(sp, slot, fnamex);
						if (err_f) { cmd_msg = "can't download from slot:" + slot_num; }
						if (!(slot % 4)) { cout << "." << flush; }  // show activity
						slot++;
					}
					dir_proc(dir_all, dir_dlp, dir_view, view_f);  // update file list
					dir_idx = xy_move_idx(dir_idx, 0, 0, dir_view.size()-1, DIR_ROWS);
				}
				// editor => slot
				else if ((last_token == "etos") || (last_token == "etos_sys")) {
					cmd_f = true;
					int32_t slot = 0;
					err_f = (tokens != 2); 
					if (err_f) { cmd_msg = "expecting slot"; }
					if (!err_f) {  // check slot limits
						if (!str_is_int(vtoken[0])) { err_f = true; }
						if (err_f) { cmd_msg = "expecting a number"; }
					}
					if (!err_f) {
						slot = strtol(vtoken[0].c_str(), NULL, 0);  // STRTOL!
						if (last_token == "etos_sys") {
							slot += 128;
							err_f = !slot_sys(slot);
						}
						else {
							err_f = slot_error(slot);
						}
						if (err_f) { cmd_msg = "slot out of range"; }
					}
					if (!err_f) { 
						err_f = ed_slot_tx(sp, slot);
						if (err_f) { cmd_msg = "can't upload to slot:" + to_string(slot); }
					}
				}
				// slot => editor
				else if ((last_token == "stoe") || (last_token == "stoe_sys")) {
					cmd_f = true;
					int32_t slot = 0;
					err_f = (tokens != 2); 
					if (err_f) { cmd_msg = "expecting slot"; }
					if (!err_f) {  // check slot limits
						err_f = !str_is_int(vtoken[0]);
						if (err_f) { cmd_msg = "expecting a number"; }
					}
					if (!err_f) {
						slot = strtol(vtoken[0].c_str(), NULL, 0);  // STRTOL!
						if (last_token == "stoe_sys") {
							slot += 128;
							err_f = !slot_sys(slot);
						}
						else {
							err_f = slot_error(slot);
						}
						if (err_f) { cmd_msg = "slot out of range"; }
					}
					if (!err_f) { 
						err_f = ed_slot_rx(sp, slot);
						if (err_f) { cmd_msg = "can't download from slot:" + to_string(slot); }
					}
					if (!err_f) { 
						mode = edit;  // edit mode
					}
				}
				// slots => list
				else if (last_token == "stol") {
					cmd_f = true;
					err_f = (tokens > 3);
					if (err_f) { cmd_msg = "expecting 0, 1, or 2 numbers"; }
					if (!err_f && (tokens > 1)) {  // 1 or 2 slots
						err_f = !str_is_int(vtoken[0]);
						if (err_f) { cmd_msg = "expecting a number"; }
						if (!err_f) { dl_min = strtol(vtoken[0].c_str(), NULL, 0); }  // STRTOL!
					}
					if (!err_f && (tokens == 3)) {  // 2 slots
						err_f = !str_is_int(vtoken[1]);
						if (err_f) { cmd_msg = "expecting a number"; }
						if (!err_f) { dl_max = strtol(vtoken[1].c_str(), NULL, 0); }  // STRTOL!
					}
					if (!err_f && (tokens == 2)) {  // 1 slot
						dl_max = dl_min;
					}
					if (!err_f && (tokens == 3)) {  // 2 slots
						int32_t tmp_min = min(dl_min, dl_max);
						dl_max = max(dl_min, dl_max);
						dl_min = tmp_min;
					}
					if (!err_f) {  // check slot limits
						err_f = slot_error(dl_min);
						err_f |= slot_error(dl_max);
						if (err_f) { cmd_msg = "slot(s) out of range"; }
					}
					if (!err_f) {  // dl slots & match contents with file names
						err_f = slots_sv_names_rx(sp, dl_min, dl_max, dir_dlp, dl_pre_sv, dl_names_sv);
						if (err_f) { cmd_msg = "can't download slot(s)"; }
					}
					if (!err_f) {
						mode = slots;  // slots mode
						dl_idx = 0;
					}
				}
				// list => bank file
				else if (last_token == "ltob") {
					cmd_f = true;
					err_f = (tokens != 2); 
					if (err_f) { cmd_msg = "expecting file"; }
					if (!err_f) { 
						string fnamex = file_ext_add(vtoken[0], BNK_EXT);  // add extension
						err_f = sv_to_file(fnamex, dl_names_sv, false);
						if (err_f) { cmd_msg = "can't write file:" + fnamex;}
					}
					dir_proc(dir_all, dir_dlp, dir_view, view_f);  // update file list
					dir_idx = xy_move_idx(dir_idx, 0, 0, dir_view.size()-1, DIR_ROWS);
				}
				// bank file => slots
				else if (last_token == "btos") {
					cmd_f = true;
					int32_t slot = 0;
					uint32_t idx = 0;
					vector<string> bnk_sv;
					err_f = (tokens != 3); 
					if (err_f) { cmd_msg = "expecting file & slot"; }
					if (!err_f) {  // check slot input
						if (!str_is_int(vtoken[1])) { err_f = true; }
						if (err_f) { cmd_msg = "expecting a number"; }
					}
					if (!err_f) {  // check slot range
						slot = strtol(vtoken[1].c_str(), NULL, 0);  // STRTOL!
						err_f = slot_error(slot);
						if (err_f) { cmd_msg = "slot out of range"; }
					}
					if (!err_f) {  // read bank file
						string rd_fnamex = file_ext_add(vtoken[0], BNK_EXT);  // add extension
						err_f = file_to_sv(rd_fnamex, bnk_sv, true);
						if (err_f) { cmd_msg = "can't read file:" + rd_fnamex;}
					}
					while ((!err_f) && (idx < bnk_sv.size()) && !slot_error(slot)) { 
						if (bnk_sv[idx].find(PRE_NO_MATCH) == string::npos) {  // if not " -?- "
							string dlp_fnamex = file_ext_add(bnk_sv[idx], DLP_EXT);  // add extension
							err_f = dlp_slot_tx(sp, slot, dlp_fnamex);
							if (err_f) { cmd_msg = "can't upload file:" + dlp_fnamex + " to slot:" + to_string(slot); }
						}
						idx++;
						if (!(idx % 2)) { cout << "." << flush; }  // show activity
						if (slot < 0) { slot--; } // go neg if starting neg
						else { slot++; }
					}
				}
				// eeprom => file
				else if ((last_token == "spi_dump") || (last_token == "eeprom_dump") || (last_token == "pre_dump")  || (last_token == "pro_dump")) {
					cmd_f = true;
					err_f = (tokens != 2); 
					if (err_f) { cmd_msg = "expecting a filename"; }
					if (!err_f) {
						string fext = SPI_EXT;  // spi_dump defaults
						uint32_t addr = EEPROM_SW_ADDR; 
						uint32_t bytes = EEPROM_SW_BYTES;
						if (last_token == "eeprom_dump") {
							fext = EEPROM_EXT;
							addr = 0;
							bytes = EEPROM_SW_ADDR + EEPROM_SW_BYTES;
						}
						else if (last_token == "pre_dump") {
							fext = PRE_EXT;
							addr = EEPROM_PRE_ADDR;
							bytes = EEPROM_PRE_BYTES;
						}
						else if (last_token == "pro_dump") {
							fext = PRO_EXT;
							addr = EEPROM_PRO_ADDR;
							bytes = EEPROM_PRO_BYTES;
						}
						string fnamex = file_ext_add(vtoken[0], fext);  // add extension
						err_f = dump_rx(sp, fnamex, addr, bytes);  // do dump
						if (err_f) { cmd_msg = "can't download to file:" + fnamex; }
					}
					dir_proc(dir_all, dir_dlp, dir_view, view_f);  // update file list
					dir_idx = xy_move_idx(dir_idx, 0, 0, dir_view.size()-1, DIR_ROWS);
				}
				// file => eeprom
				else if ((last_token == "spi_pump") || (last_token == "eeprom_pump") || (last_token == "pre_pump") || (last_token == "pro_pump")) {
					cmd_f = true;
					err_f = (tokens != 2); 
					if (err_f) { cmd_msg = "expecting a filename"; }
					if (!err_f) {
						string fext = SPI_EXT;  // spi_pump defaults
						uint32_t addr = EEPROM_SW_ADDR;
						bool rst = true;
						if (last_token == "eeprom_pump") {
							fext = EEPROM_EXT;
							addr = 0;
							rst = true;
						}
						else if (last_token == "pre_pump") {
							fext = PRE_EXT;
							addr = EEPROM_PRE_ADDR;
							rst = false;
						}
						else if (last_token == "pro_pump") {
							fext = PRO_EXT;
							addr = EEPROM_PRO_ADDR;
							rst = false;
						}
						string fnamex = file_ext_add(vtoken[0], fext);  // add extension
						err_f = pump_tx(sp, fnamex, addr, rst);  // do pump
						if (err_f) { cmd_msg = "can't upload file:" + fnamex; }
					}
				}
/*
				// preset file update
				else if (last_token == "up") {
					cmd_f = true;
					err_f = (tokens != 2);
					if (err_f) { cmd_msg = "expecting 1 filename"; }
					if (!err_f) {
						string fnamex = file_ext_add(vtoken[0], DLP_EXT);
						err_f = dlp_update(fnamex);
						if (err_f) { cmd_msg = "can't update file:" + fnamex; }
					}
				}
*/
				// do command line stuff
				if (cmd_f) { 
					cmd_prev = cmd_str;  // keep old
					cmd_hist.wr(cmd_str);  // store & clear cmd
					cursor_pos = cmd_str.size();  // reset cursor
					if (err_f) { cmd_msg = "<?> " + cmd_msg; }  
					else { cmd_msg = "<OK> " + cmd_msg; }
					cmd_prev += cmd_msg;  // keep old cmd
					cmd_msg = "";  // clear msg
					hl_f = false; // clear hilite
					con_cls();  // clear screen
				}
			}
		}

		//////////
		// draw //
		//////////

		// editor view
		if (mode == edit) {
			con_xy(0, 0);
			ed_draw(ed_x, ed_y, hl_f);
		}

		// directory view
		if (mode == files) {
			con_xy(0, 0);
			dir_draw(dir_idx, hl_f, dir_view);
		}

		// slots view
		if (mode == slots) {
			con_xy(0, 0);
			slots_draw(dl_idx, hl_f, dl_min, dl_names_sv);
		}

		// draw preset file name
		con_font('b');  // bold
		con_xy(0, FILE_Y);
		cout << setw(CMD_W - 40) << left << ("preset: " + pre_fname);
		cout << setw(40) << right << ("print: " + prn_fname);
		cout << endl << setw(CMD_W) << "";

		// draw cmd line
		if (err_f) { con_font('e'); }  // hilite errors
		else { con_font('n'); }  // normal
		con_xy(0, PREV_Y); cout << setw(CMD_W) << left << cmd_prev;
		con_font('d');  // default
		con_xy(0, CMD_Y); cout << setw(CMD_W) << left << cmd_str;
		con_xy(0+cursor_pos, CMD_Y);  // position cursor on cmd line

		// clear flag
		first_f = false;

	}


} // end of main

