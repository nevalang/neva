import { Connection, Module, Network, Node } from "./entity"



interface App {
  commands: {
    addNode(node: Node, network: Network)
    addConnection(connection: Connection, network: Network)
    setRoot(nodeName: string)
  }
}

const addNode()
