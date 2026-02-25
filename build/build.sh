#!/bin/bash

declare -A oses=(
	['windows']='exe'
	['darwin']='app'
	['linux']='bin'
	['android']='adr'
	['ios']='ios'
)
archs=(
	'amd64'
	'arm64'
)

main="${1:?needs main location}"
path="${2:?needs output location}"
name="${3:-app}"

for arch in "${archs[@]}"; do
	for os in "${!oses[@]}"; do
		if [[ "$os" = 'android' && "$(go env GOOS)" != 'android' ]] || [[ "$os" = 'ios' && "$(go env GOOS)" != 'ios' ]]; then
			continue
		else
			[[ "$arch" = 'amd64' && "$os" = 'ios' ]] && continue

			if [[ "$arch" = 'arm64' && "$os" = 'ios' ]] || [[ "$arch" = 'amd64' && "$os" = 'android' ]]; then
				export CGO_ENABLED=1
			else
				unset CGO_ENABLED
			fi
		fi

		ext="${oses[$os]}"
		app="${name}_$os-$arch.$ext"

		GOOS="$os" GOARCH="$arch" go build -ldflags='-s -w' -buildvcs=false -o "${path%/}/$app" "$main"
		echo "$app OK"
	done
done
