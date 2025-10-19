# History

## Frontend and Functional Programming

Around 2017, when I got experienced enough in JavaScript, I often caught myself imagining what an ideal programming language would look like. I started to write down these ideas without a clear plan in mind. At first those ideas were mostly about syntax.

Later in 2018 the company I was working for started to use TypeScript. This was my first experience with static typing and it blew my mind. I realized how much time a good compiler can save. It affected my notes and focus started to move from syntax to static semantic analysis. I came up with a name "Postulat" meaning strictness as core idea.

It worth noting that all this time I was also interesting in functional programming. It all started with React/Redux and brought me to the book called "An Introduction to Functional Programming Through Lambda Calculus". React introduced me to concepts such as immutability, higher-order components and reactive programming, from the book I learned about lambda reduction and how one can represent almost anything as a function, even numbers. I was kinda obsessed with idea of pure functions and expressing everything in forms functional composition which sometimes I had to say didn't contribute to readability of the code that I wrote.

## Go

Around 2019 I learned Go because of the opportunity in my company to switch into full-stack development, which I always wanted but didn't get the chance. As soon as I become middle frontend developer I started to feel like it's not enough for me to be responsible only for the client, I felt kinda inside the cage, captured in the browser sandbox. Thankfully the company started to adopt Go and started educational program where we learned about algorithms, hardware, OS and Go. Meanwhile I was reading "The Go Programming Language" book by Alan A. A. Donovan and Brian Kernighan and a ton stuff on the internet.

I felt in love with this language and started to spend almost all my time writing in it. I was shocked by simplicity, how minimal it is, how easier it is to work with it compared to TypeScript. I liked it's feature-rich standard library, builtin formatted, single linter everybody uses and the the fact that "it just works". It didn't have lot of features back in the day e.g. there was on generics, yet I felt safier than when I was working with the typescript.

Go was a cure from functional pureness, I've learned more about pointers, garbade-collection, CPUs and kinda embrased the efficiency of imperative programming. Perhaps started to do some Leetcode with Go which also helped.

Finally Go was my first language with real parallelism. I learned about goroutines, channels, mutexes, data-races, deadlocks and race-conditions, and how to avoid them, about concurrency patterns such as fan-in/out, deadline, etc. What I didn't know is that concurrency model Go uses called Communicating-Sequentional Processing (CSP) is actually a form of dataflow.

## Flow-Based Programming

Meanwhile, I don't remember how exactly, but I found out that there is an old programming paradigm nobody knows about. Nobody I know, I mean. It was created in the early 1970s by the british/canadian computer scientist J. Paul Morrison which I later had honor to talk and even work with, not that much unfortunatelly though. It's called "Flow-Based Programming". It's the idea that instead of programming in step-by-step instructions it's much more natural to express computation as data factory - like a set of pipelines that transform the data. The program composed this way is much easier to reason about, especially when it comes down to visual programming and parallelism. 

I had multiple conversations with J.P.Morrison on the FBP discord server, I was lucky enough to contribute to his FBP-Go port a little bit, Egon Elbge contributed much more though. I must mention that mr.Morrison was aware of Neva, he was kind enough to asnwer my stupid questions and provide his opinion. He considered Neva to be "FBP-like" system, rather than pure FBP, because of the reasons that are outside of the scope of this text. I am agree with him on this point.

Flow-based programming (or shorty "FBP") models program like message passing graph where you have nodes, that are instances of components. Nodes have their input and output ports, that are connected to eachother through what's called connections, that are techincally buffered queues. The nodes themselves operate completely concurrently to each-other.

Guess what else operate concurrently and what else is a buffered queue? A goroutine and a channel! Go provides primitives that are perfect for building a flow-based runtime. FBP and Go were a match made in heaven. This was turns out no coinsidence. When I found a book "Dataflow and Reactive Programming Systems: A Practical Guide" by Matt Carkci, I've found out that both CSP and FBP are forms of dataflow.

Finally I was fully aware that there are only two programming paradigms, and all the other are descedants from them: controlflow and dataflow. Go and other languages like Erlang are mixing controlflow with dataflow. What I wanted to create is pure dataflow programming language, using Go as lower-level pepresentation.

