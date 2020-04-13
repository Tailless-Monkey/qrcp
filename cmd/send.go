package cmd

import (
	"github.com/claudiodangelis/qrcp/config"
	"github.com/claudiodangelis/qrcp/logger"
	"github.com/claudiodangelis/qrcp/payload"
	"github.com/claudiodangelis/qrcp/qr"

	"github.com/claudiodangelis/qrcp/server"
	"github.com/spf13/cobra"
)

func sendCmdFunc(command *cobra.Command, args []string) error {
	log := logger.New(quietFlag)
	payload, err := payload.FromArgs(args, zipFlag)
	if err != nil {
		return err
	}
	cfg := config.New(command.Flags())
	srv, err := server.New(cfg.Interface, cfg.Port, keepaliveFlag)
	if err != nil {
		return err
	}
	if err := srv.Send(payload); err != nil {
		return err
	}
	log.Print("Scan the following URL with a QR reader to start the file transfer:")
	log.Print(srv.SendURL)
	qr.RenderString(srv.SendURL)
	if err := srv.Wait(); err != nil {
		return err
	}
	log.Print("Done")
	return nil
}

var sendCmd = &cobra.Command{
	Use:     "send",
	Short:   "Send a file(s) or directories from this host",
	Long:    "Send a file(s) or directories from this host",
	Aliases: []string{"s", "transfer"},
	Example: `# Send /path/file.gif. Webserver listens on a random port
qrcp send /path/file.gif
# Shorter version:
qrcp /path/file.gif
# Zip file1.gif and file2.gif, then send the zip package
qrcp /path/file1.gif /path/file2.gif
# Zip teh content of directory, then send the zip package
qrcp /path/directory
# Send file.gif by creating a webserver on port 8080
qrcp --port 8080 /path/file.gif
`,
	Args: cobra.MinimumNArgs(1),
	RunE: sendCmdFunc,
}
