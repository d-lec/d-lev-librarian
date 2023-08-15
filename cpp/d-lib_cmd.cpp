// d-lib_cmd.cpp
#pragma once

/////////////
// DEFINES //
/////////////

#ifdef _MSWIN
	// MSWin key codes //
	// concatenated byte key codes (printable ASCII: 0x20 thru 0x7e)
	#define KBD_TAB			0x9
	#define KBD_ENT			0xd
	#define KBD_ESC			0x1b
	#define KBD_BKSP		0x8
	//
	#define KBD_UP			0x48e0
	#define KBD_DN			0x50e0
	#define KBD_RT			0x4de0
	#define KBD_LT			0x4be0
	#define KBD_PGUP		0x49e0
	#define KBD_PGDN		0x51e0
	#define KBD_HOME		0x47e0
	#define KBD_END			0x4fe0
	#define KBD_INS			0x52e0
	#define KBD_DEL			0x53e0
	#define KBD_F1			0x3b00
	#define KBD_F2			0x3c00
	#define KBD_F3			0x3d00
	#define KBD_F4			0x3e00
	#define KBD_F5			0x3f00
	#define KBD_F6			0x4000
	#define KBD_F7			0x4100
	#define KBD_F8			0x4200
	#define KBD_F9			0x4300
	#define KBD_F10			0x4400
	#define KBD_F11			0x8500
	#define KBD_F12			0x8600
	//
	#define KBD_CTRL_UP		0x8de0
	#define KBD_CTRL_DN		0x91e0
	#define KBD_CTRL_RT		0x74e0
	#define KBD_CTRL_LT		0x73e0
	#define KBD_CTRL_PGUP	0x86e0
	#define KBD_CTRL_PGDN	0x76e0
	#define KBD_CTRL_HOME	0x77e0
	#define KBD_CTRL_END	0x75e0
	#define KBD_CTRL_DEL	0x93e0
	#define KBD_CTRL_F1		0x5e00
	#define KBD_CTRL_F2		0x5f00
	#define KBD_CTRL_F3		0x6000
	#define KBD_CTRL_F4		0x6100
	#define KBD_CTRL_F5		0x6200
	#define KBD_CTRL_F6		0x6300
	#define KBD_CTRL_F7		0x6400
	#define KBD_CTRL_F8		0x6500
	#define KBD_CTRL_F9		0x6600
	#define KBD_CTRL_F10	0x6700
	#define KBD_CTRL_F11	0x8900
	#define KBD_CTRL_F12	0x8a00
	//
	#define KBD_SHFT_F1		0x5400
	#define KBD_SHFT_F2		0x5500
	#define KBD_SHFT_F3		0x5600
	#define KBD_SHFT_F4		0x5700
	#define KBD_SHFT_F5		0x5800
	#define KBD_SHFT_F6		0x5900
	#define KBD_SHFT_F7		0x5a00
	#define KBD_SHFT_F8		0x5b00
	#define KBD_SHFT_F9		0x5c00
	#define KBD_SHFT_F10	0x5d00
	#define KBD_SHFT_F11	0x8700
	#define KBD_SHFT_F12	0x8800
	//
	#define KBD_ALT_UP		0x9800
	#define KBD_ALT_DN		0xa000
	#define KBD_ALT_RT		0x9d00
	#define KBD_ALT_LT		0x9b00
	#define KBD_ALT_PGUP	0x9900
	#define KBD_ALT_PGDN	0xa100
	#define KBD_ALT_HOME	0x9700
	#define KBD_ALT_END		0x9f00
	#define KBD_ALT_INS		0xa200
	#define KBD_ALT_DEL		0xa300
