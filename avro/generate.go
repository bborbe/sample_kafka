package generate

//go:generate mkdir -p ./avro
//go:generate $GOPATH/bin/gogen-avro --containers ./avro schema.avsc
