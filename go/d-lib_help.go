package main

var help_str = `
Usage: d-* [ command ] [ -flag <option> -flag <option> ... ]
Where: d-* = d-win (Windows); d-mac (Mac Intel); d-mm1 (Mac M1);
             d-lin (Linux Intel); d-arm (Linux ARM64); d-a32 (Linux ARM32)

COMMANDS & FLAGS:
  <command> -h                                               Help with individual command flags
  help  <-v>                                                 Command help with optional examples
  menu  <-d dir>                                             Interactive update menu (default)
  ports <-p port>                                            List ports / set port
  view  <-k|-s slot|-f file.ext> <-pro>                      View knobs|slot|DLP file or slot in PRE|PRO|EEPROM file
  match <-d dir> <-s|-d2 dir|-f file.ext> <-hdr> <-g> <-pro> Match DLP files with slots|DLP|PRE|PRO|EEPROM files
  diff  <-f file> <-k|-s slot|-f2 file> <-pro>               Compare DLP file to knobs|slot|DLP file2
  diff  <-s slot> <-k|-s2 slot> <-pro>                       Compare slot to knobs|slot2
  dump  <-f file.ext> <-k|-s slot> <-pro> <-y>               Download to DLP|SPI|PRE|PRO|EEPROM file
  pump  <-f file.ext> <-k|-s slot> <-pro>                    Upload from DLP|SPI|PRE|PRO|EEPROM file
  bank  <-f file> <-s slot> <-pro>                           Upload BNK DLP files
  split <-f file.ext> <-y>                                   Split PRE|PRO|EEPROM container file into sub files
  join  <-f file.ext> <-y>                                   Join DLP|PRE|PRO|SPI sub files into container file
  morph <-k|-s slot|-f file> <-mo|n|e|f|r|> <-i seed>        Morph knobs|slot|DLP file to knobs
  batch <-d dir> <-d2 dir> <-m|u|r> <pro> <-y>               Batch convert DLP files (mono, update, rob)
  knob  <-k page:knob|all> <-o offset|-s set|-min> <-v>      Read/set/offset/min knob value, min all page knobs
  ver   <-f file.ext>                                        Read software version & check CRC
  hcl   <ver|crc|acal|22 rk|...>                             Issue HCL command**
  loop  <"some text to loop back">                           Serial port loop back test***
  acal                                                       Issue immediate ACAL (no timer)
  reset                                                      Issue processor reset
`

