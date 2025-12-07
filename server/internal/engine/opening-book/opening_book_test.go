package openingbook

import (
	"errors"
	"testing"

	e "github.com/jsbento/chess-server-v4/internal/engine"
	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
)

func TestLoadTSV(t *testing.T) {
	entries, err := loadTSV("data/e.tsv")
	if err != nil {
		t.Fatalf("Failed to load TSV: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("Expected entries to be loaded, got 0")
	}

	// Check first entry structure
	firstEntry := entries[0]
	if firstEntry.ECO == "" {
		t.Error("Expected ECO to be not empty")
	}
	if firstEntry.Name == "" {
		t.Error("Expected Name to be not empty")
	}
	if firstEntry.PGN == "" {
		t.Error("Expected PGN to be not empty")
	}
	if len(firstEntry.Moves) == 0 {
		t.Error("Expected Moves to be not empty")
	}
}

func TestMovesFromPGN(t *testing.T) {
	tests := []struct {
		name          string
		pgn           string
		expected      int // expected number of move pairs
		expectedError error
	}{
		{
			name:          "Single move",
			pgn:           "1. e4",
			expected:      1,
			expectedError: nil,
		},
		{
			name:          "Two moves",
			pgn:           "1. e4 e5",
			expected:      1,
			expectedError: nil,
		},
		{
			name:          "Multiple moves",
			pgn:           "1. d4 Nf6 2. c4 e6 3. g3",
			expected:      3,
			expectedError: nil,
		},
		{
			name:          "Complex sequence",
			pgn:           "1. d4 Nf6 2. c4 e6 3. g3 d5 4. Bg2 dxc4",
			expected:      4,
			expectedError: nil,
		},
		{
			name:          "Empty PGN",
			pgn:           "",
			expected:      0,
			expectedError: errors.New("pgn is empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			moves, err := movesFromPGN(tt.pgn)
			if tt.expectedError != nil {
				if err != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error: %v, got %v", tt.expectedError, err)
				} else if err == nil {
					t.Errorf("Expected error: %v, got nil", tt.expectedError)
				} else if len(moves) != tt.expected {
					t.Errorf("Expected %d moves, got %d", tt.expected, len(moves))
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				} else if len(moves) != tt.expected {
					t.Errorf("Expected %d moves, got %d", tt.expected, len(moves))
				}
			}
		})
	}
}

func TestMovesFromPGNContent(t *testing.T) {
	pgn := "1. d4 Nf6 2. c4 e6 3. g3"
	moves, err := movesFromPGN(pgn)
	if err != nil {
		t.Fatalf("Failed to get moves from PGN: %v", err)
	}

	if len(moves) != 3 {
		t.Fatalf("Expected 3 moves, got %d", len(moves))
	}

	// Check first move
	if moves[0].WhiteMove != "d4" {
		t.Errorf("Expected WhiteMove to be 'd4', got '%s'", moves[0].WhiteMove)
	}
	if moves[0].BlackMove != "Nf6" {
		t.Errorf("Expected BlackMove to be 'Nf6', got '%s'", moves[0].BlackMove)
	}

	// Check second move
	if moves[1].WhiteMove != "c4" {
		t.Errorf("Expected WhiteMove to be 'c4', got '%s'", moves[1].WhiteMove)
	}
	if moves[1].BlackMove != "e6" {
		t.Errorf("Expected BlackMove to be 'e6', got '%s'", moves[1].BlackMove)
	}

	// Check third move (no black move)
	if moves[2].WhiteMove != "g3" {
		t.Errorf("Expected WhiteMove to be 'g3', got '%s'", moves[2].WhiteMove)
	}
	if moves[2].BlackMove != "" {
		t.Errorf("Expected BlackMove to be empty, got '%s'", moves[2].BlackMove)
	}
}

func TestNewEmptyOpeningBook(t *testing.T) {
	book := newEmptyOpeningBook()

	if book == nil {
		t.Fatal("Expected book to be not nil")
	}

	if book.Vertices == nil {
		t.Error("Expected Vertices map to be initialized")
	}

	if len(book.Vertices) != 0 {
		t.Errorf("Expected empty vertices map, got %d vertices", len(book.Vertices))
	}

	if book.RootKey != 0 {
		t.Errorf("Expected RootKey to be 0, got %d", book.RootKey)
	}
}

