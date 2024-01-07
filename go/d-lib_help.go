package main

var help_commands_str = `
Usage: d-* [ command ] [ -flag <option> -flag <option> ... ]
Where: d-* = d-win (Windows); d-mac (Mac Intel); d-mm1 (Mac M1);
             d-lin (Linux Intel); d-arm (Linux ARM64); d-a32 (Linux ARM32)

COMMANDS & FLAGS:
  <command> -h                                                       Help with individual command flags
  help  <-v>                                                         Command line help with optional examples
  menu  <-d dir>                                                     Interactive update menu (default)
  port  <-l> <-s port>                                               List ports / set port
  view  <-k|-s slot|-f file.ext> <-pro>                              View knobs|slot|DLP file or slot in PRE|PRO|EEPROM file
  match <-d dir> <-s|-d2 dir|-f file.ext> <-hdr> <-g> <-pro>         Match DLP files with slots|DLP|PRE|PRO|EEPROM files
  diff  <-k> <-f file> <-s slot> <-f2 file> <-s2 slot> <-pro>        Compare single slot|knobs|DLP|PRE|PRO|EEPROM files
  dump  <-f file.ext> <-k|-s slot> <-pro> <-y>                       Download to DLP|SPI|PRE|PRO|EEPROM file
  pump  <-f file.ext> <-k|-s slot> <-pro>                            Upload from DLP|SPI|PRE|PRO|EEPROM file
  bank  <-f file> <-s slot> <-pro>                                   Upload BNK DLP files
  split <-f file.ext> <-y>                                           Split PRE|PRO|EEPROM container file into sub files
  join  <-f file.ext> <-y>                                           Join DLP|PRE|PRO|SPI sub files into container file
  morph <-k|-s slot|-f file> <-pkv page:knob[0:6]|all:val> <-i seed> Morph knobs|slot|DLP|PRE|EEPROM file to knobs
  batch <-d dir> <-d2 dir> <-m|u|r> <-pro> <-y> <-dry>               Batch convert DLP files (mono, update, rob)
  knob  <-pkv page:knob|all:val> <-v>                                Read/set knob or all knobs on page
  ver   <-f file.ext> <-crc> <-pre> <-pro>                           Read software version & check CRC
  stats <-p> <-v> <-e>                                               Read stats (pitch & volume field Hz, Hive errors)
  hcl   <ver|crc|acal|22 rk|...>                                     Issue HCL command**
  loop  <"some text to loop back">                                   Serial port loop back test***
  acal                                                               Issue immediate ACAL (no timer)
  reset                                                              Issue processor reset
`

var help_notes_str = `
NOTES:
- Flags may be entered in any order.
- Flag prefix either "-" or "--" (e.g. -s=5; --s=5).
- Flags / values separator either space or "=" (e.g. -s 5; -s=5).
- Many commands require a file extension to know what to do.
- If not provided, the *.dlp file extension may be added automatically.
- If provided, an incorrect file extension flags an error.
- If the specified target directory doesn't exist it will be created.
- A <y|ENTER|q> user prompt precedes most file overwrites (-y flag overrides prompt).
- The -pkv "knob" parameter is [0:6]: 0 @ upper left, 1 @ upper right, etc.
- The "cpypg" command is UI page based and can accept multiple -p flags.
- The "morph" command is UI page based and can accept multiple -pkv flags.
- The "morph" command negative val allows the sign to flip w/ 50/50 probability.
- The "morph" "knob" and "cpypg" commands use fuzzy matching for the page name.
- The "bank" command uses the bank file path to locate all files.
- The "bank" command skips over lines in *.bnk files with C / C++ type comments.
- Preset and profile files share the same *.dlp file extension.
- The serial port number is stored in the config file "d-lib.cfg".
- The "ports" command updates the config file if a port number is given via -p.
- If missing, the config file will be automatically generated.
- If "view" output doesn't fit in the window, resize it or change the font/layout.
- Linux & Mac terminals require executable files to be prefaced with: "./" e.g. "./d-mac".
- Windows terminal may require files to be prefaced with: ".\" e.g. ".\d-win".
- Linux users may need to join the "dialout" group for serial port access.
- If the librarian hangs, CTRL-C will usually kill a terminal program.
- ** For HCL commands, consult files HCL.txt, REGS.txt, and KNOBS.txt.
- *** Loopback requires USB dongle RX and TX wires to be connected together.
`

