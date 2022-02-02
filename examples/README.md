This example demonstrates a practical usage of gorpheus. In order to compile it please issue the following command:

```
go build -o test
```

you will end up with a binary. There are two sources of migrations here:

- `migrations` directory. These will not be embedded in go and gorpheus will reference them as file paths
- `embedded` directory. These will be embedded into the binary thanks to go's embed package.

here's an list of example executions:

### migrate everything to the latest revision
```
DATABASE_URL=sqlite://example.sqlite ./test
```

Of course if you want to try other options you'd have to remove the database file since it will already be migrated to the very latest version, unless the example specifies otherwise.

### migrate users namespace (including dependencies) to version 2
```
DATABASE_URL=sqlite://example.sqlite ./test users 2
```

### roll back all migrations and drop gorpheus state table
```
DATABASE_URL=sqlite://example.sqlite ./test --vaccuum
```

### roll back all migrations inside users namespace
```
DATABASE_URL=sqlite://example.sqlite ./test users zero
```

### observe backwards migrations
```
DATABASE_URL=sqlite://example.sqlite ./test users
DATABASE_URL=sqlite://example.sqlite ./test users 1
```

the first command would migrate users to revision 2 (as defined in example code) and the second one would roll it back to version 1.