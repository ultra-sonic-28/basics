package lexer

var Keywords = map[string]bool{
	// Contrôle
	"FOR": true, "TO": true, "STEP": true, "NEXT": true,
	"IF": true, "THEN": true, "ELSE": true,
	"GOTO": true, "GOSUB": true, "RETURN": true,
	"END": true, "STOP": true,

	// Variables & logique
	"LET": true, "DIM": true,
	"REM": true,

	// I/O
	"PRINT": true, "INPUT": true,
	"GET": true,

	// Math
	"SIN": true, "COS": true, "TAN": true,
	"INT": true, "ABS": true, "RND": true,
	"SGN": true,

	// Graphique / écran
	"GR": true, "HGR": true, "TEXT": true,
	"PLOT": true, "HPLOT": true,
	"COLOR": true, "HCOLOR": true,
	"HOME": true,

	// DATA
	"DATA": true, "READ": true, "RESTORE": true,

	// Autres
	"POKE": true, "PEEK": true, "CALL": true,
	"TAB": true, "VTAB": true, "HTAB": true,
	"INVERSE": true, "NORMAL": true, "FLASH": true,

	// Extension
	"SLEEP": true,
}
