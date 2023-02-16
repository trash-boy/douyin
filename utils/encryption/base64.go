package encryption

import (
	b64 "encoding/base64"
)

func Encoding(data string)(string){
	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	return sEnc
}
