
Run with 

`go run *.go` 

```
  "UnsetUserArray": null,
  "UnsetIntArray": null,
```

contrast the output with 

`. environment; go run *.go`

```
  "UnsetUserArray": [
    "tom",
    "dick",
    "harry"
  ],
  "UnsetIntArray": [
    9,
    8,
    0,
    5,
    4,
    3,
    2,
    1
  ],

```

