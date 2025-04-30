#!/bin/bash

# Pasta de saída
OUTPUT_DIR="build"
APP_NAME="multiglass"
ICON_PATH="bgp-lg.ico"
MAC_ICON_PATH="iconbuilder.icns"

# Limpar a pasta de build
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Gerar rsrc para Windows
rsrc -ico "$ICON_PATH" -o rsrc.syso

# Função para compilar
build() {
  local os=$1
  local arch=$2
  local ext=$3

  echo "Building for $os/$arch..."

  # Definir o arquivo de saída
  output_file="$OUTPUT_DIR/${APP_NAME}-${os}-${arch}${ext}"

  # Executar o build
  GOOS=$os GOARCH=$arch go build -ldflags "-s -w" -o "$output_file"

  if [ $? -eq 0 ]; then
    echo "Successfully built: $output_file"

    # Criar pacote .app para macOS
    if [ "$os" == "darwin" ]; then
      appify --name "$APP_NAME" --icon "$MAC_ICON_PATH" "$output_file"
    fi
  else
    echo "Failed to build: $os/$arch"
  fi
}

# Compilar para macOS, Linux e Windows (amd64 e arm)
build "darwin" "amd64" ""
build "darwin" "arm64" ""
build "linux" "amd64" ""
build "linux" "arm64" ""
build "windows" "amd64" ".exe"
build "windows" "arm64" ".exe"

# Compactar os arquivos para economizar espaço
cd "$OUTPUT_DIR"
for file in *; do
  echo "Compressing $file..."
  zip "${file}.zip" "$file" && rm "$file"
done
cd ..

echo "Builds completed and saved in $OUTPUT_DIR"
