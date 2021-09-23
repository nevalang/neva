package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/emil14/neva/internal/app"
	"github.com/emil14/neva/internal/runtime"
	cli "github.com/urfave/cli/v2"
)

func main() {
	application := app.MustNew()

	app := cli.App{
		Name: "neva",
		Commands: []*cli.Command{
			{
				Name: "compile",
				Action: func(ctx *cli.Context) error {
					bb, err := ioutil.ReadFile(ctx.Args().First())
					if err != nil {
						return err
					}

					compiled, err := application.Compile(bb)
					if err != nil {
						return err
					}

					return ioutil.WriteFile(
						`C:\projects\refactored-garbanzo\examples\arr.json`, compiled, 0644,
					)
				},
			},
			{
				Name: "run",
				Action: func(cliCtx *cli.Context) error {
					bb, err := ioutil.ReadFile(cliCtx.Args().First())
					if err != nil {
						return err
					}

					io, err := application.Run(bb)
					if err != nil {
						return err
					}

					ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

					// reading the input
					go func() {
						reader := bufio.NewReader(os.Stdin)

						for {
							select {
							case <-ctx.Done():
								return
							default:
							}

							fmt.Print("PORT: ")
							inportName, err := reader.ReadString('\n')
							if err != nil {
								fmt.Println(err)
								continue
							}

							inportChan, err := io.In.Port(strings.TrimSuffix(inportName, "\n"), 0)
							if err != nil {
								fmt.Println(err)
								continue
							}

							fmt.Print("MESSAGE: ")
							msg, err := reader.ReadString('\n')
							if err != nil {
								fmt.Println(err)
								continue
							}

							intValue, err := strconv.ParseInt(strings.TrimSuffix(msg, "\n"), 10, 64)
							if err != nil {
								fmt.Println(err)
								continue
							}

							go func() {
								inportChan <- runtime.NewIntMsg(int(intValue))
							}()
						}
					}()

					// Printing the output
					for addr, ch := range io.Out {
						go func(addr runtime.PortAddr, ch chan runtime.Msg) {
							for v := range ch {
								fmt.Print(addr, v)
							}
						}(addr, ch)
					}

					<-ctx.Done()

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
