package rebackdb

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"
)

type dateFormat string

const (
	// FormatShort is a date in YYYYMMDDHHSS format.
	FormatShort dateFormat = "200601021504"
	// FormatISO is a date in YYYY-MM-DD-HH-SS format.
	FormatISO dateFormat = "2006-01-02-15-04"
)

// DumpOptions is options for rethink-dump command app.
type DumpOptions struct {
	Connection        string
	OutputFile        string
	DatabasesToExport []string
	TablesToExport    []string
	Password          string
	PasswordFile      string
	TLSCert           string
	Clients           int
	TempDir           string
}

// ResultFile is the resulting tar.gz of the dumping command
type ResultFile struct {
	Path string
	MIME string
}

// Backup fires the dump command with respective options and creates a tar.gz file from
// rethinkDB cluster
func Backup(options DumpOptions, dateFormat dateFormat) (*ResultFile, error) {
	cmdOptions, err := options.Validate()
	result := &ResultFile{MIME: "application/x-tar"}
	result.Path = fmt.Sprintf(`%s_backup_%s.tar.gz`, time.Now().Format(string(dateFormat)), options.OutputFile)

	out, err := exec.Command("rethinkdb", cmdOptions...).Output()
	if err != nil {
		return nil, err
	}

	log.Printf("Backup result: %s", out)

	return result, err
}

// Move takes resulting file from backup and moves it to an specific destination path
func (f *ResultFile) Move(folder string) error {
	out, err := exec.Command("mv", f.Path, folder+f.FileName()).Output()
	if err != nil {
		return err
	}

	log.Printf("To result: %s", out)
	return nil
}

// Validate verifies if the dump options struct has it's parameters correctly specified
func (options DumpOptions) Validate() ([]string, error) {
	var err error

	if options.Connection == "" {
		err = errors.New("RethinkDB server was not specified")
		return nil, err
	}

	if options.OutputFile == "" {
		err = errors.New("Output file name was not specified")
		return nil, err
	}

	var cmdOptions []string
	cmdOptions = append(cmdOptions, fmt.Sprintf(`-c %s`, options.Connection))
	cmdOptions = append(cmdOptions, fmt.Sprintf(`-f %s`, options.OutputFile))

	if options.TablesToExport != nil && len(options.TablesToExport) > 0 {
		if options.DatabasesToExport != nil && len(options.DatabasesToExport) > 0 {
			for _, table := range options.TablesToExport {
				for _, database := range options.DatabasesToExport {
					cmdOptions = append(cmdOptions, fmt.Sprintf(`-e %s.%s`, database, table))
				}
			}
		} else {
			for _, table := range options.TablesToExport {
				cmdOptions = append(cmdOptions, fmt.Sprintf(`-e %s`, table))
			}
		}
	} else {
		if options.DatabasesToExport != nil {
			if len(options.DatabasesToExport) > 0 {
				for _, database := range options.DatabasesToExport {
					cmdOptions = append(cmdOptions, fmt.Sprintf(`-e %s`, database))
				}
			}
		}
	}

	if options.Password != "" && options.PasswordFile != "" {
		err = errors.New("A password or a password file could be specified, but not both")
		return nil, err
	}

	if options.Password != "" {
		cmdOptions = append(cmdOptions, fmt.Sprintf(`-p %s`, options.Password))
	}

	if options.PasswordFile != "" {
		cmdOptions = append(cmdOptions, fmt.Sprintf(`--password-file %s`, options.PasswordFile))
	}

	if options.TLSCert != "" {
		cmdOptions = append(cmdOptions, fmt.Sprintf(`--tls-cert %s`, options.TLSCert))
	}

	if options.Clients != 0 {
		cmdOptions = append(cmdOptions, fmt.Sprintf(`--clients %d`, options.Clients))
	}

	if options.TempDir != "" {
		cmdOptions = append(cmdOptions, fmt.Sprintf(`--temp-dir %s`, options.TempDir))
	}

	return cmdOptions, err
}

// FileName returns the just filename component of the `Path` attribute
func (f ResultFile) FileName() string {
	_, filename := filepath.Split(f.Path)
	return filename
}
