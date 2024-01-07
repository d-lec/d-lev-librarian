package main

/*
 * Librarian for the D-Lev Theremin
 * See file "CHANGE_LOG.txt" for details
*/

import (
	"flag"
	"os"
	// for fuzzy testing:
	//"strings"
	//"fmt"
)

func main() {

	// do update if no args
	if len(os.Args) < 2 {
		menu_cmd(WORK_DIR)
    } else {
		// parse subcommands
		cmd := os.Args[1]
		switch cmd {
			case "menu"  : menu()
			case "help", "-help", "-h", "/h" : help()
			case "port"  : port()
			case "view"  : view()
			case "match" : match()
			case "diff"  : diff()
			case "bank"  : bank()
			case "dump"  : dump()
			case "pump"  : pump()
			case "split" : split()
			case "join"  : join()
			case "morph" : morph()
			case "batch" : batch()
			case "knob"  : knob()
			case "ver"   : ver()
			case "hcl"   : hcl_cmd()
			case "loop"  : loop_cmd()
			case "acal"  : acal_cmd()
			case "reset" : reset_cmd()
			case "stats" : stats()
//			case "dev"   : dev()  // dev stuff
			default : cmd_hint(cmd)
		}
	}
}  // end of main()


////////////////////
// main functions //
////////////////////

// need these for custom flags
type flag_strs_t []string

func (i *flag_strs_t) String() string { return "" }

func (i *flag_strs_t) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// show help
func help() {
	sub := flag.NewFlagSet("help", flag.ExitOnError)
	verbose := sub.Bool("v", false, "verbose mode")
	sub.Parse(os.Args[2:])
	//
	help_cmd(*verbose)
}

// list free serial ports / set port
func port() {
	sub := flag.NewFlagSet("ports", flag.ExitOnError)
	set := sub.String("s", "", "`set` port number")
	list := sub.Bool("l", false, "list ports")
	sub.Parse(os.Args[2:])
	//
	port_cmd(*set, *list)
}

