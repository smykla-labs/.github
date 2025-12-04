package main

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/smykla-labs/.github/pkg/logger"
)

var version = "dev"

// CLI defines the command-line interface structure.
type CLI struct {
	LogLevel  string     `help:"Log level (trace|debug|info|warn|error)" default:"info" enum:"trace,debug,info,warn,error"`
	UseGHAuth bool       `help:"Use 'gh auth token' for authentication"`
	DryRun    bool       `help:"Preview changes without applying them"`
	Org       string     `help:"GitHub organization" default:"smykla-labs"`
	Version   VersionCmd `cmd:"" help:"Show version information"`
	Labels    LabelsCmd  `cmd:"" help:"Label synchronization commands"`
	Files     FilesCmd   `cmd:"" help:"File synchronization commands"`
	Smyklot   SmyklotCmd `cmd:"" help:"Smyklot version synchronization commands"`
	Repos     ReposCmd   `cmd:"" help:"Repository listing commands"`
}

// VersionCmd shows version information.
type VersionCmd struct{}

// Run executes the version command.
func (*VersionCmd) Run(_ context.Context) error {
	fmt.Printf("orgsync version %s\n", version)

	return nil
}

// LabelsCmd contains label sync subcommands.
type LabelsCmd struct {
	Sync LabelsSyncCmd `cmd:"" help:"Sync labels to a repository"`
}

// LabelsSyncCmd syncs labels to a repository.
type LabelsSyncCmd struct {
	Repo       string `help:"Target repository (e.g., 'myrepo')" required:""`
	LabelsFile string `help:"Path to labels YAML file" required:""`
	Config     string `help:"JSON sync config (optional)"`
}

// Run executes the label sync command.
//
//nolint:unparam // placeholder implementation, will return errors in future
func (c *LabelsSyncCmd) Run(ctx context.Context, cli *CLI) error {
	log := logger.FromContext(ctx)
	log.Info("label sync not yet implemented",
		"repo", c.Repo,
		"labels_file", c.LabelsFile,
		"dry_run", cli.DryRun,
	)

	return nil
}

// FilesCmd contains file sync subcommands.
type FilesCmd struct {
	Sync FilesSyncCmd `cmd:"" help:"Sync files to a repository"`
}

// FilesSyncCmd syncs files to a repository.
type FilesSyncCmd struct {
	Repo         string `help:"Target repository (e.g., 'myrepo')" required:""`
	FilesConfig  string `help:"JSON files config" required:""`
	Config       string `help:"JSON sync config (optional)"`
	BranchPrefix string `help:"Branch name prefix" default:"chore/org-sync"`
	PRLabels     string `help:"Comma-separated PR labels" default:"ci/skip-all"`
}

// Run executes the file sync command.
//
//nolint:unparam // placeholder implementation, will return errors in future
func (c *FilesSyncCmd) Run(ctx context.Context, cli *CLI) error {
	log := logger.FromContext(ctx)
	log.Info("file sync not yet implemented",
		"repo", c.Repo,
		"branch_prefix", c.BranchPrefix,
		"dry_run", cli.DryRun,
	)

	return nil
}

// SmyklotCmd contains smyklot sync subcommands.
type SmyklotCmd struct {
	Sync SmyklotSyncCmd `cmd:"" help:"Sync smyklot version to a repository"`
}

// SmyklotSyncCmd syncs smyklot version to a repository.
type SmyklotSyncCmd struct {
	Repo    string `help:"Target repository (e.g., 'myrepo')" required:""`
	Version string `help:"Smyklot version (e.g., '1.9.2')" required:""`
	Tag     string `help:"Smyklot tag (e.g., 'v1.9.2')" required:""`
	Config  string `help:"JSON sync config (optional)"`
}

// Run executes the smyklot sync command.
//
//nolint:unparam // placeholder implementation, will return errors in future
func (c *SmyklotSyncCmd) Run(ctx context.Context, cli *CLI) error {
	log := logger.FromContext(ctx)
	log.Info("smyklot sync not yet implemented",
		"repo", c.Repo,
		"version", c.Version,
		"tag", c.Tag,
		"dry_run", cli.DryRun,
	)

	return nil
}

// ReposCmd contains repository listing subcommands.
type ReposCmd struct {
	List ReposListCmd `cmd:"" help:"List organization repositories"`
}

// ReposListCmd lists organization repositories.
type ReposListCmd struct{}

// Run executes the repos list command.
//
//nolint:unparam // placeholder implementation, will return errors in future
func (*ReposListCmd) Run(ctx context.Context, cli *CLI) error {
	log := logger.FromContext(ctx)
	log.Info("repos list not yet implemented", "org", cli.Org)

	return nil
}

func main() {
	var cli CLI

	ctx := kong.Parse(&cli,
		kong.Name("orgsync"),
		kong.Description("Organization sync tool for labels, files, and smyklot versions"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": version,
		},
	)

	log := logger.New(cli.LogLevel)
	ctx.BindTo(log, (*logger.Logger)(nil))

	appCtx := logger.WithContext(context.Background(), log)

	err := ctx.Run(appCtx, &cli)
	ctx.FatalIfErrorf(err)
}
