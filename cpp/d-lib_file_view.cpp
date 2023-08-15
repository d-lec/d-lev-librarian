// d-lib_file_view.cpp
// text file viewer
#pragma once

// view string vector
void sv_view(const vector<string>& sv_i, int32_t conw_i, int32_t conh_i) {
	int32_t line_min = 0;
	int32_t view_h = conh_i - 2;
	int32_t lines = sv_i.size();
	int32_t num_w = log10(lines) + 1;
	bool num_f = false;
	bool done_f = false;
	while (!done_f) {
		con_cls();  // clear screen
		for (int32_t line = line_min; line < line_min+view_h; line++) {
			if (line < lines) { 
				if (num_f) { cout << setw(num_w) << right << line + 1 << " "; }
				cout << sv_i[line]; 
			}
			cout << endl;
		}
		cout << endl << "- UP, DN, PGUP, PGDN, HOME, END, n - ESC to quit -";
		int64_t key = kbd_getkey();  // blocking
		if (key == KBD_ESC) { done_f = true; }
		else if (key == KBD_UP) { line_min--; }
		else if (key == KBD_DN) { line_min++; }
		else if (key == KBD_PGUP) { line_min -= view_h; }
		else if (key == KBD_PGDN) { line_min += view_h; }
		else if (key == KBD_HOME) { line_min = 0; }
		else if (key == KBD_END) { line_min = lines - view_h; }
		else if (key == 'n') { num_f = !num_f; }
		if (line_min < 0) { line_min = 0; }
		if (line_min > lines - view_h) { line_min = lines - view_h; }
		if (line_min < 0) { line_min = 0; }
	}
	return;
}

// view text file
// return true if error
bool file_view(const string& fname_i, int32_t conw_i, int32_t conh_i) {
	vector<string> rd_sv;
	bool error_f = file_to_sv(fname_i, rd_sv, false); 
	if (!error_f) { sv_view(rd_sv, conw_i, conh_i); }
	return(error_f);
}
