# Best Practices

> This section is under heavy development

## Formatting

### Line height

Avoid lines that are longer than 80 characters.

1. Even though we have big screens these days, the space that we can focus on with our eyes is limited to about 80 characters. With this line length we can read code without moving the head and with less amount of eyes movement.
2. Professional programmers tend to use split-screens a lot. 120 characters is ok if you only have one editor and probably one sidebar but what about 2 editors? What about 2 editors and sidebar? 2 editors and two sidebars? You got the point.
3. Many people use big fonts because they care about their eyes. Long lines means horizontal scroll for them and scroll itself takes time.
4. Professional programmers sometimes uses different kinds of "code lenses" in their IDEs so e.g. the line that is technically 120 characters becomes actually 150 or 200. Think of argument names lenses for instance.

## Architecture

## Avoid Multiple Parents

Prefer _Pipe_ or _Tree_ structure over graphs where single child have multiple parents. This makes network hard to understand visually. Sometimes it's impossible to avoid that though.
