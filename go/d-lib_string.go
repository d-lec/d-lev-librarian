package main

/*
 * d-lev string functions
*/

import (
	"strings"
	"strconv"
	"regexp"
	"sort"	
)

// check char for upper case  
func ch_is_upper(ch byte) bool {
	if (ch < 65) || (ch > 90) { return false }
	return true
}

// check char for lower case  
func ch_is_lower(ch byte) bool {
	if (ch < 97) || (ch > 122) { return false }
	return true
}

// char to upper case  
func ch_to_upper(ch byte) byte {
	if ch_is_lower(ch) { return ch - 32 }
	return ch
}

// char to lower case  
func ch_to_lower(ch byte) byte {
	if ch_is_upper(ch) { return ch + 32 }
	return ch
}

// check chars for case independent equality
func chs_same(ch0, ch1 byte) bool {
	return ch_to_upper(ch0) == ch_to_upper(ch1) 
}

// trim start/end whitespace
func strs_trim(strs []string) ([]string) {
	for idx, str := range strs {
		strs[idx] = strings.TrimSpace(str)
	}
	return strs
}

// trim & split & trim
func str_split_trim(str string) (strs []string) {
	strs = strings.Split(strings.TrimSpace(str), "\n")
	strs = strs_trim(strs)
	return
}

// convert string of multi-byte hex values to slice of ints
// hex string values on separate lines
func hexs_to_ints(hex_str string, bytes int) ([]int) {
	var ints []int
	hex_strs := str_split_trim(hex_str)
	for _, str := range hex_strs {
		sh_reg, err := strconv.ParseUint(str, 16, 32); err_chk(err)
		for b:=0; b<bytes; b++ { 
			sh_byte := int(uint8(sh_reg))
			ints = append(ints, sh_byte)
			sh_reg >>= 8
		}
	}
	return ints
}

// convert slice of ints to string of multi-byte hex values
// hex string values on separate lines
func ints_to_hexs(ints []int, bytes int) (string) {
	var hex_str string
	for i:=0; i<len(ints); i+=bytes {
		var line_int int64
		for b:=0; b<bytes; b++ { 
			line_int += int64(uint8(ints[i+b])) << (b * 8)
		}
		hex_str += strconv.FormatInt(line_int, 16) + "\n"
	}
	return hex_str
}

// check for hexness
func str_is_hex(str string) bool {
	if len(str) == 0 { return false }
	for _, ch := range str {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) { return false }
	}
	return true
}

// return index of string in slice, else -1
func str_idx(strs []string, str string) (int) {
	for idx, entry := range strs { if str == entry { return idx } }
	return -1
}

// remove C style block comments
func strip_c_cmnt(str string) string {
	c_cmnt, err:= regexp.Compile(`/\*[^*]*\*+(?:[^*/][^*]*\*+)*/`); err_chk(err)
	return c_cmnt.ReplaceAllString(str, string(""))
}

// remove C++ style EOL comments
func strip_cpp_cmnt(str string) string {
	cpp_cmnt, err := regexp.Compile(`//.*`); err_chk(err)
	return cpp_cmnt.ReplaceAllString(str, string(""))
}

// remove all C/C++ style comments
func strip_cmnt(str string) string {
	return strip_cpp_cmnt(strip_c_cmnt(str))
}

// for given pattern and target, return fuzzy score and match flag
// match flag is true if all pattern chars appear in order in target
func fuzzy_score(pattern, target string) (int, bool) {
	// recursive function to exhaustively walk target vs pattern, left to right
	// looking for all chars in pattern to appear in order in target
	// for pattern[0] char test all target chars
	// for pattern[1] char test only target chars to the right of that, etc.
	// return best match flag score and match flag
	best_score := -2 * len(target)
	match_f := false
	var walk func(int, int, int)
	walk = func(p_idx_i, t_idx_i, score_i int) {
		// matched to end of pattern: update flag & score & return
		if p_idx_i >= len(pattern) { 
			match_f = true
			score_i -= len(target)  // target length penalty
			if score_i > best_score { best_score = score_i }  // save best
			return
		}
		// else: search from target position to end, recurse if char match
		for t_idx := t_idx_i; t_idx < len(target); t_idx++ {
			if chs_same(target[t_idx], pattern[p_idx_i]) {  // chars match
				score := -t_idx  // distance from start penalty
				score += 60 / (t_idx + 1)  // 1st chars bonus
				if t_idx > 0 {
					if strings.ContainsAny(string(target[t_idx-1]), "-_. ") { score += 30 }  // post separator char bonus
					if ch_is_lower(target[t_idx-1]) && ch_is_upper(target[t_idx]) { score += 30 }  // camelCase chars bonus
					if (p_idx_i > 0) && chs_same(target[t_idx-1], pattern[p_idx_i-1]) { score += 10 }  // adjacent chars bonus
				}
				walk(p_idx_i+1, t_idx+1, score_i+score)  // recurse
			}
		}
		// end of target, quit
		return
	}
	// kick off walk
	walk(0, 0, 0)
	return best_score, match_f
}

// do fuzzy matching of pattern vs list of strings
// optionally sort by highest score
func fuzzy_list(pattern string, list []string, rtn_len int, sort_en bool) ([]string) {
	type entry struct {
		e_str string
		e_score int
		e_match bool
	}
	entries := make([]entry, 0)
	for _, l_str := range list {
		score, match := fuzzy_score(pattern, l_str)
		ent := entry{e_str:l_str, e_score:score, e_match:match}
		entries = append(entries, ent)
	}
	if sort_en { sort.SliceStable(entries, func(i, j int) bool { return entries[i].e_score > entries[j].e_score }) }  // sort by score
	var rtn_strs []string
	for _, ent := range entries {
		if len(rtn_strs) >= rtn_len { break }
		if ent.e_match { rtn_strs = append(rtn_strs, ent.e_str) }
//		if ent.e_match { rtn_strs = append(rtn_strs, strconv.Itoa(ent.e_score) + " " + ent.e_str) }  // return score too for debug
	}
	return rtn_strs
}
