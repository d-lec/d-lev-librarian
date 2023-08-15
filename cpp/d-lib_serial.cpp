// d-lib_serial.cpp
#pragma once

#ifdef _MSWIN

//////////////////
// MS WIN START //
//////////////////

// serial port struct
struct sp_type {
	string dev = "";
	DCB dcb_old, dcb_new;
	COMMTIMEOUTS timeouts_old, timeouts_new;
	HANDLE fd = 0;  // handle
};

// return true if port open error
bool sp_open_err(const sp_type& sp_io) { 
	return(sp_io.fd == INVALID_HANDLE_VALUE);
}

// open port
void sp_open_port(sp_type& sp_io, const string& dev_i) {
	sp_io.dev = dev_i;  // keep dev
	sp_io.fd = CreateFile(
		sp_io.dev.c_str(),  //port name
		GENERIC_READ | GENERIC_WRITE,  //Read/Write
		0,  // No Sharing
		NULL,  // No Security
		OPEN_EXISTING,  // Open existing port only
		0,  // Non Overlapped I/O
		NULL  // Null for Comm Devices
	);
}	

// close port (true=error)
bool sp_close_port(sp_type& sp_io) {
	// restore old attributes
	bool rtn_val = !SetCommState(sp_io.fd, &sp_io.dcb_old);
	// restore old timeouts
	rtn_val |= !SetCommTimeouts(sp_io.fd, &sp_io.timeouts_old);
	// close port
	CloseHandle(sp_io.fd);
	return(rtn_val);
}

// return true if port is open
bool sp_test_port(const string& dev_i) {
	sp_type sp;
	sp_open_port(sp, dev_i);
	bool open_f = !sp_open_err(sp);
	if (open_f) { CloseHandle(sp.fd); }
	return(open_f);
}

// config port, return true if error
bool sp_config_port(sp_type& sp_io) {
	// save old dcb
	bool rtn_val = !GetCommState(sp_io.fd, &sp_io.dcb_old);
	if (rtn_val) { return(rtn_val); }
	// save old timeouts
	rtn_val = !GetCommTimeouts(sp_io.fd, &sp_io.timeouts_old);
	if (rtn_val) { return(rtn_val); }
	// zero out new dcb
	FillMemory(&sp_io.dcb_new, sizeof(sp_io.dcb_new), 0);
	sp_io.dcb_new.DCBlength = sizeof(sp_io.dcb_new);
	// build new dcb
	rtn_val = !BuildCommDCB("230400,n,8,1", &sp_io.dcb_new);
	if (rtn_val) { return(rtn_val); }
	// set new dcb
	rtn_val = !SetCommState(sp_io.fd, &sp_io.dcb_new);
	if (rtn_val) { return(rtn_val); }
	// zero out new timeouts
	sp_io.timeouts_new = { 0 };
	// build new timeouts
	sp_io.timeouts_new.ReadIntervalTimeout = MAXDWORD; 
	sp_io.timeouts_new.ReadTotalTimeoutMultiplier = 0;
	sp_io.timeouts_new.ReadTotalTimeoutConstant = 0;
	sp_io.timeouts_new.WriteTotalTimeoutMultiplier = 0;
	sp_io.timeouts_new.WriteTotalTimeoutConstant = 0;
	// set new timeouts
	rtn_val = !SetCommTimeouts(sp_io.fd, &sp_io.timeouts_new);
	if (rtn_val) { return(rtn_val); }
	return(rtn_val);
}

// str => tx port, return byte count
DWORD sp_tx(const sp_type& sp_i, const string& str_i) {
	DWORD tx_bytes = 0;
	WriteFile(sp_i.fd, str_i.c_str(), str_i.size(), &tx_bytes, NULL);
	return(tx_bytes);
}

// rx port => str, return byte count
// note: concats to string w/o clear
// note: you must externally poll
DWORD sp_rx(const sp_type& sp_i, string& str_o) {
	char rx_buf[RX_BUF_SZ];
	DWORD rx_bytes = 0;
	ReadFile(sp_i.fd, rx_buf, RX_BUF_SZ, &rx_bytes, NULL);
	str_o += string(rx_buf, rx_bytes);
	return(rx_bytes);
}

