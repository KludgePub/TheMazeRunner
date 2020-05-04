package maze

// map single props
const (
	verticalWall = '|'
	horizontalWall = '-'
    corner = '+'
    emptySpace = ' '
    startingPoint = 'S'
    endingPoint = 'E'
    keyPoint = 'K'
	// TODO doors, traps, items?
)

// map tiles
var horizontalWallTile = []byte("+---")
var horizontalOpenTile = []byte("+   ")
var verticalWallTile = []byte("|   ")
var verticalOpenTile = []byte("    ")
