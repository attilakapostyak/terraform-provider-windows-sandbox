$plan_script = Resolve-Path "$PSScriptRoot/plan_basic_configuration.ps1"
$root_dir = Resolve-Path "$PSScriptRoot/../" 
cd $root_dir
air --build.cmd "go install $root_dir/." --build.bin "powershell $plan_script" --build.include_ext "go,tf,tfvars" --build.stop_on_error true