#else
	// POSIX key codes //
	// concatenated byte key codes (printable ASCII: 0x20 thru 0x7e)
	#define KBD_TAB			0x9
	#define KBD_ENT			0xa
	#define KBD_ESC			0x1b
	#define KBD_BKSP		0x7f
	//
	#define KBD_UP			0x1b5b41
	#define KBD_DN			0x1b5b42
	#define KBD_RT			0x1b5b43
	#define KBD_LT			0x1b5b44
	#define KBD_PGUP		0x1b5b357e
	#define KBD_PGDN		0x1b5b367e
	#define KBD_HOME		0x1b5b48
	#define KBD_END			0x1b5b46
	#define KBD_INS			0x1b5b327e
	#define KBD_DEL			0x1b5b337e
	#define KBD_F1			0x1b4f50
	#define KBD_F2			0x1b4f51
	#define KBD_F3			0x1b4f52
	#define KBD_F4			0x1b4f53
	#define KBD_F5			0x1b5b31357e
	#define KBD_F6			0x1b5b31377e
	#define KBD_F7			0x1b5b31387e
	#define KBD_F8			0x1b5b31397e
	#define KBD_F9			0x1b5b32307e
	#define KBD_F12			0x1b5b32347e
	//
	#define KBD_CTRL_UP		0x1b5b313b3541
	#define KBD_CTRL_DN		0x1b5b313b3542
	#define KBD_CTRL_RT		0x1b5b313b3543
	#define KBD_CTRL_LT		0x1b5b313b3544
	#define KBD_CTRL_PGUP	0x1b5b353b357e
	#define KBD_CTRL_PGDN	0x1b5b363b357e
	#define KBD_CTRL_HOME	0x1b5b313b3548
	#define KBD_CTRL_END	0x1b5b313b3546
	#define KBD_CTRL_DEL	0x1b5b333b357e
	#define KBD_CTRL_F1		0x1b5b313b3550
	#define KBD_CTRL_F2		0x1b5b313b3551
	#define KBD_CTRL_F3		0x1b5b313b3552
	#define KBD_CTRL_F4		0x1b5b313b3553
	#define KBD_CTRL_F5		0x1b5b31353b357e
	#define KBD_CTRL_F6		0x1b5b31373b357e
	#define KBD_CTRL_F7		0x1b5b31383b357e
	#define KBD_CTRL_F8		0x1b5b31393b357e
	#define KBD_CTRL_F9		0x1b5b32303b357e
	#define KBD_CTRL_F12	0x1b5b32343b357e
	//
	#define KBD_SHFT_F1		0x1b5b313b3250
	#define KBD_SHFT_F2		0x1b5b313b3251
	#define KBD_SHFT_F3		0x1b5b313b3252
	#define KBD_SHFT_F4		0x1b5b313b3253
	#define KBD_SHFT_F5		0x1b5b31353b327e
	#define KBD_SHFT_F6		0x1b5b31373b327e
	#define KBD_SHFT_F7		0x1b5b31383b327e
	#define KBD_SHFT_F8		0x1b5b31393b327e
	#define KBD_SHFT_F9		0x1b5b32303b327e
	#define KBD_SHFT_F12	0x1b5b32343b327e
	//
	#define KBD_ALT_UP		0x1b5b313b3341
	#define KBD_ALT_DN		0x1b5b313b3342
	#define KBD_ALT_RT		0x1b5b313b3343
	#define KBD_ALT_LT		0x1b5b313b3344
	#define KBD_ALT_PGUP	0x1b5b353b337e
	#define KBD_ALT_PGDN	0x1b5b363b337e
	#define KBD_ALT_HOME	0x1b5b313b3348
	#define KBD_ALT_END		0x1b5b313b3346
	#define KBD_ALT_INS		0x1b5b323b337e
	#define KBD_ALT_DEL		0x1b5b333b337e
#endif