// rx port flush
void sp_rx_flush(const sp_type& sp_i) {
	ms_sleep(RX_BURST_MDLY);  // delay
	string rx_str;
	sp_rx(sp_i, rx_str);
	return;
}

// rx port => str, return byte count
// note: concats to string w/o clear
// internal polling for big bursts
// safety timeout
DWORD sp_rx_burst(const sp_type& sp_i, string& str_o) {
	DWORD rx_bytes = 0;
	DWORD rx_bytes_total = 0;
	int32_t rx_count = RX_MAX;
	do {
		ms_sleep(RX_BURST_MDLY);  // delay
		rx_bytes = sp_rx(sp_i, str_o);
		rx_bytes_total += rx_bytes;  // sum
		rx_count--;
		} while ((!rx_bytes_total || rx_bytes) && rx_count > 0);
	return(rx_bytes_total);
}

// rx port wait for '>' at end
// internal polling
// safety timeout
bool sp_rx_wait(const sp_type& sp_i) {
	string rx_str;
	int32_t rx_bytes = 0;
	int32_t rx_count = 0;
	do { 
		ms_sleep(RX_POLL_MDLY);  // delay
		rx_bytes = sp_rx(sp_i, rx_str);
		rx_count++;
	} while ((!rx_bytes || (rx_str.back() != '>')) && rx_count < RX_MAX);
	return(!(rx_count < RX_MAX));
}

////////////////
// MS WIN END //
////////////////




#else

/////////////////
// POSIX START //
/////////////////

// serial port struct
struct sp_type {
	string dev = "";
	struct termios t_old, t_new;
	int32_t fd = 0;  // handle
};

// return true if port open error
bool sp_open_err(const sp_type& sp_i) { return(sp_i.fd < 0); }

// open port
void sp_open_port(sp_type& sp_io, const string& dev_i) {
	sp_io.dev = dev_i;  // keep dev
	// O_RDWR : read & write port access
	// O_NOCTTY : no terminal control over port
	// O_NDELAY : open port immediately
	sp_io.fd = open(sp_io.dev.c_str(), O_RDWR | O_NOCTTY | O_NDELAY);
}

// close port (-1 : error)
bool sp_close_port(sp_type& sp_io) {
	// restore old attributes
	tcsetattr(sp_io.fd, TCSANOW, &sp_io.t_old);
	return(close(sp_io.fd));
}

// return true if port is open
bool sp_test_port(const string& dev_i) {
	sp_type sp;
	sp_open_port(sp, dev_i);
	bool open_f = !sp_open_err(sp);
	if (open_f) { close(sp.fd); }
	return(open_f);
}

