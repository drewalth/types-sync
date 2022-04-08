# types-sync

A small CLI to copy all TypeScript `type`, `enum`, and `interface` declarations in a directory to a single `types.ts`
file.

In one example use case, you are working with a Nodejs API and a Frontend framework that both use TypeScript.
With `types-sync` you can let the API define your models and easily sync them with the client-side.

Inspired by [graphql-code-generator](https://www.graphql-code-generator.com/), an awesome tool that generates TS types
from GraphQL schema.

### Requirements

- [Go](https://go.dev/)

### Setup

Make sure this is in your `~/.bash_profile` or `~/.zshrc`:

```shell
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

### Install Command

```shell
$ go build types-sync.go
$ go install
```

### Usage

By default `types-sync` will look for a `src` directory where the command is run and generate a `types.ts` file in the
current working directory too.

```shell
$ types-sync
```

You can provide your own source directory and destination.

```shell
$ types-sync --src ~/my-project/server/src --output ~/my-project/client/src/types.ts
```

In addition, you can specify which types to exclude. For example, if you are working with a Fastify server, you can
ignore all the Fastify specific typings by adding the `--excludedTypes` flag.

```shell
$ types-sync --excludedTypes Fastify,Express
```

*Fastify types are excluded by default.

### Example

A folder structure like this:

```shell
>  src/bar.ts
>  src/nested-01/nested-02/nested.ts
>  src/nested-01/routes.ts
```

Produces:

```ts
export type MyType = {
    label: string
    nestedOne: {
        val: string
        nestedTwo: {
            foo: boolean
        }
    }
}


export interface MyInterface {
    value: string
}


export enum MyEnum {
    FOO = "FOO"
}

export type NestedType = {
    value: boolean
}

export type RouteType = {
    path: string
    nests: NestedType[]
}
```

