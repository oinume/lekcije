TOPTARGETS := all clean
SUBDIRS := backend frontend

PROTO_GEN_GO_DIR = backend/proto_gen

$(TOPTARGETS): $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)

.PHONY: $(TOPTARGETS) $(SUBDIRS)

setup:
	$(MAKE) -C backend $(MAKECMDGOALS)

.PHONY: proto/go
proto/go:
	rm -rf $(PROTO_GEN_GO_DIR) && mkdir -p $(PROTO_GEN_GO_DIR)
	protoc -I/usr/local/include -I. \
  		-I$(GOPATH)/src \
  		-Iproto \
  		-Iproto/third_party \
  		--plugin=$(GOBIN)/protoc-gen-twirp \
  		--plugin=$(GOBIN)/protoc-gen-go \
  		--go_out=paths=source_relative:$(PROTO_GEN_GO_DIR) \
  		--twirp_out=paths=source_relative:$(PROTO_GEN_GO_DIR) \
  		proto/api/v1/*.proto
