TAG=`git describe --tags`
VERSION ?= `git describe --tags`
LDFLAGS=-ldflags "-s -extldflags \"--static\" -w -X main.version=${VERSION}"

build = echo "\n\nBuilding $(1)-$(2)" && GOOS=$(1) GOARCH=$(2) go build ${LDFLAGS} -o dist/gitlab-pipe-cleaner_${VERSION}_$(1)_$(2) \
	&& bzip2 dist/gitlab-pipe-cleaner_${VERSION}_$(1)_$(2)

gitlab-pipe-cleaner: *.go
	go build ${LDFLAGS} -o gitlab-pipe-cleaner

clean:
	rm -f gitlab-pipe-cleaner

release:
	mkdir -p dist
	rm -f dist/gitlab-pipe-cleaner_${VERSION}_*
	$(call build,linux,amd64)
	$(call build,linux,386)
	$(call build,linux,arm)
	$(call build,linux,arm64)
	$(call build,darwin,amd64)
	$(call build,darwin,386)
	$(call build,windows,386)
	$(call build,windows,amd64)


