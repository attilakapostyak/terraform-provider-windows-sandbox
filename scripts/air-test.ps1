$root_dir = Resolve-Path "$PSScriptRoot/../" 
cd $root_dir
$env:TF_ACC="1"
air --build.bin "go test -v $root_dir/..." --build.include_ext "go" --build.stop_on_error true