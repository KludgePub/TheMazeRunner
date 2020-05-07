package asset

// map single props
const (
	VerticalWall   = '|'
	HorizontalWall = '-'
	Corner         = '+'
	EmptySpace     = ' '
	StartingPoint  = 'S'
	EndingPoint    = 'E'
	KeyPoint       = 'K'
)

// map tiles
var HorizontalWallTile = []byte("+---")
var HorizontalOpenTile = []byte("+   ")
var VerticalWallTile = []byte("|   ")
var VerticalOpenTile = []byte("    ")
