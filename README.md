# nina-batch

- go get
```shell
go get -u github.com/ningenMe/mami-interface@v0.x.0
```

```shell
PAT=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
go run . -t ${PAT} -s 2022-01-01 -e 2023-01-01
go run . -t ${PAT} -s 2022-01-01 -e 2023-01-01 -r `cat repo.txt`
```