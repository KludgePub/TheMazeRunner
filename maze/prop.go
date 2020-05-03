package maze

const verticalWall = '|'
const horizontalWall = '-'
const corner = '+'
const emptySpace = ' '
const startingPoint = 'S'
const endingPoint = 'E'
// TODO doors, traps, items?

var horizontalWallTile = []byte("+---")
var horizontalOpenTile = []byte("+   ")
var verticalWallTile = []byte("|   ")
var verticalOpenTile = []byte("    ")
