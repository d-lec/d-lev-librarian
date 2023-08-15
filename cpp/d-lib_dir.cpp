// d-lib_dir.cpp
#pragma once


/////////////
// DEFINES //
/////////////

#define DIR_ROWS 25
#define DIR_COLS 4
#define DIR_W 21

#define SLOTS_ROWS 20
#define SLOTS_COLS 4
#define SLOTS_W 16

// dump directory file listing to stringvec
// this code from: http://www.martinbroadhurst.com/list-the-files-in-a-directory-in-c.html
#ifdef _MSWIN
	#include <windows.h>
	void dir_rd(const string& dir_i, vector<string>& sv_o) {
		string pattern(dir_i);
		pattern.append("\\*");
		WIN32_FIND_DATA data;
		HANDLE hFind;
		if ((hFind = FindFirstFile(pattern.c_str(), &data)) != INVALID_HANDLE_VALUE) {
			do {
				sv_o.push_back(data.cFileName);
			} while (FindNextFile(hFind, &data) != 0);
			FindClose(hFind);
		}
	}
#else
	#include <dirent.h>
	void dir_rd(const string& dir_i, vector<string>& sv_o){
		DIR* dirp = opendir(dir_i.c_str());
		struct dirent * dp;
		while ((dp = readdir(dirp)) != NULL) {
			sv_o.push_back(dp->d_name);
		}
		closedir(dirp);
	}
#endif

// alpha sort file names
void dir_sort(vector<string>& dir_io) {
	sort(dir_io.begin(), dir_io.end());
}

// get file name extension
string file_ext_rd(const string& fname_i)
{
	size_t dot_pos = fname_i.find_last_of(".");
    if(dot_pos != string::npos) { return(fname_i.substr(++dot_pos)); }
    else { return(""); }
}

// kill file name extension
string file_ext_kill(const string& fname_i)
{
	size_t dot_pos = fname_i.find_last_of(".");
    if(dot_pos != string::npos) { return(fname_i.substr(0, dot_pos)); }
    else { return(fname_i); }
}

// add file name extension (if it doesn't have an extension already)
string file_ext_add(const string& fname_i, const string ext_i)
{
	if (file_ext_rd(fname_i) == "") { return(fname_i + "." + ext_i); }
	else { return(fname_i); }
}

// return only file names with given extension
vector<string> dir_ext_filter(const vector<string>& dir_i, const string ext_i) {
	vector<string> sv;
	for (uint32_t i=0; i < dir_i.size(); i++) {
		if (file_ext_rd(dir_i[i]) == ext_i) { sv.push_back(dir_i[i]); }
	}
	return(sv);
}

// kill all file name extensions
vector<string> dir_ext_kill(const vector<string>& dir_i) {
	vector<string> sv;
	for (uint32_t i=0; i < dir_i.size(); i++) {
		sv.push_back(file_ext_kill(dir_i[i]));
	}
	return(sv);
}

// autocomplete input string
string dir_acomp(const string str_i, const vector<string>& dir_i) {
	string str = str_i;  // set to input
	// check for empty
	uint32_t len = str.size();
	if (!len) { return(str); }  // empty
	// check for no and single match
	uint32_t idx = 0;
	uint32_t matches = sv_str_match(str, dir_i, len, idx);
	if (!matches) { return(str); }  // no match
	str = dir_i[idx];  // set to match
	if (matches == 1) { return(str); }  // single match
	// do multiple matches
	uint32_t matches_old = 0;
	do {
		len++;
		matches_old = matches;
		matches = sv_str_match(str, dir_i, len, idx);
	} while (matches == matches_old);
	return(str.substr(0, len-1)); // safe b/c len>1
}

// autocomplete input string - return just the delta
string dir_acomp_delta(const string str_i, const vector<string>& dir_i) {
	string str = dir_acomp(str_i, dir_i);
	return(str.substr(str_i.size(), str.size() - str_i.size()));
}

// return last word in string
string str_last_word(const string& str_i) {
	string word;
	istringstream( { str_i.rbegin(), str_i.rend() } ) >> word;
	return { word.rbegin(), word.rend() };
}

