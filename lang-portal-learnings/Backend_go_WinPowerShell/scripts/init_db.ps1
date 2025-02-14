# Refresh environment variables
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

# Enable CGO
$env:CGO_ENABLED=1

# Run mage db:init
mage db:init
