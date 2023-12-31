/*
Copyright © 2023 rares

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type listConfig struct {
	url    string
	active bool
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:           "list",
	Short:         "list todo items",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiUrl := viper.GetString("api-root")
		isActive := viper.GetBool("active")
		cfg := listConfig{
			url:    apiUrl,
			active: isActive,
		}
		return listAction(os.Stdout, &cfg)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCmd.Flags().BoolP("active", "a", false, "display only active tasks")
	viper.BindPFlag("active", listCmd.Flags().Lookup("active"))

}

func listAction(out io.Writer, cfg *listConfig) error {
	items, err := getAll(cfg)
	if err != nil {
		return err
	}

	return printAll(out, items)
}

func printAll(out io.Writer, items []item) error {
	w := tabwriter.NewWriter(out, 3, 2, 0, ' ', 0)
	for k, v := range items {
		done := ""
		if v.Done {
			done = "X"
		}

		fmt.Fprintf(w, "%s - \t%d\t%s\t\n", done, k+1, v.Task)
	}

	return w.Flush()
}
