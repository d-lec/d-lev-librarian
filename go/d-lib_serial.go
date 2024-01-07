package main

/*
 * d-lev support functions
*/

import (
	"fmt"
	"bytes"
	"strconv"
	"strings"
	"go.bug.st/serial"
)

// return a list of serial ports
func sp_list() ([]string) {
	ports, err := serial.GetPortsList(); err_chk(err)
	return ports
}

// open enumerated serial port
func sp_open() (serial.Port) {
	port := cfg_get("port")
	if port == "" { error_exit("Current port is not assigned") }
	// config as 230400bps N81
	mode := &serial.Mode{
		BaudRate: 230400,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	// open port
	sp, err := serial.Open(port, mode); 
	if err != nil { error_exit(fmt.Sprint("Can't open port: ", port)) }
	return sp
}

// write & read serial port, string i/o, optionally show activity
func sp_tx_rx(sp serial.Port, wr_str string, act_f bool) (string) {
	err := sp.ResetInputBuffer(); err_chk(err)
	err = sp.ResetOutputBuffer(); err_chk(err)
	// write to port
	_, err = sp.Write([]byte(wr_str)); err_chk(err)
	// read port, concat to buffer
	var rd_bytes bytes.Buffer
	var chars int
	for {
		rd_buf := make([]byte, RX_BUF_BYTES)
		n, err := sp.Read(rd_buf); err_chk(err)
		rd_bytes.Write(rd_buf[:n])  // concat
		if bytes.Contains(rd_buf[:n], []byte(">")) { break }  // read until prompt
		chars += n
		if act_f { chars = dots(chars) }
	}
	// done
	if act_f { fmt.Println(" download done") }
	return rd_bytes.String()
}

// read knobs data string
func sp_rx_knobs_str() (string) {
	sp := sp_open()
	rx_str := sp_tx_rx(sp, "0 " + strconv.Itoa(KNOBS_TOTAL-1) + " rk ", false)
	sp.Close()
	rx_str = decruft_hcl(rx_str)
	if strings.Count(rx_str, "\n") != KNOBS_TOTAL-1 { error_exit("Bad knob info") }
	return rx_str
}	

// read knobs, convert to pints
func sp_rx_knobs_pints(mode string) ([]int) {
	kints := hexs_to_ints(sp_rx_knobs_str(), 1)
	return kints_pints_order(kints, mode)
}

// write knobs pint data
func sp_tx_knobs_pints(pints []int, mode string) {
	sp := sp_open()
	switch mode {
	case "pre" :
		for pidx:=0; pidx<len(pre_params); pidx++ {
			_, _, _, kidx, _ := pname_lookup(pre_params[pidx].pname)
			wr_str := fmt.Sprint(kidx, " ", pints[pidx], " wk ")
			sp_tx_rx(sp, wr_str, false)
		}
	case "pro" :
		for pidx:=0; pidx<len(pro_params); pidx++ {
			_, _, _, kidx, _ := pname_lookup(pro_params[pidx].pname)
			wr_str := fmt.Sprint(kidx, " ", pints[pidx], " wk ")
			sp_tx_rx(sp, wr_str, false)
		}
	case "not" :
		for pidx:=0; pidx<len(not_params); pidx++ {
			_, _, _, kidx, _ := pname_lookup(not_params[pidx].pname)
			wr_str := fmt.Sprint(kidx, " ", pints[pidx], " wk ")
			sp_tx_rx(sp, wr_str, false)
		}
	}
	sp.Close()
}