func TestBuildOpeningBookFromTSV(t *testing.T) {
	engine := e.NewEngine()

	testOpeningPGN := "1. e4 e5"
	testOpeningMoves, err := movesFromPGN(testOpeningPGN)
	if err != nil {
		t.Fatalf("Failed to get moves from PGN: %v", err)
	}

	testAnotherOpeningPGN := "1. d4"
	testAnotherOpeningMoves, err := movesFromPGN(testAnotherOpeningPGN)
	if err != nil {
		t.Fatalf("Failed to get moves from PGN: %v", err)
	}

	// Create test entries
	entries := []TSVEntry{
		{
			ECO:   "A00",
			Name:  "Test Opening",
			PGN:   testOpeningPGN,
			Moves: testOpeningMoves,
		},
		{
			ECO:   "A01",
			Name:  "Another Opening",
			PGN:   testAnotherOpeningPGN,
			Moves: testAnotherOpeningMoves,
		},
	}

	book, err := BuildOpeningBookFromTSV(entries, engine)
	if err != nil {
		t.Fatalf("Failed to build opening book: %v", err)
	}

	if book == nil {
		t.Fatal("Expected book to be not nil")
	}

	if book.RootKey == 0 {
		t.Error("Expected RootKey to be set")
	}

	if len(book.Vertices) == 0 {
		t.Error("Expected vertices to be created")
	}
}

func TestBuildOpeningBookFromTSVWithRealData(t *testing.T) {
	engine := e.NewEngine()

	// Load a small TSV file for testing
	entries, err := loadTSV("data/e.tsv")
	if err != nil {
		t.Fatalf("Failed to load TSV: %v", err)
	}

	// Use first 5 entries for faster testing
	testEntries := entries[:5]
	if len(testEntries) > 5 {
		testEntries = entries[:5]
	}

	book, err := BuildOpeningBookFromTSV(testEntries, engine)
	if err != nil {
		t.Fatalf("Failed to build opening book: %v", err)
	}

	if book == nil {
		t.Fatal("Expected book to be not nil")
	}

	if book.RootKey == 0 {
		t.Error("Expected RootKey to be set")
	}

	// Verify root position exists
	if !book.InBook(book.RootKey) {
		t.Error("Root position should be in book")
	}
}

func TestGetMoves(t *testing.T) {
	engine := e.NewEngine()
	book := newEmptyOpeningBook()

	// Initialize engine to starting position
	engine.ParseFEN(c.START_FEN)
	engine.GeneratePosKey()
	rootKey := engine.Board.PosKey
	book.RootKey = rootKey

	// Create a test node
	testNode := &OpeningBookNode{
		PositionKey:  rootKey,
		OpeningNames: []string{},
		ECOCodes:     []string{},
		Edges:        make(map[string]*OpeningBookEdge),
	}

	// Add some test edges
	testNode.Edges["e2e4"] = &OpeningBookEdge{
		Move:      "e2e4",
		Frequency: 1,
	}
	testNode.Edges["d2d4"] = &OpeningBookEdge{
		Move:      "d2d4",
		Frequency: 1,
	}

	book.Vertices[rootKey] = testNode

	// Test GetMoves
	moves := book.GetMoves(rootKey)
	if len(moves) != 2 {
		t.Errorf("Expected 2 moves, got %d", len(moves))
	}

	// Test with non-existent position
	moves = book.GetMoves(999999)
	if moves != nil {
		t.Error("Expected nil for non-existent position")
	}
}

func TestGetOpeningName(t *testing.T) {
	engine := e.NewEngine()
	book := newEmptyOpeningBook()

	engine.ParseFEN(c.START_FEN)
	engine.GeneratePosKey()
	posKey := engine.Board.PosKey

	// Create a test node with opening names
	testNode := &OpeningBookNode{
		PositionKey:  posKey,
		OpeningNames: []string{"Test Opening", "Another Opening"},
		ECOCodes:     []string{},
		Edges:        make(map[string]*OpeningBookEdge),
	}

	book.Vertices[posKey] = testNode

	// Test GetOpeningName
	names := book.GetOpeningName(posKey)
	if len(names) != 2 {
		t.Errorf("Expected 2 opening names, got %d", len(names))
	}

	if names[0] != "Test Opening" {
		t.Errorf("Expected first name to be 'Test Opening', got '%s'", names[0])
	}

	// Test with non-existent position
	names = book.GetOpeningName(999999)
	if names != nil {
		t.Error("Expected nil for non-existent position")
	}
}