// command line autocomplete
void cmd_acomp(string& cmd_str_io, uint32_t& cursor_pos_io, const vector<string>& dir_i) {
	string cmd_str_edit = cmd_str_io.substr(0, cursor_pos_io);
	string cmd_str_last = cmd_str_io.substr(cursor_pos_io--);
	string cmd_str_word = str_last_word(cmd_str_edit);
	string cmd_str_delta = dir_acomp_delta(cmd_str_word, dir_i);
	cmd_str_io = cmd_str_edit;
	cmd_str_io += cmd_str_delta;
	cmd_str_io += cmd_str_last;
	cursor_pos_io += cmd_str_delta.size() + 1;
}

// process directory
// read, alpha sort, etc.
// dir_view_io : dir_all_io or (*.eeprom, *.spi, *.pre, *.pro, *.bnk, *.dlp (w/o ext))
void dir_proc(vector<string>& dir_all_io, vector<string>& dir_dlp_io, vector<string>& dir_view_io, bool all_f) {
	dir_all_io.clear();
	dir_rd(".", dir_all_io);  // read directory listing
	dir_sort(dir_all_io);  // alpha sort directory listing
	// *.dlp processing
	dir_dlp_io = dir_ext_filter(dir_all_io, DLP_EXT);  // filter out if not *.dlp
	dir_dlp_io = dir_ext_kill(dir_dlp_io);  // kill file extensions
	// generate views
	if (all_f) { dir_view_io = dir_all_io; }
	else { 
		dir_view_io = dir_ext_filter(dir_all_io, EEPROM_EXT);  // filter out if not *.eeprom
		vector<string> dir_spi = dir_ext_filter(dir_all_io, SPI_EXT);  // filter out if not *.spi
		vector<string> dir_pre = dir_ext_filter(dir_all_io, PRE_EXT);  // filter out if not *.pre
		vector<string> dir_pro = dir_ext_filter(dir_all_io, PRO_EXT);  // filter out if not *.pro
		vector<string> dir_bnk = dir_ext_filter(dir_all_io, BNK_EXT);  // filter out if not *.bnk
		dir_view_io.insert(dir_view_io.end(), dir_spi.begin(), dir_spi.end());  // append *.spi
		dir_view_io.insert(dir_view_io.end(), dir_pre.begin(), dir_pre.end());  // append *.pre
		dir_view_io.insert(dir_view_io.end(), dir_pro.begin(), dir_pro.end());  // append *.pro
		dir_view_io.insert(dir_view_io.end(), dir_bnk.begin(), dir_bnk.end());  // append *.bnk
		dir_view_io.insert(dir_view_io.end(), dir_dlp_io.begin(), dir_dlp_io.end());  // append *.dlp (w/o ext)
	}
}

// return idx given x & y
int32_t xy_to_idx(int32_t x_i, int32_t y_i, int32_t rows_i) {
	return(y_i + (x_i * rows_i));
}

// return x given idx
int32_t idx_to_x(int32_t idx_i, int32_t rows_i) {
	return(div_int(idx_i, rows_i));
}

// return y given idx
int32_t idx_to_y(int32_t idx_i, int32_t rows_i) {
	return(mod_int(idx_i, rows_i));
}

// return new idx given current idx and delta x & y
// stop at individual x & y limits
int32_t xy_move_idx(int32_t idx_i, int32_t dx_i, int32_t dy_i, int32_t idx_max_i, int32_t rows_i) {
	int32_t x_max = div_int(idx_max_i, rows_i);
	int32_t x = idx_to_x(idx_i, rows_i);
	int32_t y = idx_to_y(idx_i, rows_i);
	// offset
	x += dx_i;
	y += dy_i;
	// limit
	if (x < 0) { x = 0; }
	if (x > x_max) { x = x_max; }
	if (y < 0) { y = 0; }
	if (y >= rows_i) { y = rows_i - 1; }
	int32_t idx = xy_to_idx(x, y, rows_i);
	if (idx > idx_max_i) { idx = idx_max_i; }
	if (idx < 0) { idx = 0; }
	return(idx);
}

