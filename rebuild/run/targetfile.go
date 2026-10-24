package run

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"

	"gopkg.in/yaml.v3"
)

// TargetsDir is the directory of the values checkout that holds one file per
// target, named after it: targets/<name>.yaml.
const TargetsDir = "targets"

// TargetFile is what a target's file in the values checkout holds. Today that
// is only whether the target is live or a drill:
//
//	kind: drill
//
// Proposal: definitions, values and secrets adds the rest of the file.
type TargetFile struct {
	// Kind says whether the target is live or a drill. The file must say
	// which.
	Kind Kind `yaml:"kind"`
}

// ErrInvalidTargetFile is returned for a target's file that can't be read
// strictly: one that is empty, holds a key it shouldn't, more than one
// document, or a kind that is neither live nor drill, or says no kind at all.
var ErrInvalidTargetFile = errors.New("invalid target file")

// ErrUnknownTarget is returned for a target the values checkout has no file
// for, and for a name that can't be one.
var ErrUnknownTarget = errors.New("no such target")

// targetName is what a target's name may be: one path element of letters,
// digits, dots, dashes and underscores, which doesn't start with a dot, so it
// names a file in TargetsDir and nowhere else.
var targetName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// TargetFilePath returns the path of the named target's file in the values
// checkout, targets/<name>.yaml. It refuses an empty name with ErrNoTarget,
// since there is no default target, and a name that would reach outside
// TargetsDir with ErrUnknownTarget.
func TargetFilePath(name string) (string, error) {
	if name == "" {
		return "", ErrNoTarget
	}
	if !targetName.MatchString(name) {
		return "", fmt.Errorf("%w: %q is not a target name: use letters, digits, dots, dashes and underscores, not starting with a dot",
			ErrUnknownTarget, name)
	}
	return path.Join(TargetsDir, name+".yaml"), nil
}

// ParseTargetFile reads a target's file strictly. It returns
// ErrInvalidTargetFile when the file is empty, holds a key TargetFile doesn't
// have, holds more than one YAML document, or doesn't say kind: live or
// kind: drill.
func ParseTargetFile(data []byte) (TargetFile, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)

	var f TargetFile
	err := dec.Decode(&f)
	switch {
	case errors.Is(err, io.EOF):
		return TargetFile{}, fmt.Errorf("%w: it is empty, and must say kind: live or kind: drill", ErrInvalidTargetFile)
	case err != nil:
		return TargetFile{}, fmt.Errorf("%w: %w", ErrInvalidTargetFile, err)
	case f.Kind == 0:
		return TargetFile{}, fmt.Errorf("%w: it must say kind: live or kind: drill", ErrInvalidTargetFile)
	}

	var extra yaml.Node
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return TargetFile{}, fmt.Errorf("%w: it must hold one YAML document", ErrInvalidTargetFile)
	}
	return f, nil
}

// UnmarshalText reads "live" or "drill", and refuses anything else.
func (k *Kind) UnmarshalText(text []byte) error {
	switch string(text) {
	case "live":
		*k = Live
	case "drill":
		*k = Drill
	default:
		return fmt.Errorf("kind %q is neither live nor drill", text)
	}
	return nil
}
