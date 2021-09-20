package options

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "dir",
		Aliases:  []string{"d"},
		Usage:    "path to directory containing directory path, defaults to current directory",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "out",
		Aliases:  []string{"o"},
		Usage:    "output directory. will be created if does not exist",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "include",
		Aliases:  []string{"i"},
		Value:    "",
		Usage:    "patterns of file paths to include, comma delimited, may contain any glob pattern",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "exclude",
		Aliases:  []string{"e"},
		Value:    "",
		Usage:    "patterns of file paths to exclude, comma delimited, may contain any glob pattern",
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "verbose",
		Aliases:  []string{"vv"},
		Value:    false,
		Usage:    "verbose logging",
		Required: false,
	},
	&cli.IntFlag{
		Name:     "max-size",
		Value:    6,
		Usage:    "maximal file size, in MB",
		Required: false,
	},
}

type Options struct {
	CodePath         string
	IncludePatterns  []string
	ExcludePatterns  []string
	VerboseLogging   bool
	MaxFileSizeBytes int64
}

func splitListFlag(flag string) []string {
	if len(flag) == 0 {
		return []string{}
	}
	return strings.Split(flag, ",")
}

func validateDirectory(dirPath string, createIfNotExist bool) error {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		if !createIfNotExist {
			return fmt.Errorf("directory does not exist at %v", dirPath)
		}
		err = os.MkdirAll(dirPath, 0777)
		if err != nil {
			return fmt.Errorf("failed to create directory at %v: %w", dirPath, err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("directory error at %v: %w", dirPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("directory is actually a file at %v", dirPath)
	}
	return nil
}

func ParseOptions(c *cli.Context) (*Options, error) {
	opts := &Options{
		CodePath:         c.String("dir"),
		IncludePatterns:  splitListFlag(c.String("include")),
		ExcludePatterns:  splitListFlag(c.String("exclude")),
		VerboseLogging:   c.Bool("verbose"),
		MaxFileSizeBytes: int64(c.Int("max-size")) * 1024 * 1024,
	}
	var err error
	if len(opts.CodePath) == 0 {
		opts.CodePath, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %v", err)
		}
	} else {
		err = validateDirectory(opts.CodePath, false)
		if err != nil {
			return nil, fmt.Errorf("directory path '%v' is not valid: %v", opts.CodePath, err)
		}
	}

	return opts, nil
}
