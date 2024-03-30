// Code generated by ffjson <https://github.com/pquerna/ffjson>. DO NOT EDIT.
// source: relay.go

package models

import (
	"bytes"
	"fmt"
	fflib "github.com/pquerna/ffjson/fflib/v1"
	"time"
)

// MarshalJSON marshal bytes to json - template
func (j *Relay) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if j == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := j.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONBuf marshal buff to json - template
func (j *Relay) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	if j.Payload != nil {
		buf.WriteString(`{"payload":`)

		{

			err = j.Payload.MarshalJSONBuf(buf)
			if err != nil {
				return err
			}

		}
	} else {
		buf.WriteString(`{"payload":null`)
	}
	if j.Metadata != nil {
		buf.WriteString(`,"meta":`)

		{

			err = j.Metadata.MarshalJSONBuf(buf)
			if err != nil {
				return err
			}

		}
	} else {
		buf.WriteString(`,"meta":null`)
	}
	if j.RelayProof != nil {
		buf.WriteString(`,"proof":`)

		{

			err = j.RelayProof.MarshalJSONBuf(buf)
			if err != nil {
				return err
			}

		}
	} else {
		buf.WriteString(`,"proof":null`)
	}
	buf.WriteByte('}')
	return nil
}

const (
	ffjtRelaybase = iota
	ffjtRelaynosuchkey

	ffjtRelayPayload

	ffjtRelayMetadata

	ffjtRelayRelayProof
)

var ffjKeyRelayPayload = []byte("payload")

var ffjKeyRelayMetadata = []byte("meta")

var ffjKeyRelayRelayProof = []byte("proof")

// UnmarshalJSON umarshall json - template of ffjson
func (j *Relay) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return j.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

