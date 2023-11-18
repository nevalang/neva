# ARCHITECTURE

```mermaid
flowchart LR
    language-server-->vscode
    subgraph language-server
        indexer
    end
    subgraph vscode
        webview-->extension
        extension-->webview
    end
    webview-->vscode
    vscode-->language-server
```
