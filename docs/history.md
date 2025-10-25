# History

## Frontend and Functional Programming

Around 2017, as I grew more fluent in JavaScript, I kept imagining my ideal programming language. I began jotting down ideas without a clear plan. At first, they were mostly about syntax.

In 2018 my company adopted TypeScript. It was my first brush with static typing, and it blew my mind. I realized how much time a good compiler can save. My notes shifted from syntax to static semantics. I even gave the idea a name: "Postulat" — strictness as a core value.

It's worth noting that during this time I was also deeply interested in functional programming. It started with React/Redux and led me to the book "An Introduction to Functional Programming Through Lambda Calculus." React introduced immutability, higher‑order components, and reactive patterns; the book taught me lambda reduction and how almost anything — even numbers — can be represented as functions. I became a bit obsessed with pure functions and composition, sometimes at the expense of readability.

## Go

Around 2019 I learned Go thanks to an opportunity at work to move into full‑stack development. As soon as I became a mid‑level frontend developer, I felt that focusing only on the client wasn't enough — like living in a cage, trapped by the browser sandbox. Luckily, the company was adopting Go and launched an educational program on algorithms, hardware, operating systems, and Go. Meanwhile I read "The Go Programming Language" by Donovan and Kernighan and a lot more online.

I fell in love with Go and started spending most of my time writing it. Its simplicity shocked me — minimal, pragmatic, and easier to work with than TypeScript. I loved its feature‑rich standard library, the built‑in formatter, the single linter everyone uses, and the feeling that "it just works." Back then it lacked many features — no generics, for example — yet I still felt safer than in TypeScript.

Go cured my functional purism. I learned more about pointers, garbage collection, and CPUs, and I embraced the efficiency of imperative programming. I also did a lot of LeetCode in Go, which helped.

Go was my first language with true concurrency. I learned about goroutines, channels, mutexes, data races, deadlocks, and race conditions — and how to avoid them — plus patterns like fan‑in/out and deadlines. I later realized that Go's concurrency model, Communicating Sequential Processes (CSP), is actually a form of dataflow.

## Flow‑Based Programming

Somewhere along the way I discovered an old paradigm almost nobody around me knew: Flow‑Based Programming (FBP), created in the early 1970s by the British‑Canadian computer scientist J. Paul Morrison — whom I later had the honor to talk to and even work with, though only a little. The key idea is that instead of step‑by‑step instructions, computation is a data factory — a set of pipelines that transform data. Programs composed this way are much easier to reason about, especially for visual programming and parallelism.

I had multiple conversations with J. P. Morrison on the FBP Discord server. I was lucky to contribute a bit to his FBP‑Go port (Egon Elbre contributed far more). Mr. Morrison was aware of Neva; he kindly answered my questions and shared his perspective. He considered Neva "FBP‑like" rather than pure FBP — and I agree.

Flow‑based programming models a program as a message‑passing graph: nodes (instances of components) with input and output ports, connected by "connections" that are buffered queues. Nodes run fully concurrently.

Guess what else runs concurrently and uses a buffered queue? Goroutines and channels. Go provides perfect primitives for building a flow‑based runtime. FBP and Go were a match made in heaven. This was no coincidence. In the book "Dataflow and Reactive Programming Systems: A Practical Guide" by Matt Carkci, I learned that both CSP and FBP are forms of dataflow.

Eventually I came to believe there are two root paradigms — control flow and dataflow — and most others descend from them. Languages like Go or Erlang mix the two. I wanted to build a pure dataflow language, with Go as the lower‑level representation.

Many of Neva's ideas come from Morrison's FBP: the core abstractions, streaming semantics, and the vision of independent nodes running in full asynchrony. There would be no Neva without FBP.

## Visual Programming

At some point everything clicked. I realized I wanted a purely dataflow, statically typed language that compiles to Go. And I wanted a visual editor.

After learning about FBP, text‑based programming no longer felt like the way to go. It's strange that we still program by reading and editing text, as in the 1950s. Control‑flow code doesn't scale well for visual programming or parallelism — it's hard to visualize.

So the idea became: the language would be primarily visual, with source code stored in JSON or something similar — much like how n8n works today.

I believed the language could be simple, even without generics, yet safer than Go. As much as I love Go, it has flaws — see "50 Shades of Go," for example. I wanted Neva to be, among other things, a "fixed Go."

The plan was a visual language as powerful as Go, yet approachable even to non‑programmers. I asked friends with no programming background for feedback. I even built a couple of working visual editor prototypes, and they were able to run Neva code.

Over time, I realized that was the wrong direction. A "real" programming language will always be too complex for non‑technical people. I'm aware of vibe‑coding, and I appreciate it; I vibe‑code sometimes myself. But you can't expect someone with no programming background to understand how a concurrent, statically typed program that does stream processing actually works.

Eventually I moved to a hybrid: text and visuals, the best of both worlds. Text remains the single source of truth — easy to version, review, and use with IDE IntelliSense — and it can also be represented visually as connected nodes and flows. The syntax is minimal, C‑like, and reminiscent of Go.

I still believe in visual programming. We, as a species, reason better with boxes and arrows than with instruction streams and function calls, and dataflow fits that shape perfectly. It would be a shame not to use that advantage.

## Developer Experience

I've been programming professionally for more than ten years, and I've always obsessed over developer experience. I've spent countless hours refining every detail of my setup — IDE, Linux, everything. On teams, I tried to optimize the process of shipping features: integrate linting and type checking, organize code review, and so on. You could call me a DX freak. No surprise I ended up fantasizing about — and now building — my dream language.

I have strong opinions, forged by years of pain in software engineering practices. I believe it's possible to strike the right balance between strictness and flexibility: a compiler you don't fight — one that helps you and that you can trust. Human‑readable, friendly errors. Easy debugging. No‑hassle dependency management. And effortless concurrency. The dataflow paradigm, in an FBP flavor, combined with specific language design choices, makes this possible.

Programming is hard enough. There's no need to make it harder. Developers should focus on engineering — designing systems, modeling processes, building algorithms — not wrestling with dependencies, broken builds, and mysterious errors. I wanted to go one step further and create a language so natural that developers would think:

> “Damn, I wish I were writing this in Neva.”

## AI-friendly Language

Neva is not an AI‑first language. You can write it in nano with zero setup. I don't prescribe — or even encourage — AI‑generated code.

That said, I'm obsessed with process optimization. If AI can accelerate development, I'm interested. Neva wasn't designed around AI — the core decisions predate coding agents — but it happens to fit them well:

- Minimal core and a small, orthogonal set of abstractions
- Static typing with relatively strict, predictable compilation rules
- Human‑readable, deterministic compiler errors

This makes AI output easy to validate: we don't rely on the model's "intuition"; the compiler provides precise feedback loops. Generate code, compile, fix, repeat.

## Neva Programming Language

This is how the Neva project was born. I named it after the river in the city where I lived — Saint Petersburg. The river flows, just like data in FBP.

Why am I doing this? I've asked myself that many times. There were periods of monk mode when I asked: was I crazy, working so hard on a project that might never pay?

I think it's worth it. Simply because there's no other language like it. I feel like I reached a hilltop from which a specific perspective opens, and my mission is to share it. When I think about Neva, I see exactly how it should look.

I can see how much better it can be than what exists today.