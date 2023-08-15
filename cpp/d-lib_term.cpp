// d-lib_term.cpp
// micro terminal
#pragma once

// interactive serial port terminal
void term(const sp_type& sp_i) {
	bool done_f = false;
	int64_t key = 0;
	string tx_str = "";
	string rx_str = "";
	con_cls();
	cout << "- TERMINAL MODE - ESC to quit -" << endl;
	while (!done_f) {
		ms_sleep(RX_POLL_MDLY);  // delay for polling
		// rx
		rx_str = "";  // clear
		sp_rx(sp_i, rx_str);
		cout << flush << rx_str;
		// tx
		tx_str = "";  // clear
		key = kbd_rdbuf();
		if (key) { 
			if (key == 0x7f) { key = 0x08; }  // xlate del => bs
			tx_str = key; 
		}
		// exit conditions / tx
		if (key == KBD_ESC) { done_f = true; }
		else { sp_tx(sp_i, tx_str); }
	}
}
