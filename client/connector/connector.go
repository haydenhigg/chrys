package connector

import "net/url"

type Payload struct {
	Query url.Values
	Body  url.Values
}
