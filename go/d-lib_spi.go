package main

/*
 * d-lev support functions
*/

import (
	"fmt"
	"strings"
	"strconv"
	"go.bug.st/serial"
	"time"
)

// show activity via printed dots
func dots(chars int) (int) {
	if chars > 0 {
		chars -= CHARS_PER_DOT
		fmt.Print(".") 
	}
	return chars
}

// read SPI port to string, trim cruft, optionally show activity
func spi_rd(addr int, addr_end int, act_f bool) (string) {
	sp := sp_open()
	rx_str := sp_tx_rx(sp, strconv.Itoa(addr) + " " + strconv.Itoa(addr_end) + " rs ", act_f)
	sp.Close()
	rx_str = decruft_hcl(rx_str)
	if len(strings.Split(rx_str, "\n")) != 1 + (addr_end - addr) / 4 { error_exit("Bad SPI read") }
	return rx_str
}

// SPI write enable
func spi_wr_en(sp serial.Port) {
	sp_tx_rx(sp, "6 6 wr ", false)
	sp_tx_rx(sp, "6 0x100 wr ", false)  // csn hi
}

// SPI write & wait
func spi_wr_wait(sp serial.Port) {
	sp_tx_rx(sp, "6 0x100 wr ", false)  // csn hi
	time.Sleep(EE_WR_MS * time.Millisecond)
}

// SPI write protect & unprotect
func spi_wr_prot(sp serial.Port, prot_f bool) {
	spi_wr_en(sp)
	sp_tx_rx(sp, "6 1 wr ", false)  // wrsr reg
	if prot_f { sp_tx_rx(sp, "6 0xc wr ", false)
	} else { sp_tx_rx(sp, "6 0 wr ", false)	}
	spi_wr_wait(sp)
}

// write string to SPI port, optionally show activity
func spi_wr(addr int, wr_str string, act_f bool) {
	sp := sp_open()
	spi_wr_prot(sp, false)
	wr_strs := str_split_trim(wr_str)
	var chars int
	for _, wr_line := range wr_strs {
		var cmd string
		if addr % EE_PG_BYTES == 0 {  // page boundary
			spi_wr_wait(sp)
			spi_wr_en(sp)
			cmd = strconv.Itoa(addr) + " "
		}
		if wr_line != "0" { cmd += "0x" }  // no 0x for zero data
		cmd += wr_line + " ws "
		sp_tx_rx(sp, cmd, false)
		chars += len(cmd)
		addr += EE_RW_BYTES;
		if act_f { chars = dots(chars) }
	}
	// done
	spi_wr_wait(sp);
	spi_wr_prot(sp, true);
	sp.Close()
	if act_f { fmt.Println(" upload done") }
}

// return spi bulk addresses
func spi_bulk_addrs(mode string) (addr int, end int) {
	switch mode {
		case ".pre" :
			addr = EE_PRE_ADDR
			end = EE_PRE_END
		case ".pro" :
			addr = EE_PRO_ADDR
			end = EE_PRO_END
		case ".spi" :
			addr = EE_SPI_ADDR
			end = EE_SPI_END
		case ".eeprom" :
			addr = EE_START
			end = EE_END
		default :
			error_exit(fmt.Sprint("Unknown mode: ", mode))
	}
	return
}

// return spi slot addr
func spi_slot_addr(slot int, pro bool) (int) {
	slot_lim_chk(slot, pro)
	if pro { return (slot + PRE_SLOTS) * EE_PG_BYTES }
	return slot * EE_PG_BYTES
}

// trim command, address, and prompt cruft from hcl read string
func decruft_hcl(str_i string) (str_o string) {
	strs_i := str_split_trim(str_i)
	for idx, str := range strs_i {
		// kill cmd & prompt & addr
		if (idx != 0) && (idx != len(strs_i) - 1) {  
			addr_end := strings.Index(str, "]")
			str_o += str[addr_end+1:] + "\n" 
		}
	}
	return strings.TrimSpace(str_o)
}

// get single slot data string
func spi_rd_slot_str(slot int, pro bool) (string) {
	addr := spi_slot_addr(slot, pro)
	rx_str := spi_rd(addr, addr + EE_PG_BYTES - 1, false)
	return rx_str
}

// get slots data strings
func spi_rd_slots_strs(pro bool) ([]string) {
	ext := ".pre"
	if pro { ext = ".pro" }
	addr, end := spi_bulk_addrs(ext)
	rx_str := spi_rd(addr, end - 1, true)
	strs := split_pre_pro_str(rx_str)
	slots := PRE_SLOTS
	if pro { slots = PRO_SLOTS }
	if len(strs) < slots/EE_RW_BYTES { error_exit("Bad slots info") }
	return strs
}
