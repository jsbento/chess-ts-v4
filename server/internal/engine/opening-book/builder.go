package openingbook

import (
	"strings"

	e "github.com/jsbento/chess-server-v4/internal/engine"
	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
	t "github.com/jsbento/chess-server-v4/internal/engine/types"
	"github.com/jsbento/chess-server-v4/internal/engine/utils"
)

// parseSANMove converts a SAN move string (e.g., "d4", "Nf6") to engine move format
// This is a simplified parser that handles most common cases
func parseSANMove(san string, engine *e.Engine) int {
	if len(san) == 0 {
		return c.NOMOVE
	}

	// Handle special moves
	if san == "O-O" || san == "0-0" || san == "Kingside" {
		return parseCastleMove(engine, true)
	}
	if san == "O-O-O" || san == "0-0-0" || san == "Queenside" {
		return parseCastleMove(engine, false)
	}

	// Generate all legal moves
	moveList := t.NewMoveList()
	engine.GenerateAllMoves(moveList)

	// Try to match the SAN move
	for i := 0; i < moveList.Count; i++ {
		move := moveList.Moves[i].Move
		if !engine.MakeMove(move) {
			continue
		}
		engine.TakeMove()

		// Convert move to SAN-like format and compare
		if matchesSAN(san, move, engine) {
			return move
		}
	}

	return c.NOMOVE
}

// parseCastleMove finds the castling move
func parseCastleMove(engine *e.Engine, kingside bool) int {
	moveList := t.NewMoveList()
	engine.GenerateAllMoves(moveList)

	for i := 0; i < moveList.Count; i++ {
		move := moveList.Moves[i].Move
		if move&int(c.MFLAGCA) == 0 {
			continue
		}
		if !engine.MakeMove(move) {
			continue
		}
		engine.TakeMove()

		to := utils.ToSq(move)
		if kingside && (to == c.G1 || to == c.G8) {
			return move
		}
		if !kingside && (to == c.C1 || to == c.C8) {
			return move
		}
	}

	return c.NOMOVE
}

// matchesSAN checks if a move matches the SAN notation
func matchesSAN(san string, move int, engine *e.Engine) bool {
	to := utils.ToSq(move)
	piece := engine.Board.Pieces[utils.FromSq(move)]

	// Extract destination square from SAN (usually last 2 characters before promotion/castling)
	destSquare := ""
	sanUpper := strings.ToUpper(san)

	// Handle castling separately
	if sanUpper == "O-O" || sanUpper == "0-0" {
		return move&int(c.MFLAGCA) != 0 && (to == c.G1 || to == c.G8)
	}
	if sanUpper == "O-O-O" || sanUpper == "0-0-0" {
		return move&int(c.MFLAGCA) != 0 && (to == c.C1 || to == c.C8)
	}

	// Find destination square (last 2 alphanumeric characters before = or # or +)
	destIdx := len(san) - 1
	for destIdx >= 0 && !isSquareChar(san[destIdx]) {
		destIdx--
	}
	if destIdx >= 1 {
		destSquare = strings.ToLower(san[destIdx-1 : destIdx+1])
	}

	// Check if destination matches
	toFile := c.FilesBrd[to]
	toRank := c.RanksBrd[to]
	expectedSquare := string(rune('a'+toFile)) + string(rune('1'+toRank))
	if destSquare != expectedSquare {
		return false
	}

	// Check piece type (first character if uppercase)
	if len(san) > 0 {
		firstChar := string(san[0])
		pieceChar := getPieceChar(piece, engine.Board.Side)

		// If first char is uppercase letter (excluding file letters), it should match piece type
		firstCharUpper := strings.ToUpper(firstChar)
		if firstCharUpper >= "A" && firstCharUpper <= "Z" && firstCharUpper != "O" {
			// Check if it's a piece notation (N, B, R, Q, K) or a file letter (a-h for pawn moves)
			if firstChar >= "A" && firstChar <= "Z" {
				// It's a piece indicator (uppercase), should match piece type
				if firstCharUpper != pieceChar {
					return false
				}
			} else if firstChar >= "a" && firstChar <= "h" {
				// It's a file letter (lowercase), must be a pawn move
				if pieceChar != "P" {
					return false
				}
			}
		}
	}

	// Handle captures
	if strings.Contains(san, "x") {
		captured := utils.Captured(move)
		if captured == c.EMPTY && move&int(c.MFLAGEP) == 0 {
			return false
		}
	}

	// Handle promotion
	if strings.Contains(san, "=") {
		promoted := utils.Promoted(move)
		if promoted == c.EMPTY {
			return false
		}
		// Check promotion piece
		if strings.Contains(san, "=Q") && promoted != c.WQ && promoted != c.BQ {
			return false
		}
		if strings.Contains(san, "=R") && !utils.IsRQ(int(promoted)) {
			return false
		}
		if strings.Contains(san, "=B") && !utils.IsBQ(int(promoted)) {
			return false
		}
		if strings.Contains(san, "=N") && !utils.IsKn(int(promoted)) {
			return false
		}
	}

	return true
}

