package run

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"slices"
	"strings"

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

// ErrMoreThanOneLive is returned when more than one target in the values
// checkout is live. Only one target is.
var ErrMoreThanOneLive = errors.New("more than one target is live")

// TargetKind reads the named target's kind from the values checkout's target
// files, each keyed by its target's name: committed holds them as the
// checkout's commit has them, and working as the checkout stands, uncommitted
// changes included.
//
// A target is live when either version of its file says so, so an uncommitted
// edit can never turn the live target into a drill, and lift the rules that
// keep a drill off it. Every file is read strictly, in both versions, and
// TargetKind stops with ErrMoreThanOneLive when more than one target is live.
// It stops with ErrUnknownTarget when the checkout as it stands has no file for
// the target.
func TargetKind(name string, committed, working map[string][]byte) (Kind, error) {
	path, err := TargetFilePath(name)
	if err != nil {
		return 0, err
	}
	names := make([]string, 0, len(working)+len(committed))
	for n := range working {
		names = append(names, n)
	}
	for n := range committed {
		if _, ok := working[n]; !ok {
			names = append(names, n)
		}
	}
	slices.Sort(names)

	kinds := make(map[string]Kind, len(names))
	var live []string
	for _, n := range names {
		kind := Drill
		for _, v := range []struct {
			files map[string][]byte
			as    string
		}{{committed, "as committed"}, {working, "as it stands"}} {
			data, ok := v.files[n]
			if !ok {
				continue
			}
			f, err := ParseTargetFile(data)
			if err != nil {
				return 0, fmt.Errorf("target %s: %s/%s.yaml, %s: %w", n, TargetsDir, n, v.as, err)
			}
			if f.Kind == Live {
				kind = Live
			}
		}
		kinds[n] = kind
		if kind == Live {
			live = append(live, n)
		}
	}
	if len(live) > 1 {
		return 0, fmt.Errorf("%w: %s each say kind: live, in the values checkout as committed or as it stands",
			ErrMoreThanOneLive, strings.Join(live, ", "))
	}
	if _, ok := working[name]; !ok {
		return 0, fmt.Errorf("%w: %s: the values checkout has no %s", ErrUnknownTarget, name, path)
	}
	return kinds[name], nil
}