// UnmarshalJSONFFLexer fast json unmarshall - template ffjson
func (j *Relay) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error
	currentKey := ffjtRelaybase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffjtRelaynosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'm':

					if bytes.Equal(ffjKeyRelayMetadata, kn) {
						currentKey = ffjtRelayMetadata
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'p':

					if bytes.Equal(ffjKeyRelayPayload, kn) {
						currentKey = ffjtRelayPayload
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyRelayRelayProof, kn) {
						currentKey = ffjtRelayRelayProof
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.SimpleLetterEqualFold(ffjKeyRelayRelayProof, kn) {
					currentKey = ffjtRelayRelayProof
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyRelayMetadata, kn) {
					currentKey = ffjtRelayMetadata
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyRelayPayload, kn) {
					currentKey = ffjtRelayPayload
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffjtRelaynosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffjtRelayPayload:
					goto handle_Payload

				case ffjtRelayMetadata:
					goto handle_Metadata

				case ffjtRelayRelayProof:
					goto handle_RelayProof

				case ffjtRelaynosuchkey:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}

handle_Payload:

	/* handler: j.Payload type=models.Payload kind=struct quoted=false*/

	{
		if tok == fflib.FFTok_null {

			j.Payload = nil

		} else {

			if j.Payload == nil {
				j.Payload = new(Payload)
			}

			err = j.Payload.UnmarshalJSONFFLexer(fs, fflib.FFParse_want_key)
			if err != nil {
				return err
			}
		}
		state = fflib.FFParse_after_value
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Metadata:

	/* handler: j.Metadata type=models.RelayMeta kind=struct quoted=false*/

	{
		if tok == fflib.FFTok_null {

			j.Metadata = nil

		} else {

			if j.Metadata == nil {
				j.Metadata = new(RelayMeta)
			}

			err = j.Metadata.UnmarshalJSONFFLexer(fs, fflib.FFParse_want_key)
			if err != nil {
				return err
			}
		}
		state = fflib.FFParse_after_value
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_RelayProof:

	/* handler: j.RelayProof type=models.RelayProof kind=struct quoted=false*/

	{
		if tok == fflib.FFTok_null {

			j.RelayProof = nil

		} else {

			if j.RelayProof == nil {
				j.RelayProof = new(RelayProof)
			}

			err = j.RelayProof.UnmarshalJSONFFLexer(fs, fflib.FFParse_want_key)
			if err != nil {
				return err
			}
		}
		state = fflib.FFParse_after_value
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:

	return nil
}

// MarshalJSON marshal bytes to json - template
func (j *RelayMeta) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if j == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := j.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONBuf marshal buff to json - template
func (j *RelayMeta) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"block_height":`)
	fflib.FormatBits2(buf, uint64(j.BlockHeight), 10, false)
	buf.WriteByte('}')
	return nil
}

const (
	ffjtRelayMetabase = iota
	ffjtRelayMetanosuchkey

	ffjtRelayMetaBlockHeight
)

var ffjKeyRelayMetaBlockHeight = []byte("block_height")

// UnmarshalJSON umarshall json - template of ffjson
func (j *RelayMeta) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return j.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

// UnmarshalJSONFFLexer fast json unmarshall - template ffjson
func (j *RelayMeta) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error
	currentKey := ffjtRelayMetabase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffjtRelayMetanosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'b':

					if bytes.Equal(ffjKeyRelayMetaBlockHeight, kn) {
						currentKey = ffjtRelayMetaBlockHeight
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.EqualFoldRight(ffjKeyRelayMetaBlockHeight, kn) {
					currentKey = ffjtRelayMetaBlockHeight
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffjtRelayMetanosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffjtRelayMetaBlockHeight:
					goto handle_BlockHeight

				case ffjtRelayMetanosuchkey:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}

handle_BlockHeight:

	/* handler: j.BlockHeight type=uint kind=uint quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for uint", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseUint(fs.Output.Bytes(), 10, 64)

			if err != nil {
				return fs.WrapErr(err)
			}

			j.BlockHeight = uint(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:

	return nil
}

// MarshalJSON marshal bytes to json - template
func (j *RelayProof) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if j == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := j.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONBuf marshal buff to json - template
func (j *RelayProof) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"entropy":`)
	fflib.FormatBits2(buf, uint64(j.Entropy), 10, false)
	buf.WriteString(`,"session_block_height":`)
	fflib.FormatBits2(buf, uint64(j.SessionBlockHeight), 10, false)
	buf.WriteString(`,"servicer_pub_key":`)
	fflib.WriteJsonString(buf, string(j.ServicerPubKey))
	buf.WriteString(`,"blockchain":`)
	fflib.WriteJsonString(buf, string(j.Blockchain))
	if j.AAT != nil {
		buf.WriteString(`,"aat":`)

		{

			err = j.AAT.MarshalJSONBuf(buf)
			if err != nil {
				return err
			}

		}
	} else {
		buf.WriteString(`,"aat":null`)
	}
	buf.WriteString(`,"signature":`)
	fflib.WriteJsonString(buf, string(j.Signature))
	buf.WriteString(`,"request_hash":`)
	fflib.WriteJsonString(buf, string(j.RequestHash))
	buf.WriteByte('}')
	return nil
}

const (
	ffjtRelayProofbase = iota
	ffjtRelayProofnosuchkey

	ffjtRelayProofEntropy

	ffjtRelayProofSessionBlockHeight

	ffjtRelayProofServicerPubKey

	ffjtRelayProofBlockchain

	ffjtRelayProofAAT

	ffjtRelayProofSignature

	ffjtRelayProofRequestHash
)

var ffjKeyRelayProofEntropy = []byte("entropy")

var ffjKeyRelayProofSessionBlockHeight = []byte("session_block_height")

var ffjKeyRelayProofServicerPubKey = []byte("servicer_pub_key")

var ffjKeyRelayProofBlockchain = []byte("blockchain")

var ffjKeyRelayProofAAT = []byte("aat")

var ffjKeyRelayProofSignature = []byte("signature")

var ffjKeyRelayProofRequestHash = []byte("request_hash")

// UnmarshalJSON umarshall json - template of ffjson
func (j *RelayProof) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return j.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

// UnmarshalJSONFFLexer fast json unmarshall - template ffjson
func (j *RelayProof) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error
	currentKey := ffjtRelayProofbase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffjtRelayProofnosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'a':

					if bytes.Equal(ffjKeyRelayProofAAT, kn) {
						currentKey = ffjtRelayProofAAT
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'b':

					if bytes.Equal(ffjKeyRelayProofBlockchain, kn) {
						currentKey = ffjtRelayProofBlockchain
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'e':

					if bytes.Equal(ffjKeyRelayProofEntropy, kn) {
						currentKey = ffjtRelayProofEntropy
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'r':

					if bytes.Equal(ffjKeyRelayProofRequestHash, kn) {
						currentKey = ffjtRelayProofRequestHash
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 's':

					if bytes.Equal(ffjKeyRelayProofSessionBlockHeight, kn) {
						currentKey = ffjtRelayProofSessionBlockHeight
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyRelayProofServicerPubKey, kn) {
						currentKey = ffjtRelayProofServicerPubKey
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyRelayProofSignature, kn) {
						currentKey = ffjtRelayProofSignature
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.EqualFoldRight(ffjKeyRelayProofRequestHash, kn) {
					currentKey = ffjtRelayProofRequestHash
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyRelayProofSignature, kn) {
					currentKey = ffjtRelayProofSignature
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyRelayProofAAT, kn) {
					currentKey = ffjtRelayProofAAT
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyRelayProofBlockchain, kn) {
					currentKey = ffjtRelayProofBlockchain
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyRelayProofServicerPubKey, kn) {
					currentKey = ffjtRelayProofServicerPubKey
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyRelayProofSessionBlockHeight, kn) {
					currentKey = ffjtRelayProofSessionBlockHeight
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyRelayProofEntropy, kn) {
					currentKey = ffjtRelayProofEntropy
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffjtRelayProofnosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffjtRelayProofEntropy:
					goto handle_Entropy

				case ffjtRelayProofSessionBlockHeight:
					goto handle_SessionBlockHeight

				case ffjtRelayProofServicerPubKey:
					goto handle_ServicerPubKey

				case ffjtRelayProofBlockchain:
					goto handle_Blockchain

				case ffjtRelayProofAAT:
					goto handle_AAT

				case ffjtRelayProofSignature:
					goto handle_Signature

				case ffjtRelayProofRequestHash:
					goto handle_RequestHash

				case ffjtRelayProofnosuchkey:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}

handle_Entropy:

	/* handler: j.Entropy type=uint64 kind=uint64 quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for uint64", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseUint(fs.Output.Bytes(), 10, 64)

			if err != nil {
				return fs.WrapErr(err)
			}

			j.Entropy = uint64(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_SessionBlockHeight:

	/* handler: j.SessionBlockHeight type=uint kind=uint quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for uint", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseUint(fs.Output.Bytes(), 10, 64)

			if err != nil {
				return fs.WrapErr(err)
			}

			j.SessionBlockHeight = uint(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_ServicerPubKey:

	/* handler: j.ServicerPubKey type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.ServicerPubKey = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Blockchain:

	/* handler: j.Blockchain type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.Blockchain = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_AAT:

	/* handler: j.AAT type=models.AAT kind=struct quoted=false*/

	{
		if tok == fflib.FFTok_null {

			j.AAT = nil

		} else {

			if j.AAT == nil {
				j.AAT = new(AAT)
			}

			err = j.AAT.UnmarshalJSONFFLexer(fs, fflib.FFParse_want_key)
			if err != nil {
				return err
			}
		}
		state = fflib.FFParse_after_value
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Signature:

	/* handler: j.Signature type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.Signature = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_RequestHash:

	/* handler: j.RequestHash type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.RequestHash = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:

	return nil
}

// MarshalJSON marshal bytes to json - template
func (j *SendRelayRequest) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if j == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := j.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONBuf marshal buff to json - template
func (j *SendRelayRequest) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	if j.Payload != nil {
		buf.WriteString(`{"Payload":`)

		{

			err = j.Payload.MarshalJSONBuf(buf)
			if err != nil {
				return err
			}

		}
	} else {
		buf.WriteString(`{"Payload":null`)
	}
	if j.Signer != nil {
		buf.WriteString(`,"Signer":`)

		{

			err = j.Signer.MarshalJSONBuf(buf)
			if err != nil {
				return err
			}

		}
	} else {
		buf.WriteString(`,"Signer":null`)
	}
	buf.WriteString(`,"Chain":`)
	fflib.WriteJsonString(buf, string(j.Chain))
	buf.WriteString(`,"SelectedNodePubKey":`)
	fflib.WriteJsonString(buf, string(j.SelectedNodePubKey))
	if j.Session != nil {
		buf.WriteString(`,"Session":`)

		{

			err = j.Session.MarshalJSONBuf(buf)
			if err != nil {
				return err
			}

		}
	} else {
		buf.WriteString(`,"Session":null`)
	}
	if j.Timeout != nil {
		buf.WriteString(`,"Timeout":`)
		fflib.FormatBits2(buf, uint64(*j.Timeout), 10, *j.Timeout < 0)
	} else {
		buf.WriteString(`,"Timeout":null`)
	}
	buf.WriteByte('}')
	return nil
}

const (
	ffjtSendRelayRequestbase = iota
	ffjtSendRelayRequestnosuchkey

	ffjtSendRelayRequestPayload

	ffjtSendRelayRequestSigner

	ffjtSendRelayRequestChain

	ffjtSendRelayRequestSelectedNodePubKey

	ffjtSendRelayRequestSession

	ffjtSendRelayRequestTimeout
)

var ffjKeySendRelayRequestPayload = []byte("Payload")

var ffjKeySendRelayRequestSigner = []byte("Signer")

var ffjKeySendRelayRequestChain = []byte("Chain")

var ffjKeySendRelayRequestSelectedNodePubKey = []byte("SelectedNodePubKey")

var ffjKeySendRelayRequestSession = []byte("Session")

var ffjKeySendRelayRequestTimeout = []byte("Timeout")

// UnmarshalJSON umarshall json - template of ffjson
func (j *SendRelayRequest) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return j.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

// UnmarshalJSONFFLexer fast json unmarshall - template ffjson
func (j *SendRelayRequest) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error
	currentKey := ffjtSendRelayRequestbase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffjtSendRelayRequestnosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'C':

					if bytes.Equal(ffjKeySendRelayRequestChain, kn) {
						currentKey = ffjtSendRelayRequestChain
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'P':

					if bytes.Equal(ffjKeySendRelayRequestPayload, kn) {
						currentKey = ffjtSendRelayRequestPayload
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'S':

					if bytes.Equal(ffjKeySendRelayRequestSigner, kn) {
						currentKey = ffjtSendRelayRequestSigner
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeySendRelayRequestSelectedNodePubKey, kn) {
						currentKey = ffjtSendRelayRequestSelectedNodePubKey
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeySendRelayRequestSession, kn) {
						currentKey = ffjtSendRelayRequestSession
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'T':

					if bytes.Equal(ffjKeySendRelayRequestTimeout, kn) {
						currentKey = ffjtSendRelayRequestTimeout
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.SimpleLetterEqualFold(ffjKeySendRelayRequestTimeout, kn) {
					currentKey = ffjtSendRelayRequestTimeout
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeySendRelayRequestSession, kn) {
					currentKey = ffjtSendRelayRequestSession
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeySendRelayRequestSelectedNodePubKey, kn) {
					currentKey = ffjtSendRelayRequestSelectedNodePubKey
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeySendRelayRequestChain, kn) {
					currentKey = ffjtSendRelayRequestChain
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeySendRelayRequestSigner, kn) {
					currentKey = ffjtSendRelayRequestSigner
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeySendRelayRequestPayload, kn) {
					currentKey = ffjtSendRelayRequestPayload
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffjtSendRelayRequestnosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffjtSendRelayRequestPayload:
					goto handle_Payload

				case ffjtSendRelayRequestSigner:
					goto handle_Signer

				case ffjtSendRelayRequestChain:
					goto handle_Chain

				case ffjtSendRelayRequestSelectedNodePubKey:
					goto handle_SelectedNodePubKey

				case ffjtSendRelayRequestSession:
					goto handle_Session

				case ffjtSendRelayRequestTimeout:
					goto handle_Timeout

				case ffjtSendRelayRequestnosuchkey:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}

handle_Payload:

	/* handler: j.Payload type=models.Payload kind=struct quoted=false*/

	{
		if tok == fflib.FFTok_null {

			j.Payload = nil

		} else {

			if j.Payload == nil {
				j.Payload = new(Payload)
			}

			err = j.Payload.UnmarshalJSONFFLexer(fs, fflib.FFParse_want_key)
			if err != nil {
				return err
			}
		}
		state = fflib.FFParse_after_value
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Signer:

	/* handler: j.Signer type=models.Ed25519Account kind=struct quoted=false*/

	{
		if tok == fflib.FFTok_null {

			j.Signer = nil

		} else {

			if j.Signer == nil {
				j.Signer = new(Ed25519Account)
			}

			err = j.Signer.UnmarshalJSONFFLexer(fs, fflib.FFParse_want_key)
			if err != nil {
				return err
			}
		}
		state = fflib.FFParse_after_value
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Chain:

	/* handler: j.Chain type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.Chain = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_SelectedNodePubKey:

	/* handler: j.SelectedNodePubKey type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.SelectedNodePubKey = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Session:

	/* handler: j.Session type=models.Session kind=struct quoted=false*/

	{
		if tok == fflib.FFTok_null {

			j.Session = nil

		} else {

			if j.Session == nil {
				j.Session = new(Session)
			}

			err = j.Session.UnmarshalJSONFFLexer(fs, fflib.FFParse_want_key)
			if err != nil {
				return err
			}
		}
		state = fflib.FFParse_after_value
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Timeout:

	/* handler: j.Timeout type=time.Duration kind=int64 quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for Duration", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

			j.Timeout = nil

		} else {

			tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 64)

			if err != nil {
				return fs.WrapErr(err)
			}

			ttypval := time.Duration(tval)
			j.Timeout = &ttypval

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:

	return nil
}

// MarshalJSON marshal bytes to json - template
func (j *SendRelayResponse) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if j == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := j.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONBuf marshal buff to json - template
func (j *SendRelayResponse) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"response":`)
	fflib.WriteJsonString(buf, string(j.Response))
	buf.WriteByte('}')
	return nil
}

const (
	ffjtSendRelayResponsebase = iota
	ffjtSendRelayResponsenosuchkey

	ffjtSendRelayResponseResponse
)

var ffjKeySendRelayResponseResponse = []byte("response")

// UnmarshalJSON umarshall json - template of ffjson
func (j *SendRelayResponse) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return j.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

// UnmarshalJSONFFLexer fast json unmarshall - template ffjson
func (j *SendRelayResponse) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error
	currentKey := ffjtSendRelayResponsebase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffjtSendRelayResponsenosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'r':

					if bytes.Equal(ffjKeySendRelayResponseResponse, kn) {
						currentKey = ffjtSendRelayResponseResponse
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.EqualFoldRight(ffjKeySendRelayResponseResponse, kn) {
					currentKey = ffjtSendRelayResponseResponse
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffjtSendRelayResponsenosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffjtSendRelayResponseResponse:
					goto handle_Response

				case ffjtSendRelayResponsenosuchkey:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}

handle_Response:

	/* handler: j.Response type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.Response = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:

	return nil
}
