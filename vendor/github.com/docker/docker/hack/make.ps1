<#
.NOTES
    Author:  @jhowardmsft

    Summary: Windows native build script. This is similar to functionality provided
             by hack\make.sh, but uses native Windows PowerShell semantics. It does
             not support the full set of options provided by the Linux counterpart.
             For example:

             - You can't cross-build Linux docker binaries on Windows
             - Hashes aren't generated on binaries
             - 'Releasing' isn't supported.
             - Integration tests. This is because they currently cannot run inside a container,
               and require significant external setup.

             It does however provided the minimum necessary to support parts of local Windows
             development and Windows to Windows CI.

             Usage Examples (run from repo root):
                "hack\make.ps1 -Client" to build docker.exe client 64-bit binary (remote repo)
                "hack\make.ps1 -TestUnit" to run unit tests
                "hack\make.ps1 -Daemon -TestUnit" to build the daemon and run unit tests
                "hack\make.ps1 -All" to run everything this script knows about that can run in a container
                "hack\make.ps1" to build the daemon binary (same as -Daemon)
                "hack\make.ps1 -Binary" shortcut to -Client and -Daemon

.PARAMETER Client
     Builds the client binaries.

.PARAMETER Daemon
     Builds the daemon binary.

.PARAMETER Binary
     Builds the client and daemon binaries. A convenient shortcut to `make.ps1 -Client -Daemon`.

.PARAMETER Race
     Use -race in go build and go test.

.PARAMETER Noisy
     Use -v in go build.

.PARAMETER ForceBuildAll
     Use -a in go build.

.PARAMETER NoOpt
     Use -gcflags -N -l in go build to disable optimisation (can aide debugging).

.PARAMETER CommitSuffix
     Adds a custom string to be appended to the commit ID (spaces are stripped).

.PARAMETER DCO
     Runs the DCO (Developer Certificate Of Origin) test (must be run outside a container).

.PARAMETER PkgImports
     Runs the pkg\ directory imports test (must be run outside a container).

.PARAMETER GoFormat
     Runs the Go formatting test (must be run outside a container).

.PARAMETER TestUnit
     Runs unit tests.

.PARAMETER All
     Runs everything this script knows about that can run in a container.


TODO
- Unify the head commit
- Add golint and other checks (swagger maybe?)

#>