// config port, return true if error
bool sp_config_port(sp_type& sp_io) {
	// enable FTDI low latency mode
	struct serial_struct ser_info;
	ioctl(sp_io.fd, TIOCGSERIAL, &ser_info);
	ser_info.flags |= ASYNC_LOW_LATENCY;
	ioctl(sp_io.fd, TIOCSSERIAL, &ser_info);	

	// read old attributes
	bool rtn_val = tcgetattr(sp_io.fd, &sp_io.t_old);
	if (rtn_val) { return(rtn_val); }
	// copy old attributes to freely modify them
	sp_io.t_new = sp_io.t_old;
	// set new attributes
	cfsetispeed(&sp_io.t_new, B230400);  // in speed
	cfsetospeed(&sp_io.t_new, B230400);  // out speed
	// control modes
	sp_io.t_new.c_cflag &= ~PARENB;   // no parity
	sp_io.t_new.c_cflag &= ~CSTOPB;  // 1 stop bit
	sp_io.t_new.c_cflag &= ~CSIZE;  // clear data bits field
	sp_io.t_new.c_cflag |=  CS8;   // 8 data bits
	sp_io.t_new.c_cflag &= ~CRTSCTS;  // no HW flow control
	sp_io.t_new.c_cflag |= CREAD;  // enable rx
	sp_io.t_new.c_cflag |= CLOCAL;  // ignore modem control lines
	// input modes
	sp_io.t_new.c_iflag &= ~(IXON | IXOFF);  // no input XON/XOFF flow control
//	sp_io.t_new.c_iflag &= ~IXANY;  // only START char can restart output (?)
//	sp_io.t_new.c_iflag |= INLCR;  // xlate NL to CR
//	sp_io.t_new.c_iflag |= ICRNL;  // xlate CR to NL
	sp_io.t_new.c_iflag &= ~INLCR;  // disable xlate NL to CR
	sp_io.t_new.c_iflag &= ~ICRNL;  // disable xlate CR to NL
//	sp_io.t_new.c_iflag |= IGNCR;  // ignore CR
	// output modes
//	sp_io.t_new.c_oflag |= OCRNL;  // xlate CR to NL
//	sp_io.t_new.c_oflag &= ~OCRNL;  // disable xlate CR to NL
//	sp_io.t_new.c_oflag &= ~ONLCR;  // disable xlate NL to CR-NL
//	sp_io.t_new.c_oflag |= ONOCR;  // don't output CR at column 0
//	sp_io.t_new.c_oflag |= ONLRET;  // don't output CR
	// local modes
	sp_io.t_new.c_lflag &= ~ICANON;  // non-canonical mode
	sp_io.t_new.c_lflag &= ~ECHO;  // disable local echo
//		t_new.c_lflag &= ~ECHOE;  // disable echo erase
	sp_io.t_new.c_lflag &= ~ISIG;  // disable signal gen from some chars
	// special modes
	sp_io.t_new.c_cc[VMIN] = 0;  // 0 min bytes
	sp_io.t_new.c_cc[VTIME] = 0;  // 0/10 sec inter byte time
	// write new attributes, return nz @ error
	return(tcsetattr(sp_io.fd, TCSANOW, &sp_io.t_new));
}

// str => tx port, return byte count
int32_t sp_tx(const sp_type& sp_i, const string& str_i) {
	int32_t tx_bytes = write(sp_i.fd, str_i.c_str(), str_i.size());
	return(tx_bytes);
}

// rx port => str, return byte count
// note: concats to string w/o clear
// you must externally poll
int32_t sp_rx(const sp_type& sp_i, string& str_o) {
	char rx_buf[RX_BUF_SZ];
	int32_t rx_bytes = 0;
	rx_bytes = read(sp_i.fd, &rx_buf, RX_BUF_SZ);		
	str_o += string(rx_buf, rx_bytes);
	return(rx_bytes);
}

// rx port flush
void sp_rx_flush(const sp_type& sp_i) {
	ms_sleep(RX_BURST_MDLY);  // delay
	string rx_str;
	sp_rx(sp_i, rx_str);
	return;
}

// rx port => str, return byte count
// note: concats to string w/o clear
// internal polling for big bursts
// safety timeout
int32_t sp_rx_burst(const sp_type& sp_i, string& str_o) {
	int32_t rx_bytes = 0;
	int32_t rx_bytes_total = 0;
	int32_t rx_count = 0;
	do {
		ms_sleep(RX_BURST_MDLY);  // delay
		rx_bytes = sp_rx(sp_i, str_o);
		rx_bytes_total += rx_bytes;  // sum
		rx_count++;
		} while ((!rx_bytes_total || rx_bytes) && rx_count < RX_MAX);
	return(rx_bytes_total);
}

// rx port wait for '>' at end
// internal polling
// safety timeout
bool sp_rx_wait(const sp_type& sp_i) {
	string rx_str;
	int32_t rx_bytes = 0;
	int32_t rx_count = 0;
	do { 
		ms_sleep(RX_POLL_MDLY);  // delay
		rx_bytes = sp_rx(sp_i, rx_str);
		rx_count++;
	} while ((!rx_bytes || (rx_str.back() != '>')) && rx_count < RX_MAX);
	return(!(rx_count < RX_MAX));
}

///////////////
// POSIX END //
///////////////

#endif

