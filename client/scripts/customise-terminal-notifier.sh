#!/bin/bash

readonly program="$(basename "${0}")"

function syntax_error {
  echo "${program}: ${1}" >&2
  echo "Try \`${program} --help\` for more information." >&2
  exit 1
}

# instructions
function usage {
  echo "
    usage: ${program} -i <file> -b <bundle_id>

    options:
      -i <file>, --icon <file>               Set the icon (can be either an 'icns' or an image format, like 'png')
      -b <bundle_id>, --bundle <bundle_id>   Set the bundle id
      -h, --help                             Show this message
  " | sed -E 's/^ {4}//'
}

# set options
while [[ "${1}" ]]; do
  case "${1}" in
    -h | --help)
      usage
      exit 0
      ;;
    -i | --icon)
      icon="${2}"
      shift
      ;;
    -b | --bundle)
      id="${2}"
      shift
      ;;
    -*)
      syntax_error "unrecognized option: ${1}"
      ;;
  esac
  shift
done

function make_icns {
  local file="${1}"
  local iconset="$(mktemp -d)"
  local output_icon="$(mktemp).icns"

  for size in {16,32,64,128,256,512}; do
    sips --resampleHeightWidth "${size}" "${size}" "${file}" --out "${iconset}/icon_${size}x${size}.png" &> /dev/null
    sips --resampleHeightWidth "$((size * 2))" "$((size * 2))" "${file}" --out "${iconset}/icon_${size}x${size}@2x.png" &> /dev/null
  done

  mv "${iconset}" "${iconset}.iconset"
  iconutil --convert icns "${iconset}.iconset" --output "${output_icon}"

  echo "${output_icon}" # so its path is returned when the function ends
}

# stop executing if an option is missing
if [[ -z "${icon}" || -z "${id}" ]]; then usage && exit 1; fi

# set variables
tn_version='2.0.0'
tmp_dir="$(mktemp -d -t 'terminal-notifier')"
#app="${tmp_dir}/terminal-notifier-${tn_version}/terminal-notifier.app"
app="${tmp_dir}/terminal-notifier.app"

# get terminal notifier
curl --progress-bar --location "https://github.com/julienXX/terminal-notifier/releases/download/${tn_version}/terminal-notifier-${tn_version}.zip" | ditto -xk - "${tmp_dir}"

# set icon and bundle id
[[ "${icon}" != *'.icns' ]] && icon=$(make_icns "${icon}") # convert icon, if not already 'icns'
cp "${icon}" "${app}/Contents/Resources/Terminal.icns"
sed -i '' "s/fr.julienxx.oss.terminal-notifier/${id}/" "${app}/Contents/Info.plist"

# move it to the Desktop
mv "${app}" "${HOME}/Desktop"