// report active or file version
func ver() {
	sub := flag.NewFlagSet("ver", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	crc := sub.Bool("crc", false, "crc check")
	pre := sub.Bool("pre", false, "preset check")
	pro := sub.Bool("pro", false, "profile check")
	sub.Parse(os.Args[2:])
	//
	ver_cmd(*file, *crc, *pre, *pro)
}

// report stats
func stats() {
	sub := flag.NewFlagSet("stats", flag.ExitOnError)
	p_hz := sub.Bool("p", false, "pitch field Hz")
	v_hz := sub.Bool("v", false, "volume field Hz")
	h_er := sub.Bool("e", false, "hive errors")
	sub.Parse(os.Args[2:])
	//
	stats_cmd(*p_hz, *v_hz, *h_er)
}

// view knobs, DLP file, slot
func view() {
	sub := flag.NewFlagSet("view", flag.ExitOnError)
	file := sub.String("f", "", "view `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	knobs := sub.Bool("k", false, "view knobs")
	slot := sub.String("s", "", "view `slot` number")
	sub.Parse(os.Args[2:])
	//
	view_cmd(*file, *pro, *knobs, *slot, -1)
}

// twiddle knob
func knob() {
	sub := flag.NewFlagSet("knob", flag.ExitOnError)
	pkv := sub.String("pkv", "", "page:all|knob[0:6]:val")
	view := sub.Bool("v", false, "view all knobs")
	sub.Parse(os.Args[2:])
	//
	knob_cmd(*pkv, *view)
}

// diff DLP file(s) / slot(s) / knobs
func diff() {
	sub := flag.NewFlagSet("diff", flag.ExitOnError)
	file := sub.String("f", "", "compare `file` name")
	file2 := sub.String("f2", "", "compare `file2` name")
	pro := sub.Bool("pro", false, "profile mode")
	knobs := sub.Bool("k", false, "compare knobs")
	slot := sub.String("s", "", "compare `slot` number")
	slot2 := sub.String("s2", "", "compare `slot2` number")
	sub.Parse(os.Args[2:])
	//
	diff_cmd(*file, *file2, *pro, *knobs, *slot, *slot2)
}

// match slots / DLP files w/ DLP files & list
func match() {
	sub := flag.NewFlagSet("match", flag.ExitOnError)
	dir := sub.String("d", "", "`directory` name")
	dir2 := sub.String("d2", "", "`directory` name")
	file := sub.String("f", "", "compare `file` name")
	pro := sub.Bool("pro", false, "profile mode")
	hdr := sub.Bool("hdr", false, "header format")
	guess := sub.Bool("g", false, "guess")
	slots := sub.Bool("s", false, "slots")
	sub.Parse(os.Args[2:])
	//
	match_cmd(*dir, *dir2, *file, *pro, *hdr, *guess, *slots)
}


////////////////////////////
// file upload & download //
////////////////////////////

// dump to file
func dump() {
	sub := flag.NewFlagSet("dump", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	slot := sub.String("s", "", "source `slot` number")
	knobs := sub.Bool("k", false, "source knobs")
	pro := sub.Bool("pro", false, "profile mode")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	dump_cmd(*file, *slot, *knobs, *pro, *yes)
}

// pump from file
func pump() {
	sub := flag.NewFlagSet("pump", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	slot := sub.String("s", "", "target `slot` number")
	knobs := sub.Bool("k", false, "target knobs")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	pump_cmd(*file, *slot, *knobs, *pro)
}

// *.bnk => *.dlps => slots
func bank() {
	sub := flag.NewFlagSet("btos", flag.ExitOnError)
	slot := sub.String("s", "", "starting `slot` number")
	file := sub.String("f", "", "bank `file` name")
	file2 := sub.String("f2", "", "target `file2` name")
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	bank_cmd(*slot, *file, *file2, *pro)
}


////////////////////////
// split / join files //
////////////////////////

// split bulk files
func split() {
	sub := flag.NewFlagSet("split", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	split_cmd(*file, *yes)
}

// join bulk files
func join() {
	sub := flag.NewFlagSet("join", flag.ExitOnError)
	file := sub.String("f", "", "target `file` name")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	join_cmd(*file, *yes)
}


//////////////////
// morph & copy //
//////////////////

// morph knobs|slot|file to knobs
func morph() {
	var pkv flag_strs_t
	sub := flag.NewFlagSet("morph", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	knobs := sub.Bool("k", false, "source knobs")
	slot := sub.String("s", "", "source `slot` number")
	seed := sub.Int64("i", 0, "initial seed")
	sub.Var(&pkv, "pkv", "page:all|knob[0:6]:val")
	sub.Parse(os.Args[2:])
	if *seed == 0 { *seed = timeseed() }
	morph_cmd(*file, *knobs, *slot, *seed, pkv)
}


//////////////
// updating //
//////////////

// read, processs, write all *.dlp in dir => dir2
func batch() {
	sub := flag.NewFlagSet("batch", flag.ExitOnError)
	dir := sub.String("d", "", "source `directory` name")
	dir2 := sub.String("d2", "", "target `directory` name")
	pro := sub.Bool("pro", false, "profile mode")
	mono := sub.Bool("m", false, "to mono")
	update := sub.Bool("u", false, "update")
	robs := sub.Bool("r", false, "Rob S. PP")
	yes := sub.Bool("y", false, "overwrite files")
	dry := sub.Bool("dry", false, "dry run")
	sub.Parse(os.Args[2:])
	//
	process_dlps(*dir, *dir2, *pro, *mono, *update, *robs, *yes, *dry)
}

// do a bunch of update stuff via interactive menu
func menu() {
	sub := flag.NewFlagSet("menu", flag.ExitOnError)
	dir := sub.String("d", WORK_DIR, "work `directory` name")
	sub.Parse(os.Args[2:])
	//
	menu_cmd(*dir)	
}


/////////
// dev //
/////////

/*
func dev() {
	var pm flag_strs_t
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	sub.Var(&pm, "pm", "page:morph")
	sub.Parse(os.Args[2:])
	//
	dev_cmd(pm)
}
*/

/*
	// fuzz testing
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	pattern := sub.String("p", "", "pattern string")
	file := sub.String("f", "", "source `file` name")
	sub.Parse(os.Args[2:])
	//
	file_str := file_read_str(*file)
	str_split := strings.Split(strings.TrimSpace(file_str), "\n")
	list := fuzzy_list(*pattern, str_split, 30, true)
	for _, word := range list { fmt.Println(word) }
*/


/*
	// fuzz testing
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	pattern := sub.String("p", "", "pattern string")
	sub.Parse(os.Args[2:])
	//
	fmt.Println(fuzzy_top(*pattern, lib_cmds))
*/


/*
	// fuzz testing
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	pattern := sub.String("p", "", "pattern string")
	file := sub.String("f", "", "source `file` name")
	sub.Parse(os.Args[2:])
	//
	fuzzy_test(*pattern, *file)
*/
/*
	// fuzz testing
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	pattern := sub.String("p", "", "pattern string")
	entry := sub.String("e", "", "entry string")
	sub.Parse(os.Args[2:])
	//
	fmt.Println(fuzzy_score(*pattern, *entry))
*/

/*
	// check comment removal from file
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	sub.Parse(os.Args[2:])
	//
	str := file_read_str(*file)
	fmt.Println(strip_cmnt(str))
*/

/*
	// find DLP files with certain characteristics
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	dir := sub.String("d", ".", "`directory` name")
	sub.Parse(os.Args[2:])
	//
	find_dlp(*dir)
	//dev_cmd()
*/
/*
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	str := sub.String("s", "", "test string")
	sub.Parse(os.Args[2:])
	//
	fmt.Print("_", *str, "_\n")
*/

