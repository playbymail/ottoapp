#!/bin/bash
# tools/build.sh - create a production distribution tarball in ./dist/prod

# Pre-flight checks
command -v gtar >/dev/null 2>&1 || {
  echo "❌ error: gtar not found. Install GNU tar."
  exit 2
}

[ -d "backend" -a -d "frontend" -a -d "dist" ] || {
  echo "❌ error: must be run from root of the repository"
  exit 2
}

## fetch the version
VERSION=$( go run ./cmd/ottoapp version )
[ "${VERSION}" == "" ] && {
  echo "❌ error: unable to fetch version information from ottoapp"
  exit 2
}
echo "📦  info: building version '${VERSION}'"
backendArtifact="ottoapp-${VERSION}"
frontendArtifact="emberjs-${VERSION}"
tarballArtifact="ottoapp-${VERSION}.tgz"
prodBackend="dist/prod/${backendArtifact}"
prodFrontend="dist/prod/${frontendArtifact}"
prodTarball="dist/prod/${tarballArtifact}"

## remove and recreate the production deployment directory
echo "📦  info: clearing out dist/prod"
rm -rf dist/prod || {
  echo "❌ error: could not clear out dist/prod"
  exit 2
}
mkdir -p dist/prod || {
  echo "❌ error: could not rebuild dist/prod"
  exit 2
}

## build the executable for linux
echo "🛠️  info: building '${prodBackend}'"
CGO_ENABLED=0    # make the executable as static as possible
GOOS=linux
GOARCH=amd64
GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} go build -o "${prodBackend}" ./cmd/ottoapp || {
  echo "❌ error: Go build failed"
  exit 2
}
echo "✅  info: created backend executable: '${prodBackend}'"

## build the ember deployment
echo "🛠️  info: building '${prodFrontend}'"
cd frontend || {
  echo "❌ error: could not set def to frontend"
  exit 2
}
ember build --environment=production || {
  echo "❌ error: ember build failed"
  exit 2
}
cd .. || {
  echo "❌ error: could not set def to repo root"
  exit 2
}
mv frontend/dist "${prodFrontend}"
echo "✅  info: created ember build: '${prodFrontend}'"

## build the deployment tarball
echo "🛠️  info: building '${prodTarball}'"
cd dist/prod || {
  echo "❌ error: failed to set def to dist/prod"
  exit 2
}
gtar -cz -f ${tarballArtifact} --exclude=".DS_Store" ${backendArtifact} ${frontendArtifact} || {
  echo "❌ error: failed to create tarball"
  exit 2
}
echo "✅  info: created tarball: ${prodTarball}"

exit 0
