package main

/*
 * Librarian for the D-Lev Theremin
 * See file "CHANGE_LOG.txt" for details
*/

import (
	"flag"
	"os"
	"fmt"
)

func main() {

	// do update if no args
	if len(os.Args) < 2 {
		menu_cmd(WORK_DIR)
    } else {
		// parse subcommands
		switch os.Args[1] {
			case "menu" : menu()
			case "help", "-help", "-h", "/h" : help()
			case "ports" : ports()
			case "view" : view()
			case "match" : match()
			case "diff" : diff()
			case "bank" : bank()
			case "dump" : dump()
			case "pump" : pump()
			case "split" : split()
			case "join" : join()
			case "morph" : morph()
			case "batch" : batch()
			case "knob" : knob()
			case "ver" : ver()
			case "hcl" : hcl_cmd()
			case "loop" :loop_cmd()
			case "acal" : acal_cmd()
			case "reset" : reset_cmd()
			case "dev" : dev()  // dev stuff
			default : error_exit(fmt.Sprint("Unknown command: ", os.Args[1]))
		}
	}
}  // end of main()


////////////////////
// main functions //
////////////////////

// show help
func help() {
	sub := flag.NewFlagSet("help", flag.ExitOnError)
	verbose := sub.Bool("v", false, "verbose mode")
	sub.Parse(os.Args[2:])
	//
	help_cmd(*verbose)
}

// list free serial ports / set port
func ports() {
	sub := flag.NewFlagSet("ports", flag.ExitOnError)
	port := sub.String("p", "", "`port` number")
	sub.Parse(os.Args[2:])
	//
	ports_cmd(*port)
}

// report active or file version
func ver() {
	sub := flag.NewFlagSet("ver", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	sub.Parse(os.Args[2:])
	//
	ver_cmd(*file, false)
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
	knob := sub.String("k", "", "page:knob[0:6]")
	ofs := sub.String("o", "", "knob offset `value`")
	set := sub.String("s", "", "knob set `value`")
	min := sub.Bool("min", false, "set knob min")
	view := sub.Bool("v", false, "view all knobs")
	sub.Parse(os.Args[2:])
	//
	knob_cmd(*knob, *ofs, *set, *min, *view)
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
	pro := sub.Bool("pro", false, "profile mode")
	sub.Parse(os.Args[2:])
	//
	bank_cmd(*slot, *file, *pro)
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


///////////
// morph //
///////////

func morph() {
	sub := flag.NewFlagSet("morph", flag.ExitOnError)
	file := sub.String("f", "", "source `file` name")
	knobs := sub.Bool("k", false, "source knobs")
	slot := sub.String("s", "", "source `slot` number")
	seed := sub.Int("i", timeseed(), "initial seed")
	mo := sub.Int("mo", 0, "oscillator multiplier")
	mn := sub.Int("mn", 0, "noise multiplier")
	me := sub.Int("me", 0, "eq (bass & treble) multiplier")
	mf := sub.Int("mf", 0, "filter multiplier")
	mr := sub.Int("mr", 0, "resonator multiplier")
	sub.Parse(os.Args[2:])
	morph_cmd(*file, *knobs, *slot, *seed, *mo, *mn, *me, *mf, *mr)
}


//////////////
// updating //
//////////////

// read, processs, write all *.dlp in dir => dir2
func batch() {
	sub := flag.NewFlagSet("batch", flag.ExitOnError)
	dir := sub.String("d", ".", "source `directory` name")
	dir2 := sub.String("d2", ".", "target `directory` name")
	pro := sub.Bool("pro", false, "profile mode")
	mono := sub.Bool("m", false, "to mono")
	update := sub.Bool("u", false, "update")
	robs := sub.Bool("r", false, "Rob S. PP")
	yes := sub.Bool("y", false, "overwrite files")
	sub.Parse(os.Args[2:])
	//
	process_dlps(*dir, *dir2, *pro, *mono, *update, *robs, *yes)
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

func dev() {

	// find DLP files with certain characteristics
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	dir := sub.String("d", ".", "`directory` name")
	sub.Parse(os.Args[2:])
	//
	find_dlp(*dir)
	//dev_cmd()

/*
	sub := flag.NewFlagSet("dev", flag.ExitOnError)
	str := sub.String("s", "", "test string")
	sub.Parse(os.Args[2:])
	//
	fmt.Print("_", *str, "_\n")
*/

}