// draw files listing
void dir_draw(int32_t idx_i, bool ed_i, const vector<string>& dir_i) {
	static int32_t x_min = 0;  // first col to display, static for hysteresis
	int32_t last_col = (dir_i.size() / DIR_ROWS);  // last col (total)
	int32_t x_max = min(x_min + DIR_COLS - 1, last_col);  // last col to display
	int32_t x_idx = idx_to_x(idx_i, DIR_ROWS);
	// drawing limits set by hilite location
	if (x_idx < x_min) { 
		x_min = x_idx; 
		x_max = min(x_min + DIR_COLS - 1, last_col);
	}
	if (x_idx > x_max) { 
		x_max = x_idx; 
		x_min = max(0, x_max - DIR_COLS + 1);
	}
	// draw
	con_font('n');
	for (int32_t y=0; y<DIR_ROWS; y++) {
		for (int32_t x=x_min; x<=x_max; x++) {
			int32_t idx = xy_to_idx(x, y, DIR_ROWS);
			if (uint32_t(idx) < dir_i.size()) {  // draw entry
				if (ed_i && (idx == idx_i)) { con_font('s'); }
				cout << setw(DIR_W) << left << dir_i[idx];
				con_font('n');
			}
			else { cout << setw(DIR_W) << ""; }  // else draw blank
			if (x == x_max) { cout << endl; }  // new line
		}
	}
	con_font('o'); cout << setw(49) << right << "F2 - FILES" << setw(36) << ""; 
	con_font('n'); 
	cout << endl;  // ending divider
}

// print files listing to string
string dir_print_to_str(const vector<string>& dir_i) {
	stringstream ss;
	for (uint32_t idx=0; idx<dir_i.size(); idx++) {
		ss << dir_i[idx] << endl;
	}
	ss << endl;
	return(ss.str());
}

// print files screen => txt file
// uses str_to_file()
bool dir_print_to_file(const string& tfname_i, const vector<string>& dir_i, bool append_i) {
	string str = dir_print_to_str(dir_i);
	return(str_to_file(tfname_i, str, append_i));
}

// delete preset file @ idx
bool dlp_del(uint32_t idx_i, const vector<string>& dir_i) {
	string filen = dir_i[idx_i] + "." + DLP_EXT;
	return(file_delete(filen));
}

// draw slots listing
void slots_draw(int32_t idx_i, bool ed_i, int32_t min_i, const vector<string>& names_sv_i) {
	static int32_t x_min = 0;  // first col to display, static for hysteresis
	int32_t last_col = (names_sv_i.size() / SLOTS_ROWS);  // last col (total)
	int32_t x_max = min(x_min + SLOTS_COLS - 1, last_col);  // last col to display
	int32_t x_idx = idx_to_x(idx_i, SLOTS_ROWS);
	// drawing limits set by hilite location
	if (x_idx < x_min) { 
		x_min = x_idx; 
		x_max = min(x_min + SLOTS_COLS - 1, last_col);
	}
	if (x_idx > x_max) { 
		x_max = x_idx; 
		x_min = max(0, x_max - SLOTS_COLS + 1);
	}
	// draw
	con_font('n');
	for (int32_t y=0; y<SLOTS_ROWS; y++) {
		for (int32_t x=x_min; x<=x_max; x++) {
			int32_t idx = xy_to_idx(x, y, SLOTS_ROWS);
			if (uint32_t(idx) < names_sv_i.size()) {
				int32_t slot = min_i + idx;
				if (ed_i && (idx == idx_i)) { con_font('s'); }
				cout << setw(4) << right << slot << " " << setw(SLOTS_W) << left << names_sv_i[idx];
				con_font('n');
			}
			else { cout << setw(5 + SLOTS_W) << ""; }  // else draw blank
			if (x == x_max) { cout << endl; }  // new line
		}
	}
	con_font('o'); cout << setw(49) << right << "F3 - SLOTS" << setw(36) << ""; 
	con_font('n'); 
	cout << endl;  // ending divider
}

// print slots listing to string
string slots_print_to_str(int32_t min_i, const vector<string>& names_sv_i) {
	uint32_t last_col = names_sv_i.size() / SLOTS_ROWS;
	stringstream ss;
	for (uint32_t y=0; y<SLOTS_ROWS; y++) {
		for (uint32_t x=0; x<=last_col; x++) {
			uint32_t idx = xy_to_idx(x, y, SLOTS_ROWS);
			if (idx < names_sv_i.size()) {
				int32_t slot = min_i + idx;
				ss << setw(4) << right << slot << " " << setw(SLOTS_W) << left << names_sv_i[idx];
			}
			if (x == last_col) { ss << endl; }
		}
	}
	ss << endl;
	return(ss.str());
}

// print slots screen => txt file
// uses str_to_file()
bool slots_print_to_file(const string& tfname_i, int32_t min_i, const vector<string>& slot_names_i, bool append_i) {
	string str = slots_print_to_str(min_i, slot_names_i);
	return(str_to_file(tfname_i, str, append_i));
}

