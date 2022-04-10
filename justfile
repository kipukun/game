alias b := build
defaultOS := 'windows'
out := 'build.exe'

build OS=defaultOS:
	GOOS="{{OS}}" go build -o bin/{{out}}