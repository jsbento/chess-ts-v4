package openingbook

type OpeningBook struct {
	Vertices map[uint64]*OpeningBookNode // Key: position key (PosKey)
	RootKey  uint64                      // Starting position key
}

type OpeningBookNode struct {
	PositionKey  uint64                      // Position identifier
	OpeningNames []string                    // All opening names leading here
	ECOCodes     []string                    // All ECO codes
	Edges        map[string]*OpeningBookEdge // Key: move in UCI format
}

type OpeningBookEdge struct {
	Move         string // Move in UCI format (e.g., "e2e4")
	ToNode       *OpeningBookNode
	Frequency    int      // How many times this move appears
	OpeningNames []string // Opening names using this move
}

// newEmptyOpeningBook creates a new empty opening book
func newEmptyOpeningBook() *OpeningBook {
	return &OpeningBook{
		Vertices: make(map[uint64]*OpeningBookNode),
		RootKey:  0,
	}
}

// GetMoves returns available moves from a position
func (ob *OpeningBook) GetMoves(positionKey uint64) []string {
	node, exists := ob.Vertices[positionKey]
	if !exists {
		return nil
	}

	moves := make([]string, 0, len(node.Edges))
	for move := range node.Edges {
		moves = append(moves, move)
	}
	return moves
}

// GetOpeningName returns opening names for a position
func (ob *OpeningBook) GetOpeningName(positionKey uint64) []string {
	node, exists := ob.Vertices[positionKey]
	if !exists {
		return nil
	}
	return node.OpeningNames
}

// GetECO returns ECO codes for a position
func (ob *OpeningBook) GetECO(positionKey uint64) []string {
	node, exists := ob.Vertices[positionKey]
	if !exists {
		return nil
	}
	return node.ECOCodes
}

// InBook checks if a position exists in the opening book
func (ob *OpeningBook) InBook(positionKey uint64) bool {
	_, exists := ob.Vertices[positionKey]
	return exists
}