var help_verbose_str = `
NOTES:
- Flags may be entered in any order.
- Flag prefix either "-" or "--" (e.g. -s=5; --s=5).
- Flags / values separator either space or "=" (e.g. -s 5; -s=5).
- If not provided, the *.dlp file extension is added automatically.
- If provided, an incorrect file extension flags an error.
- If the specified target directory doesn't exist it will be created.
- Many commands require a file extension to know what to do.
- A <y|ENTER|q> user prompt precedes most file overwrites (-y flag overrides prompt).
- The "knob" command page name matches first chars, and is case agnostic.
- The "knob" command knob number [0:6]: 0 @ upper left, 1 @ upper right, etc.
- The "bank" command uses the file path to locate all files.
- The "bank" command skips over lines in *.bnk files that begin with "//".
- Preset and profile files share the same *.dlp file extension.
- The serial port number is stored in the config file "d-lib.cfg".
- The "ports" command updates the config file if a port number is given via -p.
- If missing, the config file will be automatically generated.
- If "view" output doesn't fit in the window, resize it or change the font/layout.
- Linux & Mac terminals require executable files to be prefaced with: "./" e.g. "./d-mac".
- Windows terminal may require files to be prefaced with: ".\" e.g. ".\d-win".
- Linux users may need to join the "dialout" group for serial port access.
- If the librarian hangs, CTRL-C will usually kill a terminal program.
- ** See files HCL.txt, REGS.txt, and KNOBS.txt for details.
- *** Requires USB dongle RX and TX wires to be connected together.

USAGE EXAMPLES: (e.g. Windows build)
- List flags for a command (e.g. morph):
    d-win morph -h
- Show librarian version & compact help:
    d-win help
- Show librarian version & verbose help:
    d-win help -v
- Interactive update menu (temp directory "_WORK_"):
    d-win
- Interactive update menu with temp directory "my_work_dir":
    d-win menu -d my_work_dir
- List all serial ports & current port:
    d-win ports
- List all serial ports & set port to 5:
    d-win ports -p 5
- View all current knob values:
    d-win view -k
- View preset in slot 20:
    d-win view -s 20
- View profile in slot 2:
    d-win view -s 2 -pro
- View preset file "bassoon.dlp":
    d-win view -f bassoon
- View profile file "some_pro.dlp":
    d-win view -f some_pro -pro 
- View preset 55 in file "my_old_presets.pre"
    d-win view -s 198 -f my_old_presets.pre
- View profile 0 in file "some_stuff.pro"
    d-win view -s 0 -f some_stuff.pro
- View preset 198 in file "my_backup.eeprom"
    d-win view -s 198 -f my_backup.eeprom
- View profile 3 in file "my_backup.eeprom"
    d-win view -s 3 -f my_backup.eeprom -pro
- Match slots with DLP files in "_ALL_" directory:
    d-win match -s -d _ALL_
- Match slots with DLP files in "_ALL_" directory with best guess:
    d-win match -s -d _ALL_ -g
- Match DLP files in "_OLD_" directory with DLP files in "_ALL_" directory:
    d-win match -d2 _OLD_ -d _ALL_
- Match presets in "my.pre" with DLP files in "_ALL_" directory:
    d-win match -f my.pre -d _ALL_
- Match profiles in "my.pro" with DLP files in "_ALL_/sys" directory:
    d-win match -f my.pro -d _ALL_/sys
- Match presets in "my.eeprom" with DLP files in "_ALL_" directory:
    d-win match -f my.eeprom -d _ALL_
- Match profiles in "my.eeprom" with DLP files in "_ALL_/sys" directory:
    d-win match -f my.eeprom -d _ALL_/sys -pro
- Compare current knob values to preset file "mimi.dlp":
    d-win diff -f mimi -k
- Compare preset in slot 7 to file "saw.dlp":
    d-win diff -f saw -s 7
- Compare preset file "trixie.dlp" to file "patsy.dlp":
    d-win diff -f patsy -f2 trixie
- Compare profile file "_sys_3.dlp" to file "_sys_0.dlp":
    d-win diff -f _sys_0 -f2 _sys_3 -pro
- Compare preset in slot 20 to preset in slot 45:
    d-win diff -s 45 -s2 20
- Compare current knob values to preset in slot 3:
    d-win diff -s 3 -k
- Compare profile in slot 3 to profile in slot 0:
    d-win diff -s 0 -s2 3 -pro
- Download preset knobs to preset file "him_her.dlp":
    d-win dump -k -f him_her
- Upload preset file "flute.dlp" to preset knobs:
    d-win pump -k -f flute
- Download profile knobs to profile file "my_prof_4.dlp":
    d-win dump -k -f my_prof_4 -pro
- Upload profile file "some_prof.dlp" to profile knobs:
    d-win pump -k -f some_prof -pro
- Download preset slot 5 to preset file "female7.dlp":
    d-win dump -s 5 -f female7
- Upload preset file "cello8.dlp" to preset slot 9:
    d-win pump -f cello8 -s 9
- Download profile slot 0 to profile file "my_sys.dlp":
    d-win dump -s 0 -f my_sys -pro
- Upload profile file "_sys_9.dlp" to profile slot 3
    d-win pump -f _sys_9 -s 3 -pro
- Upload bank of presets in bank file "mybank.bnk" to preset slots 10, 11, 12, etc.:
    d-win bank -f mybank -s 10
- Download software & all presets & profiles to file "2022-01-23.eeprom":
    d-win dump -f 2022-01-23.eeprom
- Upload software & all presets & profiles from file "factory.eeprom":
    d-win pump -f factory.eeprom
- Download software to file "sw_backup.spi":
    d-win dump -f sw_backup.spi
- Upload software from file "f9e1c5c7.spi":
    d-win pump -f f9e1c5c7.spi
- Download all presets to file "old_presets.pre":
    d-win dump -f old_presets.pre
- Upload all presets from file "my_dlev.pre":
    d-win pump -f my_dlev.pre
- Download all profiles to file "my_setup.pro":
    d-win dump -f my_setup.pro
- Upload all profiles from file "your_setup.pro":
    d-win pump -f your_setup.pro
- Split file "some.eeprom" into "some.pre", "some.pro", some.spi":
    d-win split -f some.eeprom
- Split file "my_setup.pro" into "pro_000.dlp" thru "pro_005.dlp":
    d-win split -f my_setup.pro
- Split file "my_new.pre" into "000.dlp" thru "249.dlp":
    d-win split -f my_new.pre
- Join files "some.pre", "some.pro" and "some.spi" into "some.eeprom":
    d-win join -f some.eeprom
- Join files "pro_000.dlp" thru "pro_005.dlp" to "stuff.pro":
    d-win join -f stuff.pro
- Join files "000.dlp" thru "249.dlp" to "some.pre":
    d-win join -f some.pre
- Morph knobs (osc):
    d-win morph -mo 12
- Morph slot 23 (filters, resonator, seed):
    d-win morph -s 23 -mf 5 -mr 20 -i 9
- Morph file "cello_8" (osc, filters, resonator):
    d-win morph -f cello_8 -mo 10 -mf 10 -mr 10
- Batch convert all presets in the _ALL_ directory to mono in the _MONO_ directory:
    d-win batch -d _ALL_ -d2 _MONO_ -m
- Batch update all presets in the _ALL_ directory and overwrite them:
    d-win batch -d _ALL_ -d2 _ALL_ -u
- Read knob RESON:mode:
    d-win knob -k re:4
- Set knob RESON:mode to 10 & view all knobs:
    d-win knob -k re:4 -s 10 -v
- Inc knob 1_FORM:reso by 4:
    d-win knob -k 1_f:6 -o 4
- Dec knob FLT_OSC:xmix by -2:
    d-win knob -k flt_o:5 -o -2
- Minimize all FLT_NOISE knobs:
    d-win knob -k flt_n:all -min
- Read the installed software version & check CRC:
    d-win ver
- Read file "2020.eeprom" software version & check CRC:
    d-win ver -f 2020.eeprom
- Read file "some.spi" software version & check CRC:
    d-win ver -f some.spi
- Read processor registers 0 thru 9:
    d-win hcl 0 9 rr
- Loop back serial port text "testing 123":
    d-win loop "testing 123"
- Perform an ACAL:
    d-win acal
- Reset the processor:
    d-win reset
`