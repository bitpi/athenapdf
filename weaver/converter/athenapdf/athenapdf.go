package athenapdf

import (
	"github.com/arachnys/athenapdf/weaver/converter"
	"github.com/arachnys/athenapdf/weaver/gcmd"
	"github.com/arachnys/athenapdf/weaver/util"
	"log"
	"os"
	"strings"
)

// AthenaPDF represents a conversion job for athenapdf CLI.
// AthenaPDF implements the Converter interface with a custom Convert method.
type AthenaPDF struct {
	// AthenaPDF inherits properties from UploadConversion, and as such,
	// it supports uploading of its results to S3
	// (if the necessary credentials are given).
	// See UploadConversion for more information.
	converter.UploadConversion
	// CMD is the base athenapdf CLI command that will be executed.
	// e.g. 'athenapdf -S -T 120'
	CMD string
	// Aggressive will alter the athenapdf CLI conversion behaviour by passing
	// an '-A' command-line flag to indicate aggressive content extraction
	// (ideal for a clutter-free reading experience).
	Aggressive bool
}

// constructCMD returns a string array containing the AthenaPDF command to be
// executed by Go's os/exec Output. It does this using a base command, and path
// string.
// It will set an additional '-A' flag if aggressive is set to true.
// See athenapdf CLI for more information regarding the aggressive mode.
func constructCMD(base string, path string, aggressive bool) []string {
	args := strings.Fields(base)
	args = append(args, path)
	if aggressive {
		args = append(args, "-A")
	}
	return args
}

// Convert returns a byte slice containing a PDF converted from HTML
// using athenapdf CLI.
// See the Convert method for Conversion for more information.
func (c AthenaPDF) Convert(done <-chan struct{}) ([]byte, error) {
	// TODO: Check content type
	p, err := util.HandleOctetStream(c.Path)
	// TODO: Distinguish between "REAL" errors, and a bad URL
	if err != nil {
		return nil, err
	}
	if p != "" {
		// GC
		defer os.Remove(p)
		c.Path = p
	}

	log.Printf("[AthenaPDF] converting to PDF: %s\n", c.Path)

	// Construct the command to execute
	cmd := constructCMD(c.CMD, c.Path, c.Aggressive)
	out, err := gcmd.Execute(cmd, done)
	if err != nil {
		return nil, err
	}

	return out, nil
}
