package main

import (
	"bytes"
	"fmt"
	"github.com/shihyuho/go-gitignore/pkg/gitignoreio"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl-gitignore/pkg/formatter"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"text/template"
)

const (
	help = `
建立或增加更多的 .gitignore

從 gitignore.io 取得指定類別 .gitignore 的內容

	$ gitignore go intellij+all

若不指定 ignore 類別, 會有一個互動式選單讓你更方便的做選擇

	$ gitignore

使用 '--outout' 指定路徑, 可以將內容儲存成檔案而不是印在 console
如果指定的路徑檔案已經存在, 將會將內容接續在原本的檔案中

	$ gitignore -o /path/to/.gitignore

如果指定的路徑是目錄, 將會以 '.gitignore' 為預設檔名儲存檔案

	$ gitignore -o .

Environment:

  ${{.ENV}}	set the api url for integrating with gitignore.io
`
)

var (
	// global flags
	offline, _ = strconv.ParseBool(os.Getenv("SL_OFFLINE"))
	verbose, _ = strconv.ParseBool(os.Getenv("SL_VERBOSE"))
	token      = os.Getenv("SL_TOKEN")

	optionSize int
	output     string
	client     = &gitignoreio.Client{}
)

func main() {
	if err := newRootCmd(os.Args[1:]).Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitignore [TYPE...]",
		Short: "fetches a .gitignore from gitignore.io",
		Long:  format(help),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if offline {
				return fmt.Errorf("can not run the command in offline mode")
			}
			client.Log = logrus.New()
			client.Log.SetOutput(cmd.OutOrStdout())
			client.Log.SetFormatter(&formatter.PlainFormatter{})
			if verbose {
				client.Log.SetLevel(logrus.DebugLevel)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				types, err := client.Prompt(optionSize)
				if err != nil {
					return err
				}
				args = append(args, types...)
			}
			content, err := client.Retrieve(args)
			if err != nil {
				return err
			}
			return out(content)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&offline, "offline", offline, "work offline, Overrides $SL_OFFLINE")
	f.BoolVarP(&verbose, "verbose", "v", verbose, "enable verbose output, Overrides $SL_VERBOSE")
	f.StringVar(&token, "token", token, "github access token. Overrides $SL_TOKEN")
	f.StringVar(&client.API, "api", gitignoreio.GetAPI(), "specify api url for integrating with gitignore.io, Overrides $"+gitignoreio.ENV)
	f.StringVarP(&output, "output", "o", "", "specify the output file or directory to save instead of stdout")
	f.IntVar(&optionSize, "option-size", gitignoreio.DefaultPromptOptions, "specify the showing option size for prompt")
	f.Parse(args)

	return cmd
}

func out(content []byte) error {
	if output != "" {
		return client.Save(content, output)
	}
	fmt.Println(string(content))
	return nil
}

func format(tpl string) string {
	var buf bytes.Buffer
	parsed := template.Must(template.New("").Parse(tpl))
	data := make(map[string]string)
	data["ENV"] = gitignoreio.ENV
	err := parsed.Execute(&buf, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return buf.String()
}
