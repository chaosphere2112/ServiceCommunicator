package servicecomm

import (
	"bytes"
	"github.com/chaosphere2112/StateMachine/state"
	"strconv"
	"strings"
)

const (
	state_done             = state.FinishedState
	state_start            = 0
	state_lineNumber       = 1
	state_finishLineNumber = 2
	state_startLineLength  = 3
	state_lineLength       = 4
	state_finishLineLength = 5
	state_line             = 6
	state_finishLine       = 7
)

const (
	//48 = '0' 57 = '9'
	ascii_num_start       = 48
	ascii_num_end         = 57
	ascii_newline         = 10
	ascii_carriage_return = 13
	ascii_dollar          = 36
	ascii_star            = 42
)

const (
	line_num_delimiter    = ascii_star
	line_length_delimiter = ascii_dollar
)

type ServiceDecoder struct {
	engine            *state.StateEngine
	states            []func(*ServiceDecoder, byte) (string, int)
	valid             bytes.Buffer
	numberLines       int64
	currentLineLength int64
	LastMessage       []string
}

func (machine *ServiceDecoder) TakeData(data byte, state int) (string, int) {
	return machine.states[state](machine, data)
}

func (machine *ServiceDecoder) addState(state func(*ServiceDecoder, byte) (string, int)) {
	machine.states = append(machine.states, state)
}

func NewDecoder() *ServiceDecoder {

	machine := new(ServiceDecoder)

	machine.engine = new(state.StateEngine)

	machine.engine.Machine = machine

	machine.addState((*ServiceDecoder).startState)
	machine.addState((*ServiceDecoder).readNumbers)
	machine.addState((*ServiceDecoder).readNewline)
	machine.addState((*ServiceDecoder).startLineLength)
	machine.addState((*ServiceDecoder).readNumbers)
	machine.addState((*ServiceDecoder).readNewline)
	machine.addState((*ServiceDecoder).readLine)
	machine.addState((*ServiceDecoder).readNewline)

	return machine

}

func (machine *ServiceDecoder) assignValue(read_data string) error {
	switch machine.engine.CurrentState {

	case state_finishLineNumber:
		val, err := strconv.ParseInt(read_data, 10, 64)
		if err != nil {
			return err
		}
		machine.numberLines = val

	case state_finishLineLength:
		val, err := strconv.ParseInt(read_data, 10, 64)
		if err != nil {
			return err
		}
		machine.currentLineLength = val

	}

	return nil
}

func (machine *ServiceDecoder) startState(data byte) (string, int) {

	if data == line_num_delimiter {
		return "", state_lineNumber
	}

	return "", state_start

}

func (machine *ServiceDecoder) startLineLength(data byte) (string, int) {

	if data == line_length_delimiter {
		return "", state_lineLength
	}

	return "", state_start

}

func (machine *ServiceDecoder) readNumbers(data byte) (string, int) {

	if data >= ascii_num_start && data <= ascii_num_end {
		machine.valid.WriteByte(data)
		return "", machine.engine.CurrentState
	}

	if data == ascii_carriage_return {
		return "", machine.engine.CurrentState + 1
	}

	return "", state_start

}

func (machine *ServiceDecoder) readLine(data byte) (string, int) {

	if int64(machine.valid.Len()) < machine.currentLineLength {
		machine.valid.WriteByte(data)
		return "", machine.engine.CurrentState
	}

	if data == ascii_carriage_return {
		return "", state_finishLine
	}

	return "", state_start

}

func (machine *ServiceDecoder) DecodeMessage(str string) {
	machine.engine.ConsumeBytes([]byte(str))
}

func (machine *ServiceDecoder) Done() {
	machine.LastMessage = machine.engine.Output
}

func (machine *ServiceDecoder) readNewline(data byte) (string, int) {

	stringvalue := machine.valid.String()
	machine.valid.Truncate(0)

	if data == ascii_newline {
		err := machine.assignValue(stringvalue)
		if err != nil {
			return "", state_start
		} else {

			if machine.engine.CurrentState == state_finishLine {
				if int64(len(machine.engine.Output))+1 < machine.numberLines {
					return stringvalue, state_startLineLength
				} else {
					return stringvalue, state_done
				}
			} else {
				return "", machine.engine.CurrentState + 1
			}
		}
	}

	return "", state_start
}

func EncodeMessage(str string) string {
	writer := new(bytes.Buffer)

	words := strings.Split(str, " ")
	delim := "\r\n"
	writer.WriteString("*")
	writer.WriteString(strconv.FormatInt(int64(len(words)), 10))
	writer.WriteString(delim)

	for _, word := range words {
		writer.WriteString("$")
		writer.WriteString(strconv.FormatInt(int64(len(word)), 10))
		writer.WriteString(delim)
		writer.WriteString(word)
		writer.WriteString(delim)
	}

	return writer.String()
}
