// d-lib_preset_upd.cpp
// Preset update functions
#pragma once

///////////////////////
// ENCODER RESCALING //
///////////////////////

// encoder => exp use
uint32_t enc_exp(uint64_t enc_i, uint64_t width_i, uint64_t db12_i, bool down_i) {
	uint64_t mult = db12_i << 28;
	uint64_t offs = (16 - db12_i) << 28;
	uint64_t enc_fs = enc_i << (32 - width_i);  // full scale
	uint64_t enc_mult = (enc_fs * mult) >> 32;  // * mult
	uint64_t enc_offs = enc_mult + offs;  // + offs
	double_t enc_5_27 = double_t(enc_offs) / exp2(27);  // => 5.27 decimal
	uint64_t exp_o = exp2(enc_5_27);
	if (down_i) {
		uint64_t dn = 1 << ((-2 * db12_i) + 32);
		exp_o -= dn;
	}
	return(exp_o);
}

// encoder => 16.16 inv use
uint32_t enc_inv(uint64_t enc_i, uint64_t width_i) {
	uint64_t enc_fs = enc_i << (32 - width_i);  // full scale
	uint64_t enc_frac = enc_fs >> 16;  // 16.16 fractional
	double_t enc_16_16 = double_t(enc_frac) / exp2(16);  // => 16.16 decimal
	uint64_t inv_o = exp2(16) / enc_16_16;
	return(inv_o);
}

// knee update: [0:63]48db_dn => [0:31]inv_rev
uint8_t knee_upd(uint64_t enc_i) {
	if (!enc_i) return (0);  // zero out for zero in to turn off
	if (enc_i == 63) return (31);  // max out for max in
	uint32_t use_exp = (enc_exp(enc_i, 6, 4, true)) >> 20;  // old use
	uint32_t error = 1 << 31;
	uint32_t enc_o = 0;
	for(int i=0; i<32; i++) {
		uint32_t use_inv = (enc_inv(31-i, 5) - 0xc000) >> 12;  // new use
		uint32_t error_tmp = labs(int64_t(use_exp) - int64_t(use_inv));
		if (error_tmp < error) {
			error = error_tmp;
			enc_o = i;
		}
	}
	return(enc_o);
}

// span update: [0:31]48db_dn => [0:31]inv_rev
uint8_t span_upd(uint64_t enc_i) {
	if (!enc_i) return (0);  // zero out for zero in to turn off
	if (enc_i == 31) return (31);  // max out for max in
	uint32_t use_exp = ((enc_exp(enc_i, 5, 4, true)) >> 22) + 4;  // old use (+ 2 bits res)
	uint32_t error = 1 << 31;
	uint32_t enc_o = 0;
	for(int i=0; i<32; i++) {
		uint32_t use_inv = enc_inv(31-i, 5) >> 14;  // new use (+ 2 bits res)
		uint32_t error_tmp = labs(int64_t(use_exp) - int64_t(use_inv));
		if (error_tmp < error) {
			error = error_tmp;
			enc_o = i;
		}
	}
	return(enc_o);
}

// fall update: [0:63]48db_rev => [0:63]84db_rev_dn
uint8_t fall_upd(uint64_t enc_i) {
	if (!enc_i) return (0);  // zero out for zero in to turn off
	uint32_t use_exp = (enc_exp(63-enc_i, 6, 4, false)) >> 13;  // old use
	uint32_t error = 1 << 31;
	uint32_t enc_o = 0;
	for(int i=0; i<64; i++) {
		uint32_t use_exp2 = (enc_exp(63-i, 6, 7, true)) >> 7;  // new use
		uint32_t error_tmp = labs(int64_t(use_exp) - int64_t(use_exp2));
		if (error_tmp < error) {
			error = error_tmp;
			enc_o = i;
		}
	}
	return(enc_o);
}

// rise update: [0:63]48db_rev => [0:63]72db_rev
uint8_t rise_upd(uint64_t enc_i) {
	if (!enc_i) return (0);  // zero out for zero in to turn off
	uint32_t use_exp = (enc_exp(63-enc_i, 6, 4, false)) >> 10;  // old use
	uint32_t error = 1 << 31;
	uint32_t enc_o = 0;
	for(int i=0; i<64; i++) {
		uint32_t use_exp2 = (enc_exp(63-i, 6, 6, false)) >> 7;  // new use
		uint32_t error_tmp = labs(int64_t(use_exp) - int64_t(use_exp2));
		if (error_tmp < error) {
			error = error_tmp;
			enc_o = i;
		}
	}
	return(enc_o);
}


// freq update based on associated pmod and now global oct
// freq_i is [0:192] with 24 steps per octave
// pmod_i is [-63:63] full scale squared, with 32 == +1 octave
// oct_i is [-7:7]
uint8_t freq_upd(uint8_t freq_i, int8_t pmod_i, int8_t oct_i) {
	if (!oct_i) return (freq_i);  // zero octave : no change
	if (!pmod_i) return (freq_i);  // zero pmod: no change
	int32_t sgn = (pmod_i < 0) ? 1 : -1;
	int32_t offs = sgn * ((pmod_i * pmod_i * oct_i * 24) / (32 * 32));
	int32_t freq_o = freq_i + offs;
	if (freq_o < 0) { freq_o = 0; }
	if (freq_o > 192) { freq_o = 192; }
	return(freq_o);
}


// preset file => preset file
// read, update, write back preset file (after any significant encoder changes)
// return true if error
bool dlp_update(const string& dlp_i) {
	bool error_f = false;
	string pre_str;
	uint8_t pre_uba[PRE_PARAMS] = { 0 };
	uint8_t out_uba[PRE_PARAMS] = { 0 };
	
	// read from file
	if (!error_f) {	error_f = file_to_str(dlp_i, pre_str); }
	if (!error_f) { error_f = hex_str_to_uba(pre_str, pre_uba, 4, PRE_PARAMS); }

	// default: output = input, zero out high fluff
	for (uint32_t i=0; i<PRE_PARAMS; i++) {
		if (i > 96) { out_uba[i] = 0; }
		else { out_uba[i] = pre_uba[i]; }
	}

	////////////////////////
	// special cases here //
	////////////////////////
/*
	// resonator harm/2
	out_uba[35] = int8_t(pre_uba[35]) / 2;
*/
	// pitch correction slew knob reversal
	uint8_t temp = (pre_uba[70]) ? 32 - pre_uba[70] : 31;
	out_uba[70] = temp;

	// write back to file
	if (!error_f) { uba_to_str(out_uba, pre_str, PRE_PARAMS); }
	if (!error_f) { error_f = str_to_file(dlp_i, pre_str, false); }
	return(error_f);
}
