package basicauth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// SetWWWAuthenticate set an auth error to the response
func SetWWWAuthenticate(res http.ResponseWriter, typ, realm string) {
	realmStr := ""
	if realm != "" {
		realmStr = fmt.Sprintf(`realm="%s"`, realm)
	}

	authStr := fmt.Sprintf("%s %s", typ, realmStr)
	res.Header().Set("WWW-Authenticate", authStr)
}

// ParseAuthHeader returns a username and password from a list of Auth header
// Supported type are "basic" and "password". Example
// 		basic user:password
//		password base64EncodedPassword
// TODO(melvin): support multiple realm
func ParseAuthHeader(auths []string, typeWanted, realmWanted string) (string, string, error) {
	for _, auth := range auths {
		data := strings.Split(auth, " ")
		nbArg := len(data)

		// the string can contain the type, the realm, and the encoded data.
		// The type is always first param, the realm is optional
		// Ex: basic real="secretArea" base64EncodedString
		if nbArg > 1 {
			// We check the type is right
			providedType := strings.ToLower(data[0])
			if providedType != typeWanted {
				continue
			}

			//  If a realm is wanted, we need at least 3 params
			// (type, realm, and the encoded string)
			if nbArg != 3 && realmWanted != "" {
				continue
			}

			//  if a realm is not wanted we should not have more that 2 params
			// (the type and the encoded string)
			if realmWanted == "" && nbArg != 2 {
				continue
			}
			realmWantedStr := fmt.Sprintf(`realm="%s"`, realmWanted)
			encoded := data[1] // let's assume we don't want a realm

			// We want a realm
			if realmWanted != "" {
				switch realmWantedStr {
				case data[1]:
					encoded = data[2]
				case data[2]:
					encoded = data[1]
				default:
					// wanted realm not provided
					continue
				}
			}

			decodedString, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				return "", "", err
			}
			decoded := string(decodedString[:])

			switch strings.ToLower(data[0]) {
			case "basic":
				cred := strings.Split(decoded, ":")
				return cred[0], cred[1], nil
			case "password":
				return "", decoded, nil
			}
		}
	}

	return "", "", nil
}
