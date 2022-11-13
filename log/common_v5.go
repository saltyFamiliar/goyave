package log

import (
	"fmt"
	"net"
	"strconv"
)

const (
	// Format is the fmt format code for common logs
	Format string = "%s %s %s [%s] \"%s %s %s\" %d %d"

	// TimestampFormat is the time format used in common logs
	TimestampFormat string = "02/Jan/2006:15:04:05 -0700"
)

// CommonLogFormatter build a log entry using the Common Log Format.
func CommonLogFormatterV5(ctx *LogContext) string {
	req := ctx.Request.Request()
	url := ctx.Request.URL()

	username := "-"
	if url.User != nil {
		if name := url.User.Username(); name != "" {
			username = name
		}
	}

	host, _, err := net.SplitHostPort(req.RemoteAddr)

	if err != nil {
		host = req.RemoteAddr
	}

	uri := req.RequestURI

	// Requests using the CONNECT method over HTTP/2.0 must use
	// the authority field (aka r.Host) to identify the target.
	// Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
	if req.ProtoMajor == 2 && req.Method == "CONNECT" {
		uri = req.Host
	}
	if uri == "" {
		uri = url.RequestURI()
	}

	return fmt.Sprintf(Format,
		host,
		"-",
		username,
		ctx.Now.Format(TimestampFormat),
		req.Method,
		strconv.QuoteToASCII(uri),
		req.Proto,
		ctx.Response.GetStatus(),
		ctx.Length,
	)
}

// CombinedLogFormatter build a log entry using the Combined Log Format.
func CombinedLogFormatterV5(ctx *LogContext) string {
	return fmt.Sprintf("%s \"%s\" \"%s\"", CommonLogFormatterV5(ctx), ctx.Request.Referrer(), ctx.Request.UserAgent())
}
