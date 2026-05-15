package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/tliron/commonlog"
	cli "github.com/urfave/cli/v2"

	"github.com/nevalang/neva/pkg/graphdoc"
	"github.com/nevalang/neva/pkg/indexer"
)

// newVisualCmd serves a read-only standalone visual explorer for Neva graphs.
func newVisualCmd(workdir string) *cli.Command {
	return &cli.Command{
		Name:      "visual",
		Usage:     "Run readonly visual graph explorer for a Neva workspace",
		ArgsUsage: "Provide optional workspace path",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "listen",
				Usage: "HTTP listen address",
				Value: "127.0.0.1:7777",
			},
			&cli.BoolFlag{
				Name:  "open",
				Usage: "Open browser automatically",
				Value: true,
			},
		},
		Action: func(cliCtx *cli.Context) error {
			workspacePath := workdir
			if cliCtx.Args().Len() > 0 {
				workspacePath = cliCtx.Args().Get(0)
			}

			idx, err := indexer.NewDefault(commonlog.GetLoggerf("neva.visual"))
			if err != nil {
				return fmt.Errorf("create indexer: %w", err)
			}

			build, found, scanErr := idx.FullScan(cliCtx.Context, workspacePath)
			if scanErr != nil {
				return fmt.Errorf("scan workspace: %w", scanErr)
			}
			if !found {
				return errors.New("no Neva module found in workspace")
			}

			doc := graphdoc.ProjectBuild(build, workspacePath)
			listenAddr := cliCtx.String("listen")

			mux := http.NewServeMux()
			mux.HandleFunc("/api/graph/workspace", func(w http.ResponseWriter, _ *http.Request) {
				writeJSON(w, doc)
			})
			mux.HandleFunc("/api/graph/file", func(w http.ResponseWriter, req *http.Request) {
				fileID := req.URL.Query().Get("id")
				for _, pkg := range doc.Packages {
					for _, file := range pkg.Files {
						if file.ID == fileID {
							writeJSON(w, file)
							return
						}
					}
				}
				http.Error(w, "file not found", http.StatusNotFound)
			})
			mux.HandleFunc("/api/graph/component", func(w http.ResponseWriter, req *http.Request) {
				componentID := req.URL.Query().Get("id")
				for _, pkg := range doc.Packages {
					for _, file := range pkg.Files {
						for _, component := range file.Components {
							if component.ID == componentID {
								writeJSON(w, component)
								return
							}
						}
					}
				}
				http.Error(w, "component not found", http.StatusNotFound)
			})
			mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				if _, err := w.Write([]byte(visualIndexHTML)); err != nil {
					panic(err)
				}
			})

			server := &http.Server{
				Addr:    listenAddr,
				Handler: mux,
			}

			url := "http://" + listenAddr
			fmt.Printf("visual explorer running at %s\n", url)
			if cliCtx.Bool("open") {
				if openErr := openBrowser(cliCtx.Context, url); openErr != nil {
					fmt.Printf("unable to open browser: %v\n", openErr)
				}
			}

			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return fmt.Errorf("serve visual explorer: %w", err)
			}
			return nil
		},
	}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		panic(err)
	}
}

func openBrowser(ctx context.Context, url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.CommandContext(ctx, "open", url)
	case "windows":
		cmd = exec.CommandContext(ctx, "rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.CommandContext(ctx, "xdg-open", url)
	}
	return cmd.Start()
}

const visualIndexHTML = `<!doctype html>
<html>
<head>
  <meta charset="utf-8" />
  <title>Neva Visual (Read-only)</title>
  <style>
    body { font-family: sans-serif; margin: 0; display: grid; grid-template-columns: 380px 1fr; height: 100vh; }
    #left { border-right: 1px solid #ddd; overflow: auto; padding: 12px; }
    #right { overflow: auto; padding: 12px; }
    h2,h3 { margin: 8px 0; }
    button { margin: 4px 0; width: 100%; text-align: left; }
    pre { background: #f6f6f6; padding: 10px; overflow: auto; }
  </style>
</head>
<body>
  <div id="left">
    <h2>Program Explorer</h2>
    <div id="tree"></div>
  </div>
  <div id="right">
    <h2>Component View</h2>
    <div id="details">Select a component.</div>
  </div>
<script>
(async function () {
  const resp = await fetch('/api/graph/workspace');
  const doc = await resp.json();
  const tree = document.getElementById('tree');
  const details = document.getElementById('details');

  for (const pkg of doc.packages) {
    const pkgTitle = document.createElement('h3');
      pkgTitle.textContent = 'Package: ' + pkg.name;
    tree.appendChild(pkgTitle);

    for (const file of pkg.files) {
      const fileTitle = document.createElement('div');
      fileTitle.textContent = 'File: ' + file.name;
      fileTitle.style.fontWeight = 'bold';
      tree.appendChild(fileTitle);

      for (const component of file.components) {
        const btn = document.createElement('button');
        btn.textContent = 'Component ' + component.name;
        btn.onclick = () => {
          details.innerHTML = '';
          const title = document.createElement('h3');
          title.textContent = component.name;
          details.appendChild(title);

          const meta = document.createElement('pre');
          meta.textContent = JSON.stringify(component, null, 2);
          details.appendChild(meta);
        };
        tree.appendChild(btn);
      }

      for (const iface of file.interfaces) {
        const btn = document.createElement('button');
        btn.textContent = 'Interface ' + iface.name;
        btn.onclick = () => {
          details.innerHTML = '';
          const title = document.createElement('h3');
          title.textContent = iface.name;
          details.appendChild(title);

          const meta = document.createElement('pre');
          meta.textContent = JSON.stringify(iface, null, 2);
          details.appendChild(meta);
        };
        tree.appendChild(btn);
      }
    }
  }
})();
</script>
</body>
</html>`
