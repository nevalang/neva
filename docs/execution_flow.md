# Execution Flow

Program is network, it consist of nodes. Nodes of that network are called _ports_.
There are _input ports_ that receives data and _output ports_ that sends it.
Every input port (_inport_) is connected to some output port (_outport_).
And Every outport is connected to some inport.
Every inport and every outport is unique.
Every port (in and out) has corresponding `source`, `name` and `type` characteristics.
There are 2 types of `source` for _port_: io and worker.
IO means port belongs to network's io and it serves it as input or output.
Worker means port belongs to what called `worker` - an instance of _component_ (_module_ or _operator_).
The process of data flowing through connection beteween such two ports is what called _stream_.
Objects that flows are _messages_.

## Messages

Under the hood messages are objects of special `Msg` interface.
There are several types of these objects (like `int`, `str`, etc) but they all satisfy the interface.
Interface designed in a way where operators can access the data in a type-safe way.
This makes runtime crashes impossible.

Technically `Msg` is an interface where only getters defined.
There is no "behaviour" in `Msg`, it's just data.
This data is once created and then can never be changes.
That's why don't care about race conditions and other mutations hustle.
