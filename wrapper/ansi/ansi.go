package wrapper_ansi

const ResetAttrs = "\x1b[0m"
const ResetCursor = "\x1b[H"
const CleanConsole = "\x1B[2J\x1B[H"

const EraseLine = "\r\033[K"

const HideCursor = "\x1b[?25l"
const ShowCursor = "\x1b[?25h"
