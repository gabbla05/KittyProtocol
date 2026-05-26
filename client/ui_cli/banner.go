package ui_cli

import "fmt"

func PrintBanner() {

	cat := []string{
		ColorPink2 + "                                             \\`*-.                 ",
		ColorPink2 + "                                              )  _`-.                ",
		ColorPink2 + "           ^ ^                               .  : `. .               ",
		ColorPink2 + "░█▄ ▄█░█▀▀░█▀█░█░░░█░                        : _   '  \\              ",
		ColorPink2 + "░█░▀░█░█▀▀░█░█░█▄▀▄█░                        ; *` _.   `*-._         ",
		ColorPink1 + "░▀░░░▀░▀▀▀░▀▀▀░▀░░░▀░                        `-.-'          `-.      ",
		ColorPink1 + "░█▀▀░█▀▀░█▀▀░█▀█░█▀▀░█▀▀░█▀▄                   ;       `       `.    ",
		ColorPink3 + "░▀▀█░▀▀█░█▀▀░█░█░█░█░█▀▀░█▀▄                   :.       .        \\   ",
		ColorPink3 + "░▀▀▀░▀▀▀░▀▀▀░▀░▀░▀▀▀░▀▀▀░▀░▀                   . \\  .   :   .-'   .  ",
		ColorPink2 + "Powered by KITTYPROTOCOL                       '  `+.;  ;  '      :  ",
		ColorPink2 + "                _            _        __       :  '  |    ;       ;-.",
		ColorPink2 + "__ _____ _ _ __(_)___ _ _   / |      /  \\      ; '   : :`-:     _.`* ;",
		ColorPink2 + "\\ V / -_) '_(_-< / _ \\ ' \\  | |  _  | () |  .*' /  .*' ; .*`- +'  `*'",
		ColorPink2 + " \\_/\\___|_| /__/_\\___/_||_| |_| (_)  \\__/   `*-*   `*-*  `*-*'" + ColorReset,
	}

	for _, line := range cat {
		fmt.Println(line)
	}
}
