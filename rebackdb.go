package rebackdb

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type dateFormat string

const (
	// FormatShort is a date in YYYYMMDDHHSS format.
	FormatShort dateFormat = "200601021504"
	// FormatISO is a date in YYYY-MM-DD-HH-SS format.
	FormatISO dateFormat = "2006-01-02-15-04"
)

type _os string

const (
	// Unix is for MAC and Linux
	Unix _os = "mv"
	// Windows is for Microsoft Windows
	Windows _os = "move"
)

// DumpOptions is options for rethink-dump command app.
type DumpOptions struct {
	Connection        string
	OutputFileName    string
	DatabasesToExport []string
	TablesToExport    []string
	PasswordFile      string
	TLSCert           string
	Clients           int
	TempDir           string
	DateFormat        dateFormat
	OperativeSystem   _os
}

// ResultFile is the resulting tar.gz of the dumping command
type ResultFile struct {
	Path            string
	MIME            string
	OSMoveCommand   string
	CommandExecuted string
	CommandOutput   string
}

// Backup fires the dump command with respective options and creates a tar.gz file from
// rethinkDB cluster
func Backup(options DumpOptions) (*ResultFile, error) {
	binary, err := exec.LookPath("rethinkdb")
	if err != nil {
		return nil, err
	}

	cmdOptions, err := options.Validate()
	if err != nil {
		return nil, err
	}

	result := &ResultFile{MIME: "application/x-tar"}
	result.Path = fmt.Sprintf(`%s_backup_%s.tar.gz`, time.Now().Format(string(options.DateFormat)), options.OutputFileName)
	result.OSMoveCommand = string(options.OperativeSystem)

	cmd := exec.Command(binary, cmdOptions...)
	env := os.Environ()
	cmd.Env = env

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err = cmd.Start()
	if err != nil {
		result.CommandOutput = errbuf.String()
		return result, err
	}

	err = cmd.Wait()
	if err != nil {
		result.CommandOutput = errbuf.String()
		return result, err
	}

	result.CommandExecuted = fmt.Sprintf("%s %s", binary, strings.Join(cmdOptions, " "))
	result.CommandOutput = outbuf.String()

	return result, err
}

// Move takes resulting file from backup and moves it to an specific destination path
func (f *ResultFile) Move(folder string) (*ResultFile, error) {
	binary, err := exec.LookPath(string(f.OSMoveCommand))
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(binary, f.Path, folder+f.FileName())
	env := os.Environ()
	cmd.Env = env

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	f.CommandExecuted = fmt.Sprintf("%s %s %s", binary, f.Path, folder+f.FileName())

	err = cmd.Start()
	if err != nil {
		f.CommandOutput = errbuf.String()
		return f, err
	}

	err = cmd.Wait()
	if err != nil {
		f.CommandOutput = errbuf.String()
		return f, err
	}

	f.CommandOutput = outbuf.String()

	return f, nil
}

// Validate verifies if the dump options struct has it's parameters correctly specified
func (options DumpOptions) Validate() ([]string, error) {
	var err error

	if options.Connection == "" {
		err = errors.New("RethinkDB server was not specified")
		return nil, err
	}

	if options.OutputFileName == "" {
		err = errors.New("Output file name was not specified")
		return nil, err
	}

	if options.DateFormat == "" {
		err = errors.New("A date format must be specified")
		return nil, err
	}

	if options.OperativeSystem == "" {
		err = errors.New("An operative system must be specified")
		return nil, err
	}

	var cmdOptions []string
	cmdOptions = append(cmdOptions, "dump")
	cmdOptions = append(cmdOptions, fmt.Sprintf(`-c %s`, options.Connection))
	cmdOptions = append(cmdOptions, fmt.Sprintf(`-f %s_backup_%s.tar.gz`, time.Now().Format(string(options.DateFormat)), options.OutputFileName))

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
