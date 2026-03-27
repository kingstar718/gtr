package main

type InputType int

const (
	TypeUnknown InputType = iota
	TypeCoordinate
	TypeHTTP
	TypeTimestamp
	TypeText
)

func detectInputType(input string) InputType {
	if isHTTPURL(input) {
		return TypeHTTP
	}

	if isTimeFormat(input) {
		return TypeTimestamp
	}

	if isCoordinateFormat(input) {
		return TypeCoordinate
	}

	return TypeText
}

func handleAutoConvert(input string) error {
	inputType := detectInputType(input)

	switch inputType {
	case TypeCoordinate:
		return handleCoordinateConvert(input)
	case TypeHTTP:
		return handleHTTPRequest(input)
	case TypeTimestamp:
		return handleTimeConvert(input)
	case TypeText:
		return handleAllConversions(input)
	default:
		return handleAllConversions(input)
	}
}
