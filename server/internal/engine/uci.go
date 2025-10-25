package engine

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
	t "github.com/jsbento/chess-server-v4/internal/engine/types"
	"github.com/jsbento/chess-server-v4/pkg/utils"
)

func (e *Engine) ParseGo(line string, info *t.SearchInfo) string {
	depth, movesToGo, moveTime := -1, 30, -1
	time, inc := -1, 0

	info.TimeSet = false

	tokens := strings.Split(line, " ")
	for i := 0; i < len(tokens); i++ {
		switch tokens[i] {
		case "infinite":
			continue
		case "binc":
			if e.Board.Side == c.BLACK {
				inc, _ = strconv.Atoi(tokens[i+1])
				i++
			}
		case "winc":
			if e.Board.Side == c.WHITE {
				inc, _ = strconv.Atoi(tokens[i+1])
				i++
			}
		case "wtime":
			if e.Board.Side == c.WHITE {
				time, _ = strconv.Atoi(tokens[i+1])
				i++
			}
		case "btime":
			if e.Board.Side == c.BLACK {
				time, _ = strconv.Atoi(tokens[i+1])
				i++
			}
		case "movestogo":
			movesToGo, _ = strconv.Atoi(tokens[i+1])
			i++
		case "depth":
			depth, _ = strconv.Atoi(tokens[i+1])
			i++
		case "movetime":
			moveTime, _ = strconv.Atoi(tokens[i+1])
			i++
		default:
			continue
		}
	}

	if moveTime != -1 {
		time = moveTime
		movesToGo = 1
	}
	info.StartTime = utils.GetTimeMs()
	info.Depth = depth

	if time != -1 {
		info.TimeSet = true
		time /= movesToGo
		time -= 50
		info.StopTime = info.StartTime + int64(time) + int64(inc)
	}

	if depth == -1 {
		info.Depth = c.MAX_DEPTH
	}

	// fmt.Printf("time:%d start:%d stop:%d depth:%d timeset:%t\n", time, info.StartTime, info.StopTime, info.Depth, info.TimeSet)
	return e.SearchPosition(info)
}

func (e *Engine) ParsePosition(line string) error {
	commands := strings.Split(line, " moves ")
	posCommand := commands[0][9:]

	if strings.Contains(posCommand, "startpos") {
		e.ParseFEN(c.START_FEN)
	} else if strings.Contains(posCommand, "fen") {
		posCommand = posCommand[4:]
		e.ParseFEN(posCommand)
	} else {
		return errors.New("postion must be followed by 'startpos' or 'fen'")
	}

	if len(commands) > 1 {
		for _, move := range strings.Split(commands[1], " ") {
			m := e.ParseMove(move)
			if m == c.NOMOVE {
				break
			}
			e.MakeMove(m)
			e.Board.Ply = 0
		}
	}

	return nil
}

// modify to run in sockets
func (e *Engine) UCILoop() {
	in := bufio.NewReader(os.Stdin)
	fmt.Println("uciok")

	info := &t.SearchInfo{}

	for {
		line, err := in.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %s\n", err.Error())
			continue
		}
		line = strings.TrimSpace(strings.Trim(line, "\r\n\t"))

		if line == "" {
			continue
		} else if line == "isready" {
			fmt.Println("readyok")
		} else if strings.Contains(line, "position") {
			err := e.ParsePosition(line)
			if err != nil {
				panic(err)
			}
		} else if line == "ucinewgame" {
			err := e.ParsePosition("position startpos")
			if err != nil {
				panic(err)
			}
		} else if strings.Contains(line, "go") {
			e.ParseGo(line, info)
		} else if line == "quit" {
			info.Quit = true
			break
		} else if line == "uci" {
			fmt.Println("id name ViceGo 1.0")
			fmt.Println("id author jsbento")
			fmt.Println("uciok")
		}

		if info.Quit {
			break
		}
	}
}

func (e *Engine) ParseUCICommand(command string, info *t.SearchInfo) (string, error) {
	line := strings.TrimSpace(strings.Trim(command, "\r\n\t"))

	if line == "" {
		return "", nil
	} else if line == "isready" {
		return "readyok", nil
	} else if strings.Contains(line, "position") {
		err := e.ParsePosition(line)
		if err != nil {
			return "", fmt.Errorf("error parsing position: %s", err.Error())
		}
		return "Position received", nil
	} else if line == "ucinewgame" {
		err := e.ParsePosition("position startpos")
		if err != nil {
			return "", fmt.Errorf("error starting new game: %s", err.Error())
		}
		return "New game started", nil
	} else if strings.Contains(line, "go") {
		return e.ParseGo(line, info), nil
	} else if line == "quit" {
		info.Quit = true
		return "Game Stopped", nil
	} else if line == "uci" {
		return "id name ViceGo 1.0\nid author jsbento\nuciok", nil
	}
	return "No command found", nil
}
