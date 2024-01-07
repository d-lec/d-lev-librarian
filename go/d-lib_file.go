package main

/*
 * d-lib support functions
*/

import (
	"os"
	"path/filepath"
	"errors"
	"strings"
	"fmt"
)

// check for blank <file> name
func file_blank_chk(file string) {
	if file == "" { error_exit("Missing file name") }
}

// check for blank <dir> name
func dir_blank_chk(dir string) {
	if dir == "" { error_exit("Missing directory name") }
}

// check if dir or file exists
func path_exists(path string) (bool) {
    _, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// check if dir exists, quit if not
func dir_chk(dir string) {
	dir_blank_chk(dir)
	if !path_exists(dir) { error_exit(fmt.Sprint("Directory ", dir, " does not exist")) }
}

// check if <file> exists, quit if not
func file_chk(file string) {
	file_blank_chk(file)
	if !path_exists(file) { error_exit(fmt.Sprint("File ", file, " does not exist")) }
}

// create directory for file if directory does not exist
func file_make_dir(file string) {
	dir, _ := filepath.Split(file)
	if dir == "" { return }
	err_chk(os.MkdirAll(dir, os.ModePerm))
}

// check <file> <ext>, add if missing, return it
func file_ext_chk(file string, ext string) (string) {
	file_blank_chk(file)
    f_ext := filepath.Ext(file)
    if f_ext == "" { 
		file += ext
		return file 
	}
	if f_ext != ext { error_exit(fmt.Sprint("Wrong file extension: ", f_ext, " (expecting: ", ext, " or none)")) }
	return file
}

// read trimmed string from file
func file_read_str(file string) (string) {
	file_chk(file)
	file_bytes, err := os.ReadFile(file); err_chk(err)
	return strings.TrimSpace(string(file_bytes))
}

// check if <file> exists, prompt to overwrite
func file_write_chk(file string, yes bool) (bool) {
	file_blank_chk(file)
    if path_exists(file) { return user_prompt("Overwrite file " + file + "?", yes, false) }
	return true
}

// write trimmed string to checked file
func file_write_str(file, data string, yes bool) (bool) {
	if file_write_chk(file, yes) {
		file_make_dir(file)
		err_chk(os.WriteFile(file, []byte(strings.TrimSpace(data)), 0666))
		return true
	}
	return false
}

// get name and contents of all *<ext> files in <dir>
// optionally trim file extensions
func dir_read_strs(dir, ext string, trim bool) (name_strs, data_strs []string) {
	files, err := os.ReadDir(dir); err_chk(err)
	for _, file := range files {
		if (filepath.Ext(file.Name()) == ext) && (file.IsDir() == false) {
			file_str := file_read_str(filepath.Join(dir, file.Name()))
			data_strs = append(data_strs, file_str)
			if trim { 
				name_strs = append(name_strs, strings.TrimSuffix(file.Name(), ext))
			} else {
				name_strs = append(name_strs, file.Name())
			}
		}
    }
    return
}


////////////
// CONFIG //
////////////

// set key value in config file
// create file if it doesn't exist
func cfg_set(key, value string) {
	if !path_exists(CFG_FILE) {	 // create file
		file_write_str(CFG_FILE, "", true) 
	}
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	file_str := file_read_str(CFG_FILE)
	lines := strings.Split(file_str, "\n")
	str := ""
	found := false
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == key {
			str += key + " " + value + "\n"
			found = true
		} else { str += line + "\n"	}
	}
	if !found { str += key + " " + value + "\n" }
	file_write_str(CFG_FILE, str, true)
}

// get first matching key value from config file
// return "" if file doesn't exist or no key match
func cfg_get(key string) (string) {
	if !path_exists(CFG_FILE) { return "" }
	key = strings.TrimSpace(key)
	file_str := file_read_str(CFG_FILE)
	lines := strings.Split(file_str, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == key {
			return fields[1]
		}
	}
	return ""
}