Many, many things about Neva are stolen from the J.Paul's FBP. Fundamental abstractions, the way streaming works, the the idea of a program consisting of independent nodes running in full asynchronism, etc. There would be no Neva without FBP.

## Visual Programming

As you probably understand yourself now, at some point, everything clicked. I realized I want pure dataflow staticly typed programming language that compiles to Go. But there was one more thing I didn't mention yet - visual editor.

After I learned about FBP I also decided that text-based programming is not the way to go. I started to think that it's weird that we still program by reading and editing text, just like our grandparents did back in the 50's. The reason for that was that controlflow paradigm doesn't scale well for two things - parallelism and visual programming. It's just difficult to visualize.

Do the idea of the language at a time became that it will be first of all visual, and the source code will be stored in JSON or something like that. This is basically like n8n works now.

I beieved that the language can be really simple, for example without generics. But at the same time I wanted it to be safier than Go. Despite being my beloved language Go had a lot of flaws, you can read about them in articles like 50 shades of Go. I wanted to Neva to be, among many other things, a "fixed Go" in some terms.

The idea was that it will be a visual language that is at the same time as powerful as e.g. Go and also is actually easy for everyone, even person with no programming experience. Those days I was asking my non-programmers friends for the feedback. I managed to write a couple working visual editor prototypes even, and they were able to actually run some nevalang code.

However, as the time went by I realized that it's wrong direction. The "real" programming language will always be too complicated for non-tech people. Please note I am aware of vibe-coding and I'm ok with it. I'm somewhat a vibe-coder myself. What I'm talking about is that you will not make a person who doesn't know anything about programming to understand how concurrent statically typed program that does some stream processing actually work.

Eventually, I moved away from the idea of a purely visual language to a hybrid of textual and visual to get best of two worlds: text remains the single source of truth — easy to version, review, and work with in IDE intellisence — but it can also be represented visually as connected nodes and flows. The syntax is minimal, C-like, similar to Go.

I still do believe in visual programming though, because we as species are much more capable of thinking about computations in terms of boxes and arrows, rather than instructions and function calls, and because dataflow is actually really well suited for that, it would be a shame not to use this advantage.

## Developer Experience

I was programming professionally for more than 10 years now and I always was a person obsessed with developer experience. I've spend who knows how many hours tweaking every little detail in my setup: IDE, Linux. When it came to actual development I always tried to optimize the processes of shipping new features to production. Integrate linting, type-checking, organize code-review, etc. You perhaps might call me DX freak. My twisted brain was constantly thinking what else could be improved? Not the actual product development, but the optimization of the development itself, that was my passion. No surprize I found myself fantasizing and now implementing my dream language.

I have my own very opinionated set of views on how things should be, gained from the years of pain in software engineering practises. I do believe that it is possible to get to perfect balance betwen strictness and flexibility, having a compiler that you don't have to fight but that yet helps you, and that you can really trust. Human-readable, friendly errors. The ease of debug. No hustle dependency management. Even more, effortless concurrency! Dataflow paradigm in FBP-flawor combined with specific ideas from programming language design makes it possible.

Programming is hard enough already. There’s no need to make it harder. Developers should focus on solving engineering problems — designing systems, modeling processes, building algorithms — not wrestling with dependencies, broken builds, and mysterious errors. I wanted to go one step further — to create a language so natural that developers would think:

> “Damn, I wish I were writing this in Neva.”

## Neva Programming Language

This is how Neva project was born. I called it by the name of the river in the sity I was living - Saint-Petersburg. Because you know the river flows just like data in FBP.

Why am I doing it? I've asked this question myself who knows how many times. There were a times when I went completely monk mode and asked my self serious question - is it possible that I am crazy? Working so hard on a project that probably won't gain me any money?

I think it worth it. Simply because there is no other language like it. I feel like I got to the heel from which a specific perspective opens, and my mission is to show it to the other people. Because when I think about Nevalang, I see exactly how it should look.

I can see how much better it can be than everything that exists today.