// SPDX-FileCopyrightText: SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

// Package cmd
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tashima42/helm-questions/genquestions"
)

// version is set at build time.
var version = "dev"

var (
	valuesFilename    string
	schemaFilename    string
	questionsFilename string
	header            string
	headerFilename    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helm questions",
	Short: "Generate the Rancher's native chart questions.yml file",
	Long: `Generate the Rancher's native chart questions.yml file based on the
x-rancher-question annotations in values.schema.json.

The schema is usually generated from the values.yaml comments with
the helm-schema plugin, run it first so that the schema is up to date.

Annotation fields (an object, or a list of objects):
  label, order        required
  group               required, except for subquestions
  variable            path relative to the annotated key, needed for keys using $ref
  subquestion_of      variable of the parent question
  type, description, default
                      override the values inferred from the schema and values.yaml
  empty_option_label  label of the "" enum option, "None" by default
Any other field is copied as is (required, min, max, show_if, show_subquestion_if...)`,
	Example: `  helm questions
  helm questions -s charts/app/values.schema.json -v charts/app/values.yaml -q charts/app/questions.yml
  helm questions --header-file HEADER.txt
  helm questions -q - --header ""`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("header") && headerFilename != "" {
			return errors.New("--header and --header-file are mutually exclusive")
		}
		text := header
		if headerFilename != "" {
			data, err := os.ReadFile(headerFilename)
			if err != nil {
				return err
			}
			text = string(data)
		} else if !cmd.Flags().Changed("header") {
			text = genquestions.DefaultHeader
		}

		values := valuesFilename
		if !cmd.Flags().Changed("values-file") {
			values = detectValues(valuesFilename)
		}
		return genquestions.Generate(genquestions.Options{
			ValuesFilename:    values,
			SchemaFilename:    schemaFilename,
			QuestionsFilename: questionsFilename,
			Header:            text,
		})
	},
}

// detectValues falls back to values.yml when the default values.yaml does not exist.
func detectValues(def string) string {
	if _, err := os.Stat(def); err == nil {
		return def
	}
	if _, err := os.Stat("values.yml"); err == nil {
		return "values.yml"
	}
	return def
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&valuesFilename, "values-file", "v", "values.yaml", "helm chart values file, values.yml is used if values.yaml does not exist")
	rootCmd.Flags().StringVarP(&schemaFilename, "schema-file", "s", "values.schema.json", "values json schema file")
	rootCmd.Flags().StringVarP(&questionsFilename, "questions-file", "q", "questions.yml", "questions file to write, - for stdout")
	rootCmd.Flags().StringVar(&header, "header", "", "comment added at the top of the questions file, empty for none (default: generated file notice)")
	rootCmd.Flags().StringVar(&headerFilename, "header-file", "", "file with the comment added at the top of the questions file")
}
