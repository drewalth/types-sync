# types-sync

A small CLI to copy all TypeScript `type`, `enum`, and `interface` declarations in a directory to a single `types.ts`
file.

Great for when you are working with a Nodejs API and a frontend that both use TypeScript. With `types-sync` you can
write your models once then easily sync them with the client-side helping to reduce duplication.

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

Running `types-sync` on the [`example/`](./example) directory:

```shell
$ types-sync --src ./example
```

Produces a single TS file with consolidated declarations:

```ts
export enum MyEnum {
    VALUE = "VALUE"
}

export interface MyInterface {
    value: string
}

export type MyType = {
    value: boolean
    field: NestedType
}

export type NestedType = {
    value: number
    obj: {
        foo: {
            bar: string[]
        }
    }
}
```

### Git Hooks

Use `types-sync` in a Git hook with [Husky](https://www.npmjs.com/package/husky) from your frontend project.

```shell
$ npx husky add .husky/pre-push "types-sync --src ../server/src --output ./src/types.ts && npm run lint && npm run test && npm run build"
```