package lexer

var Keywords = map[string]bool{
	// Contrôle
	"FOR": true, "TO": true, "STEP": true, "NEXT": true,
	"IF": true, "THEN": true, "ELSE": true,
	"GOTO": true, "GOSUB": true, "RETURN": true,
	"END": true, "STOP": true,

	// Variables & logique
	"LET": true, "DIM": true,
	"AND": true,
	"OR":  true,
	"NOT": true,
	"REM": true, "CLEAR": true,

	// I/O
	"PRINT": true,
	"INPUT": true, "GET": true,

	// Opérations sur les chaines de caractères
	"LEFT$": true, "RIGHT$": true, "MID$": true,
	"STR$": true,
	"CHR$": true,
	"ASC":  true, "VAL": true,
	"LEN": true,

	// Math
	"SIN": true, "COS": true, "TAN": true, "ATN": true,
	"INT": true, "ABS": true, "RND": true,
	"SGN": true, "SQR": true,

	// Graphique / écran
	"GR": true, "HGR": true, "TEXT": true,
	"PLOT": true, "HPLOT": true,
	"COLOR": true, "HCOLOR": true,
	"HOME": true,

	// DATA
	"DATA": true, "READ": true, "RESTORE": true,

	// Autres
	"POKE": true, "PEEK": true, "CALL": true,
	"PR":  true,
	"TAB": true, "VTAB": true, "HTAB": true,
	"SPC":     true,
	"INVERSE": true, "NORMAL": true, "FLASH": true,

	// Extension
	"WAIT": true,
}
