# mameshiba

A tiny CLI that prints a random fun fact in a speech bubble, said by a Mameshiba-style bean dog.

```
 ╭─────────────────────────────────────────╮
 │ ねえ、知ってる？                        │
 │ It's against the law to burp, or sneeze │
 │ in a church in Nebraska, USA.           │
 ╰─────────────────────────────────────────╯
   ╲
  (◕ ᴥ ◕ ᗜ   )
```

## Build and run

Facts are fetched live, so you need an internet connection.

### Using Go

```
git clone https://github.com/gktgtrkmn/mameshiba
cd mameshiba
go run .
```

To build a binary:

```
go build -o mameshiba .
./mameshiba
```

### With Nix

Run without installing:

```
nix run github:gktgtrkmn/mameshiba
```

Install into your profile:

```
nix profile install github:gktgtrkmn/mameshiba#mameshiba
```

For development, `nix develop` gives you a shell with Go, gopls, and gotools.

## Plan

- Facts currently come from the [Useless Facts API](https://uselessfacts.jsph.pl/).
- Later: a local table of trivia in the style of the original Mameshiba commercials, so the tool works offline.