// isSquareChar checks if a character is part of a square notation (a-h, 1-8)
func isSquareChar(c byte) bool {
	return (c >= 'a' && c <= 'h') || (c >= '1' && c <= '8')
}

// getPieceChar returns the character representation of a piece
func getPieceChar(piece c.Piece, side c.Side) string {
	switch piece {
	case c.WP, c.BP:
		return "P"
	case c.WN, c.BN:
		return "N"
	case c.WB, c.BB:
		return "B"
	case c.WR, c.BR:
		return "R"
	case c.WQ, c.BQ:
		return "Q"
	case c.WK, c.BK:
		return "K"
	default:
		return "P"
	}
}

// BuildOpeningBookFromTSV builds an opening book graph from TSV entries
func BuildOpeningBookFromTSV(entries []TSVEntry, engine *e.Engine) (*OpeningBook, error) {
	book := newEmptyOpeningBook()

	// Initialize engine to starting position
	if err := engine.ParseFEN(c.START_FEN); err != nil {
		return nil, err
	}
	engine.GeneratePosKey()
	book.RootKey = engine.Board.PosKey

	// Process each entry
	for _, entry := range entries {
		// Reset to starting position for each entry
		if err := engine.ParseFEN(c.START_FEN); err != nil {
			continue // Skip invalid entries
		}
		engine.GeneratePosKey()

		// Traverse moves in this entry
		for moveIdx, movePair := range entry.Moves {
			// Parse and make white move
			if movePair.WhiteMove != "" {
				fromPosKey := engine.Board.PosKey
				whiteMove := parseSANMove(movePair.WhiteMove, engine)
				if whiteMove == c.NOMOVE {
					break // Invalid move, skip this entry
				}
				if !engine.MakeMove(whiteMove) {
					break // Illegal move, skip this entry
				}
				engine.GeneratePosKey()
				toPosKey := engine.Board.PosKey

				// Create or update nodes and edge
				book.addNodeIfNotExists(fromPosKey)
				book.addNodeIfNotExists(toPosKey)

				// Add edge with move in UCI format
				moveUCI := utils.PrintMove(whiteMove)
				book.addEdge(fromPosKey, toPosKey, moveUCI, entry.ECO, entry.Name)

				// If this is the last move and there's no black move, store metadata
				if moveIdx == len(entry.Moves)-1 && movePair.BlackMove == "" {
					book.addMetadata(toPosKey, entry.ECO, entry.Name)
				}
			}

			// Parse and make black move if present
			if movePair.BlackMove != "" {
				fromPosKey := engine.Board.PosKey
				blackMove := parseSANMove(movePair.BlackMove, engine)
				if blackMove == c.NOMOVE {
					break // Invalid move, skip this entry
				}
				if !engine.MakeMove(blackMove) {
					break // Illegal move, skip this entry
				}
				engine.GeneratePosKey()
				toPosKey := engine.Board.PosKey

				// Create or update nodes and edge
				book.addNodeIfNotExists(fromPosKey)
				book.addNodeIfNotExists(toPosKey)

				// Add edge with move in UCI format
				moveUCI := utils.PrintMove(blackMove)
				book.addEdge(fromPosKey, toPosKey, moveUCI, entry.ECO, entry.Name)

				// If this is the last move, store metadata
				if moveIdx == len(entry.Moves)-1 {
					book.addMetadata(toPosKey, entry.ECO, entry.Name)
				}
			}
		}
	}

	return book, nil
}

// addNodeIfNotExists creates a node if it doesn't exist
func (ob *OpeningBook) addNodeIfNotExists(posKey uint64) {
	if _, exists := ob.Vertices[posKey]; !exists {
		ob.Vertices[posKey] = &OpeningBookNode{
			PositionKey:  posKey,
			OpeningNames: []string{},
			ECOCodes:     []string{},
			Edges:        make(map[string]*OpeningBookEdge),
		}
	}
}

// addEdge adds or updates an edge between two nodes
func (ob *OpeningBook) addEdge(fromPosKey, toPosKey uint64, moveUCI, eco, name string) {
	fromNode, exists := ob.Vertices[fromPosKey]
	if !exists {
		return
	}

	toNode, exists := ob.Vertices[toPosKey]
	if !exists {
		return
	}

	// Check if edge already exists
	if edge, exists := fromNode.Edges[moveUCI]; exists {
		// Update existing edge
		edge.Frequency++
		if !contains(edge.OpeningNames, name) {
			edge.OpeningNames = append(edge.OpeningNames, name)
		}
	} else {
		// Create new edge
		fromNode.Edges[moveUCI] = &OpeningBookEdge{
			Move:         moveUCI,
			ToNode:       toNode,
			Frequency:    1,
			OpeningNames: []string{name},
		}
	}
}

// addMetadata adds ECO code and opening name to a node
func (ob *OpeningBook) addMetadata(posKey uint64, eco, name string) {
	node, exists := ob.Vertices[posKey]
	if !exists {
		return
	}

	if !contains(node.ECOCodes, eco) {
		node.ECOCodes = append(node.ECOCodes, eco)
	}
	if !contains(node.OpeningNames, name) {
		node.OpeningNames = append(node.OpeningNames, name)
	}
}

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
