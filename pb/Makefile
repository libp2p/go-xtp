PB = xtp-ctl.proto
GO = $(PB:.proto=.pb.go)

all: $(GO)

%.pb.go: %.proto
		protoc --gogo_out=. --proto_path=../../../../../../:/usr/local/opt/protobuf/include:. $<

clean:
		rm $(GO)
