// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package binder centralizes the "Cobra + Viper + env + config file" glue.
//
// Adapted from Carolyn Van Slyck's excellent write-up
// "The Sting of the Viper":
// https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/
//
// Goal: for every command, make sure Cobra's flags end up with the final
// value after considering:
//  1. command-line flags
//  2. environment variables (with a prefix, dashes -> underscores)
//  3. config file values
//  4. Cobra's built-in defaults
package binder

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Options describes how to load configuration and env vars for a command.
type Options struct {
	DebugConfig bool

	// EnvPrefix is the prefix for environment variables, e.g. "OTTOAPP".
	// So a flag named "favorite-color" becomes OTTOAPP_FAVORITE_COLOR.
	EnvPrefix string

	// ConfigFile, if set, is an explicit path to a config file.
	// If empty, we will try ConfigName + ConfigPaths.
	ConfigFile string

	// ConfigName is the base name of the config without extension,
	// e.g. "ottoapp".
	ConfigName string

	// ConfigType is the file type, e.g. "json", "yaml", "toml".
	// If empty, Viper will try to guess from the file extension when
	// ConfigFile is set.
	ConfigType string

	// ConfigPaths is a list of directories to search (e.g. $HOME, .).
	ConfigPaths []string

	// Viper is the instance to use. If nil, the package-level viper is used.
	Viper *viper.Viper
}

// Bind wires Viper to the given Cobra command using the provided options.
// Call this from your root command's PersistentPreRunE so it runs for every
// subcommand:
//
//	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
//	    return binder.Bind(cmd, binder.Options{
//	        EnvPrefix:  "OTTOAPP",
//	        ConfigName: "ottoapp",
//	        ConfigPaths: []string{homeDir},
//	        ConfigType: "json",
//	        ConfigFile: cfgFileFlag, // can be empty
//	    })
//	}

func Bind(cmd *cobra.Command, opts Options) error {
	v := opts.Viper
	if v == nil {
		v = viper.GetViper()
	}

	// 1. config file
	if opts.ConfigFile != "" {
		v.SetConfigFile(opts.ConfigFile)
	} else {
		for _, p := range opts.ConfigPaths {
			v.AddConfigPath(p)
		}
		if opts.ConfigName != "" {
			v.SetConfigName(opts.ConfigName)
		}
		if opts.ConfigType != "" {
			v.SetConfigType(opts.ConfigType)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// 2. env
	if opts.EnvPrefix != "" {
		v.SetEnvPrefix(opts.EnvPrefix)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// 3. apply globals to ALL flags (local + inherited)
	applyViperToAllFlags(cmd, v, opts.DebugConfig)

	// 4. apply hierarchical per-command config (commands.foo.bah)
	applyCommandHierarchy(cmd, v, opts.DebugConfig)

	return nil
}

// visitAllFlags walks local + inherited flags.
func visitAllFlags(cmd *cobra.Command, fn func(f *pflag.Flag)) {
	// inherited (persistent from parents)
	cmd.InheritedFlags().VisitAll(fn)
	// local to this command
	cmd.Flags().VisitAll(fn)
}

func applyViperToAllFlags(cmd *cobra.Command, v *viper.Viper, debug bool) {
	visitAllFlags(cmd, func(f *pflag.Flag) {
		// allow env for this key
		_ = v.BindEnv(f.Name)

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			_ = f.Value.Set(fmt.Sprintf("%v", val))
			if debug {
				fmt.Printf("[binder] %s.%s = %v (config/env)\n", cmd.CommandPath(), f.Name, val)
			}
			return
		}

		if debug {
			source := "default"
			if f.Changed {
				source = "cli"
			}
			fmt.Printf("[binder] %s.%s = %v (%s)\n", cmd.CommandPath(), f.Name, f.Value.String(), source)
		}
	})
}

func applyCommandHierarchy(cmd *cobra.Command, v *viper.Viper, debug bool) {
	commandsNode := v.Sub("commands")
	if commandsNode == nil {
		return
	}

	// build chain of names under root: foo, bah
	var chain []string
	for c := cmd; c != nil && c.Parent() != nil; c = c.Parent() {
		if c.Parent() == nil {
			break
		}
		chain = append([]string{c.Name()}, chain...)
	}

	current := commandsNode
	for _, name := range chain {
		next := current.Sub(name)
		if next == nil {
			return
		}
		applySubConfigToAllFlags(cmd, next, debug)
		current = next
	}
}

func applySubConfigToAllFlags(cmd *cobra.Command, sub *viper.Viper, debug bool) {
	visitAllFlags(cmd, func(f *pflag.Flag) {
		if f.Changed {
			return // CLI wins
		}
		if sub.IsSet(f.Name) {
			val := sub.Get(f.Name)
			_ = f.Value.Set(fmt.Sprintf("%v", val))
			if debug {
				fmt.Printf("[binder] %s.%s = %v (command-config)\n", cmd.CommandPath(), f.Name, val)
			}
		}
	})
}

// DumpResolved builds a map of "what the command will actually use":
// 1. start with viper's settings
// 2. overwrite with any cobra flag values (CLI) from this command
//
// This way you see the real, final values.
func DumpResolved(cmd *cobra.Command, v *viper.Viper) map[string]interface{} {
	out := map[string]interface{}{}

	// start with what Viper knows
	if v != nil {
		for k, vval := range v.AllSettings() {
			out[k] = vval
		}
	}

	// overwrite with ALL flags (inherited + local)
	visitAllFlags(cmd, func(f *pflag.Flag) {
		out[f.Name] = f.Value.String()
	})

	return out
}
