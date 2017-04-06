build:
	go install -v -x -a -buildmode=shared runtime sync/atomic #构建核心基本库
	go install -v -x -a -buildmode=shared -linkshared #构建GO动态库

example:
	go build -v -x -linkshared _example/md_example.go
	go build -v -x -linkshared _example/trader_example.go
