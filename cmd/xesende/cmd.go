package main

import (
	"flag"
	"log"
	"os"

	"github.com/gobs/pretty"
	"hawx.me/code/hadfield"
	"hawx.me/code/xesende"
)

var (
	account  = flag.String("account", "", "")
	username = flag.String("username", "", "")
	password = flag.String("password", "", "")
)

const pageSize = 20

const authHelp = `Authentication required:
  Either pass the --username USER and --password PASS options, or set the
  ESENDEX_USERNAME and ESENDEX_PASSWORD environment variables.
`

func pageOpts(page int) xesende.Option {
	startIndex := (page - 1) * pageSize

	return xesende.Page(startIndex, pageSize)
}

var templates = hadfield.Templates{
	Help: `usage: xesende [command] [arguments]

  A command line client for the Esendex REST API.

  Options:
    --username USER    # Username to authenticate with
    --password PASS    # Password to authenticate with
    --account REF      # Account reference to use
    --help             # Display this message

  Commands: {{range .}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}
`,
	Command: `usage: xesende {{.Usage}}
{{.Long}}
`,
}

func main() {
	flag.Parse()

	if *username == "" {
		*username = os.Getenv("ESENDEX_USERNAME")
	}

	if *password == "" {
		*password = os.Getenv("ESENDEX_PASSWORD")
	}

	if *account == "" {
		*account = os.Getenv("ESENDEX_ACCOUNT")
	}

	if *username == "" || *password == "" {
		log.Fatal(authHelp)
	}

	client := xesende.New(*username, *password)

	commands := hadfield.Commands{
		receivedCmd(client),
		sentCmd(client),
		messageCmd(client),
		accountsCmd(client),
	}

	if *account == "" {
		commands = append(commands,
			receivedCmd(client),
			sentCmd(client))
	} else {
		accountClient := client.Account(*account)

		commands = append(commands,
			receivedCmd(accountClient),
			sentCmd(accountClient),
			sendCmd(accountClient))
	}

	hadfield.Run(commands, templates)
}

type receivedClient interface {
	Received(opts ...xesende.Option) (*xesende.ReceivedMessagesResponse, error)
}

func receivedCmd(client receivedClient) *hadfield.Command {
	var page int

	cmd := &hadfield.Command{
		Usage: "received [options]",
		Short: "lists received messages",
		Long: `
  Received displays a list of received messages.

    --page NUM       # Display given page
`,
		Run: func(cmd *hadfield.Command, args []string) {
			resp, err := client.Received()
			if err != nil {
				log.Fatal(err)
			}

			pretty.PrettyPrint(resp.Messages)
		},
	}

	cmd.Flag.IntVar(&page, "page", 0, "")

	return cmd
}

type sentClient interface {
	Sent(opts ...xesende.Option) (*xesende.SentMessagesResponse, error)
}

func sentCmd(client sentClient) *hadfield.Command {
	var page int

	cmd := &hadfield.Command{
		Usage: "sent [options]",
		Short: "lists sent messages",
		Long: `
  Sent displays a list of sent messages.

    --page NUM       # Display given page
`,
		Run: func(cmd *hadfield.Command, args []string) {
			resp, err := client.Sent(pageOpts(page))
			if err != nil {
				log.Fatal(err)
			}

			pretty.PrettyPrint(resp.Messages)
		},
	}

	cmd.Flag.IntVar(&page, "page", 1, "")

	return cmd
}

func messageCmd(client *xesende.Client) *hadfield.Command {
	return &hadfield.Command{
		Usage: "message MESSAGEID",
		Short: "displays a message",
		Long: `
  Message displays the details for a message.
`,
		Run: func(cmd *hadfield.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Require MESSAGEID parameter")
			}

			resp, err := client.Message(args[0])
			if err != nil {
				log.Fatal(err)
			}

			pretty.PrettyPrint(resp)
		},
	}
}

func accountsCmd(client *xesende.Client) *hadfield.Command {
	return &hadfield.Command{
		Usage: "accounts",
		Short: "list accounts",
		Long: `
  List accounts available to the user.
`,
		Run: func(cmd *hadfield.Command, args []string) {
			resp, err := client.Accounts()
			if err != nil {
				log.Fatal(err)
			}

			pretty.PrettyPrint(resp.Accounts)
		},
	}
}

func sendCmd(client *xesende.AccountClient) *hadfield.Command {
	return &hadfield.Command{
		Usage: "send TO BODY",
		Short: "send messages",
		Long: `
  Send a single sms message.
`,
		Run: func(cmd *hadfield.Command, args []string) {
			resp, err := client.Send([]xesende.Message{
				{
					To:   args[0],
					Body: args[1],
				},
			})

			if err != nil {
				log.Fatal(err)
			}

			pretty.PrettyPrint(resp)
		},
	}
}