var help_examples_str = `
USAGE EXAMPLES:
- List flags for a command (e.g. morph):
    LIB_EXE morph -h
- Show librarian version & compact help:
    LIB_EXE help
- Show librarian version & verbose help:
    LIB_EXE help -v
- Interactive update menu (temp directory "_WORK_"):
    LIB_EXE
- Interactive update menu with temp directory "my_work_dir":
    LIB_EXE menu -d my_work_dir
- List all serial ports:
    LIB_EXE port -l
- Set port to 5 & list all serial ports:
    LIB_EXE port -l -s 5
- View all current knob values:
    LIB_EXE view -k
- View preset in slot 20:
    LIB_EXE view -s 20
- View profile in slot 2:
    LIB_EXE view -s 2 -pro
- View preset file "bassoon.dlp":
    LIB_EXE view -f bassoon
- View profile file "some_pro.dlp":
    LIB_EXE view -f some_pro -pro 
- View preset 55 in file "my_old_presets.pre"
    LIB_EXE view -s 198 -f my_old_presets.pre
- View profile 0 in file "some_stuff.pro"
    LIB_EXE view -s 0 -f some_stuff.pro
- View preset 198 in file "my_backup.eeprom"
    LIB_EXE view -s 198 -f my_backup.eeprom
- View profile 3 in file "my_backup.eeprom"
    LIB_EXE view -s 3 -f my_backup.eeprom -pro
- Match slots with DLP files in "_ALL_" directory:
    LIB_EXE match -s -d _ALL_
- Match slots with DLP files in "_ALL_" directory with best guess:
    LIB_EXE match -s -d _ALL_ -g
- Match DLP files in "_OLD_" directory with DLP files in "_ALL_" directory:
    LIB_EXE match -d2 _OLD_ -d _ALL_
- Match presets in "my.pre" with DLP files in "_ALL_" directory:
    LIB_EXE match -f my.pre -d _ALL_
- Match profiles in "my.pro" with DLP files in "_ALL_/sys" directory:
    LIB_EXE match -f my.pro -d _ALL_/sys
- Match presets in "my.eeprom" with DLP files in "_ALL_" directory:
    LIB_EXE match -f my.eeprom -d _ALL_
- Match profiles in "my.eeprom" with DLP files in "_ALL_/sys" directory:
    LIB_EXE match -f my.eeprom -d _ALL_/sys -pro
- Compare current knob values to preset file "mimi.dlp":
    LIB_EXE diff -f mimi -k
- Compare preset in slot 7 to file "saw.dlp":
    LIB_EXE diff -f saw -s 7
- Compare preset file "trixie.dlp" to file "patsy.dlp":
    LIB_EXE diff -f patsy -f2 trixie
- Compare profile file "_sys_3.dlp" to file "_sys_0.dlp":
    LIB_EXE diff -f _sys_0 -f2 _sys_3 -pro
- Compare preset in slot 20 to preset in slot 45:
    LIB_EXE diff -s 45 -s2 20
- Compare current knob values to preset in slot 3:
    LIB_EXE diff -s 3 -k
- Compare profile in slot 3 to profile in slot 0:
    LIB_EXE diff -s 0 -s2 3 -pro
- Compare preset preset in slot 45 to preset 36 in file "somepres.pre":
    LIB_EXE diff -s 45 -f2 somepres.pre -s2 36
- Compare profile 2 in file "jims.pro" to profile in slot 3:
    LIB_EXE diff -f jims.pro -s 2 -s2 3 -pro
- Download preset knobs to preset file "him_her.dlp":
    LIB_EXE dump -k -f him_her
- Upload preset file "flute.dlp" to preset knobs:
    LIB_EXE pump -k -f flute
- Download profile knobs to profile file "my_prof_4.dlp":
    LIB_EXE dump -k -f my_prof_4 -pro
- Upload profile file "some_prof.dlp" to profile knobs:
    LIB_EXE pump -k -f some_prof -pro
- Download preset slot 5 to preset file "female7.dlp":
    LIB_EXE dump -s 5 -f female7
- Upload preset file "cello8.dlp" to preset slot 9:
    LIB_EXE pump -f cello8 -s 9
- Download profile slot 0 to profile file "my_sys.dlp":
    LIB_EXE dump -s 0 -f my_sys -pro
- Upload profile file "_sys_9.dlp" to profile slot 3
    LIB_EXE pump -f _sys_9 -s 3 -pro
- Upload bank of presets in bank file "mybank.bnk" to preset slots 10, 11, 12, etc.:
    LIB_EXE bank -f mybank -s 10
- Write bank of presets in bank file "new.bnk" to preset slots 1, 2, etc. in file "some.pre":
    LIB_EXE bank -f newbank -s 1 -f2 some.pre
- Download software & all presets & profiles to file "2022-01-23.eeprom":
    LIB_EXE dump -f 2022-01-23.eeprom
- Upload software & all presets & profiles from file "factory.eeprom":
    LIB_EXE pump -f factory.eeprom
- Download software to file "sw_backup.spi":
    LIB_EXE dump -f sw_backup.spi
- Upload software from file "f9e1c5c7.spi":
    LIB_EXE pump -f f9e1c5c7.spi
- Download all presets to file "old_presets.pre":
    LIB_EXE dump -f old_presets.pre
- Upload all presets from file "my_dlev.pre":
    LIB_EXE pump -f my_dlev.pre
- Download all profiles to file "my_setup.pro":
    LIB_EXE dump -f my_setup.pro
- Upload all profiles from file "your_setup.pro":
    LIB_EXE pump -f your_setup.pro
- Split file "some.eeprom" into "some.pre", "some.pro", some.spi":
    LIB_EXE split -f some.eeprom
- Split file "my_setup.pro" into "pro_000.dlp" thru "pro_005.dlp":
    LIB_EXE split -f my_setup.pro
- Split file "my_new.pre" into "000.dlp" thru "249.dlp":
    LIB_EXE split -f my_new.pre
- Join files "some.pre", "some.pro" and "some.spi" into "some.eeprom":
    LIB_EXE join -f some.eeprom
- Join files "pro_000.dlp" thru "pro_005.dlp" to "stuff.pro":
    LIB_EXE join -f stuff.pro
- Join files "000.dlp" thru "249.dlp" to "some.pre":
    LIB_EXE join -f some.pre
- Morph knobs (all osc knobs & signs):
    LIB_EXE morph -k -pkv 0o:all:-12 -pkv 1o:all:-12
- Morph slot 23 (oscillator filter freq, resonator freq, seed):
    LIB_EXE morph -s 23 -pkv fo:2:5 -pkv r:0:20 -i 9
- Morph file "violin_4" (0 formant filters, resonator tap):
    LIB_EXE morph -f violin_4 -pkv 0f:2:8 -pkv 0f:4:8 -pkv r:1:10
- Batch convert all presets in the _ALL_ directory to mono in the _MONO_ directory:
    LIB_EXE batch -d _ALL_ -d2 _MONO_ -m
- Batch update all presets in the _ALL_ directory and overwrite them:
    LIB_EXE batch -d _ALL_ -d2 _ALL_ -u
- Read knob RESON:mode:
    LIB_EXE knob -pkv re:4
- Set knob RESON:mode to 10 & view all knobs:
    LIB_EXE knob -pkv re:4:10 -v
- Minimize all FLT_NOISE knobs:
    LIB_EXE knob -pkv fn:all:0
- Read the installed software version & check CRC:
    LIB_EXE ver
- Read file "2020.eeprom" software version & check CRC:
    LIB_EXE ver -f 2020.eeprom
- Read file "some.spi" software version & check CRC:
    LIB_EXE ver -f some.spi
- Report pitch & volume field frequencies (Hz):
    LIB_EXE stats -p -v
- Report processor errors:
    LIB_EXE stats -e
- Read processor registers 0 thru 9:
    LIB_EXE hcl 0 9 rr
- Loop back serial port text "testing 123":
    LIB_EXE loop "testing 123"
- Perform an ACAL:
    LIB_EXE acal
- Reset the processor:
    LIB_EXE reset
`