param(
    [Parameter(Mandatory=$False)][switch]$Client,
    [Parameter(Mandatory=$False)][switch]$Daemon,
    [Parameter(Mandatory=$False)][switch]$Binary,
    [Parameter(Mandatory=$False)][switch]$Race,
    [Parameter(Mandatory=$False)][switch]$Noisy,
    [Parameter(Mandatory=$False)][switch]$ForceBuildAll,
    [Parameter(Mandatory=$False)][switch]$NoOpt,
    [Parameter(Mandatory=$False)][string]$CommitSuffix="",
    [Parameter(Mandatory=$False)][switch]$DCO,
    [Parameter(Mandatory=$False)][switch]$PkgImports,
    [Parameter(Mandatory=$False)][switch]$GoFormat,
    [Parameter(Mandatory=$False)][switch]$TestUnit,
    [Parameter(Mandatory=$False)][switch]$All
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"
$pushed=$False  # To restore the directory if we have temporarily pushed to one.

# Utility function to get the commit ID of the repository
Function Get-GitCommit() ***REMOVED***
    if (-not (Test-Path ".\.git")) ***REMOVED***
        # If we don't have a .git directory, but we do have the environment
        # variable DOCKER_GITCOMMIT set, that can override it.
        if ($env:DOCKER_GITCOMMIT.Length -eq 0) ***REMOVED***
            Throw ".git directory missing and DOCKER_GITCOMMIT environment variable not specified."
    ***REMOVED***
        Write-Host "INFO: Git commit ($env:DOCKER_GITCOMMIT) assumed from DOCKER_GITCOMMIT environment variable"
        return $env:DOCKER_GITCOMMIT
***REMOVED***
    $gitCommit=$(git rev-parse --short HEAD)
    if ($(git status --porcelain --untracked-files=no).Length -ne 0) ***REMOVED***
        $gitCommit="$gitCommit-unsupported"
        Write-Host ""
        Write-Warning "This version is unsupported because there are uncommitted file(s)."
        Write-Warning "Either commit these changes, or add them to .gitignore."
        git status --porcelain --untracked-files=no | Write-Warning
        Write-Host ""
***REMOVED***
    return $gitCommit
***REMOVED***

# Utility function to determine if we are running in a container or not.
# In Windows, we get this through an environment variable set in `Dockerfile.Windows`
Function Check-InContainer() ***REMOVED***
    if ($env:FROM_DOCKERFILE.Length -eq 0) ***REMOVED***
        Write-Host ""
        Write-Warning "Not running in a container. The result might be an incorrect build."
        Write-Host ""
        return $False
***REMOVED***
    return $True
***REMOVED***

# Utility function to warn if the version of go is correct. Used for local builds
# outside of a container where it may be out of date with master.
Function Verify-GoVersion() ***REMOVED***
    Try ***REMOVED***
        $goVersionDockerfile=(Get-Content ".\Dockerfile" | Select-String "ENV GO_VERSION").ToString().Split(" ")[2]
        $goVersionInstalled=(go version).ToString().Split(" ")[2].SubString(2)
***REMOVED***
    Catch [Exception] ***REMOVED***
        Throw "Failed to validate go version correctness: $_"
***REMOVED***
    if (-not($goVersionInstalled -eq $goVersionDockerfile)) ***REMOVED***
        Write-Host ""
        Write-Warning "Building with golang version $goVersionInstalled. You should update to $goVersionDockerfile"
        Write-Host ""
***REMOVED***
***REMOVED***

# Utility function to get the commit for HEAD
Function Get-HeadCommit() ***REMOVED***
    $head = Invoke-Expression "git rev-parse --verify HEAD"
    if ($LASTEXITCODE -ne 0) ***REMOVED*** Throw "Failed getting HEAD commit" ***REMOVED***

    return $head
***REMOVED***

# Utility function to get the commit for upstream
Function Get-UpstreamCommit() ***REMOVED***
    Invoke-Expression "git fetch -q https://github.com/docker/docker.git refs/heads/master"
    if ($LASTEXITCODE -ne 0) ***REMOVED*** Throw "Failed fetching" ***REMOVED***

    $upstream = Invoke-Expression "git rev-parse --verify FETCH_HEAD"
    if ($LASTEXITCODE -ne 0) ***REMOVED*** Throw "Failed getting upstream commit" ***REMOVED***

    return $upstream
***REMOVED***

# Build a binary (client or daemon)
Function Execute-Build($type, $additionalBuildTags, $directory) ***REMOVED***
    # Generate the build flags
    $buildTags = "autogen"
    if ($Noisy)                     ***REMOVED*** $verboseParm=" -v" ***REMOVED***
    if ($Race)                      ***REMOVED*** Write-Warning "Using race detector"; $raceParm=" -race"***REMOVED***
    if ($ForceBuildAll)             ***REMOVED*** $allParm=" -a" ***REMOVED***
    if ($NoOpt)                     ***REMOVED*** $optParm=" -gcflags "+""""+"-N -l"+"""" ***REMOVED***
    if ($additionalBuildTags -ne "") ***REMOVED*** $buildTags += $(" " + $additionalBuildTags) ***REMOVED***

    # Do the go build in the appropriate directory
    # Note -linkmode=internal is required to be able to debug on Windows.
    # https://github.com/golang/go/issues/14319#issuecomment-189576638
    Write-Host "INFO: Building $type..."
    Push-Location $root\cmd\$directory; $global:pushed=$True
    $buildCommand = "go build" + `
                    $raceParm + `
                    $verboseParm + `
                    $allParm + `
                    $optParm + `
                    " -tags """ + $buildTags + """" + `
                    " -ldflags """ + "-linkmode=internal" + """" + `
                    " -o $root\bundles\"+$directory+".exe"
    Invoke-Expression $buildCommand
    if ($LASTEXITCODE -ne 0) ***REMOVED*** Throw "Failed to compile $type" ***REMOVED***
    Pop-Location; $global:pushed=$False
***REMOVED***


# Validates the DCO marker is present on each commit
Function Validate-DCO($headCommit, $upstreamCommit) ***REMOVED***
    Write-Host "INFO: Validating Developer Certificate of Origin..."
    # Username may only contain alphanumeric characters or dashes and cannot begin with a dash
    $usernameRegex='[a-zA-Z0-9][a-zA-Z0-9-]+'

    $dcoPrefix="Signed-off-by:"
    $dcoRegex="^(Docker-DCO-1.1-)?$dcoPrefix ([^<]+) <([^<>@]+@[^<>]+)>( \(github: ($usernameRegex)\))?$"

    $counts = Invoke-Expression "git diff --numstat $upstreamCommit...$headCommit"
    if ($LASTEXITCODE -ne 0) ***REMOVED*** Throw "Failed git diff --numstat" ***REMOVED***

    # Counts of adds and deletes after removing multiple white spaces. AWK anyone? :(
    $adds=0; $dels=0; $($counts -replace '\s+', ' ') | %***REMOVED*** 
        $a=$_.Split(" "); 
        if ($a[0] -ne "-") ***REMOVED*** $adds+=[int]$a[0] ***REMOVED***
        if ($a[1] -ne "-") ***REMOVED*** $dels+=[int]$a[1] ***REMOVED***
***REMOVED***
    if (($adds -eq 0) -and ($dels -eq 0)) ***REMOVED*** 
        Write-Warning "DCO validation - nothing to validate!"
        return
***REMOVED***

    $commits = Invoke-Expression "git log  $upstreamCommit..$headCommit --format=format:%H%n"
    if ($LASTEXITCODE -ne 0) ***REMOVED*** Throw "Failed git log --format" ***REMOVED***
    $commits = $($commits -split '\s+' -match '\S')
    $badCommits=@()
    $commits | %***REMOVED***
        # Skip commits with no content such as merge commits etc
        if ($(git log -1 --format=format: --name-status $_).Length -gt 0) ***REMOVED***
            # Ignore exit code on next call - always process regardless
            $commitMessage = Invoke-Expression "git log -1 --format=format:%B --name-status $_"
            if (($commitMessage -match $dcoRegex).Length -eq 0) ***REMOVED*** $badCommits+=$_ ***REMOVED***
    ***REMOVED***
***REMOVED***
    if ($badCommits.Length -eq 0) ***REMOVED***
        Write-Host "Congratulations!  All commits are properly signed with the DCO!"
***REMOVED*** else ***REMOVED***
        $e = "`nThese commits do not have a proper '$dcoPrefix' marker:`n"
        $badCommits | %***REMOVED*** $e+=" - $_`n"***REMOVED***
        $e += "`nPlease amend each commit to include a properly formatted DCO marker.`n`n"
        $e += "Visit the following URL for information about the Docker DCO:`n"
        $e += "https://github.com/docker/docker/blob/master/CONTRIBUTING.md#sign-your-work`n"
        Throw $e
***REMOVED***
***REMOVED***

# Validates that .\pkg\... is safely isolated from internal code
Function Validate-PkgImports($headCommit, $upstreamCommit) ***REMOVED***
    Write-Host "INFO: Validating pkg import isolation..."

    # Get a list of go source-code files which have changed under pkg\. Ignore exit code on next call - always process regardless
    $files=@(); $files = Invoke-Expression "git diff $upstreamCommit...$headCommit --diff-filter=ACMR --name-only -- `'pkg\*.go`'"
    $badFiles=@(); $files | %***REMOVED***
        $file=$_
        # For the current changed file, get its list of dependencies, sorted and uniqued.
        $imports = Invoke-Expression "go list -e -f `'***REMOVED******REMOVED*** .Deps ***REMOVED******REMOVED***`' $file"
        if ($LASTEXITCODE -ne 0) ***REMOVED*** Throw "Failed go list for dependencies on $file" ***REMOVED***
        $imports = $imports -Replace "\[" -Replace "\]", "" -Split(" ") | Sort-Object | Get-Unique
        # Filter out what we are looking for
        $imports = $imports -NotMatch "^github.com/docker/docker/pkg/" `
                            -NotMatch "^github.com/docker/docker/vendor" `
                            -Match "^github.com/docker/docker" `
                            -Replace "`n", ""
        $imports | % ***REMOVED*** $badFiles+="$file imports $_`n" ***REMOVED***
***REMOVED***
    if ($badFiles.Length -eq 0) ***REMOVED***
        Write-Host 'Congratulations!  ".\pkg\*.go" is safely isolated from internal code.'
***REMOVED*** else ***REMOVED***
        $e = "`nThese files import internal code: (either directly or indirectly)`n"
        $badFiles | %***REMOVED*** $e+=" - $_"***REMOVED***
        Throw $e
***REMOVED***
***REMOVED***

# Validates that changed files are correctly go-formatted
Function Validate-GoFormat($headCommit, $upstreamCommit) ***REMOVED***
    Write-Host "INFO: Validating go formatting on changed files..."

    # Verify gofmt is installed
    if ($(Get-Command gofmt -ErrorAction SilentlyContinue) -eq $nil) ***REMOVED*** Throw "gofmt does not appear to be installed" ***REMOVED***

    # Get a list of all go source-code files which have changed.  Ignore exit code on next call - always process regardless
    $files=@(); $files = Invoke-Expression "git diff $upstreamCommit...$headCommit --diff-filter=ACMR --name-only -- `'*.go`'"
    $files = $files | Select-String -NotMatch "^vendor/"
    $badFiles=@(); $files | %***REMOVED***
        # Deliberately ignore error on next line - treat as failed
        $content=Invoke-Expression "git show $headCommit`:$_"

        # Next set of hoops are to ensure we have LF not CRLF semantics as otherwise gofmt on Windows will not succeed.
        # Also note that gofmt on Windows does not appear to support stdin piping correctly. Hence go through a temporary file.
        $content=$content -join "`n"
        $content+="`n"
        $outputFile=[System.IO.Path]::GetTempFileName()
        if (Test-Path $outputFile) ***REMOVED*** Remove-Item $outputFile ***REMOVED***
        [System.IO.File]::WriteAllText($outputFile, $content, (New-Object System.Text.UTF8Encoding($False)))
        $currentFile = $_ -Replace("/","\")
        Write-Host Checking $currentFile
        Invoke-Expression "gofmt -s -l $outputFile"
        if ($LASTEXITCODE -ne 0) ***REMOVED*** $badFiles+=$currentFile ***REMOVED***
        if (Test-Path $outputFile) ***REMOVED*** Remove-Item $outputFile ***REMOVED***
***REMOVED***
    if ($badFiles.Length -eq 0) ***REMOVED***
        Write-Host 'Congratulations!  All Go source files are properly formatted.'
***REMOVED*** else ***REMOVED***
        $e = "`nThese files are not properly gofmt`'d:`n"
        $badFiles | %***REMOVED*** $e+=" - $_`n"***REMOVED***
        $e+= "`nPlease reformat the above files using `"gofmt -s -w`" and commit the result."
        Throw $e
***REMOVED***
***REMOVED***

# Run the unit tests
Function Run-UnitTests() ***REMOVED***
    Write-Host "INFO: Running unit tests..."
    $testPath="./..."
    $goListCommand = "go list -e -f '***REMOVED******REMOVED***if ne .Name """ + '\"github.com/docker/docker\"' + """***REMOVED******REMOVED******REMOVED******REMOVED***.ImportPath***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***' $testPath"
    $pkgList = $(Invoke-Expression $goListCommand)
    if ($LASTEXITCODE -ne 0) ***REMOVED*** Throw "go list for unit tests failed" ***REMOVED***
    $pkgList = $pkgList | Select-String -Pattern "github.com/docker/docker"
    $pkgList = $pkgList | Select-String -NotMatch "github.com/docker/docker/vendor"
    $pkgList = $pkgList | Select-String -NotMatch "github.com/docker/docker/man"
    $pkgList = $pkgList | Select-String -NotMatch "github.com/docker/docker/integration"
    $pkgList = $pkgList -replace "`r`n", " "
    $goTestCommand = "go test" + $raceParm + " -cover -ldflags -w -tags """ + "autogen daemon" + """ -a """ + "-test.timeout=10m" + """ $pkgList"
    Invoke-Expression $goTestCommand
    if ($LASTEXITCODE -ne 0) ***REMOVED*** Throw "Unit tests failed" ***REMOVED***
***REMOVED***

# Start of main code.
Try ***REMOVED***
    Write-Host -ForegroundColor Cyan "INFO: make.ps1 starting at $(Get-Date)"

    # Get to the root of the repo
    $root = $(Split-Path $MyInvocation.MyCommand.Definition -Parent | Split-Path -Parent)
    Push-Location $root

    # Handle the "-All" shortcut to turn on all things we can handle.
    # Note we expressly only include the items which can run in a container - the validations tests cannot
    # as they require the .git directory which is excluded from the image by .dockerignore
    if ($All) ***REMOVED*** $Client=$True; $Daemon=$True; $TestUnit=$True ***REMOVED***

    # Handle the "-Binary" shortcut to build both client and daemon.
    if ($Binary) ***REMOVED*** $Client = $True; $Daemon = $True ***REMOVED***

    # Default to building the daemon if not asked for anything explicitly.
    if (-not($Client) -and -not($Daemon) -and -not($DCO) -and -not($PkgImports) -and -not($GoFormat) -and -not($TestUnit)) ***REMOVED*** $Daemon=$True ***REMOVED***

    # Verify git is installed
    if ($(Get-Command git -ErrorAction SilentlyContinue) -eq $nil) ***REMOVED*** Throw "Git does not appear to be installed" ***REMOVED***

    # Verify go is installed
    if ($(Get-Command go -ErrorAction SilentlyContinue) -eq $nil) ***REMOVED*** Throw "GoLang does not appear to be installed" ***REMOVED***

    # Get the git commit. This will also verify if we are in a repo or not. Then add a custom string if supplied.
    $gitCommit=Get-GitCommit
    if ($CommitSuffix -ne "") ***REMOVED*** $gitCommit += "-"+$CommitSuffix -Replace ' ', '' ***REMOVED***

    # Get the version of docker (eg 17.04.0-dev)
    $dockerVersion="0.0.0-dev"

    # Give a warning if we are not running in a container and are building binaries or running unit tests.
    # Not relevant for validation tests as these are fine to run outside of a container.
    if ($Client -or $Daemon -or $TestUnit) ***REMOVED*** $inContainer=Check-InContainer ***REMOVED***

    # If we are not in a container, validate the version of GO that is installed.
    if (-not $inContainer) ***REMOVED*** Verify-GoVersion ***REMOVED***

    # Verify GOPATH is set
    if ($env:GOPATH.Length -eq 0) ***REMOVED*** Throw "Missing GOPATH environment variable. See https://golang.org/doc/code.html#GOPATH" ***REMOVED***

    # Run autogen if building binaries or running unit tests.
    if ($Client -or $Daemon -or $TestUnit) ***REMOVED***
        Write-Host "INFO: Invoking autogen..."
        Try ***REMOVED*** .\hack\make\.go-autogen.ps1 -CommitString $gitCommit -DockerVersion $dockerVersion -Platform "$env:PLATFORM" ***REMOVED***
        Catch [Exception] ***REMOVED*** Throw $_ ***REMOVED***
***REMOVED***

    # DCO, Package import and Go formatting tests.
    if ($DCO -or $PkgImports -or $GoFormat) ***REMOVED***
        # We need the head and upstream commits for these
        $headCommit=Get-HeadCommit
        $upstreamCommit=Get-UpstreamCommit

        # Run DCO validation
        if ($DCO) ***REMOVED*** Validate-DCO $headCommit $upstreamCommit ***REMOVED***

        # Run `gofmt` validation
        if ($GoFormat) ***REMOVED*** Validate-GoFormat $headCommit $upstreamCommit ***REMOVED***

        # Run pkg isolation validation
        if ($PkgImports) ***REMOVED*** Validate-PkgImports $headCommit $upstreamCommit ***REMOVED***
***REMOVED***

    # Build the binaries
    if ($Client -or $Daemon) ***REMOVED***
        # Create the bundles directory if it doesn't exist
        if (-not (Test-Path ".\bundles")) ***REMOVED*** New-Item ".\bundles" -ItemType Directory | Out-Null ***REMOVED***

        # Perform the actual build
        if ($Daemon) ***REMOVED*** Execute-Build "daemon" "daemon" "dockerd" ***REMOVED***
        if ($Client) ***REMOVED***
            # Get the Docker channel and version from the environment, or use the defaults.
            if (-not ($channel = $env:DOCKERCLI_CHANNEL)) ***REMOVED*** $channel = "edge" ***REMOVED***
            if (-not ($version = $env:DOCKERCLI_VERSION)) ***REMOVED*** $version = "17.06.0-ce" ***REMOVED***

            # Download the zip file and extract the client executable.
            Write-Host "INFO: Downloading docker/cli version $version from $channel..."
            $url = "https://download.docker.com/win/static/$channel/x86_64/docker-$version.zip"
            Invoke-WebRequest $url -OutFile "docker.zip"
            Try ***REMOVED***
                Add-Type -AssemblyName System.IO.Compression.FileSystem
                $zip = [System.IO.Compression.ZipFile]::OpenRead("$PWD\docker.zip")
                Try ***REMOVED***
                    if (-not ($entry = $zip.Entries | Where-Object ***REMOVED*** $_.Name -eq "docker.exe" ***REMOVED***)) ***REMOVED***
                        Throw "Cannot find docker.exe in $url"
                ***REMOVED***
                    [System.IO.Compression.ZipFileExtensions]::ExtractToFile($entry, "$PWD\bundles\docker.exe", $true)
            ***REMOVED***
                Finally ***REMOVED***
                    $zip.Dispose()
            ***REMOVED***
        ***REMOVED***
            Finally ***REMOVED***
                Remove-Item -Force "docker.zip"
        ***REMOVED***
    ***REMOVED***
***REMOVED***

    # Run unit tests
    if ($TestUnit) ***REMOVED*** Run-UnitTests ***REMOVED***

    # Gratuitous ASCII art.
    if ($Daemon -or $Client) ***REMOVED***
        Write-Host
        Write-Host -ForegroundColor Green " ________   ____  __."
        Write-Host -ForegroundColor Green " \_____  \ `|    `|/ _`|"
        Write-Host -ForegroundColor Green " /   `|   \`|      `<"
        Write-Host -ForegroundColor Green " /    `|    \    `|  \"
        Write-Host -ForegroundColor Green " \_______  /____`|__ \"
        Write-Host -ForegroundColor Green "         \/        \/"
        Write-Host
***REMOVED***
***REMOVED***
Catch [Exception] ***REMOVED***
    Write-Host -ForegroundColor Red ("`nERROR: make.ps1 failed:`n$_")

    # More gratuitous ASCII art.
    Write-Host
    Write-Host -ForegroundColor Red  "___________      .__.__             .___"
    Write-Host -ForegroundColor Red  "\_   _____/____  `|__`|  `|   ____   __`| _/"
    Write-Host -ForegroundColor Red  " `|    __) \__  \ `|  `|  `| _/ __ \ / __ `| "
    Write-Host -ForegroundColor Red  " `|     \   / __ \`|  `|  `|_\  ___// /_/ `| "
    Write-Host -ForegroundColor Red  " \___  /  (____  /__`|____/\___  `>____ `| "
    Write-Host -ForegroundColor Red  "     \/        \/             \/     \/ "
    Write-Host

    Throw $_
***REMOVED***
Finally ***REMOVED***
    Pop-Location # As we pushed to the root of the repo as the very first thing
    if ($global:pushed) ***REMOVED*** Pop-Location ***REMOVED***
    Write-Host -ForegroundColor Cyan "INFO: make.ps1 ended at $(Get-Date)"
***REMOVED***