// edit command line, return true if key_i used.
// enter inserts space @ end and puts cursor there.
// tab & other white space ignored.
// ctrl l&r arrows postion cursor at next word or beginning / end.
bool cmd_edit(string& cmd_str_io, uint32_t& cursor_pos_io, uint64_t key_i) {
	bool rtn_val = true;  // default true
	// control chars:
	if (key_i == KBD_ESC) {
		cmd_str_io.clear();
		cursor_pos_io = 0;
	}
	else if (key_i == KBD_LT) {
		if (cursor_pos_io) { cursor_pos_io--; }
	}
	else if (key_i == KBD_RT) {
		if (cursor_pos_io < cmd_str_io.size()) { cursor_pos_io++; }
	}
	else if (key_i == KBD_HOME) {
		cursor_pos_io = 0;
	}
	else if (key_i == KBD_END) {
		cursor_pos_io = cmd_str_io.size();
	}
	else if (key_i == KBD_BKSP) {
		if (cursor_pos_io) { cmd_str_io.erase(--cursor_pos_io, 1); }
	}
	else if (key_i == KBD_DEL) {
		if (cursor_pos_io < cmd_str_io.size()) { cmd_str_io.erase(cursor_pos_io, 1); }
	}
	else if (key_i == KBD_ENT) {
		cmd_str_io += " ";
		cursor_pos_io = cmd_str_io.size();
	}
	else if (key_i == KBD_CTRL_RT) {
		bool done_f = false;
		while (!done_f) {
			if (cursor_pos_io < cmd_str_io.size()) {  // not @ end
				cursor_pos_io++;  // inc
				if ((cmd_str_io[cursor_pos_io-1] == ' ') && (cmd_str_io[cursor_pos_io] != ' ')) { done_f = true; }
			}
			else { done_f = true; }  // @ end
		}
	}
	else if (key_i == KBD_CTRL_LT) {
		bool done_f = false;
		while (!done_f) {
			if (cursor_pos_io > 1) {  // not @ start
				cursor_pos_io--;  // dec
				if ((cmd_str_io[cursor_pos_io-1] == ' ') && (cmd_str_io[cursor_pos_io] != ' ')) { done_f = true; }
			}
			else {  // @ start
				cursor_pos_io = 0;
				done_f = true; 
			} 
		}
	}
	// non-white or space char:
	else if ((ch_is_nonwhite(key_i)) || (key_i == ' ')) { 
		cmd_str_io.insert(cursor_pos_io++, 1, char(key_i)); 
	}
	else { rtn_val = false; }
	return(rtn_val);
}


// command line history circular buffer
// inc write pointer for each write, zero out read pointer & input command string
// inc / dec read pointer for each read (limited to 0 and writes depth).
template <uint32_t CMD_DEPTH>  // bs just to set a constant
class cmd_mem {
private:

	// state and such
	string mem[CMD_DEPTH];
	uint32_t wr_ptr = 0;
	uint32_t rd_ptr = 0;
	uint32_t writes = 0;

	// level & modulo (safe) indices
	uint32_t level() { return(wr_ptr - rd_ptr); }
	uint32_t rd_idx() { return(rd_ptr % CMD_DEPTH); }
	uint32_t wr_idx() { return(wr_ptr % CMD_DEPTH); }

public:

	// comand string => mem, clear cmd
	void wr(string& cmd_str_io) { 
		str_trim(cmd_str_io);  // trim whitespace at ends
		mem[wr_idx()] = cmd_str_io;  // sto cmd
		wr_ptr++;  // inc wr ptr
		rd_ptr = wr_ptr;  // reset rd ptr
		if (writes < CMD_DEPTH-1) { writes++; }  // inc write count
		mem[wr_idx()] = "";  // clear next slot
		cmd_str_io = mem[rd_idx()];  // clear cmd_str
	}

	// mem => comand string
	bool rd(string& cmd_str_io, bool deeper_f) { 
		bool rd_f = false;
		if (deeper_f) {  // go deeper
			if (level() < writes) { rd_ptr--; rd_f = true; }
		}
		else {  // come back
			if (level() > 0) { rd_ptr++; rd_f = true; }
		}
		if (rd_f) { cmd_str_io = mem[rd_idx()]; }
		return(rd_f);
	}
};