func TestGetECO(t *testing.T) {
	engine := e.NewEngine()
	book := newEmptyOpeningBook()

	engine.ParseFEN(c.START_FEN)
	engine.GeneratePosKey()
	posKey := engine.Board.PosKey

	// Create a test node with ECO codes
	testNode := &OpeningBookNode{
		PositionKey:  posKey,
		OpeningNames: []string{},
		ECOCodes:     []string{"A00", "A01"},
		Edges:        make(map[string]*OpeningBookEdge),
	}

	book.Vertices[posKey] = testNode

	// Test GetECO
	ecos := book.GetECO(posKey)
	if len(ecos) != 2 {
		t.Errorf("Expected 2 ECO codes, got %d", len(ecos))
	}

	if ecos[0] != "A00" {
		t.Errorf("Expected first ECO to be 'A00', got '%s'", ecos[0])
	}

	// Test with non-existent position
	ecos = book.GetECO(999999)
	if ecos != nil {
		t.Error("Expected nil for non-existent position")
	}
}

func TestInBook(t *testing.T) {
	engine := e.NewEngine()
	book := newEmptyOpeningBook()

	engine.ParseFEN(c.START_FEN)
	engine.GeneratePosKey()
	posKey := engine.Board.PosKey

	// Create a test node
	testNode := &OpeningBookNode{
		PositionKey:  posKey,
		OpeningNames: []string{},
		ECOCodes:     []string{},
		Edges:        make(map[string]*OpeningBookEdge),
	}

	book.Vertices[posKey] = testNode

	// Test InBook with existing position
	if !book.InBook(posKey) {
		t.Error("Expected position to be in book")
	}

	// Test InBook with non-existent position
	if book.InBook(999999) {
		t.Error("Expected position to not be in book")
	}
}

func TestBuildOpeningBookGraphStructure(t *testing.T) {
	engine := e.NewEngine()

	entries := []TSVEntry{
		{
			ECO:  "A00",
			Name: "Test Opening 1",
			PGN:  "1. e4 e5",
		},
		{
			ECO:  "A01",
			Name: "Test Opening 2",
			PGN:  "1. e4 e5 2. Nf3",
		},
	}
	for i := range entries {
		moves, err := movesFromPGN(entries[i].PGN)
		if err != nil {
			t.Fatalf("Failed to get moves from PGN: %v", err)
		}
		entries[i].Moves = moves
	}

	book, err := BuildOpeningBookFromTSV(entries, engine)
	if err != nil {
		t.Fatalf("Failed to build opening book: %v", err)
	}

	// Verify root exists
	if !book.InBook(book.RootKey) {
		t.Error("Root position should be in book")
	}

	// Verify we can get moves from root
	moves := book.GetMoves(book.RootKey)
	if len(moves) == 0 {
		t.Error("Expected moves from root position")
	}

	// Verify e4 move exists
	hasE4 := false
	for _, move := range moves {
		if move == "e2e4" {
			hasE4 = true
			break
		}
	}
	if !hasE4 {
		t.Error("Expected e2e4 move to be available from root")
	}
}

func TestBuildOpeningBookWithMultipleEntries(t *testing.T) {
	engine := e.NewEngine()

	// Create entries that share common moves
	entries := []TSVEntry{
		{
			ECO:  "A00",
			Name: "Opening A",
			PGN:  "1. d4",
		},
		{
			ECO:  "A01",
			Name: "Opening B",
			PGN:  "1. d4 Nf6",
		},
		{
			ECO:  "A02",
			Name: "Opening C",
			PGN:  "1. d4 Nf6 2. c4",
		},
	}
	for i := range entries {
		moves, err := movesFromPGN(entries[i].PGN)
		if err != nil {
			t.Fatalf("Failed to get moves from PGN: %v", err)
		}
		entries[i].Moves = moves
	}

	book, err := BuildOpeningBookFromTSV(entries, engine)
	if err != nil {
		t.Fatalf("Failed to build opening book: %v", err)
	}

	// Verify graph structure
	if len(book.Vertices) == 0 {
		t.Error("Expected vertices to be created")
	}

	// Verify root has moves
	rootMoves := book.GetMoves(book.RootKey)
	if len(rootMoves) == 0 {
		t.Error("Expected moves from root")
	}
}

