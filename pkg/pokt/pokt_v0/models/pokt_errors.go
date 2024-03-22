//go:generate ffjson $GOFILE
package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const modulePocketCore = "pocketcore"
const moduleRoot = "sdk"

var (
	ErrPocketCoreInvalidBlockHeight = PocketSdkError{Codespace: modulePocketCore, Code: 60}
	ErrPocketCoreOverService        = PocketSdkError{Codespace: modulePocketCore, Code: 71}
	ErrPocketEvidenceSealed         = PocketSdkError{Codespace: modulePocketCore, Code: 90}
	ErrorServicerNotFound           = PocketSdkError{Codespace: moduleRoot, Message: "Failed to find correct servicer PK"}
)

var (
	sdkErrorRegexCodespace = regexp.MustCompile(`codespace: (\w+)`)
	sdkErrorRegexCode      = regexp.MustCompile(`code: (\d+)`)
	sdkErrorRegexMessage   = regexp.MustCompile(`message: \\"(.+?)\\"`)
)

type PocketRPCError struct {
	HttpCode int    `json:"code"`
	Message  string `json:"message"`
}

func (r PocketRPCError) Error() string {
	return fmt.Sprintf(`ERROR: HttpCode: %d Message: %s`, r.HttpCode, r.Message)
}

type PocketSdkError struct {
	Codespace string
	Code      uint64
	Message   string // Fallback if code does not exist
}

func (r PocketSdkError) Error() string {
	return fmt.Sprintf(`ERROR: Codespace: %s Code: %d`, r.Codespace, r.Code)
}

func (r PocketRPCError) ToSdkError() *PocketSdkError {
	msgLower := strings.ToLower(r.Message)
	codespaceMatches := sdkErrorRegexCodespace.FindStringSubmatch(msgLower)
	if len(codespaceMatches) != 2 {
		return nil
	}

	sdkError := PocketSdkError{
		Codespace: codespaceMatches[1],
	}

	codeMatches := sdkErrorRegexCode.FindStringSubmatch(msgLower)
	if len(codeMatches) == 2 {
		code, err := strconv.ParseUint(codeMatches[1], 10, 0)
		if err == nil {
			sdkError.Code = code
		}
	}

	// If the code parsed is zero, then pocket core did not return a code
	// Internal errors do not have a code attached to it, so we should parse message
	if sdkError.Code == 0 {
		matchesMessage := sdkErrorRegexMessage.FindStringSubmatch(msgLower)
		if len(matchesMessage) == 2 {
			sdkError.Message = matchesMessage[1]
		}
	}

	return &sdkError
}