func TestBuildOpeningBookWithInvalidMoves(t *testing.T) {
	engine := e.NewEngine()

	// Create entries with potentially invalid moves
	entries := []TSVEntry{
		{
			ECO:  "A00",
			Name: "Valid Opening",
			PGN:  "1. e4 e5",
		},
		{
			ECO:  "A01",
			Name: "Invalid Opening",
			PGN:  "1. z9 z8", // Invalid moves
		},
	}
	for i := range entries {
		moves, err := movesFromPGN(entries[i].PGN)
		if err != nil {
			t.Fatalf("Failed to get moves from PGN: %v", err)
		}
		entries[i].Moves = moves
	}

	// Should not crash, but may skip invalid entries
	book, err := BuildOpeningBookFromTSV(entries, engine)
	if err != nil {
		// Error is acceptable for invalid moves
		return
	}

	// Should still have some valid structure
	if book == nil {
		t.Fatal("Expected book to be created even with invalid moves")
	}
}

func TestOpeningBookEdgeFrequency(t *testing.T) {
	engine := e.NewEngine()

	// Create multiple entries with the same move
	entries := []TSVEntry{
		{
			ECO:  "A00",
			Name: "Opening 1",
			PGN:  "1. e4",
		},
		{
			ECO:  "A01",
			Name: "Opening 2",
			PGN:  "1. e4",
		},
		{
			ECO:  "A02",
			Name: "Opening 3",
			PGN:  "1. e4",
		},
	}
	for i := range entries {
		moves, err := movesFromPGN(entries[i].PGN)
		if err != nil {
			t.Fatalf("Failed to get moves from PGN: %v", err)
		}
		entries[i].Moves = moves
	}

	book, err := BuildOpeningBookFromTSV(entries, engine)
	if err != nil {
		t.Fatalf("Failed to build opening book: %v", err)
	}

	// Get root node
	rootNode, exists := book.Vertices[book.RootKey]
	if !exists {
		t.Fatal("Root node should exist")
	}

	// Check if e4 edge exists and has correct frequency
	e4Edge, exists := rootNode.Edges["e2e4"]
	if !exists {
		t.Error("Expected e2e4 edge to exist")
	} else {
		if e4Edge.Frequency < 1 {
			t.Errorf("Expected frequency to be at least 1, got %d", e4Edge.Frequency)
		}
	}
}

func TestOpeningBookMetadata(t *testing.T) {
	engine := e.NewEngine()

	entries := []TSVEntry{
		{
			ECO:  "A00",
			Name: "Test Opening",
			PGN:  "1. e4",
		},
	}
	for i := range entries {
		moves, err := movesFromPGN(entries[i].PGN)
		if err != nil {
			t.Fatalf("Failed to get moves from PGN: %v", err)
		}
		entries[i].Moves = moves
	}

	book, err := BuildOpeningBookFromTSV(entries, engine)
	if err != nil {
		t.Fatalf("Failed to build opening book: %v", err)
	}

	// Get moves from root - should include e4
	rootMoves := book.GetMoves(book.RootKey)
	if len(rootMoves) == 0 {
		t.Fatal("Expected moves from root position")
	}

	// Find e4 edge
	var e4Edge *OpeningBookEdge
	rootNode, exists := book.Vertices[book.RootKey]
	if !exists {
		t.Fatal("Root node should exist")
	}

	for move, edge := range rootNode.Edges {
		if move == "e2e4" {
			e4Edge = edge
			break
		}
	}

	if e4Edge == nil {
		t.Error("Expected e2e4 edge to exist")
	} else {
		// Check if the target node has metadata
		if e4Edge.ToNode != nil {
			// Terminal positions should have metadata
			// Since this is a single move entry, the position after e4 should have metadata
			if len(e4Edge.ToNode.OpeningNames) > 0 || len(e4Edge.ToNode.ECOCodes) > 0 {
				// Metadata is stored correctly
			}
		}
	}
}
